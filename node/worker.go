package node

import (
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"sync"
	"time"
)

type IWorker interface {
	GenProof(req rpc.ProofRequest) (rpc.ProofResponse, error)
	Add()
	Del()
	ParallelNums() int
	CurrentNums() int
}

//type ISyncCommitteeWorker interface {
//	GenProof(req rpc.SyncCommitteeProofRequest) (rpc.SyncCommitteeProofResponse, error)
//	Add()
//	Del()
//	ParallelNums() int
//	CurrentNums() int
//}

var _ IWorker = (*Worker)(nil)

type Worker struct {
	client      rpc.IProof
	maxNums     int
	currentNums int
	lock        sync.Mutex
}

func (w *Worker) ParallelNums() int {
	return w.maxNums
}

func (w *Worker) CurrentNums() int {
	return w.currentNums
}

func (w *Worker) GenProof(req rpc.ProofRequest) (rpc.ProofResponse, error) {
	return w.client.GenZkProof(req)
}

func NewWorker(client rpc.IProof, parallelNums int) *Worker {
	return &Worker{
		client:      client,
		maxNums:     parallelNums,
		currentNums: 0,
	}
}

func (w *Worker) Del() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.currentNums--
}

func (w *Worker) Add() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.currentNums++
}

var _ IWorker = (*LocalWorker)(nil)

type LocalWorker struct {
	parallelNums int
	currentNums  int
	lock         sync.Mutex
}

func NewLocalWorker(maxNums int) IWorker {
	return &LocalWorker{
		parallelNums: maxNums,
		currentNums:  0,
	}
}

func (l *LocalWorker) ParallelNums() int {
	return l.parallelNums
}

func (l *LocalWorker) CurrentNums() int {
	return l.currentNums
}

func (l *LocalWorker) GenProof(req rpc.ProofRequest) (rpc.ProofResponse, error) {
	// todo
	logger.Info("local worker gen proof now: %v %v", req.TxId, req.ProofType)
	time.Sleep(6 * time.Second)
	response := rpc.ProofResponse{}
	err := objParse(req, &response)
	if err != nil {
		return response, nil
	}
	return response, nil
}

func (l *LocalWorker) Add() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.currentNums--
}

func (l *LocalWorker) Del() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.currentNums++
}
