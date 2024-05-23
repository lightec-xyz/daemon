package node

import (
	"context"
	"fmt"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/beacon"
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
	beaconClient   *beacon.Client
	store          store.IStore
	memory         store.IStore
	genesisPeriod  uint64
	state          *State
	btcProofResp   chan *common.ZkProofResponse
	ethProofResp   chan *common.ZkProofResponse
	syncCommitResp chan *common.ZkProofResponse
	lock           sync.Mutex
}

func NewManager(btcClient *bitcoin.Client, ethClient *ethereum.Client, beaconClient *beacon.Client, btcProofResp, ethProofResp, syncCommitteeProofResp chan *common.ZkProofResponse,
	store, memory store.IStore, schedule *Schedule, fileStore *FileStorage, genesisPeriod uint64, state *State) (IManager, error) {
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
		beaconClient:   beaconClient,
		genesisPeriod:  genesisPeriod,
		state:          state,
	}, nil
}

func (m *manager) Init() error {

	return nil
}

func (m *manager) ReceiveRequest(requestList []*common.ZkProofRequest) error {
	for _, req := range requestList {
		logger.Info("queue receive gen Proof request:%v", req.Id())
		m.proofQueue.Push(req)
		err := m.UpdateProofStatus(req, common.ProofQueued)
		if err != nil {
			logger.Error("update Proof status error:%v %v", req.Id(), err)
		}
	}
	return nil
}

func (m *manager) GetProofRequest() (*common.ZkProofRequest, bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	logger.Warn("current proof queue length: %v", m.proofQueue.Len())
	if m.proofQueue.Len() == 0 {
		return nil, false, nil
	}
	m.proofQueue.Iterator(func(index int, value *common.ZkProofRequest) error {
		logger.Debug("queryQueueReq: %v", value.Id())
		return nil
	})
	request, ok := m.proofQueue.Pop()
	if !ok {
		logger.Error("should never happen,parse Proof request error")
		return nil, false, fmt.Errorf("parse Proof request error")
	}
	logger.Debug("get proof request:%v", request.Id())

	request.StartTime = time.Now()
	m.pendingQueue.Push(request.Id(), request)
	err := m.UpdateProofStatus(request, common.ProofGenerating)
	if err != nil {
		logger.Error("update Proof status error:%v %v", request.Id(), err)
	}
	return request, true, nil
}

func (m *manager) UpdateProofStatus(req *common.ZkProofRequest, status common.ProofStatus) error {
	// todo
	if req.ReqType == common.DepositTxType || req.ReqType == common.RedeemTxType {
		err := UpdateProof(m.store, req.TxHash, "", req.ReqType, status)
		if err != nil {
			logger.Error("update Proof status error:%v %v", req.Id(), err)
			return err
		}
	}
	return nil
}

func (m *manager) SendProofResponse(responses []*common.ZkProofResponse) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, response := range responses {
		chanResponse, err := m.getChanResponse(response.ZkProofType)
		if err != nil {
			logger.Error("get chan response error:%v", err)
			return err
		}
		chanResponse <- response
		logger.Info("send Proof response:%v %v %v", response.ZkProofType.String(), response.Period, response.TxHash)
		proofId := response.Id()
		logger.Info("delete pending request:%v", proofId)
		m.pendingQueue.Delete(proofId)
		// todo
		err = m.waitUpdateProofStatus(response)
		if err != nil {
			logger.Error("wait update Proof status error:%v", err)
			return err
		}
	}
	return nil
}

