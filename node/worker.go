package node

import (
	"github.com/lightec-xyz/daemon/rpc"
	"sync"
	"time"
)

var _ rpc.IProof = (*Worker)(nil)

type Worker struct {
	client      rpc.IProof
	maxNums     int
	currentNums int
	lock        sync.Mutex
}

func NewWorker(client rpc.IProof, parallelNums int) *Worker {
	return &Worker{
		client:      client,
		maxNums:     parallelNums,
		currentNums: 0,
	}
}

func (w *Worker) ProofInfo(proofId string) (rpc.ProofInfo, error) {
	return w.client.ProofInfo(proofId)
}

func (w *Worker) GenDepositProof(req rpc.DepositRequest) (rpc.DepositResponse, error) {
	return w.client.GenDepositProof(req)
}

func (w *Worker) GenRedeemProof(req rpc.RedeemRequest) (rpc.RedeemResponse, error) {
	return w.GenRedeemProof(req)
}

func (w *Worker) GenVerifyProof(req rpc.VerifyRequest) (rpc.VerifyResponse, error) {
	return w.client.GenVerifyProof(req)
}

func (w *Worker) GenSyncCommGenesisProof(req rpc.SyncCommGenesisRequest) (rpc.SyncCommGenesisResponse, error) {
	return w.client.GenSyncCommGenesisProof(req)
}

func (w *Worker) GenSyncCommitUnitProof(req rpc.SyncCommUnitsRequest) (rpc.SyncCommUnitsResponse, error) {
	return w.client.GenSyncCommitUnitProof(req)
}

func (w *Worker) GenSyncCommRecursiveProof(req rpc.SyncCommRecursiveRequest) (rpc.SyncCommRecursiveResponse, error) {
	return w.client.GenSyncCommRecursiveProof(req)
}

func (w *Worker) MaxNums() int {
	return w.maxNums
}

func (w *Worker) CurrentNums() int {
	return w.currentNums
}

func (w *Worker) DelReqNum() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.currentNums--
}

func (w *Worker) AddReqNum() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.currentNums++
}

func NewLocalWorker(maxNums int) rpc.IProof {
	return &LocalWorker{
		maxNums:     maxNums,
		currentNums: 0,
	}
}

var _ rpc.IProof = (*LocalWorker)(nil)

type LocalWorker struct {
	maxNums     int
	currentNums int
	lock        sync.Mutex
}

func (l *LocalWorker) ProofInfo(proofId string) (rpc.ProofInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LocalWorker) GenRedeemProof(req rpc.RedeemRequest) (rpc.RedeemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LocalWorker) GenVerifyProof(req rpc.VerifyRequest) (rpc.VerifyResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LocalWorker) GenSyncCommGenesisProof(req rpc.SyncCommGenesisRequest) (rpc.SyncCommGenesisResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LocalWorker) GenSyncCommitUnitProof(req rpc.SyncCommUnitsRequest) (rpc.SyncCommUnitsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LocalWorker) GenSyncCommRecursiveProof(req rpc.SyncCommRecursiveRequest) (rpc.SyncCommRecursiveResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LocalWorker) MaxNums() int {
	return l.maxNums
}

func (l *LocalWorker) CurrentNums() int {
	return l.currentNums
}

func (l *LocalWorker) GenDepositProof(req rpc.DepositRequest) (rpc.DepositResponse, error) {
	// todo
	time.Sleep(6 * time.Second)
	response := rpc.DepositResponse{}
	err := objParse(req, &response)
	if err != nil {
		return response, nil
	}
	return response, nil
}

func (l *LocalWorker) AddReqNum() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.currentNums--
}

func (l *LocalWorker) DelReqNum() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.currentNums++
}
