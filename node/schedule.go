package node

import (
	"encoding/json"
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

func (m *Schedule) GenZKProof(worker IWorker, req ProofRequest) (ProofResponse, error) {
	worker.Add()
	defer worker.Del()
	//todo
	proofResp := ProofResponse{}
	rpcReq := rpc.ProofRequest{}
	err := objParse(req, &rpcReq)
	if err != nil {
		return proofResp, err
	}
	proofResponse, err := worker.GenProof(rpcReq)
	if err != nil {
		logger.Error("gen zk proof error:%v", err)
		return proofResp, err
	}
	err = objParse(proofResponse, &proofResp)
	if err != nil {
		return proofResp, nil
	}
	return proofResp, nil

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
