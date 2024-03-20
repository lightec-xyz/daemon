package node

import (
	"sync"
	"time"

	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
)

var _ rpc.IWorker = (*LocalWorker)(nil)

func NewLocalWorker(maxNums int) rpc.IWorker {
	return &LocalWorker{
		maxNums:     maxNums,
		currentNums: 0,
		wid:         UUID(),
	}
}

type LocalWorker struct {
	maxNums     int
	currentNums int
	lock        sync.Mutex
	wid         string
}

func (l *LocalWorker) Id() string {
	return l.wid
}

func (l *LocalWorker) ProofInfo(proofId string) (rpc.ProofInfo, error) {
	logger.Debug("Proof info")
	time.Sleep(10 * time.Second)
	return rpc.ProofInfo{
		Status: 0,
		Proof:  "",
		TxId:   proofId,
	}, nil
}

func (l *LocalWorker) GenDepositProof(req rpc.DepositRequest) (rpc.DepositResponse, error) {
	logger.Debug("gen deposit Proof")
	time.Sleep(6 * time.Second)
	return rpc.DepositResponse{
		Proof: common.ZkProof([]byte("deposit Proof")),
	}, nil
}

func (l *LocalWorker) GenRedeemProof(req rpc.RedeemRequest) (rpc.RedeemResponse, error) {
	logger.Debug("gen redeem Proof")
	time.Sleep(10 * time.Second)
	return rpc.RedeemResponse{
		Proof: common.ZkProof([]byte("redeem Proof")),
	}, nil
}

func (l *LocalWorker) GenVerifyProof(req rpc.VerifyRequest) (rpc.VerifyResponse, error) {
	logger.Debug("verify Proof")
	time.Sleep(10 * time.Second)
	return rpc.VerifyResponse{
		Proof: common.ZkProof([]byte("verify Proof")),
	}, nil
}

func (l *LocalWorker) GenSyncCommGenesisProof(req rpc.SyncCommGenesisRequest) (rpc.SyncCommGenesisResponse, error) {
	logger.Debug("gen genesis Proof")
	time.Sleep(10 * time.Second)
	return rpc.SyncCommGenesisResponse{
		Proof: common.ZkProof([]byte("genesis Proof")),
	}, nil
}

func (l *LocalWorker) GenSyncCommitUnitProof(req rpc.SyncCommUnitsRequest) (rpc.SyncCommUnitsResponse, error) {
	logger.Debug("gen units Proof")
	time.Sleep(10 * time.Second)
	return rpc.SyncCommUnitsResponse{
		Proof: common.ZkProof([]byte("units Proof")),
	}, nil
}

func (l *LocalWorker) GenSyncCommRecursiveProof(req rpc.SyncCommRecursiveRequest) (rpc.SyncCommRecursiveResponse, error) {
	logger.Debug("gen recursive Proof")
	time.Sleep(10 * time.Second)
	return rpc.SyncCommRecursiveResponse{
		Proof: common.ZkProof([]byte("recursive proof")),
	}, nil
}

func (l *LocalWorker) MaxNums() int {
	return l.maxNums
}

func (l *LocalWorker) CurrentNums() int {
	return l.currentNums
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

var _ rpc.IWorker = (*Worker)(nil)

type Worker struct {
	client      rpc.IProof
	maxNums     int
	currentNums int
	lock        sync.Mutex
	wid         string
}

func (w *Worker) Id() string {
	return w.wid
}

func NewWorker(client rpc.IProof, parallelNums int) *Worker {
	return &Worker{
		client:      client,
		maxNums:     parallelNums,
		currentNums: 0,
		wid:         UUID(),
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
