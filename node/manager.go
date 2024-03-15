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
		if req.reqType == SyncComUnitType || req.reqType == SyncComRecursiveType {
			// sync commit Proof Has higher priority
			m.proofQueue.PushBack(req)
		} else {
			m.proofQueue.PushFront(req)
		}
	}
	return nil
}

func (m *manager) genProof() error {
	if m.proofQueue.Len() == 0 {
		time.Sleep(1 * time.Second)
		return nil
	}
	element := m.proofQueue.Back()
	request, ok := element.Value.(ZkProofRequest)
	if !ok {
		logger.Error("should never happen,parse Proof request error")
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
		logger.Warn(" no find best worker to gen Proof")
		time.Sleep(10 * time.Second)
		return nil
	}
	m.proofQueue.Remove(element)
	proofSubmitted, err := m.CheckProofStatus(request)
	if err != nil {
		logger.Error("check Proof error:%v", err)
		return err
	}
	if proofSubmitted {
		logger.Info("Proof already submitted:%v", request.String())
		return nil
	}
	logger.Debug("worker %v start generate Proof type: %v Period: %v", worker.Id(), request.reqType.String(), request.period)
	chanResponse := m.getChanResponse(request.reqType)
	go func() {
		err := m.workerGenProof(worker, request, chanResponse)
		if err != nil {
			logger.Error("worker %v gen Proof error:%v %v %v", worker.Id(), request.reqType, request.period, err)
			//  take fail request to queue again
			m.proofQueue.PushBack(request)
			logger.Info("add Proof request type: %v ,Period: %v to queue again", request.reqType.String(), request.period)
			return
		}
	}()

	return nil
}

func (m *manager) workerGenProof(worker rpc.IWorker, request ZkProofRequest, resp chan ZkProofResponse) error {
	worker.AddReqNum()
	defer worker.DelReqNum()
	var zkbProofResponse ZkProofResponse
	switch request.reqType {
	case DepositTxType:
		depositProofParam, ok := request.data.(DepositProofParam)
		if !ok {
			return fmt.Errorf("not deposit Proof param")
		}
		depositRpcRequest := rpc.DepositRequest{
			Version: depositProofParam.Version,
		}
		proofResponse, err := worker.GenDepositProof(depositRpcRequest)
		if err != nil {
			logger.Error("gen deposit Proof error:%v", err)
			return err
		}
		zkbProofResponse = NewZkProofResp(request.reqType, request.period, proofResponse.Body)
	case VerifyTxType:
		verifyProofParam, ok := request.data.(VerifyProofParam)
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
		zkbProofResponse = NewZkProofResp(request.reqType, request.period, proofResponse.Body)
	case RedeemTxType:
		redeemProofParam, ok := request.data.(RedeemProofParam)
		if !ok {
			return fmt.Errorf("not deposit Proof param")
		}
		redeemRpcRequest := rpc.RedeemRequest{
			Version: redeemProofParam.Version,
		}
		proofResponse, err := worker.GenRedeemProof(redeemRpcRequest)
		if err != nil {
			logger.Error("gen redeem Proof error:%v", err)
			return err
		}
		zkbProofResponse = NewZkProofResp(request.reqType, request.period, proofResponse.Body)

	case SyncComGenesisType:
		genesisReq, ok := request.data.(GenesisProofParam)
		if !ok {
			logger.Error("parse sync comm genesis request error")
			return fmt.Errorf("parse sync comm genesis request error")
		}
		genesisRpcRequest := rpc.SyncCommGenesisRequest{
			Version: genesisReq.Version,
		}
		proofResponse, err := worker.GenSyncCommGenesisProof(genesisRpcRequest)
		if err != nil {
			logger.Error("gen sync comm genesis Proof error:%v", err)
			return err
		}
		zkbProofResponse = NewZkProofResp(request.reqType, request.period, proofResponse.Body)

	case SyncComUnitType:
		updateWithVersion, ok := request.data.(UnitProofParam)
		if !ok {
			return fmt.Errorf("parse sync comm unit request error")
		}
		commUnitsRequest := rpc.SyncCommUnitsRequest{
			Version: updateWithVersion.Version,
		}
		proofResponse, err := worker.GenSyncCommitUnitProof(commUnitsRequest)
		if err != nil {
			logger.Error("gen sync comm unit Proof error:%v", err)
			return err
		}
		zkbProofResponse = NewZkProofResp(request.reqType, request.period, proofResponse.Body)

	case SyncComRecursiveType:
		recursiveProofParam, ok := request.data.(RecursiveProofParam)
		if !ok {
			return fmt.Errorf("parse sync comm recursive request error")
		}
		recursiveRequest := rpc.SyncCommRecursiveRequest{
			Version: recursiveProofParam.Version,
		}
		proofResponse, err := worker.GenSyncCommRecursiveProof(recursiveRequest)
		if err != nil {
			logger.Error("gen sync comm recursive Proof error:%v", err)
			return err
		}
		zkbProofResponse = NewZkProofResp(request.reqType, request.period, proofResponse.Body)
	default:
		logger.Error("never should happen Proof type:%v", request.reqType)
		return fmt.Errorf("never should happen Proof type:%v", request.reqType)

	}
	resp <- zkbProofResponse
	return nil

}

func (m *manager) getChanResponse(reqType ZkProofType) chan ZkProofResponse {
	switch reqType {
	case DepositTxType, VerifyTxType:
		return m.btcProofResp
	case RedeemTxType:
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

type Queue struct {
	list     *list.List
	lock     sync.Mutex
	capacity uint64
}

func NewZkProofResp(reqType ZkProofType, period uint64, body []byte) ZkProofResponse {
	return ZkProofResponse{
		ZkProofType: reqType,
		Period:      period,
		Proof:       body,
		Status:      ProofSuccess,
	}
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
