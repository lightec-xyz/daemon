package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
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
	err := m.state.RemoveProofRequest(id)
	if err != nil {
		logger.Error("remove proof request error:%v", err)
		return err
	}
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
		timeout := time.Now().Sub(request.StartTime) >= request.ProofType.Timeout()
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
	var err error
	if len(proofTypes) == 0 {
		request, ok = m.proofQueue.Pop()
	} else {
		request, ok, err = m.proofQueue.PopFn(func(request *common.ZkProofRequest) (bool, error) {
			if len(proofTypes) == 0 {
				return true, nil
			}
			for _, req := range proofTypes {
				if request.ProofType == req {
					return true, nil
				}
			}
			return false, nil
		})
	}
	if !ok {
		logger.Warn("no find match proof task")
		return nil, false, nil
	}
	exists, err := CheckProof(m.fileStore, request.ProofType, request.Index, request.SIndex, request.Hash)
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
		chanResponse, err := m.getChanResponse(response.ProofType)
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

func (m *manager) DistributeRequest() error {
	_, find, err := m.proofQueue.PopFn(func(req *common.ZkProofRequest) (bool, error) {
		proofSubmitted, err := m.CheckProofStatus(req)
		if err != nil {
			logger.Error("check Proof error:%v", err)
			return false, err
		}
		if proofSubmitted {
			logger.Info("Proof already submitted:%v", req.String())
			return true, nil
		}
		worker, ok, err := m.schedule.findWorker(req.ProofType)
		if err != nil {
			logger.Error("find worker error:%v", err)
			return false, err
		}
		if !ok {
			return false, fmt.Errorf("no find worker") // skip proofQueue loop
		}
		chanResp, err := m.getChanResponse(req.ProofType)
		if err != nil {
			logger.Error("get chan response error:%v", err)
			return false, err
		}
		worker.AddReqNum()
		m.pendingQueue.Push(req.RequestId(), req)
		go func(req *common.ZkProofRequest, chaResp chan *common.ZkProofResponse) {
			defer func() {
				worker.DelReqNum()
				m.pendingQueue.Delete(req.RequestId())
			}()
			logger.Debug("worker %v start generate Proof type: %v", worker.Id(), req.RequestId())
			err := m.fileStore.StoreRequest(req)
			if err != nil {
				logger.Error("store Proof error:%v %v %v", req.ProofType.String(), req.Index, err)
				return
			}
			count := 0
			for {
				// todo
				if count >= 1 {
					logger.Error("gen Proof error:%v %v %v", req.ProofType.String(), req.Index, count)
					m.proofQueue.Push(req)
					return
				}
				count++
				zkProofResponse, err := WorkerGenProof(worker, req)
				if err != nil {
					logger.Error("worker %v gen Proof error:%v %v %v %v", worker.Id(), req.ProofType.String(), req.Index, count, err)
					continue
				}
				for _, item := range zkProofResponse {
					logger.Debug("complete generate Proof type: %v", item.RespId())
					if chaResp != nil {
						chaResp <- item
						logger.Debug("chan send -- %v", item.RespId())
					}
					err := StoreZkProof(m.fileStore, item.ProofType, item.Index, item.SIndex, item.Hash, item.Proof, item.Witness)
					if err != nil {
						logger.Error("store Proof error:%v %v", item.RespId(), err)
						return
					}
				}
				return
			}
		}(req, chanResp)
		return true, nil
	})
	if err != nil {
		logger.Warn("find worker error:%v", err)
		time.Sleep(10 * time.Second)
		return err
	}
	if !find {
		logger.Warn("no find match proof task")
		time.Sleep(10 * time.Second)
		return nil
	}
	return nil
}

// todo
func (m *manager) CacheRequest(request *common.ZkProofRequest) {
	m.proofQueue.Push(request)
	m.pendingQueue.Delete(request.RequestId())
}

func (m *manager) getChanResponse(reqType common.ZkProofType) (chan *common.ZkProofResponse, error) {
	switch reqType {
	case common.BtcDepositType, common.BtcChangeType, common.BtcBulkType, common.BtcPackedType:
		return m.btcProofResp, nil
	case common.RedeemTxType, common.TxInEth2, common.BeaconHeaderType, common.BeaconHeaderFinalityType: // todo
		return m.ethProofResp, nil
	case common.SyncComGenesisType, common.SyncComUnitType, common.SyncComOuterType, common.SyncComRecursiveType:
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
	if req.ProofType == common.BtcDepositType || req.ProofType == common.RedeemTxType {
		err := UpdateProof(s.store, req.Hash, "", req.ProofType, status)
		if err != nil {
			logger.Error("update Proof status error:%v %v", req.RequestId(), err)
			return err
		}
	}
	return nil
}
