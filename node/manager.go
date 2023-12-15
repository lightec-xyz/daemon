package node

import (
	"container/list"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/store"
	"sync"
	"time"
)

type Manager struct {
	proofQueue   *SafeList
	schedule     *Schedule
	store        store.IStore
	memory       store.IStore
	proofRequest chan []ProofRequest
	btcProofResp chan ProofResponse
	ethProofResp chan ProofResponse
	exit         chan struct{}
}

func NewManager(proofRequest chan []ProofRequest, btcProofResp, ethProofResp chan ProofResponse, store, memory store.IStore, schedule *Schedule) *Manager {
	return &Manager{
		proofQueue:   NewSafeList(),
		schedule:     schedule,
		store:        store,
		memory:       memory,
		proofRequest: proofRequest,
		btcProofResp: btcProofResp,
		ethProofResp: ethProofResp,
		exit:         make(chan struct{}, 1),
	}
}

func (m *Manager) run() {
	logger.Info("run manager proof queue")
	for {
		select {
		case requestList := <-m.proofRequest:
			for _, req := range requestList {
				logger.Info("queue receive gen proof request:%v", req.String())
				m.proofQueue.PushBack(req)
			}
		case <-m.exit:
			logger.Debug("manager proof queue exit ...")
			return

		}
	}
}

func (m *Manager) genProof() {
	//todo
	logger.Info("start gen proof goroutine")
	for {
		select {
		case <-m.exit:
			logger.Debug("gen proof goroutine exit ...")
			return
		default:
			if m.proofQueue.Len() == 0 {
				//logger.Debug("no proof need to do,wait now ....")
				time.Sleep(1 * time.Second)
				continue
			}
			worker, find, err := m.schedule.findBestWorker()
			if err != nil {
				logger.Error("find best worker error:%v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			if !find {
				logger.Warn(" no find best worker to gen proof")
				time.Sleep(1 * time.Second)
				continue
			}
			frontElement := m.proofQueue.Front()
			proofRequest, ok := frontElement.Value.(ProofRequest)
			if !ok {
				logger.Error("should never happen,parse proof request error")
				continue
			}
			logger.Info("start gen proof:%v", proofRequest.String())
			m.proofQueue.Remove(frontElement)
			go func() {
				proofResponse, err := m.schedule.GenZKProof(worker, proofRequest)
				if err != nil {
					//todo add queue again or cli retry ?
					logger.Error("gen proof error:%v %v", proofRequest.TxId, err)
				}
				logger.Info("worker response proof: %v", proofResponse.String())
				switch proofResponse.ProofType {
				case Deposit:
					m.btcProofResp <- proofResponse
				case Redeem:
					m.ethProofResp <- proofResponse
				default:
					logger.Error("never should happen proof type:%v", proofResponse.ProofType)
				}
			}()
		}
	}
}

func (m *Manager) Close() {
	close(m.exit)
}

type SafeList struct {
	list *list.List
	mu   sync.Mutex
}

func NewSafeList() *SafeList {
	return &SafeList{
		list: list.New(),
		mu:   sync.Mutex{},
	}
}

func (sl *SafeList) PushBack(value interface{}) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.list.PushBack(value)
}

func (sl *SafeList) PushFront(value interface{}) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.list.PushFront(value)
}

func (sl *SafeList) Front() *list.Element {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	return sl.list.Front()

}

func (sl *SafeList) Len() int {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	return sl.list.Len()
}
func (sl *SafeList) Remove(e *list.Element) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.list.Remove(e)
}
