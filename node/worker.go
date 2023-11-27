package node

import (
	"github.com/lightec-xyz/daemon/rpc"
	"sync"
)

type IWorker interface {
	GenProof(req rpc.ProofRequest) (rpc.ProofResponse, error)
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
	lock         sync.Locker
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

var _ IWorker = (*LocalWorker)(nil)

type LocalWorker struct {
	parallelNums int
	currentNums  int
	lock         sync.Locker
}

func (l *LocalWorker) ParallelNums() int {
	//TODO implement me
	panic("implement me")
}

func (l *LocalWorker) CurrentNums() int {
	//TODO implement me
	panic("implement me")
}

func (l *LocalWorker) GenProof(req rpc.ProofRequest) (rpc.ProofResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LocalWorker) Add() {
	//TODO implement me
	panic("implement me")
}

func (l *LocalWorker) Del() {
	//TODO implement me
	panic("implement me")
}
