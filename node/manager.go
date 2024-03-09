package node

import (
	"container/list"
	"fmt"
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
	schedule       *Schedule
	btcClient      *bitcoin.Client
	ethClient      *ethereum.Client
	store          store.IStore
	memory         store.IStore
	btcProofResp   chan ZkProofResponse
	ethProofResp   chan ZkProofResponse
	syncCommitResp chan ZkProofResponse
}

func NewManager(cfg NodeConfig, btcProofResp, ethProofResp, syncCommitteeProofResp chan ZkProofResponse, store, memory store.IStore, schedule *Schedule) (*manager, error) {
	btcClient, err := bitcoin.NewClient(cfg.BtcUrl, cfg.BtcUser, cfg.BtcPwd, cfg.BtcNetwork)
	if err != nil {
		logger.Error("new bitcoin rpc client error:%v", err)
		return nil, err
	}
	ethClient, err := ethereum.NewClient(cfg.EthUrl, cfg.ZkBridgeAddr, cfg.ZkBtcAddr)
	if err != nil {
		logger.Error("new ethereum rpc client error:%v", err)
		return nil, err
	}
	return &manager{
		proofQueue:     NewQueue(),
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
	//	logger.Error("read un gen proof error:%v", err)
	//	return err
	//}
	//for _, req := range dbRequests {
	//	submitted, err := m.CheckProofStatus(req)
	//	if err != nil {
	//		logger.Error("check proof error:%v", err)
	//		return err
	//	}
	//	if !submitted {
	//		logger.Info("add un gen proof request:%v", req.String())
	//		m.cacheQueue.PushBack(req)
	//	} else {
	//		err := DeleteUnGenProof(m.store, getChainByProofType(req), req.TxHash)
	//		if err != nil {
	//			logger.Error("delete un gen proof error:%v", err)
	//			return err
	//		}
	//	}
	//}
	return nil
}

func (m *manager) run(requestList []ZkProofRequest) error {
	for _, req := range requestList {
		logger.Info("queue receive gen proof request:%v %v", req.period, req.reqType.String())
		if req.reqType == SyncComUnitType || req.reqType == SyncComRecursiveType {
			// sync commit proof Has higher priority
			m.proofQueue.PushBack(req)
		} else {
			m.proofQueue.PushFront(req)
		}
	}
	return nil
}

func (m *manager) genProof() error {
	if m.proofQueue.Len() == 0 {
		//logger.Debug("no proof need to do,wait now ....")
		time.Sleep(1 * time.Second)
		return nil
	}
	element := m.proofQueue.Back()
	request, ok := element.Value.(ZkProofRequest)
	if !ok {
		logger.Error("should never happen,parse proof request error")
		time.Sleep(1 * time.Second)
		return nil
	}
	worker, find, err := m.schedule.findBestWorker(request.reqType)
	if err != nil {
		logger.Error("find best worker error:%v", err)
		time.Sleep(1 * time.Second)
		return err
	}
	if !find {
		logger.Warn(" no find best worker to gen proof")
		time.Sleep(10 * time.Second)
		return nil
	}
	// todo
	m.proofQueue.Remove(element)
	proofSubmitted, err := m.CheckProofStatus(request)
	if err != nil {
		logger.Error("check proof error:%v", err)
		return err
	}
	if proofSubmitted {
		//logger.Info("proof already submitted:%v", request.String())
		return nil
	}
	//logger.Info("start gen proof:%v", request.String())
	go workerGenProof(worker, request, m.syncCommitResp)

	return nil
}

func workerGenProof(worker rpc.IWorker, request ZkProofRequest, resp chan ZkProofResponse) error {
	worker.AddReqNum()
	defer worker.DelReqNum()
	var zkbProofResponse ZkProofResponse
	switch request.reqType {
	case DepositTxType:
		// todo
		zkbProofResponse.zkProofType = DepositTxType
		proofResponse, err := worker.GenDepositProof(rpc.DepositRequest{})
		if err != nil {
			logger.Error("gen deposit proof error:%v", err)
			return err
		}
		zkbProofResponse.period = request.period
		zkbProofResponse.proof = proofResponse.Body
		zkbProofResponse.Status = ProofSuccess

	case RedeemTxType:
		proofResponse, err := worker.GenRedeemProof(rpc.RedeemRequest{})
		if err != nil {
			// todo
			logger.Error("gen redeem proof error:%v", err)
			return err
		}
		zkbProofResponse.period = request.period
		zkbProofResponse.zkProofType = RedeemTxType
		zkbProofResponse.proof = proofResponse.Body
		zkbProofResponse.Status = ProofSuccess

	case SyncComGenesisType:
		proofResponse, err := worker.GenSyncCommGenesisProof(rpc.SyncCommGenesisRequest{})
		if err != nil {
			//todo
			logger.Error("gen sync comm genesis proof error:%v", err)
			return err
		}
		zkbProofResponse.period = request.period
		zkbProofResponse.zkProofType = SyncComGenesisType
		zkbProofResponse.proof = proofResponse.Body
		zkbProofResponse.Status = ProofSuccess

	case SyncComUnitType:
		proofResponse, err := worker.GenSyncCommitUnitProof(rpc.SyncCommUnitsRequest{})
		if err != nil {
			//todo
			logger.Error("gen sync comm unit proof error:%v", err)
			return err
		}
		zkbProofResponse.period = request.period
		zkbProofResponse.zkProofType = SyncComUnitType
		zkbProofResponse.proof = proofResponse.Body
		zkbProofResponse.Status = ProofSuccess

	case SyncComRecursiveType:
		proofResponse, err := worker.GenSyncCommRecursiveProof(rpc.SyncCommRecursiveRequest{})
		if err != nil {
			// todo
			logger.Error("gen sync comm recursive proof error:%v", err)
			return err
		}
		zkbProofResponse.period = request.period
		zkbProofResponse.zkProofType = SyncComRecursiveType
		zkbProofResponse.proof = proofResponse.Body
		zkbProofResponse.Status = ProofSuccess
	default:
		logger.Error("never should happen proof type:%v", request.reqType)
		return fmt.Errorf("never should happen proof type:%v", request.reqType)

	}
	resp <- zkbProofResponse
	return nil

}

func (m *manager) CheckProofStatus(request ZkProofRequest) (bool, error) {
	// todo check proof
	return false, nil
}

func (m *manager) Close() {

}

type Queue struct {
	list     *list.List
	lock     sync.Mutex
	capacity uint64
}

func NewQueue() *Queue {
	return &Queue{
		list: list.New(),
		lock: sync.Mutex{},
	}
}

func NewQueueWithCapacity(capacity uint64) *Queue {
	return &Queue{
		list:     list.New(),
		lock:     sync.Mutex{},
		capacity: capacity,
	}
}

func (sl *Queue) CanPush() bool {
	if sl.capacity == 0 {
		return true
	}
	sl.lock.Lock()
	defer sl.lock.Unlock()
	return sl.list.Len() < int(sl.capacity)
}

func (sl *Queue) PushBack(value interface{}) {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	sl.list.PushBack(value)
}

func (sl *Queue) PushFront(value interface{}) {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	sl.list.PushFront(value)
}

func (sl *Queue) Front() *list.Element {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	return sl.list.Front()

}
func (sl *Queue) Back() *list.Element {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	return sl.list.Back()

}

func (sl *Queue) Len() int {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	return sl.list.Len()
}
func (sl *Queue) Remove(e *list.Element) {
	sl.lock.Lock()
	defer sl.lock.Unlock()
	sl.list.Remove(e)
}
