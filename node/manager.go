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
	genesisPeriod  uint64
	cache          *Cache
	btcProofResp   chan *common.ZkProofResponse
	ethProofResp   chan *common.ZkProofResponse
	syncCommitResp chan *common.ZkProofResponse
	preparedData   *PreparedData
	lock           sync.Mutex
	state          *State
}

func NewManager(btcClient *bitcoin.Client, ethClient *ethereum.Client, prep *PreparedData, btcProofResp, ethProofResp, syncCommitteeProofResp chan *common.ZkProofResponse,
	store, memory store.IStore, schedule *Schedule, fileStore *FileStorage, genesisPeriod uint64, cache *Cache) (IManager, error) {
	queue := NewArrayQueue()
	state, err := NewState(queue, fileStore, store, cache, prep, genesisPeriod, genesisPeriod, btcClient, ethClient)
	if err != nil {
		logger.Error("new state error:%v", err)
		return nil, err
	}
	return &manager{
		proofQueue:     queue,
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
		genesisPeriod:  genesisPeriod,
		cache:          cache,
		state:          state,
	}, nil
}

func (m *manager) Init() error {
	logger.Debug("manger load db cache now ...")
	allPendingRequests, err := ReadAllPendingRequests(m.store)
	if err != nil {
		logger.Error("read all pending requests error: %v", err)
		return err
	}
	for _, request := range allPendingRequests {
		logger.Info("load pending request:%v", request.Id())
		m.pendingQueue.Push(request.Id(), request)
		m.cache.Store(request.Id(), nil)
		err = DeletePendingRequest(m.store, request.Id())
		if err != nil {
			logger.Error("delete pending request error:%v", err)
		}
	}
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

func (m *manager) GetProofRequest(proofTypes []common.ZkProofType) (*common.ZkProofRequest, bool, error) {
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

	request, ok := m.proofQueue.PopFn(func(request *common.ZkProofRequest) bool {
		if len(proofTypes) == 0 {
			return true
		}
		for _, req := range proofTypes {
			if request.ReqType == req {
				return true
			}
		}
		return false
	})
	if !ok {
		logger.Warn("no find match proof task")
		return nil, false, nil
	}
	// todo
	exists, err := CheckProof(m.fileStore, request.ReqType, request.Index, 0, request.TxHash)
	if err != nil {
		logger.Error("check Proof error:%v %v", request.Id(), err)
		return nil, false, err
	}
	if exists {
		return nil, false, nil
	}

	logger.Debug("get proof request:%v", request.Id())
	request.StartTime = time.Now()
	m.pendingQueue.Push(request.Id(), request)
	err = m.UpdateProofStatus(request, common.ProofGenerating)
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
	m.lock.Lock() // todo
	defer m.lock.Unlock()
	for _, response := range responses {
		chanResponse, err := m.getChanResponse(response.ZkProofType)
		if err != nil {
			logger.Error("get chan response error:%v", err)
			return err
		}
		chanResponse <- response
		proofId := response.Id()
		logger.Info("delete pending request:%v", proofId)
		m.pendingQueue.Delete(proofId)
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
			logger.Error("get redeem request error:%v %v", resp.Id(), err)
			return nil, false, err
		}
		return []*common.ZkProofRequest{request}, ok, nil
	case common.BeaconHeaderType:
		txes, err := ReadAllTxBySlot(m.store, resp.Index)
		if err != nil {
			logger.Error("get redeem request error:%v %v", resp.Id(), err)
			return nil, false, err
		}
		var result []*common.ZkProofRequest
		for _, tx := range txes {
			request, ok, err := m.GetRedeemRequest(tx.TxHash)
			if err != nil {
				logger.Error("get redeem request error:%v %v", resp.Id(), err)
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
		txes, err := ReadAllTxByFinalizedSlot(m.store, resp.Index)
		if err != nil {
			logger.Error("get redeem request error:%v %v", resp.Id(), err)
			return nil, false, err
		}
		var result []*common.ZkProofRequest
		for _, tx := range txes {
			request, ok, err := m.GetRedeemRequest(tx.TxHash)
			if err != nil {
				logger.Error("get redeem request error:%v %v", resp.Id(), err)
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
	data, ok, err := m.preparedData.GetRedeemRequestData(m.genesisPeriod, txSlot, txHash)
	if err != nil {
		logger.Error("get redeem request data error: %v %v", txHash, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	request := common.NewZkProofRequest(common.RedeemTxType, data, txSlot, 0, txHash)
	return request, true, nil
}

// todo
func (m *manager) DistributeRequest() error {
	request, ok, err := m.GetProofRequest(nil)
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
			defer worker.DelReqNum()
			logger.Debug("worker %v start generate Proof type: %v", worker.Id(), req.Id())
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

func (m *manager) getChanResponse(reqType common.ZkProofType) (chan *common.ZkProofResponse, error) {
	switch reqType {
	case common.DepositTxType, common.VerifyTxType, common.BtcBulkType, common.BtcPackedType, common.BtcWrapType:
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

func (m *manager) CheckState() error {
	logger.Debug("check pending request now")
	m.pendingQueue.Iterator(func(request *common.ZkProofRequest) error {
		timout, err := m.checkRequestTimeout(request)
		if err != nil {
			logger.Error("check pending request error:%v", err)
			return err
		}
		if timout {
			logger.Debug("%v timeout,add proof queue again", request.Id())
			m.pendingQueue.Delete(request.Id())
			m.proofQueue.Push(request)
		}
		return nil
	})
	return nil
}
func (m *manager) checkRequestTimeout(request *common.ZkProofRequest) (bool, error) {
	if request == nil {
		return false, fmt.Errorf("request is nil")
	}
	if request.StartTime.IsZero() {
		logger.Error("request start time is zero: %v", request.Id())
		return false, fmt.Errorf("request start time is zero: %v", request.Id())
	}
	isTimeout := false
	currentTime := time.Now()
	switch request.ReqType {
	case common.SyncComUnitType:
		if currentTime.Sub(request.StartTime).Hours() >= 1.3 { // todo
			isTimeout = true
		}
	default:
		if currentTime.Sub(request.StartTime).Minutes() >= 30 { // todo
			isTimeout = true
		}
	}
	return isTimeout, nil
}

func (m *manager) GetSlotByHash(hash string) (uint64, bool, error) {
	// todo
	dbTx, err := ReadDbTx(m.store, hash)
	if err != nil {
		logger.Error("get tx receipt error: %v %v", hash, err)
		return 0, false, err
	}
	beaconSlot, ok, err := ReadBeaconSlot(m.store, dbTx.Height)
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
	logger.Debug("manager start  cache cache now ...")
	m.pendingQueue.Iterator(func(value *common.ZkProofRequest) error {
		logger.Debug("write pending request to db :%v", value.Id())
		err := WritePendingRequest(m.store, value.Id(), value)
		if err != nil {
			logger.Error("write pending request error:%v %v", value.Id(), err)
			return err
		}
		return nil
	})
	return nil

}

// todo
func (m *manager) waitUpdateProofStatus(resp *common.ZkProofResponse) error {
	switch resp.ZkProofType {
	case common.TxInEth2, common.BeaconHeaderType, common.BeaconHeaderFinalityType:
		requests, ok, err := m.checkRedeemRequest(resp)
		if err != nil {
			logger.Error("check redeem request error:%v %v", resp.Id(), err)
			return err
		}
		if !ok {
			return nil
		}
		for _, req := range requests {
			if !m.cache.Check(req.Id()) {
				logger.Debug("add redeem request:%v to queue", req.Id())
				m.cache.Store(req.Id(), nil)
				m.proofQueue.Push(req)
				err := m.UpdateProofStatus(req, common.ProofQueued)
				if err != nil {
					logger.Error("update Proof status error:%v %v", req.Id(), err)
				}
			}
		}
		return nil
	default:
	}
	return nil
}

func (m *manager) CheckProofState() error {
	err := m.state.CheckBtcState()
	if err != nil {
		logger.Error("check btc state error:%v", err)
		return err
	}
	err = m.state.CheckEthState()
	if err != nil {
		logger.Error("check eth state error:%v", err)
		return err
	}
	err = m.state.CheckBeaconState()
	if err != nil {
		logger.Error("check beacon state error:%v", err)
		return err
	}
	return nil
}
