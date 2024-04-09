package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	"sync"
	"time"
)

type manager struct {
	proofQueue     *Queue
	pendingQueue   *Queue
	schedule       *Schedule
	btcClient      *bitcoin.Client
	ethClient      *ethereum.Client
	store          store.IStore
	memory         store.IStore
	btcProofResp   chan common.ZkProofResponse
	ethProofResp   chan common.ZkProofResponse
	syncCommitResp chan common.ZkProofResponse
	lock           sync.Mutex
}

func NewManager(btcClient *bitcoin.Client, ethClient *ethereum.Client, btcProofResp, ethProofResp, syncCommitteeProofResp chan common.ZkProofResponse, store, memory store.IStore, schedule *Schedule) (*manager, error) {
	return &manager{
		proofQueue:     NewQueue(),
		pendingQueue:   NewQueue(),
		schedule:       schedule,
		store:          store,
		memory:         memory,
		btcProofResp:   btcProofResp,
		ethProofResp:   ethProofResp,
		syncCommitResp: syncCommitteeProofResp,
		btcClient:      btcClient,
		ethClient:      ethClient,
	}, nil
}

func (m *manager) init() error {
	//dbRequests, err := ReadAllUnGenProof(m.store)
	//if err != nil {
	//	logger.Error("read un gen Proof error:%v", err)
	//	return err
	//}
	//for _, req := range dbRequests {
	//	submitted, err := m.CheckProofStatus(req)
	//	if err != nil {
	//		logger.Error("check Proof error:%v", err)
	//		return err
	//	}
	//	if !submitted {
	//		logger.Info("add un gen Proof request:%v", req.FilterLogs())
	//		m.cacheQueue.PushBack(req)
	//	} else {
	//		err := DeleteUnGenProof(m.store, getChainByProofType(req), req.TxHash)
	//		if err != nil {
	//			logger.Error("delete un gen Proof error:%v", err)
	//			return err
	//		}
	//	}
	//}
	return nil
}

func (m *manager) run(requestList []*common.ZkProofRequest) error {
	for _, req := range requestList {
		logger.Info("queue receive gen Proof request:%v %v", req.ReqType.String(), req.Period)
		// Todo queue need to sort by req weight ?
		if req.ReqType == common.SyncComGenesisType || req.ReqType == common.SyncComRecursiveType {
			m.proofQueue.PushBack(req)
		} else {
			m.proofQueue.PushFront(req)
		}
	}
	return nil
}

func (m *manager) GetProofRequest() (*common.ZkProofRequest, bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.proofQueue.Len() == 0 {
		logger.Warn("current queue is empty")
		return nil, false, nil
	}
	//todo
	element := m.proofQueue.Back()
	request, ok := element.Value.(*common.ZkProofRequest)
	if !ok {
		logger.Error("should never happen,parse Proof request error")
		return nil, false, fmt.Errorf("parse Proof request error")
	}
	// todo
	m.proofQueue.Remove(element)
	logger.Info("get proof request:%v %v", request.ReqType.String(), request.Period)
	m.pendingQueue.PushBack(request)
	return request, true, nil
}

func (m *manager) SendProofResponse(response common.ZkProofResponse) error {
	chanResponse := m.getChanResponse(response.ZkProofType)
	chanResponse <- response
	logger.Info("send Proof response:%v %v", response.ZkProofType.String(), response.Period)
	return nil
}

func (m *manager) genProof() error {
	if m.proofQueue.Len() == 0 {
		time.Sleep(2 * time.Second)
		return nil
	}
	element := m.proofQueue.Back()
	request, ok := element.Value.(*common.ZkProofRequest)
	if !ok {
		logger.Error("should never happen,parse Proof request error")
		time.Sleep(5 * time.Second)
		return nil
	}
	proofSubmitted, err := m.CheckProofStatus(request)
	if err != nil {
		logger.Error("check Proof error:%v", err)
		return err
	}
	if proofSubmitted {
		logger.Info("Proof already submitted:%v", request.String())
		return nil
	}
	chanResponse := m.getChanResponse(request.ReqType)
	_, find, err := m.schedule.findBestWorker(func(worker rpc.IWorker) error {
		worker.AddReqNum()
		m.proofQueue.Remove(element)
		go func(req *common.ZkProofRequest) {
			logger.Debug("worker %v start generate Proof type: %v Period: %v", worker.Id(), req.ReqType.String(), req.Period)
			zkProofResponse, err := WorkerGenProof(worker, req)
			if err != nil {
				logger.Error("worker %v gen Proof error:%v %v %v", worker.Id(), req.ReqType.String(), req.Period, err)
				//  take fail request to queue again
				m.proofQueue.PushBack(request)
				logger.Info("add Proof request type: %v ,Period: %v to queue again", req.ReqType.String(), req.Period)
				return
			}
			chanResponse <- zkProofResponse
			logger.Info("complete generate Proof type: %v Period: %v", req.ReqType.String(), req.Period)
		}(request)
		return nil
	})
	if err != nil {
		logger.Error("find best worker error:%v", err)
		time.Sleep(1 * time.Second)
		return err
	}
	if !find {
		//logger.Warn(" no find best worker to gen Proof")
		time.Sleep(10 * time.Second)
		return nil
	}
	time.Sleep(2 * time.Second)
	return nil
}