func (m *manager) checkRedeemRequest(resp *common.ZkProofResponse) ([]*common.ZkProofRequest, bool, error) {
	switch resp.ZkProofType {
	case common.TxInEth2:
		request, ok, err := m.GetRedeemRequest(resp.TxHash)
		if err != nil {
			logger.Error("get redeem request error:%v %v", resp.Id())
			return nil, false, err
		}
		return []*common.ZkProofRequest{request}, ok, nil
	case common.BeaconHeaderType:
		txes := m.state.GetTxSlot(resp.Period)
		var result []*common.ZkProofRequest
		for _, tx := range txes {
			request, ok, err := m.GetRedeemRequest(tx)
			if err != nil {
				logger.Error("get redeem request error:%v %v", resp.Id())
				return nil, false, err
			}
			if ok {
				result = append(result, request)
			}
		}
		if len(result) == 0 {
			return nil, false, nil
		}
		return result, true, nil
	case common.BeaconHeaderFinalityType:
		txes := m.state.GetFinalizeSlot(resp.Period)
		var result []*common.ZkProofRequest
		for _, tx := range txes {
			request, ok, err := m.GetRedeemRequest(tx)
			if err != nil {
				logger.Error("get redeem request error:%v %v", resp.Id())
				return nil, false, err
			}
			if ok {
				result = append(result, request)
			}
		}
		if len(result) == 0 {
			return nil, false, nil
		}
		return result, true, nil
	default:
		return nil, false, nil
	}
}

