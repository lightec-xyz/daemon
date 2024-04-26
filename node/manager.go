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
	proofQueue     *ArrayQueue
	pendingQueue   *PendingQueue
	schedule       *Schedule
	fileStore      *FileStorage
	btcClient      *bitcoin.Client
	ethClient      *ethereum.Client
	store          store.IStore
	memory         store.IStore
	btcProofResp   chan *common.ZkProofResponse
	ethProofResp   chan *common.ZkProofResponse
	syncCommitResp chan *common.ZkProofResponse
	lock           sync.Mutex
}

func NewManager(btcClient *bitcoin.Client, ethClient *ethereum.Client, btcProofResp, ethProofResp, syncCommitteeProofResp chan *common.ZkProofResponse,
	store, memory store.IStore, schedule *Schedule, fileStore *FileStorage) (IManager, error) {
	return &manager{
		proofQueue:     NewArrayQueue(),
		pendingQueue:   NewPendingQueue(),
		schedule:       schedule,
		store:          store,
		memory:         memory,
		fileStore:      fileStore,
		btcProofResp:   btcProofResp,
		ethProofResp:   ethProofResp,
		syncCommitResp: syncCommitteeProofResp,
		btcClient:      btcClient,
		ethClient:      ethClient,
	}, nil
}

func (m *manager) Init() error {
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

func (m *manager) ReceiveRequest(requestList []*common.ZkProofRequest) error {
	for _, req := range requestList {
		logger.Info("queue receive gen Proof request:%v %v", req.ReqType.String(), req.Index)
		m.proofQueue.Push(req)
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
	request, ok := m.proofQueue.Pop()
	if !ok {
		logger.Error("should never happen,parse Proof request error")
		return nil, false, fmt.Errorf("parse Proof request error")
	}
	logger.Info("get proof request:%v %v", request.ReqType.String(), request.Index)
	request.StartTime = time.Now()
	m.pendingQueue.Push(request.ZkId, request)
	return request, true, nil
}

func (m *manager) SendProofResponse(responses []*common.ZkProofResponse) error {
	for _, response := range responses {
		chanResponse := m.getChanResponse(response.ZkProofType)
		chanResponse <- response
		logger.Info("send Proof response:%v %v %v", response.ZkProofType.String(), response.Period, response.TxHash)
		proofId := response.Id()
		logger.Info("delete pending request:%v", proofId)
		m.pendingQueue.Delete(proofId)
	}
	return nil
}

// todo
func (m *manager) DistributeRequest() error {
	logger.Debug("proofQueue len:%v", m.proofQueue.Len())
	if m.proofQueue.Len() == 0 {
		time.Sleep(10 * time.Second)
		return nil
	}
	request, ok := m.proofQueue.Pop()
	if !ok {
		logger.Error("should never happen,parse Proof request error")
		time.Sleep(10 * time.Second)
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
		go func(req *common.ZkProofRequest, chaResp chan *common.ZkProofResponse) {
			logger.Debug("worker %v start generate Proof type: %v Index: %v", worker.Id(), req.ReqType.String(), req.Index)
			err := m.fileStore.StoreRequest(req)
			if err != nil {
				logger.Error("store Proof error:%v %v %v", req.ReqType.String(), req.Index, err)
				return
			}
			count := 0
			for {
				if count >= 1 {
					// todo
					logger.Error("gen Proof error:%v %v %v", req.ReqType.String(), req.Index, count)
					//m.proofQueue.Push(request)
					return
				}
				count++
				zkProofResponse, err := WorkerGenProof(worker, req)
				if err != nil {
					logger.Error("worker %v gen Proof error:%v %v %v %v", worker.Id(), req.ReqType.String(), req.Index, count, err)
					continue
				}
				for _, item := range zkProofResponse {
					chaResp <- item
					logger.Debug("complete generate Proof type: %v Index: %v", item.ZkProofType.String(), item.Period)
				}
				return
			}
		}(request, chanResponse)
		return nil
	})
	if err != nil {
		logger.Error("find best worker error:%v", err)
		time.Sleep(10 * time.Second)
		return err
	}
	if !find {
		//logger.Warn(" no find best worker to gen Proof")
		m.proofQueue.Push(request)
		time.Sleep(10 * time.Second)
		return nil
	}

	time.Sleep(10 * time.Second)
	return nil
}

func WorkerGenProof(worker rpc.IWorker, request *common.ZkProofRequest) ([]*common.ZkProofResponse, error) {
	defer worker.DelReqNum()
	var result []*common.ZkProofResponse
	switch request.ReqType {
	case common.DepositTxType:
		var depositRpcRequest rpc.DepositRequest
		err := common.ParseObj(request.Data, &depositRpcRequest)
		if err != nil {
			return nil, fmt.Errorf("not deposit Proof param")
		}
		proofResponse, err := worker.GenDepositProof(depositRpcRequest)
		if err != nil {
			logger.Error("gen deposit Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewZkTxProofResp(request.ReqType, request.TxHash, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)
	case common.VerifyTxType:
		var verifyRpcRequest rpc.VerifyRequest
		err := common.ParseObj(request.Data, &verifyRpcRequest)
		if err != nil {
			return nil, fmt.Errorf("not verify Proof param")
		}
		proofResponse, err := worker.GenVerifyProof(verifyRpcRequest)
		if err != nil {
			logger.Error("gen verify Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewZkTxProofResp(request.ReqType, request.TxHash, proofResponse.Proof, proofResponse.Wit)
		result = append(result, zkbProofResponse)
	case common.TxInEth2:
		// todo
		var txInEth2Req rpc.TxInEth2ProveRequest
		err := common.ParseObj(request.Data, &txInEth2Req)
		if err != nil {
			logger.Error("parse txInEth2 Proof param error:%v", err)
			return nil, fmt.Errorf("not txInEth2 Proof param")
		}
		proofResponse, err := worker.TxInEth2Prove(&txInEth2Req)
		if err != nil {
			logger.Error("gen redeem Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewZkTxProofResp(request.ReqType, request.TxHash, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)
	case common.RedeemTxType:
		// todo
		//var redeemRpcRequest rpc.RedeemRequest
		redeemRpcRequest, ok := request.Data.(rpc.RedeemRequest)
		if !ok {
			logger.Error("parse redeem Proof param error:%v", request.Id())
			return nil, fmt.Errorf("not redeem Proof param")
		}
		proofResponse, err := worker.GenRedeemProof(&redeemRpcRequest)
		if err != nil {
			logger.Error("gen redeem Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewZkTxProofResp(request.ReqType, request.TxHash, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)
	case common.SyncComGenesisType:
		var genesisRpcRequest rpc.SyncCommGenesisRequest
		err := common.ParseObj(request.Data, &genesisRpcRequest)
		if err != nil {
			return nil, fmt.Errorf("not genesis Proof param")
		}
		proofResponse, err := worker.GenSyncCommGenesisProof(genesisRpcRequest)
		if err != nil {
			logger.Error("gen sync comm genesis Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewZkProofResp(request.ReqType, request.Index, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)

	case common.SyncComUnitType:
		var commUnitsRequest rpc.SyncCommUnitsRequest
		err := common.ParseObj(request.Data, &commUnitsRequest)
		if err != nil {
			return nil, fmt.Errorf("not sync comm unit Proof param")
		}
		proofResponse, err := worker.GenSyncCommitUnitProof(commUnitsRequest)
		if err != nil {
			logger.Error("gen sync comm unit Proof error:%v", err)
			return nil, err
		}
		// todo
		zkbProofResponse := NewZkProofResp(request.ReqType, request.Index, proofResponse.Proof, proofResponse.Witness)
		outerProof := NewZkProofResp(common.UnitOuter, request.Index, proofResponse.OuterProof, proofResponse.OuterWitness)
		result = append(result, zkbProofResponse)
		result = append(result, outerProof)
	case common.SyncComRecursiveType:
		var recursiveRequest rpc.SyncCommRecursiveRequest
		err := common.ParseObj(request.Data, &recursiveRequest)
		if err != nil {
			return nil, fmt.Errorf("not sync comm recursive Proof param")
		}
		proofResponse, err := worker.GenSyncCommRecursiveProof(recursiveRequest)
		if err != nil {
			logger.Error("gen sync comm recursive Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewZkProofResp(request.ReqType, request.Index, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)

	case common.BeaconHeaderType:
		// todo
		var blockHeaderRequest rpc.BlockHeaderRequest
		err := common.ParseObj(request.Data, &blockHeaderRequest)
		if err != nil {
			logger.Error("not block header Proof param")
			return nil, fmt.Errorf("not block header Proof param")
		}
		response, err := worker.BlockHeaderProve(&blockHeaderRequest)
		if err != nil {
			logger.Error("gen block header Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewZkProofResp(request.ReqType, request.Index, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)
	case common.BeaconHeaderFinalityType:
		// todo
		var finalityRequest rpc.BlockHeaderFinalityRequest
		err := common.ParseObj(request.Data, &finalityRequest)
		if err != nil {
			return nil, fmt.Errorf("not block header finality Proof param")
		}
		response, err := worker.BlockHeaderFinalityProve(&finalityRequest)
		if err != nil {
			logger.Error("gen block header finality Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewZkProofResp(request.ReqType, request.Index, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)
	default:
		logger.Error("never should happen Proof type:%v", request.ReqType)
		return nil, fmt.Errorf("never should happen Proof type:%v", request.ReqType)

	}

	for _, item := range result {
		logger.Info("send zkProof:%v %v", item.Period, item.ZkProofType.String())
	}
	return result, nil

}

func (m *manager) getChanResponse(reqType common.ZkProofType) chan *common.ZkProofResponse {
	switch reqType {
	case common.DepositTxType, common.VerifyTxType:
		return m.btcProofResp
	case common.RedeemTxType, common.TxInEth2, common.BeaconHeaderType: // todo
		return m.ethProofResp
	case common.SyncComGenesisType, common.SyncComUnitType, common.SyncComRecursiveType, common.BeaconHeaderFinalityType:
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

func (m *manager) CheckPendingRequest() error {
	logger.Debug("check pending request now")
	m.pendingQueue.Iterator(func(request *common.ZkProofRequest) error {
		if request.StartTime.IsZero() {
			logger.Error("request start time is zero")
			return fmt.Errorf("request start time is zero")
		}
		currentTime := time.Now()
		if currentTime.Sub(request.StartTime).Hours() >= 3 { // todo
			logger.Warn("gen proof request timeout:%v %v,add to queue again", request.ReqType.String(), request.Index)
			m.proofQueue.Push(request)
		}
		return nil
	})
	return nil
}

func (m *manager) Close() error {

	return nil

}

func NewZkProofResp(reqType common.ZkProofType, period uint64, proof common.ZkProof, witness []byte) *common.ZkProofResponse {
	return &common.ZkProofResponse{
		ZkProofType: reqType,
		Period:      period,
		Proof:       proof,
		Witness:     witness,
		Status:      common.ProofSuccess,
	}
}

func NewZkTxProofResp(reqType common.ZkProofType, txHash string, proof common.ZkProof, witness []byte) *common.ZkProofResponse {
	return &common.ZkProofResponse{
		ZkProofType: reqType,
		TxHash:      txHash,
		Proof:       proof,
		Witness:     witness,
		Status:      common.ProofSuccess,
	}
}
