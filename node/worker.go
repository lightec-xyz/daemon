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

type LocalWorker struct {
	circuit     *circuits.Circuit
	dataDir     string
	maxNums     int
	currentNums int
	lock        sync.Mutex
	wid         string
}

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

func (w *LocalWorker) TxInEth2Prove(req *rpc.TxInEth2ProveReq) (*rpc.TxInEth2ProveResp, error) {
	logger.Debug("local worker transaction in eth2")
	proof, err := w.circuit.TxInEth2Prove(req.TxData)
	if err != nil {
		logger.Error("TxInEth2Prove error: %v", err)
		return nil, err
	}
	hexProof, err := circuits.ProofToHexSolBytes(proof.Proof)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return &rpc.TxInEth2ProveResp{
		ProofStr: hexProof,
		//Proof:    circuits.ProofToBytes(proof.Proof),
		Witness: circuits.WitnessToBytes(proof.Wit),
	}, nil

}

func (w *LocalWorker) TxBlockIsParentOfCheckPointProve(req *rpc.TxBlockIsParentOfCheckPointProveReq) (*rpc.TxBlockIsParentOfCheckPointResp, error) {
	//TODO implement me
	panic("implement me")
}

func (w *LocalWorker) CheckPointFinalityProve(req *rpc.CheckPointFinalityProveReq) (*rpc.CheckPointFinalityProveResp, error) {
	//TODO implement me
	panic("implement me")
}

func (w *LocalWorker) Id() string {
	return w.wid
}

func (w *LocalWorker) ProofInfo(proofId string) (rpc.ProofInfo, error) {
	logger.Debug("Proof info")
	time.Sleep(10 * time.Second)
	return rpc.ProofInfo{
		Status: 0,
		Proof:  "",
		TxId:   proofId,
	}, nil
}

func (w *LocalWorker) GenDepositProof(req rpc.DepositRequest) (rpc.DepositResponse, error) {
	logger.Debug("gen deposit proof")
	proof, err := w.circuit.DepositProve(req.TxHash, req.BlockHash)
	if err != nil {
		logger.Error(err.Error())
		return rpc.DepositResponse{}, fmt.Errorf("gen deposit prove error: %v", err)
	}
	hexProof, err := circuits.ProofToHexSolBytes(proof.Proof)
	if err != nil {
		logger.Error(err.Error())
		return rpc.DepositResponse{}, nil
	}
	return rpc.DepositResponse{
		TxHash: req.TxHash,
		//Proof:    common.ZkProof(hexProof),
		ProofStr: hexProof,
		Witness:  circuits.WitnessToBytes(proof.Wit),
	}, nil
}

func (w *LocalWorker) GenRedeemProof(req rpc.RedeemRequest) (rpc.RedeemResponse, error) {
	logger.Debug("gen redeem Proof")
	time.Sleep(10 * time.Second)
	return rpc.RedeemResponse{
		Proof: common.ZkProof([]byte("redeem Proof")),
	}, nil
}

func (w *LocalWorker) GenVerifyProof(req rpc.VerifyRequest) (rpc.VerifyResponse, error) {
	logger.Debug("verify Proof")
	time.Sleep(10 * time.Second)
	return rpc.VerifyResponse{
		Proof: common.ZkProof([]byte("verify Proof")),
	}, nil
}

func (w *LocalWorker) GenSyncCommGenesisProof(req rpc.SyncCommGenesisRequest) (rpc.SyncCommGenesisResponse, error) {
	logger.Debug("start gen genesis prove %v period", req.Period)
	proof, err := w.circuit.GenesisProve(req.FirstProof, req.SecondProof, req.FirstWitness, req.SecondWitness,
		req.GenesisID, req.FirstID, req.SecondID)
	if err != nil {
		logger.Error("genesis prove error %v", err)
		return rpc.SyncCommGenesisResponse{}, err
	}
	logger.Debug("complete  genesis prove %v", req.Period)
	return rpc.SyncCommGenesisResponse{
		Version:   req.Version,
		Period:    req.Period,
		ProofType: common.SyncComGenesisType,
		Proof:     common.ZkProof(circuits.ProofToBytes(proof.Proof)),
		Witness:   circuits.WitnessToBytes(proof.Wit),
	}, nil
}

func (w *LocalWorker) GenSyncCommitUnitProof(req rpc.SyncCommUnitsRequest) (rpc.SyncCommUnitsResponse, error) {
	// todo
	logger.Debug("unit prove request: %v period", req.Period)
	var update utils.LightClientUpdateInfo
	err := deepCopy(req, &update)
	if err != nil {
		logger.Error("deep copy error %v", err)
		return rpc.SyncCommUnitsResponse{}, err
	}
	proof, err := w.circuit.UnitProve(req.Period, &update)
	if err != nil {
		logger.Error("unit prove error %v", err)
		return rpc.SyncCommUnitsResponse{}, err
	}
	logger.Debug("complete unit prove %v", req.Period)
	return rpc.SyncCommUnitsResponse{
		Version:   req.Version,
		Period:    req.Period,
		ProofType: common.SyncComUnitType,
		Proof:     common.ZkProof(circuits.ProofToBytes(proof.Proof)),
		Witness:   circuits.WitnessToBytes(proof.Wit),
	}, nil

}

func (w *LocalWorker) GenSyncCommRecursiveProof(req rpc.SyncCommRecursiveRequest) (rpc.SyncCommRecursiveResponse, error) {
	logger.Debug("recursive prove request %v period", req.Period)
	proof, err := w.circuit.RecursiveProve(req.Choice, req.FirstProof, req.SecondProof, req.FirstWitness, req.SecondWitness,
		req.BeginId, req.RelayId, req.EndId)
	if err != nil {
		logger.Error("recursive prove error %v", err)
		return rpc.SyncCommRecursiveResponse{}, err
	}
	logger.Debug("complete recursive prove %v", req.Period)
	return rpc.SyncCommRecursiveResponse{
		Version:   req.Version,
		Period:    req.Period,
		ProofType: common.SyncComRecursiveType,
		Proof:     common.ZkProof(circuits.ProofToBytes(proof.Proof)),
		Witness:   circuits.WitnessToBytes(proof.Wit),
	}, nil
}

func (w *LocalWorker) MaxNums() int {
	return w.maxNums
}

func (w *LocalWorker) CurrentNums() int {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.currentNums
}

func (w *LocalWorker) AddReqNum() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.currentNums = w.currentNums + 1
}

func (w *LocalWorker) DelReqNum() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.currentNums = w.currentNums - 1
}

var _ rpc.IWorker = (*Worker)(nil)

type Worker struct {
	client      rpc.IProof
	maxNums     int
	currentNums int
	lock        sync.Mutex
	wid         string
}

func (w *Worker) TxInEth2Prove(req *rpc.TxInEth2ProveReq) (*rpc.TxInEth2ProveResp, error) {

	//TODO implement me
	panic("implement me")
}

func (w *Worker) TxBlockIsParentOfCheckPointProve(req *rpc.TxBlockIsParentOfCheckPointProveReq) (*rpc.TxBlockIsParentOfCheckPointResp, error) {
	//TODO implement me
	panic("implement me")
}

func (w *Worker) CheckPointFinalityProve(req *rpc.CheckPointFinalityProveReq) (*rpc.CheckPointFinalityProveResp, error) {
	//TODO implement me
	panic("implement me")
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

func objToJson(obj interface{}) string {
	ojbBytes, err := json.Marshal(obj)
	if err != nil {
		return "error obj to josn"
	}
	return string(ojbBytes)

}