func WorkerGenProof(worker rpc.IWorker, request *common.ZkProofRequest) (common.ZkProofResponse, error) {
	defer worker.DelReqNum()
	var zkbProofResponse common.ZkProofResponse
	switch request.ReqType {
	case common.DepositTxType:
		var depositParam DepositProofParam
		err := ParseObj(request.Data, &depositParam)
		if err != nil {
			return zkbProofResponse, fmt.Errorf("not deposit Proof param")
		}
		depositRpcRequest := rpc.DepositRequest{
			Version:   depositParam.Version,
			TxHash:    request.TxHash,
			BlockHash: depositParam.BlockHash,
		}
		proofResponse, err := worker.GenDepositProof(depositRpcRequest)
		if err != nil {
			logger.Error("gen deposit Proof error:%v", err)
			return zkbProofResponse, err
		}
		zkbProofResponse = NewZkTxProofResp(request.ReqType, request.TxHash, proofResponse.Proof, proofResponse.Witness)
	case common.VerifyTxType:
		var verifyProofParam VerifyProofParam
		err := ParseObj(request.Data, &verifyProofParam)
		if err != nil {
			return zkbProofResponse, fmt.Errorf("not verify Proof param")
		}
		verifyRpcRequest := rpc.VerifyRequest{
			Version:   verifyProofParam.Version,
			TxHash:    verifyProofParam.TxHash,
			BlockHash: verifyProofParam.BlockHash,
		}
		proofResponse, err := worker.GenVerifyProof(verifyRpcRequest)
		if err != nil {
			logger.Error("gen verify Proof error:%v", err)
			return zkbProofResponse, err
		}
		zkbProofResponse = NewZkTxProofResp(request.ReqType, request.TxHash, proofResponse.Proof, proofResponse.Wit)

	case common.TxInEth2:
		var redeemParam RedeemProofParam
		err := ParseObj(request.Data, &redeemParam)
		if err != nil {
			return zkbProofResponse, fmt.Errorf("not txInEth2 Proof param")
		}
		txInEth2Req := &rpc.TxInEth2ProveReq{
			Version: redeemParam.Version,
			TxHash:  request.TxHash,
			TxData:  redeemParam.TxData,
		}
		proofResponse, err := worker.TxInEth2Prove(txInEth2Req)
		if err != nil {
			logger.Error("gen redeem Proof error:%v", err)
			return zkbProofResponse, err
		}
		zkbProofResponse = NewZkTxProofResp(request.ReqType, request.TxHash, proofResponse.Proof, proofResponse.Witness)

	case common.RedeemTxType:
		var redeemParam RedeemProofParam
		err := ParseObj(request.Data, &redeemParam)
		if err != nil {
			return zkbProofResponse, fmt.Errorf("not redeem Proof param")
		}
		redeemRpcRequest := rpc.RedeemRequest{
			Version: redeemParam.Version,
			TxHash:  request.TxHash,
			TxData:  redeemParam.TxData,
		}
		proofResponse, err := worker.GenRedeemProof(redeemRpcRequest)
		if err != nil {
			logger.Error("gen redeem Proof error:%v", err)
			return zkbProofResponse, err
		}
		zkbProofResponse = NewZkTxProofResp(request.ReqType, request.TxHash, proofResponse.Proof, proofResponse.Witness)

	case common.SyncComGenesisType:
		var genesisReq GenesisProofParam
		err := ParseObj(request.Data, &genesisReq)
		if err != nil {
			return zkbProofResponse, fmt.Errorf("not genesis Proof param")
		}
		genesisRpcRequest := rpc.SyncCommGenesisRequest{
			Version:       genesisReq.Version,
			Period:        request.Period,
			FirstProof:    genesisReq.FirstProof,
			FirstWitness:  genesisReq.FirstWitness,
			SecondProof:   genesisReq.SecondProof,
			SecondWitness: genesisReq.SecondWitness,
			GenesisID:     genesisReq.GenesisId,
			FirstID:       genesisReq.FirstId,
			SecondID:      genesisReq.SecondId,
			RecursiveFp:   genesisReq.RecursiveFp,
		}
		proofResponse, err := worker.GenSyncCommGenesisProof(genesisRpcRequest)
		if err != nil {
			logger.Error("gen sync comm genesis Proof error:%v", err)
			return zkbProofResponse, err
		}
		zkbProofResponse = NewZkProofResp(request.ReqType, request.Period, proofResponse.Proof, proofResponse.Witness)

	case common.SyncComUnitType:
		var unitParam UnitProofParam
		err := ParseObj(request.Data, &unitParam)
		if err != nil {
			return zkbProofResponse, fmt.Errorf("not sync comm unit Proof param")
		}
		commUnitsRequest := rpc.SyncCommUnitsRequest{
			Version:                 unitParam.Version,
			Period:                  request.Period,
			AttestedHeader:          unitParam.AttestedHeader,
			CurrentSyncCommittee:    unitParam.CurrentSyncCommittee,
			SyncAggregate:           unitParam.SyncAggregate,
			NextSyncCommittee:       unitParam.NextSyncCommittee,
			NextSyncCommitteeBranch: unitParam.NextSyncCommitteeBranch,
			FinalizedHeader:         unitParam.FinalizedHeader,
			FinalityBranch:          unitParam.FinalityBranch,
			SignatureSlot:           unitParam.SignatureSlot,
		}
		proofResponse, err := worker.GenSyncCommitUnitProof(commUnitsRequest)
		if err != nil {
			logger.Error("gen sync comm unit Proof error:%v", err)
			return zkbProofResponse, err
		}
		zkbProofResponse = NewZkProofResp(request.ReqType, request.Period, proofResponse.Proof, proofResponse.Witness)

	case common.SyncComRecursiveType:
		var recursiveParam RecursiveProofParam
		err := ParseObj(request.Data, &recursiveParam)
		if err != nil {
			return zkbProofResponse, fmt.Errorf("not sync comm recursive Proof param")
		}
		recursiveRequest := rpc.SyncCommRecursiveRequest{
			Version:       recursiveParam.Version,
			Period:        request.Period,
			Choice:        recursiveParam.Choice,
			FirstProof:    recursiveParam.FirstProof,
			FirstWitness:  recursiveParam.FirstWitness,
			SecondProof:   recursiveParam.SecondProof,
			SecondWitness: recursiveParam.SecondWitness,
			BeginId:       recursiveParam.BeginId,
			RelayId:       recursiveParam.RelayId,
			EndId:         recursiveParam.EndId,
			RecursiveFp:   recursiveParam.RecursiveFp,
		}
		proofResponse, err := worker.GenSyncCommRecursiveProof(recursiveRequest)
		if err != nil {
			logger.Error("gen sync comm recursive Proof error:%v", err)
			return zkbProofResponse, err
		}
		zkbProofResponse = NewZkProofResp(request.ReqType, request.Period, proofResponse.Proof, proofResponse.Witness)
	default:
		logger.Error("never should happen Proof type:%v", request.ReqType)
		return zkbProofResponse, fmt.Errorf("never should happen Proof type:%v", request.ReqType)

	}

	logger.Info("send zkProof:%v %v", zkbProofResponse.Period, zkbProofResponse.ZkProofType.String())
	return zkbProofResponse, nil

}

