package node

import (
	"container/list"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/store"
	"sync"
)

type Manager struct {
	proofQueue    *SafeList
	schedule      *Schedule
	store         store.IStore
	memory        store.IStore
	proofRequest  chan []ProofRequest
	proofResponse chan []ProofResponse
}

func NewManager(proofRequest chan []ProofRequest, proofResponse chan []ProofResponse, store, memory store.IStore, schedule *Schedule) *Manager {
	return &Manager{
		proofQueue:    NewSafeList(),
		schedule:      schedule,
		store:         store,
		memory:        memory,
		proofRequest:  proofRequest,
		proofResponse: proofResponse,
	}
}

func (m *Manager) Run() {
	//todo
	for {
		select {
		case requestList := <-m.proofRequest:
			for _, req := range requestList {
				logger.Debug("manager receive proof request:%v", req)
				m.proofQueue.PushBack(req)
			}
		}
	}
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