func (m *manager) GetRedeemRequest(txHash string) (*common.ZkProofRequest, bool, error) {
	// todo
	txSlot, ok, err := m.GetSlotByHash(txHash)
	if err != nil {
		logger.Error("get slot by hash error: %v %v", txHash, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	data, ok, err := GetRedeemRequestData(m.fileStore, m.genesisPeriod, txSlot, txHash, m.beaconClient, m.ethClient.Client)
	if err != nil {
		logger.Error("get redeem request data error: %v %v", txHash, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	request := common.NewZkProofRequest(common.RedeemTxType, data, txSlot, txHash)
	return request, true, nil
}

// todo
func (m *manager) DistributeRequest() error {
	request, ok, err := m.GetProofRequest()
	if err != nil {
		logger.Error("get Proof request error:%v", err)
		time.Sleep(10 * time.Second)
		return err
	}
	if !ok {
		//logger.Warn("current queue is empty")
		time.Sleep(10 * time.Second)
		return nil
	}
	proofSubmitted, err := m.CheckProofStatus(request)
	if err != nil {
		logger.Error("check Proof error:%v", err)
		m.CacheRequest(request)
		time.Sleep(10 * time.Second)
		return err
	}
	if proofSubmitted {
		logger.Info("Proof already submitted:%v", request.String())
		return nil
	}
	chanResponse, err := m.getChanResponse(request.ReqType)
	if err != nil {
		logger.Error("get chan response error:%v", err)
		m.CacheRequest(request)
		time.Sleep(10 * time.Second)
		return err
	}
	_, find, err := m.schedule.findBestWorker(func(worker rpc.IWorker) error {
		worker.AddReqNum()
		go func(req *common.ZkProofRequest, chaResp chan *common.ZkProofResponse) {
			logger.Debug("worker %v start generate Proof type: %v", worker.Id(), req.Id())
			err := m.fileStore.StoreRequest(req)
			if err != nil {
				logger.Error("store Proof error:%v %v %v", req.ReqType.String(), req.Index, err)
				return
			}
			defer worker.DelReqNum()
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
					logger.Debug("complete generate Proof type: %v", item.Id())
					chaResp <- item
					logger.Debug("chan send -- %v", item.Id())
				}
				m.pendingQueue.Delete(req.Id())
				return
			}
		}(request, chanResponse)
		return nil
	})
	if err != nil {
		logger.Error("find best worker error:%v", err)
		m.CacheRequest(request)
		time.Sleep(10 * time.Second)
		return err
	}
	if !find {
		logger.Warn(" no find best worker to gen Proof: %v", request.Id())
		m.CacheRequest(request)
		time.Sleep(10 * time.Second)
		return nil
	}
	time.Sleep(10 * time.Second)
	return nil
}

// todo
func (m *manager) CacheRequest(request *common.ZkProofRequest) {
	m.proofQueue.Push(request)
	m.pendingQueue.Delete(request.Id())
}

func WorkerGenProof(worker rpc.IWorker, request *common.ZkProofRequest) ([]*common.ZkProofResponse, error) {
	//defer worker.DelReqNum()
	var result []*common.ZkProofResponse
	switch request.ReqType {
	case common.DepositTxType:
		var depositRpcRequest rpc.DepositRequest
		err := common.ParseObj(request.Data, &depositRpcRequest)
		if err != nil {
			logger.Error("parse deposit Proof param error: %v", request.TxHash)
			return nil, fmt.Errorf("not deposit Proof param %v", request.TxHash)
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
			logger.Error("parse verify Proof param error: %v", request.TxHash)
			return nil, fmt.Errorf("not verify Proof param %v", request.TxHash)
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
		var redeemRpcRequest rpc.RedeemRequest
		err := common.ParseObj(request.Data, &redeemRpcRequest)
		if err != nil {
			logger.Error("parse redeem Proof param error:%v", request.Id())
			return nil, fmt.Errorf("not redeem Proof param")
		}
		proofResponse, err := worker.GenRedeemProof(&redeemRpcRequest)
		if err != nil {
			logger.Error("gen redeem Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.TxHash, proofResponse.Proof, proofResponse.Witness)
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
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.TxHash, response.Proof, response.Witness)
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

func (m *manager) getChanResponse(reqType common.ZkProofType) (chan *common.ZkProofResponse, error) {
	switch reqType {
	case common.DepositTxType, common.VerifyTxType:
		return m.btcProofResp, nil
	case common.RedeemTxType, common.TxInEth2, common.BeaconHeaderType, common.BeaconHeaderFinalityType: // todo
		return m.ethProofResp, nil
	case common.SyncComGenesisType, common.SyncComUnitType, common.UnitOuter, common.SyncComRecursiveType:
		return m.syncCommitResp, nil
	default:
		logger.Error("never should happen Proof type:%v", reqType.String())
		return nil, fmt.Errorf("never should happen Proof type:%v", reqType.String())
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
		if currentTime.Sub(request.StartTime).Minutes() >= 30 { // todo
			logger.Warn("gen proof request timeout:%v %v,add to queue again", request.ReqType.String(), request.Index)
			m.proofQueue.Push(request) // todo
			m.pendingQueue.Delete(request.Id())
			err := m.UpdateProofStatus(request, common.ProofQueued)
			if err != nil {
				logger.Error("update Proof status error:%v %v", request.Id(), err)
			}
		}
		return nil
	})
	return nil
}

func (m *manager) GetSlotByHash(hash string) (uint64, bool, error) {
	txHash := ethCommon.HexToHash(hash)
	receipt, err := m.ethClient.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		logger.Error("get tx receipt error: %v %v", hash, err)
		return 0, false, err
	}
	// todo
	beaconSlot, ok, err := ReadBeaconSlot(m.store, receipt.BlockNumber.Uint64())
	if err != nil {
		logger.Error("get beacon slot error: %v %v", hash, err)
		return 0, false, err
	}
	if !ok {
		return 0, false, nil
	}
	return beaconSlot, true, nil
}

func (m *manager) Close() error {

	return nil

}

// todo ,just temp use ,will remove
func (m *manager) waitUpdateProofStatus(resp *common.ZkProofResponse) error {
	switch resp.ZkProofType {
	case common.TxInEth2, common.BeaconHeaderType, common.BeaconHeaderFinalityType:
		time.Sleep(6 * time.Second)
		return nil
	default:

	}
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

func NewProofResp(reqType common.ZkProofType, period uint64, txHash string, proof common.ZkProof, witness []byte) *common.ZkProofResponse {
	return &common.ZkProofResponse{
		ZkProofType: reqType,
		Period:      period,
		Proof:       proof,
		TxHash:      txHash,
		Witness:     witness,
		Status:      common.ProofSuccess,
	}
}
