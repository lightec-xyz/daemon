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

type ISyncCommitteeWorker interface {
	rpc.I

	Add()
	Del()
	ParallelNums() int
	CurrentNums() int
}

var _ IWorker = (*Worker)(nil)

type Worker struct {
	client       rpc.ProofAPI
	parallelNums int
	currentNums  int
	lock         sync.Mutex
}

func (w *Worker) ParallelNums() int {
	return w.parallelNums
}

func (w *Worker) CurrentNums() int {
	return w.currentNums
}

func (w *Worker) GenProof(req rpc.ProofRequest) (rpc.ProofResponse, error) {
	return w.client.GenZkProof(req)
}

func NewWorker(client rpc.ProofAPI, parallelNums int) *Worker {
	return &Worker{
		client:       client,
		parallelNums: parallelNums,
		currentNums:  0,
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

var _ ISyncCommitteeWorker = (*SyncCommitteeWorker)(nil)

type SyncCommitteeWorker struct {
	client       rpc.SyncCommitteeProofAPI
	parallelNums int
	currentNums  int
	lock         sync.Mutex
}

func (w *SyncCommitteeWorker) ParallelNums() int {
	return w.parallelNums
}

func (w *SyncCommitteeWorker) CurrentNums() int {
	return w.currentNums
}

func (w *SyncCommitteeWorker) GenProof(req rpc.SyncCommitteeProofRequest) (rpc.SyncCommitteeProofResponse, error) {
	return w.client.GenZkSyncCommitteeProof(req)
}

func (w *SyncCommitteeWorker) Del() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.currentNums--
}

func (w *SyncCommitteeWorker) Add() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.currentNums++
}

func NewSyncCommitterWorker(client rpc.SyncCommitteeProofAPI, parallelNums int) *SyncCommitteeWorker {
	return &SyncCommitteeWorker{
		client:       client,
		parallelNums: parallelNums,
		currentNums:  0,
	}
}

var _ IWorker = (*LocalWorker)(nil)

type LocalWorker struct {
	parallelNums int
	currentNums  int
	lock         sync.Mutex
}

func NewLocalWorker(parallelNums int) IWorker {
	return &LocalWorker{
		parallelNums: parallelNums,
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
