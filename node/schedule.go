package node

import (
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"sort"
)

type Schedule struct {
	Workers []IWorker
}

func NewSchedule(workers ...IWorker) *Schedule {
	return &Schedule{
		Workers: workers,
	}
}

func (m *Schedule) GenZKProof(worker IWorker, req ProofRequest) (rpc.ProofResponse, error) {
	worker.Add()
	defer worker.Del()
	proofResponse, err := worker.GenProof(rpc.ProofRequest{
		TxId: req.TxId,
	})
	if err != nil {
		logger.Error("gen zk proof error:%v", err)
		return rpc.ProofResponse{}, err
	}
	return proofResponse, nil

}

func (m *Schedule) findBestWorker() (IWorker, bool, error) {
	// todo

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
