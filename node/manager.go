package node

import (
	"container/list"
	"github.com/lightec-xyz/daemon/logger"
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
	//	submitted, err := m.CheckGenProofStatus(req)
	//	if err != nil {
	//		logger.Error("check proof error:%v", err)
	//		return err
	//	}
	//	if !submitted {
	//		logger.Info("add un gen proof request:%v", req.String())
	//		m.proofQueue.PushBack(req)
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
		logger.Info("queue receive gen proof request:%v", req.reqType)
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
	proofSubmitted, err := m.CheckGenProofStatus(request)
	if err != nil {
		logger.Error("check proof error:%v", err)
		return err
	}
	if proofSubmitted {
		//logger.Info("proof already submitted:%v", request.String())
		return nil
	}
	//logger.Info("start gen proof:%v", request.String())
	go func() {
		proofResponse, err := m.schedule.GenZKProof(worker, request)
		if err != nil {
			//todo add queue again or cli retry ?
			//logger.Error("gen proof error:%v %v", request.TxHash, err)
			return
		}
		logger.Info("worker response proof: %v", proofResponse.String())
		switch proofResponse.zkProofType {
		case DepositTxType:
			m.btcProofResp <- proofResponse
		case RedeemTxType:
			m.ethProofResp <- proofResponse
		case SyncComGenesisType, SyncComUnitType, SyncComRecursiveType:
			m.syncCommitResp <- proofResponse
		default:
			logger.Error("never should happen proof type:%v", proofResponse.zkProofType)
		}
	}()
	return nil
}

func (m *manager) CheckGenProofStatus(request ZkProofRequest) (bool, error) {
	return false, nil
	//if request.ProofType == Deposit {
	//	txId := request.Utxos[0].TxId
	//	exists, err := CheckDepositDestHash(m.store, m.ethClient, txId)
	//	if err != nil {
	//		logger.Error("check deposit proof error:%v", err)
	//		return false, err
	//	}
	//	return exists, nil
	//} else if request.ProofType == Redeem {
	//	exists, err := m.btcClient.CheckTx(request.BtcTxId)
	//	if err != nil {
	//		logger.Error("check btc tx error: %v %v", request.BtcTxId, err)
	//		return false, err
	//	}
	//	return exists, nil
	//} else {
	//	//todo
	//	return false, fmt.Errorf("unknown proof type")
	//}
}

func getChainByProofType(req ProofRequest) ChainType {
	if req.ProofType == Deposit {
		return Bitcoin
	} else if req.ProofType == Redeem {
		return Ethereum
	} else {
		panic("never should happen")
	}
}

func (m *manager) Close() {

}

type Queue struct {
	list *list.List
	mu   sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		list: list.New(),
		mu:   sync.Mutex{},
	}
}

func (sl *Queue) PushBack(value interface{}) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.list.PushBack(value)
}

func (sl *Queue) PushFront(value interface{}) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.list.PushFront(value)
}

func (sl *Queue) Front() *list.Element {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	return sl.list.Front()

}
func (sl *Queue) Back() *list.Element {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	return sl.list.Back()

}

func (sl *Queue) Len() int {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	return sl.list.Len()
}
func (sl *Queue) Remove(e *list.Element) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.list.Remove(e)
}
