package node

import (
	"container/list"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
)

type Manager struct {
	proofQueue    *list.List
	schedule      *Schedule
	store         store.IStore
	memory        store.IStore
	proofRequest  chan rpc.ProofRequest
	proofResponse chan rpc.ProofResponse
}

func NewManager(store, memory store.IStore, schedule *Schedule) *Manager {
	return &Manager{
		proofQueue:    list.New(),
		schedule:      schedule,
		store:         store,
		memory:        memory,
		proofRequest:  make(chan rpc.ProofRequest, 1000),
		proofResponse: make(chan rpc.ProofResponse, 100),
	}
}

func (m *Manager) Close() {
	//todo
	for {
		select {
		case request := <-m.proofRequest:
			worker, find, err := m.schedule.findBestWorker()
			if err != nil {
				logger.Error("find best worker error:%v", err)
				continue
			}
			if !find {
				logger.Warn("no find best worker,wait now ")
				continue
			}
			go func() {
				proofResponse, err := m.schedule.GenZKProof(worker, request)
				if err != nil {
					// todo retry
					
					logger.Error("worker gen proof error:%v", err)
					return
				}
				logger.Info("success gen zk proof:%v", proofResponse)
				m.proofResponse <- proofResponse
			}()
		}
	}

}
