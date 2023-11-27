package node

import (
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/store"
	"sort"
)

type Schedule struct {
	Workers     []IWorker
	store       store.IStore
	memoryStore store.IStore
}

func NewSchedule(store, memoryStore store.IStore, worker ...IWorker) *Schedule {
	return &Schedule{
		store:       store,
		memoryStore: memoryStore,
		Workers:     worker,
	}
}

func (m *Schedule) GenZKProof(worker IWorker, req rpc.ProofRequest) (rpc.ProofResponse, error) {
	worker.Add()
	defer worker.Del()
	proofResponse, err := worker.GenProof(req)
	if err != nil {
		logger.Error("gen zk proof error:%v", err)
		return rpc.ProofResponse{}, err
	}
	return proofResponse, nil

}

func (m *Schedule) findBestWorker() (IWorker, bool, error) {
	//todo
	var tmpWorkers []IWorker
	for _, worker := range m.Workers {
		if worker.CurrentNums() < worker.ParallelNums() {
			tmpWorkers = append(tmpWorkers, worker)
		}
	}
	if len(tmpWorkers) == 0 {
		return nil, false, nil
	}
	sort.Slice(tmpWorkers, func(i, j int) bool {
		return tmpWorkers[i].CurrentNums() < tmpWorkers[j].ParallelNums()
	})
	bestWork := tmpWorkers[0]
	return bestWork, true, nil

}
