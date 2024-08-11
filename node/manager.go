package node

import (
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
	store, memory store.IStore, schedule *Schedule, fileStore *FileStorage, btcGenesisHeight, genesisPeriod, genesisSlot uint64, cache *Cache) (IManager, error) {
	queue := NewArrayQueue()
	state, err := NewState(queue, fileStore, store, cache, prep, btcGenesisHeight, genesisPeriod, genesisSlot, btcClient, ethClient)
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
		logger.Info("load pending request:%v", request.RequestId())
		m.pendingQueue.Push(request.RequestId(), request)
		m.cache.Store(request.RequestId(), nil)
		err = DeletePendingRequest(m.store, request.RequestId())
		if err != nil {
			logger.Error("delete pending request error:%v", err)
		}
	}
	return nil
}
func (m *manager) RemoveProofRequest(id string) error {
	m.proofQueue.Remove(id)
	m.pendingQueue.Delete(id)
	logger.Debug("remove pending request:%v", id)
	return nil
}

func (m *manager) CheckBtcState() error {
	err := m.state.CheckBtcState()
	if err != nil {
		logger.Error("check btc state error:%v", err)
		return err
	}
	return nil
}

func (m *manager) CheckEthState() error {
	err := m.state.CheckEthState()
	if err != nil {
		logger.Error("check eth state error:%v", err)
		return err
	}
	return nil
}

func (m *manager) CheckBeaconState() error {
	err := m.state.CheckBeaconState()
	if err != nil {
		logger.Error("check beacon state error:%v", err)
		return err
	}
	return nil
}

func (m *manager) CheckState() error {
	logger.Debug("check pending request now")
	m.pendingQueue.Iterator(func(request *common.ZkProofRequest) error {
		if request == nil {
			return nil
		}
		if request.StartTime.IsZero() {
			logger.Error("request start time is zero: %v", request.RequestId())
			return nil
		}
		timeout := time.Now().Sub(request.StartTime) >= request.ReqType.Timeout()
		if timeout {
			logger.Debug("%v timeout,add proof queue again", request.RequestId())
			m.pendingQueue.Delete(request.RequestId())
			m.proofQueue.Push(request)
		}
		return nil
	})
	return nil
}

func (m *manager) PendingProofList() []*common.ZkProofRequest {
	return m.proofQueue.List()
}

func (m *manager) ReceiveRequest(requestList []*common.ZkProofRequest) error {
	for _, req := range requestList {
		logger.Info("queue receive gen Proof request:%v", req.RequestId())
		m.proofQueue.Push(req)
		err := m.UpdateProofStatus(req, common.ProofQueued)
		if err != nil {
			logger.Error("update Proof status error:%v %v", req.RequestId(), err)
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
	//m.proofQueue.Iterator(func(index int, value *common.ZkProofRequest) error {
	//	logger.Debug("queryQueueReq: %v", value.RequestId())
	//	return nil
	//})
	var request *common.ZkProofRequest
	var ok bool
	if len(proofTypes) == 0 {
		request, ok = m.proofQueue.Pop()
	} else {
		request, ok = m.proofQueue.PopFn(func(request *common.ZkProofRequest) bool {
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
	}
	if !ok {
		logger.Warn("no find match proof task")
		return nil, false, nil
	}
	exists, err := CheckProof(m.fileStore, request.ReqType, request.Index, request.SIndex, request.TxHash)
	if err != nil {
		logger.Error("check Proof error:%v %v", request.RequestId(), err)
		return nil, false, err
	}
	if exists {
		return nil, false, nil
	}
	logger.Debug("get proof request:%v", request.RequestId())
	request.StartTime = time.Now()
	m.pendingQueue.Push(request.RequestId(), request)
	err = m.UpdateProofStatus(request, common.ProofGenerating)
	if err != nil {
		logger.Error("update Proof status error:%v %v", request.RequestId(), err)
	}
	return request, true, nil
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
		if chanResponse != nil {
			chanResponse <- response
		}
		proofId := response.RespId()
		logger.Info("delete pending request:%v", proofId)
		m.pendingQueue.Delete(proofId)
		err = m.state.CheckProofRequest(response)
		if err != nil {
			logger.Error("check proof request error:%v", err)
			return err
		}
	}
	return nil
}

// todo
func (m *manager) DistributeRequest() error {
	logger.Debug("start distribute request now")
	request, ok, err := m.GetProofRequest(nil)
	waitTime := 3 * time.Second
	if err != nil {
		logger.Error("get Proof request error:%v", err)
		time.Sleep(waitTime)
		return err
	}
	if !ok {
		//logger.Warn("current queue is empty")
		time.Sleep(waitTime)
		return nil
	}
	proofSubmitted, err := m.CheckProofStatus(request)
	if err != nil {
		logger.Error("check Proof error:%v", err)
		m.CacheRequest(request)
		time.Sleep(waitTime)
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
		time.Sleep(waitTime)
		return err
	}
	_, find, err := m.schedule.findBestWorker(func(worker rpc.IWorker) error {
		worker.AddReqNum()
		go func(req *common.ZkProofRequest, chaResp chan *common.ZkProofResponse) {
			defer worker.DelReqNum()
			logger.Debug("worker %v start generate Proof type: %v", worker.Id(), req.RequestId())
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
					logger.Debug("complete generate Proof type: %v", item.RespId())
					if chaResp != nil {
						chaResp <- item
						logger.Debug("chan send -- %v", item.RespId())
					}
					err := StoreZkProof(m.fileStore, item.ZkProofType, item.Index, item.End, item.TxHash, item.Proof, item.Witness)
					if err != nil {
						logger.Error("store Proof error:%v %v", item.RespId(), err)
						return
					}
				}
				m.pendingQueue.Delete(req.RequestId())
				return
			}
		}(request, chanResponse)
		return nil
	})
	if err != nil {
		logger.Error("find best worker error:%v", err)
		m.CacheRequest(request)
		time.Sleep(waitTime)
		return err
	}
	if !find {
		logger.Warn(" no find best worker to gen Proof: %v", request.RequestId())
		m.CacheRequest(request)
		time.Sleep(waitTime)
		return nil
	}
	time.Sleep(waitTime)
	return nil
}

// todo
func (m *manager) CacheRequest(request *common.ZkProofRequest) {
	m.proofQueue.Push(request)
	m.pendingQueue.Delete(request.RequestId())
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
		//logger.Error("never should happen Proof type:%v", reqType.String())
		return nil, nil
	}
}

func (m *manager) CheckProofStatus(request *common.ZkProofRequest) (bool, error) {
	// todo check Proof
	return false, nil
}

func (m *manager) Close() error {
	logger.Debug("manager start  cache cache now ...")
	m.pendingQueue.Iterator(func(value *common.ZkProofRequest) error {
		logger.Debug("write pending request to db :%v", value.RequestId())
		err := WritePendingRequest(m.store, value.RequestId(), value)
		if err != nil {
			logger.Error("write pending request error:%v %v", value.RequestId(), err)
			return err
		}
		return nil
	})
	return nil

}

func (s *manager) UpdateProofStatus(req *common.ZkProofRequest, status common.ProofStatus) error {
	// todo
	if req.ReqType == common.DepositTxType || req.ReqType == common.RedeemTxType {
		err := UpdateProof(s.store, req.TxHash, "", req.ReqType, status)
		if err != nil {
			logger.Error("update Proof status error:%v %v", req.RequestId(), err)
			return err
		}
	}
	return nil
}
