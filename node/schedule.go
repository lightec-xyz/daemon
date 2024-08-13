package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/ws"
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

func (s *Schedule) Close() error {
	for _, worker := range s.Workers {
		if worker != nil {
			_ = worker.Close()
		}
	}
	return nil
}
func (s *Schedule) AddWorker(endpoint string, nums int) error {
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
	s.Workers = append(s.Workers, newWorker)
	return nil
}

func (s *Schedule) AddWsWorker(opt *rpc.WsOpt) error {
	wsConn := ws.NewConn(opt.Conn, nil, nil, true)
	wsConn.Run()
	proofClient, err := rpc.NewCustomWsProofClient(wsConn)
	if err != nil {
		logger.Error("new worker error:%v", err)
		return err
	}
	worker := NewWorker(proofClient, 1)
	s.Workers = append(s.Workers, worker)
	return nil
}

func (s *Schedule) findWorker(req common.ZkProofType) (rpc.IWorker, bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var tmpWorkers []rpc.IWorker
	for _, worker := range s.Workers {
		if matchReqType(worker.SupportProofType(), req) && worker.CurrentNums() < worker.MaxNums() {
			tmpWorkers = append(tmpWorkers, worker)
		}
	}
	if len(tmpWorkers) == 0 {
		return nil, false, nil
	}
	sort.SliceStable(tmpWorkers, func(i, j int) bool {
		return tmpWorkers[i].CurrentNums() < tmpWorkers[j].CurrentNums()
	})
	bestWork := tmpWorkers[0]
	return bestWork, true, nil
}

func (s *Schedule) findBestWorker(work func(worker rpc.IWorker) error) (rpc.IWorker, bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// todo find work by Proof type
	var tmpWorkers []rpc.IWorker
	for _, worker := range s.Workers {
		if worker.CurrentNums() < worker.MaxNums() {
			tmpWorkers = append(tmpWorkers, worker)
		}
	}
	if len(tmpWorkers) == 0 {
		return nil, false, nil
	}
	sort.SliceStable(tmpWorkers, func(i, j int) bool {
		return tmpWorkers[i].CurrentNums() < tmpWorkers[j].CurrentNums()
	})
	bestWork := tmpWorkers[0]
	err := work(bestWork)
	if err != nil {
		return nil, false, err
	}
	return bestWork, true, nil
}

func matchReqType(wProofTypes []common.ZkProofType, reqType common.ZkProofType) bool {
	for _, wProofType := range wProofTypes {
		if wProofType == reqType {
			return true
		}
	}
	return false
}
