package node

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"sort"
	"strings"
	"sync"
)

type Schedule struct {
	Workers             []IWorker
	syncer              sync.Mutex
	syncCommitteeSyncer sync.Mutex
}

func NewSchedule(workers []IWorker) *Schedule {
	// todo
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

func (m *Schedule) AddSyncCommitteeWorker(endpoint string, nums int) error {
	//var client rpc.SyncCommitteeProofAPI
	//var err error
	//if strings.HasPrefix(endpoint, "http") {
	//	client, err = rpc.NewProofClient(endpoint)
	//	if err != nil {
	//		logger.Error("new worker error:%v %v", endpoint, err)
	//		return err
	//	}
	//} else if strings.HasPrefix(endpoint, "ws") {
	//	client, err = rpc.NewWsProofClient(endpoint)
	//	if err != nil {
	//		logger.Error("new worker error:%v %v", endpoint, err)
	//		return err
	//	}
	//} else {
	//	return fmt.Errorf("unSupport protocol: %v", endpoint)
	//}
	//
	//newWorker := NewWorker(client, nums)
	//m.Workers = append(m.Workers, newWorker)
	return nil
}

func (m *Schedule) GenZKProof(worker IWorker, req ZkProofRequest) (ZkProofResponse, error) {
	worker.Add()
	defer worker.Del()
	// todo
	proofResp := ZkProofResponse{}
	rpcReq := rpc.ProofRequest{}
	err := objParse(req, &rpcReq)
	if err != nil {
		logger.Error("parse proof request error:%v", err)
		return proofResp, err
	}
	proofResponse, err := worker.GenProof(rpcReq)
	if err != nil {
		logger.Error("gen zk proof error:%v", err)
		return proofResp, err
	}
	err = objParse(proofResponse, &proofResp)
	if err != nil {
		logger.Error("parse proof response error:%v", err)
		return proofResp, nil
	}
	return proofResp, nil

}

func (m *Schedule) findBestWorker(proofType ZkProofType) (IWorker, bool, error) {
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

func objParse(src, dest interface{}) error {
	marshal, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, dest)
	if err != nil {
		return err
	}
	return nil
}
