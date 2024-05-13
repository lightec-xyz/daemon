package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"sort"
	"strings"
	"sync"
)

type Schedule struct {
	Workers []rpc.IWorker
	lock    sync.Mutex
}

func NewSchedule(workers []rpc.IWorker) *Schedule {
	return &Schedule{
		Workers: workers,
	}
}

func (m *Schedule) AddWorker(endpoint string, nums int) error {
	var client rpc.IProof
	var err error
	if strings.HasPrefix(endpoint, "http") {
		client, err = rpc.NewProofClient(endpoint)
		if err != nil {
			logger.Error("new worker error:%v %v", endpoint, err)
			return err
		}
	} else if strings.HasPrefix(endpoint, "ws") {
		client, err = rpc.NewWsProofClient(endpoint)
		if err != nil {
			logger.Error("new worker error:%v %v", endpoint, err)
			return err
		}
	} else {
		return fmt.Errorf("unSupport protocol: %v", endpoint)
	}

	newWorker := NewWorker(client, nums)
	m.Workers = append(m.Workers, newWorker)
	return nil
}

func (m *Schedule) findBestWorker(work func(worker rpc.IWorker) error) (rpc.IWorker, bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	// todo find work by Proof type
	var tmpWorkers []rpc.IWorker
	for _, worker := range m.Workers {
		if worker.CurrentNums() < worker.MaxNums() {
			tmpWorkers = append(tmpWorkers, worker)
		}
	}
	if len(tmpWorkers) == 0 {
		return nil, false, nil
	}
	sort.Slice(tmpWorkers, func(i, j int) bool {
		return tmpWorkers[i].CurrentNums() < tmpWorkers[j].CurrentNums()
	})
	bestWork := tmpWorkers[0]
	err := work(bestWork)
	if err != nil {
		return nil, false, err
	}
	return bestWork, true, nil

}
