package node

import (
	"fmt"
	"time"

	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
)

type manager struct {
	txProofQueue   *Queue
	schedule       *Schedule
	btcClient      *bitcoin.Client
	ethClient      *ethereum.Client
	store          store.IStore
	memory         store.IStore
	btcProofResp   chan ZkProofResponse
	ethProofResp   chan ZkProofResponse
	syncCommitResp chan ZkProofResponse
}

func NewManager(btcClient *bitcoin.Client, ethClient *ethereum.Client, btcProofResp, ethProofResp, syncCommitteeProofResp chan ZkProofResponse, store, memory store.IStore, schedule *Schedule) (*manager, error) {
	return &manager{
		txProofQueue:   NewQueue(),
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

func (m *manager) run(requestList []ZkProofRequest) error {
	for _, req := range requestList {
		logger.Info("queue receive gen Proof request:%v %v", req.reqType.String(), req.period)
		// Todo queue need to sort by req weight ?
		if req.reqType == SyncComGenesisType || req.reqType == SyncComRecursiveType {
			m.txProofQueue.PushBack(req)
		} else {
			m.txProofQueue.PushFront(req)
		}
	}
	return nil
}

func (m *manager) genProof() error {
	if m.txProofQueue.Len() == 0 {
		time.Sleep(2 * time.Second)
		return nil
	}
	element := m.txProofQueue.Back()
	request, ok := element.Value.(ZkProofRequest)
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
	chanResponse := m.getChanResponse(request.reqType)
	_, find, err := m.schedule.findBestWorker(func(worker rpc.IWorker) error {
		worker.AddReqNum()
		m.txProofQueue.Remove(element)
		go func(req ZkProofRequest) {
			logger.Debug("worker %v start generate Proof type: %v Period: %v", worker.Id(), req.reqType.String(), req.period)
			err := m.workerGenProof(worker, req, chanResponse)
			if err != nil {
				logger.Error("worker %v gen Proof error:%v %v %v", worker.Id(), req.reqType.String(), req.period, err)
				//  take fail request to queue again
				m.txProofQueue.PushBack(request)
				logger.Info("add Proof request type: %v ,Period: %v to queue again", req.reqType.String(), req.period)
				return
			}
			logger.Info("complete generate Proof type: %v Period: %v", req.reqType.String(), req.period)
		}(request)
		return nil
	})
	if err != nil {
		logger.Error("find best worker error:%v", err)
		time.Sleep(1 * time.Second)
		return err
	}
	if !find {
		logger.Warn(" no find best worker to gen Proof")
		time.Sleep(10 * time.Second)
		return nil
	}
	time.Sleep(2 * time.Second)
	return nil
}

func (m *manager) workerGenProof(worker rpc.IWorker, request ZkProofRequest, resp chan ZkProofResponse) error {
	defer worker.DelReqNum()
	var zkbProofResponse ZkProofResponse
	switch request.reqType {
	case DepositTxType:
		depositParam, ok := request.data.(DepositProofParam)
		if !ok {
			return fmt.Errorf("not deposit Proof param")
		}
		depositRpcRequest := rpc.DepositRequest{
			Version:   depositParam.Version,
			TxHash:    request.TxHash,
			BlockHash: depositParam.BlockHash,
		}
		proofResponse, err := worker.GenDepositProof(depositRpcRequest)
		if err != nil {
			logger.Error("gen deposit Proof error:%v", err)
			return err
		}
		zkbProofResponse = NewZkTxProofResp(request.reqType, proofResponse.ProofStr, request.TxHash, proofResponse.Proof, proofResponse.Witness)
	case VerifyTxType:
		verifyProofParam, ok := request.data.(*VerifyProofParam)
		if !ok {
			return fmt.Errorf("not deposit Proof param")
		}
		verifyRpcRequest := rpc.VerifyRequest{
			Version: verifyProofParam.Version,
		}
		proofResponse, err := worker.GenVerifyProof(verifyRpcRequest)
		if err != nil {
			logger.Error("gen verify Proof error:%v", err)
			return err
		}
		zkbProofResponse = NewZkProofResp(request.reqType, request.period, proofResponse.Proof, nil)

	case TxInEth2:
		redeemParam, ok := request.data.(TxInEth2Param)
		if !ok {
			return fmt.Errorf("not txInEth2 Proof param")
		}
		txInEth2Req := &rpc.TxInEth2ProveReq{
			Version: redeemParam.Version,
			TxHash:  request.TxHash,
			TxData:  redeemParam.TxData,
		}
		proofResponse, err := worker.TxInEth2Prove(txInEth2Req)
		if err != nil {
			logger.Error("gen redeem Proof error:%v", err)
			return err
		}
		zkbProofResponse = NewZkTxProofResp(request.reqType, "", request.TxHash, proofResponse.Proof, proofResponse.Witness)

	case RedeemTxType:
		redeemParam, ok := request.data.(*RedeemProofParam)
		if !ok {
			return fmt.Errorf("not deposit Proof param")
		}
		redeemRpcRequest := rpc.RedeemRequest{
			Version: redeemParam.Version,
			TxHash:  request.TxHash,
			TxData:  redeemParam.TxData,
		}
		proofResponse, err := worker.GenRedeemProof(redeemRpcRequest)
		if err != nil {
			logger.Error("gen redeem Proof error:%v", err)
			return err
		}
		zkbProofResponse = NewZkTxProofResp(request.reqType, "", request.TxHash, proofResponse.Proof, proofResponse.Witness)

	case SyncComGenesisType:
		genesisReq, ok := request.data.(*GenesisProofParam)
		if !ok {
			logger.Error("parse sync comm genesis request error")
			return fmt.Errorf("parse sync comm genesis request error")
		}
		genesisRpcRequest := rpc.SyncCommGenesisRequest{
			Version:       genesisReq.Version,
			Period:        request.period,
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
			return err
		}
		zkbProofResponse = NewZkProofResp(request.reqType, request.period, proofResponse.Proof, proofResponse.Witness)

	case SyncComUnitType:
		unitParam, ok := request.data.(*UnitProofParam)
		if !ok {
			return fmt.Errorf("parse sync comm unit request error")
		}
		commUnitsRequest := rpc.SyncCommUnitsRequest{
			Version:                 unitParam.Version,
			Period:                  request.period,
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
			return err
		}
		zkbProofResponse = NewZkProofResp(request.reqType, request.period, proofResponse.Proof, proofResponse.Witness)

	case SyncComRecursiveType:
		recursiveParam, ok := request.data.(*RecursiveProofParam)
		if !ok {
			return fmt.Errorf("parse sync comm recursive request error")
		}
		recursiveRequest := rpc.SyncCommRecursiveRequest{
			Version:       recursiveParam.Version,
			Period:        request.period,
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
			return err
		}
		zkbProofResponse = NewZkProofResp(request.reqType, request.period, proofResponse.Proof, proofResponse.Witness)
	default:
		logger.Error("never should happen Proof type:%v", request.reqType)
		return fmt.Errorf("never should happen Proof type:%v", request.reqType)

	}
	resp <- zkbProofResponse
	logger.Info("send zkProof:%v %v", zkbProofResponse.Period, zkbProofResponse.ZkProofType.String())
	return nil

}

func (m *manager) getChanResponse(reqType ZkProofType) chan ZkProofResponse {
	switch reqType {
	case DepositTxType, VerifyTxType:
		return m.btcProofResp
	case RedeemTxType, TxInEth2:
		return m.ethProofResp
	case SyncComGenesisType, SyncComUnitType, SyncComRecursiveType:
		return m.syncCommitResp
	default:
		logger.Error("never should happen Proof type:%v", reqType)
		return nil
	}
}

func (m *manager) CheckProofStatus(request ZkProofRequest) (bool, error) {
	// todo check Proof
	return false, nil
}

func (m *manager) Close() {

}

func NewZkProofResp(reqType ZkProofType, period uint64, proof common.ZkProof, witness []byte) ZkProofResponse {
	return ZkProofResponse{
		ZkProofType: reqType,
		Period:      period,
		Proof:       proof,
		Witness:     witness,
		Status:      ProofSuccess,
	}
}

func NewZkTxProofResp(reqType ZkProofType, txHash, proofStr string, proof common.ZkProof, witness []byte) ZkProofResponse {
	return ZkProofResponse{
		ZkProofType: reqType,
		TxHash:      txHash,
		Proof:       proof,
		ProofStr:    proofStr,
		Witness:     witness,
		Status:      ProofSuccess,
	}
}
