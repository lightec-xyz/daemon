package node

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"reflect"
	"sync"
	"time"

	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
)

var _ rpc.IWorker = (*LocalWorker)(nil)

func NewLocalWorker(setupDir, dataDir string, maxNums int) (rpc.IWorker, error) {
	config := circuits.CircuitConfig{
		DataDir: setupDir,
	}
	circuit, err := circuits.NewCircuit(&config)
	if err != nil {
		return nil, err
	}
	return &LocalWorker{
		dataDir:     dataDir,
		maxNums:     maxNums,
		currentNums: 0,
		wid:         UUID(),
		circuit:     circuit,
	}, nil
}

type LocalWorker struct {
	circuit     *circuits.Circuit
	dataDir     string
	maxNums     int
	currentNums int
	lock        sync.Mutex
	wid         string
	fileStore   *FileStore
}

func (l *LocalWorker) Id() string {
	return l.wid
}

func (l *LocalWorker) ProofInfo(proofId string) (rpc.ProofInfo, error) {
	logger.Debug("Proof info")
	time.Sleep(10 * time.Second)
	return rpc.ProofInfo{
		Status: 0,
		Proof:  common.ZkProof([]byte("")),
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
	proof, err := l.circuit.GenesisProve(req.FirstProof, req.FirstWitness, req.SecondProof, req.SecondWitness,
		req.GenesisID, req.FirstID, req.SecondID, req.RecursiveFp)
	if err != nil {
		logger.Error("unit prove error", err)
		return rpc.SyncCommGenesisResponse{}, err
	}
	logger.Debug("complete %v genesis prove", req.Period)
	return rpc.SyncCommGenesisResponse{
		Version:   req.Version,
		Period:    req.Period,
		ProofType: common.SyncComGenesisType,
		Proof:     circuits.ProofToBytes(proof.Proof),
		Witness:   circuits.WitnessToBytes(proof.Wit),
	}, nil
}

func (l *LocalWorker) GenSyncCommitUnitProof(req rpc.SyncCommUnitsRequest) (rpc.SyncCommUnitsResponse, error) {
	logger.Debug("gen units Proof")
	var update utils.LightClientUpdateInfo
	err := deepCopy(req, &update)
	if err != nil {
		logger.Error("deep copy error", err)
		return rpc.SyncCommUnitsResponse{}, err
	}
	proof, err := l.circuit.UnitProve(&update)
	if err != nil {
		logger.Error("unit prove error", err)
		return rpc.SyncCommUnitsResponse{}, err
	}
	logger.Debug("complete %v unit prove", req.Period)
	return rpc.SyncCommUnitsResponse{
		Version:   req.Version,
		Period:    req.Period,
		ProofType: common.SyncComUnitType,
		Proof:     circuits.ProofToBytes(proof.Proof),
		Witness:   circuits.WitnessToBytes(proof.Wit),
	}, nil

}

func (l *LocalWorker) GenSyncCommRecursiveProof(req rpc.SyncCommRecursiveRequest) (rpc.SyncCommRecursiveResponse, error) {
	logger.Debug("gen recursive Proof")
	var update utils.LightClientUpdateInfo
	err := deepCopy(req, &update)
	if err != nil {
		logger.Error("deep copy error", err)
		return rpc.SyncCommRecursiveResponse{}, err
	}
	proof, err := l.circuit.RecursiveProve(req.Choice, req.FirstProof, req.SecondProof, req.FirstWitness, req.SecondWitness,
		req.BeginId, req.RelayId, req.EndId, req.RecursiveFp)
	if err != nil {
		logger.Error("recursive prove error", err)
		return rpc.SyncCommRecursiveResponse{}, err
	}
	logger.Debug("complete %v recursive prove", req.Period)
	return rpc.SyncCommRecursiveResponse{
		Version:   req.Version,
		Period:    req.Period,
		ProofType: common.SyncComRecursiveType,
		Proof:     circuits.ProofToBytes(proof.Proof),
		Witness:   circuits.WitnessToBytes(proof.Wit),
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

func GetGenesisProofPath(datadir string) string {
	return fmt.Sprintf("%s/genesis/genesis.proof", datadir)
}

func GetGenesisWitnessPath(datadir string, period uint64) string {
	return fmt.Sprintf("%s/%d_genesis.witness", datadir, period)
}

func GetUnitProofPath(datadir string, period uint64) string {
	return fmt.Sprintf("%s/%d_unit.proof", datadir, period)
}

func GetRecursiveProofPath(datadir string, period uint64) string {
	return fmt.Sprintf("%s/%d_recursive.proof", datadir, period)
}

func GetUnitWitnessPath(unitDir string, period uint64) string {
	return fmt.Sprintf("%s/%d_unit.witness", unitDir, period)
}

func GetRecursiveWitnessPath(recursiveDir string, period uint64) string {
	return fmt.Sprintf("%s/%d_recursive.witness", recursiveDir, period)
}

func deepCopy(src, dst interface{}) error {
	if reflect.ValueOf(dst).Kind() != reflect.Ptr {
		return fmt.Errorf("dst must be a pointer")
	}
	srcBytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(srcBytes, dst)
	if err != nil {
		return err
	}
	return nil
}