func (m *manager) getChanResponse(reqType common.ZkProofType) chan common.ZkProofResponse {
	switch reqType {
	case common.DepositTxType, common.VerifyTxType:
		return m.btcProofResp
	case common.RedeemTxType, common.TxInEth2: // todo
		return m.ethProofResp
	case common.SyncComGenesisType, common.SyncComUnitType, common.SyncComRecursiveType:
		return m.syncCommitResp
	default:
		logger.Error("never should happen Proof type:%v", reqType)
		return nil
	}
}

func (m *manager) CheckProofStatus(request *common.ZkProofRequest) (bool, error) {
	// todo check Proof
	return false, nil
}

func (m *manager) Close() {

}

func NewZkProofResp(reqType common.ZkProofType, period uint64, proof common.ZkProof, witness []byte) common.ZkProofResponse {
	return common.ZkProofResponse{
		ZkProofType: reqType,
		Period:      period,
		Proof:       proof,
		Witness:     witness,
		Status:      common.ProofSuccess,
	}
}

func NewZkTxProofResp(reqType common.ZkProofType, txHash string, proof common.ZkProof, witness []byte) common.ZkProofResponse {
	return common.ZkProofResponse{
		ZkProofType: reqType,
		TxHash:      txHash,
		Proof:       proof,
		Witness:     witness,
		Status:      common.ProofSuccess,
	}
}
