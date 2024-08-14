package node

import (
	"fmt"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/daemon/circuits"
	proverType "github.com/lightec-xyz/provers/circuits/types"
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

func (w *LocalWorker) BtcDuperRecursiveProve(req *rpc.BtcDuperRecursiveRequest) (*rpc.ProofResponse, error) {
	first, err := circuits.HexToProof(req.First)
	if err != nil {
		logger.Error("btc recursive hex to proof error: %v", err)
		return nil, err
	}
	second, err := circuits.HexToProof(req.Second)
	if err != nil {
		logger.Error("btc recursive hex to proof error: %v", err)
		return nil, err
	}
	proof, err := w.circuit.BtcDuperRecursiveProve(req.Data, first, second)
	if err != nil {
		logger.Error("btc duper recursive prove error: %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		return nil, err
	}
	return &rpc.ProofResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil

}

func (w *LocalWorker) BtcDepthRecursiveProve(req *rpc.BtcDepthRecursiveRequest) (*rpc.ProofResponse, error) {
	recursiveProof, err := circuits.HexToProof(req.Recursive)
	if err != nil {
		logger.Error("btc recursive hex to proof error: %v", err)
		return nil, err
	}
	unitProof, err := circuits.HexToProof(req.Unit)
	if err != nil {
		logger.Error("btc recursive hex to proof error: %v", err)
		return nil, err
	}
	proof, err := w.circuit.BtcDepthRecursiveProve(req.Data, recursiveProof, unitProof)
	if err != nil {
		logger.Error("btc depth recursive prove error: %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc depth recursive prove error: %v", err)
		return nil, err
	}
	return &rpc.ProofResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil

}

func (w *LocalWorker) BtcChainProve(req *rpc.BtcChainRequest) (*rpc.ProofResponse, error) {
	recursive, err := circuits.HexToProof(req.Recursive)
	if err != nil {
		logger.Error("btc recursive hex to proof error: %v", err)
		return nil, err
	}
	base, err := circuits.HexToProof(req.Base)
	if err != nil {
		logger.Error("btc recursive hex to proof error: %v", err)
		return nil, err
	}
	middle, err := circuits.HexToProof(req.MidLevel)
	if err != nil {
		logger.Error("btc recursive hex to proof error: %v", err)
		return nil, err
	}
	uppper, err := circuits.HexToProof(req.Upper)
	if err != nil {
		logger.Error("btc recursive hex to proof error: %v", err)
		return nil, err
	}
	result, err := w.circuit.BtcChainProve(req.Data, recursive, base, middle, uppper)
	if err != nil {
		logger.Error("btc chain prove error: %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(result)
	if err != nil {
		logger.Error("btc chain prove error: %v", err)
		return nil, err
	}
	return &rpc.ProofResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil
}

func (w *LocalWorker) BtcDepositProve(req *rpc.BtcDepositRequest) (*rpc.ProofResponse, error) {
	blockChain, err := circuits.HexToProof(req.BlockChain)
	if err != nil {
		logger.Error("btc deposit hex to proof error: %v", err)
		return nil, err
	}
	cpDepth, err := circuits.HexToProof(req.CpDepth)
	if err != nil {
		logger.Error("btc deposit hex to proof error: %v", err)
		return nil, err
	}
	txDepth, err := circuits.HexToProof(req.TxDepth)
	if err != nil {
		logger.Error("btc deposit hex to proof error: %v", err)
		return nil, err
	}
	proof, err := w.circuit.BtcDepositProve(req.Data, blockChain, txDepth, cpDepth, ethCommon.HexToHash(req.R),
		ethCommon.HexToHash(req.S), ethCommon.HexToAddress(req.ProverAddr))
	if err != nil {
		logger.Error("btc deposit prove error: %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc deposit prove error: %v", err)
		return nil, err
	}
	return &rpc.ProofResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil
}

func (w *LocalWorker) BtcChangeProve(req *rpc.BtcChangeRequest) (*rpc.ProofResponse, error) {
	blockChain, err := circuits.HexToProof(req.BlockChain)
	if err != nil {
		logger.Error("btc change hex to proof error: %v", err)
		return nil, err
	}
	cpDepth, err := circuits.HexToProof(req.CpDepth)
	if err != nil {
		logger.Error("btc change hex to proof error: %v", err)
		return nil, err
	}
	txDepth, err := circuits.HexToProof(req.TxDepth)
	if err != nil {
		logger.Error("btc change hex to proof error: %v", err)
		return nil, err
	}
	redeem, err := circuits.HexToProof(req.Redeem)
	if err != nil {
		logger.Error("btc change hex to proof error: %v", err)
		return nil, err
	}
	proof, err := w.circuit.BtcChangeProve(req.Data, blockChain, txDepth, cpDepth, redeem, ethCommon.HexToHash(req.R),
		ethCommon.HexToHash(req.S), ethCommon.HexToAddress(req.ProverAddr))
	if err != nil {
		logger.Error("btc change prove error: %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc change prove error: %v", err)
		return nil, err
	}
	return &rpc.ProofResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil

}

func (w *LocalWorker) BtcGenesis(req *rpc.BtcGenesisRequest) (*rpc.ProofResponse, error) {
	firstProof, err := circuits.HexToProof(req.First)
	if err != nil {
		logger.Error("btc recursive hex to proof error: %v", err)
		return nil, err
	}
	secondProof, err := circuits.HexToProof(req.Second)
	if err != nil {
		logger.Error("btc recursive hex to proof error: %v", err)
		return nil, err
	}
	proof, err := w.circuit.BtcGenesisProve(req.Data, firstProof, secondProof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	return &rpc.ProofResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil

}

func (w *LocalWorker) BtcBaseProve(req *rpc.BtcBaseRequest) (*rpc.ProofResponse, error) {
	proof, err := w.circuit.BtcBaseProve(req.Data)
	if err != nil {
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	return &rpc.ProofResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil
}

func (w *LocalWorker) BtcMiddleProve(req *rpc.BtcMiddleRequest) (*rpc.ProofResponse, error) {
	proofs, err := circuits.HexToProofs(req.Proofs)
	if err != nil {
		logger.Error("btc middle hex to proofs error: %v", err)
		return nil, err
	}
	proof, err := w.circuit.BtcMiddleProve(req.Data, proofs)
	if err != nil {
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	return &rpc.ProofResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil
}

func (w *LocalWorker) BtcUpperProve(req *rpc.BtcUpperRequest) (*rpc.ProofResponse, error) {
	proofs, err := circuits.HexToProofs(req.Proofs)
	if err != nil {
		logger.Error("btc middle hex to proofs error: %v", err)
		return nil, err
	}
	proof, err := w.circuit.BtcUpperProve(req.Data, proofs)
	if err != nil {
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	return &rpc.ProofResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil
}

func (w *LocalWorker) BtcBulkProve(req *rpc.BtcBulkRequest) (*rpc.BtcBulkResponse, error) {
	proof, err := w.circuit.BtcBulkProve(req.Data)
	if err != nil {
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	return &rpc.BtcBulkResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil
}

func (w *LocalWorker) BtcPackedRequest(req *rpc.BtcPackedRequest) (*rpc.BtcPackResponse, error) {
	recursive, err := circuits.HexToProof(req.Recursive)
	if err != nil {
		logger.Error("btc packed hex to proofs error: %v", err)
		return nil, err
	}
	bulk, err := circuits.HexToProof(req.Bulk)
	if err != nil {
		logger.Error("btc packed hex to proofs error: %v", err)
		return nil, err
	}
	proof, err := w.circuit.BtcPackProve(req.Data, recursive, bulk)
	if err != nil {
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	return &rpc.BtcPackResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil

}

func (w *LocalWorker) SupportProofType() []common.ZkProofType {
	return nil
}

func (w *LocalWorker) Close() error {
	return nil
}

func NewLocalWorker(setupDir, dataDir string, maxNums int) (rpc.IWorker, error) {
	config := circuits.CircuitConfig{
		DataDir:  setupDir,
		SetupDir: setupDir,
		Debug:    common.GetEnvDebugMode(),
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

func (w *LocalWorker) TxInEth2Prove(req *rpc.TxInEth2ProveRequest) (*rpc.TxInEth2ProveResponse, error) {
	logger.Debug("start TxInEth2Prove: %v", req.TxHash)
	proof, err := w.circuit.TxInEth2Prove(req.TxData)
	if err != nil {
		logger.Error("TxInEth2Prove error: %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	logger.Debug("complete TxInEth2Prove: %v", req.TxHash)
	return &rpc.TxInEth2ProveResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil

}

func (w *LocalWorker) BlockHeaderProve(req *rpc.BlockHeaderRequest) (*rpc.BlockHeaderResponse, error) {
	logger.Debug("start BlockHeaderProve %v", req.BeginSlot)
	var middleHeader []*proverType.BeaconHeader
	err := common.ParseObj(req.Headers, &middleHeader)
	if err != nil {
		logger.Error("deep copy error %v", err)
		return nil, err
	}
	headers := proverType.BeaconHeaderChain{
		BeginSlot:           req.BeginSlot,
		EndSlot:             req.EndSlot,
		BeginRoot:           req.BeginRoot,
		EndRoot:             req.EndRoot,
		MiddleBeaconHeaders: middleHeader,
	}
	//logger.Debug("len: %v", len(headers.MiddleBeaconHeaders))
	proof, err := w.circuit.BeaconHeaderProve(headers)
	if err != nil {
		logger.Error("BlockHeaderProve error: %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	logger.Debug("complete BlockHeaderProve %v", req.BeginSlot)
	return &rpc.BlockHeaderResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil

}

func (w *LocalWorker) BlockHeaderFinalityProve(req *rpc.BlockHeaderFinalityRequest) (*rpc.BlockHeaderFinalityResponse, error) {
	logger.Debug("start BlockHeaderFinalityProve %v", req.Index)

	recursive, err := circuits.HexToProof(req.RecursiveProof)
	if err != nil {
		logger.Error("btc recursive hex to proof error: %v", err)
		return nil, err
	}
	outer, err := circuits.HexToProof(req.OuterProof)
	if err != nil {
		logger.Error("btc recursive hex to proof error: %v", err)
		return nil, err
	}
	proof, err := w.circuit.BeaconHeaderFinalityUpdateProve(req.GenesisSCSSZRoot, recursive, outer, req.FinalityUpdate, req.ScUpdate)
	if err != nil {
		logger.Error("BeaconHeaderFinalityUpdateProve error: %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	logger.Debug("complete BlockHeaderFinalityProve %v", req.Index)
	return &rpc.BlockHeaderFinalityResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil
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

func (w *LocalWorker) GenRedeemProof(req *rpc.RedeemRequest) (*rpc.RedeemResponse, error) {
	logger.Debug("start gen redeem proof: %v", req.TxHash)
	tx, err := circuits.HexToProof(req.TxProof)
	if err != nil {
		logger.Error("gen redeem proof error: %v", err)
		return nil, fmt.Errorf("gen redeem proof error: %v", err)
	}
	bh, err := circuits.HexToProof(req.BhProof)
	if err != nil {
		logger.Error("gen redeem proof error: %v", err)
		return nil, fmt.Errorf("gen redeem proof error: %v", err)
	}
	bhf, err := circuits.HexToProof(req.BhfProof)
	if err != nil {
		logger.Error("gen redeem proof error: %v", err)
		return nil, fmt.Errorf("gen redeem proof error: %v", err)
	}
	proof, err := w.circuit.RedeemProve(tx, bh, bhf, req.BeginId, req.EndId, req.GenesisScRoot, req.CurrentSCSSZRoot, req.TxVar, req.ReceiptVar)
	if err != nil {
		logger.Error(err.Error())
		return nil, fmt.Errorf("gen redeem proof error: %v", err)
	}
	proofSolBytes, err := circuits.ProofToSolBytes(proof.Proof)
	if err != nil {
		logger.Error("TxInEth2Prove error: %v", err)
		return nil, err
	}
	witnessBytes, err := circuits.WitnessToBytes(proof.Wit)
	if err != nil {
		logger.Error("TxInEth2Prove error: %v", err)
		return nil, err
	}
	logger.Debug("complete gen redeem proof: %v", req.TxHash)
	return &rpc.RedeemResponse{
		Proof:   proofSolBytes,
		Witness: witnessBytes,
	}, nil
}

func (w *LocalWorker) GenSyncCommGenesisProof(req rpc.SyncCommGenesisRequest) (*rpc.SyncCommGenesisResponse, error) {
	logger.Debug("start gen genesis prove %v Index", req.Period)
	first, err := circuits.HexToProof(req.FirstProof)
	if err != nil {
		logger.Error("hex to proof error %v", err)
		return nil, err
	}
	second, err := circuits.HexToProof(req.SecondProof)
	if err != nil {
		logger.Error("hex to proof error %v", err)
		return nil, err
	}
	proof, err := w.circuit.GenesisProve(first, second, req.GenesisID, req.FirstID, req.SecondID)
	if err != nil {
		logger.Error("genesis prove error %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	logger.Debug("complete  genesis prove %v", req.Period)
	return &rpc.SyncCommGenesisResponse{
		Version:   req.Version,
		Period:    req.Period,
		ProofType: common.SyncComGenesisType,
		Proof:     proofBytes,
		Witness:   witnessBytes,
	}, nil
}

func (w *LocalWorker) GenSyncCommitUnitProof(req rpc.SyncCommUnitsRequest) (*rpc.SyncCommUnitsResponse, error) {
	ok, err := common.VerifyLightClientUpdate(req.Data)
	if err != nil {
		logger.Error("verify light client update error %v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("verify light client update error")
	}

	unitProof, outerProof, err := w.circuit.UnitProve(req.Index, req.Data)
	if err != nil {
		logger.Error("unit prove error %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(unitProof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}

	outerProofBytes, outerWitnessBytes, err := circuits.PlonkProofToBytes(outerProof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	logger.Debug("complete unit prove %v", req.Index)
	return &rpc.SyncCommUnitsResponse{
		Version:      req.Version,
		Period:       req.Index,
		ProofType:    common.SyncComUnitType,
		Proof:        proofBytes,
		Witness:      witnessBytes,
		OuterProof:   outerProofBytes,
		OuterWitness: outerWitnessBytes,
	}, nil

}

func (w *LocalWorker) GenSyncCommRecursiveProof(req rpc.SyncCommRecursiveRequest) (*rpc.SyncCommRecursiveResponse, error) {
	logger.Debug("start recursive prove %v", req.Period)
	first, err := circuits.HexToProof(req.FirstProof)
	if err != nil {
		logger.Error("hex to proof error %v", err)
		return nil, err
	}
	second, err := circuits.HexToProof(req.SecondProof)
	if err != nil {
		logger.Error("hex to proof error %v", err)
		return nil, err
	}
	proof, err := w.circuit.RecursiveProve(req.Choice, first, second, req.BeginId, req.RelayId, req.EndId)
	if err != nil {
		logger.Error("recursive prove error %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	return &rpc.SyncCommRecursiveResponse{
		Version:   req.Version,
		Period:    req.Period,
		ProofType: common.SyncComRecursiveType,
		Proof:     proofBytes,
		Witness:   witnessBytes,
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

func (w *Worker) BtcDuperRecursiveProve(req *rpc.BtcDuperRecursiveRequest) (*rpc.ProofResponse, error) {
	return w.client.BtcDuperRecursiveProve(req)
}

func (w *Worker) BtcDepthRecursiveProve(req *rpc.BtcDepthRecursiveRequest) (*rpc.ProofResponse, error) {
	return w.client.BtcDepthRecursiveProve(req)
}

func (w *Worker) BtcChainProve(req *rpc.BtcChainRequest) (*rpc.ProofResponse, error) {
	return w.client.BtcChainProve(req)
}

func (w *Worker) BtcDepositProve(req *rpc.BtcDepositRequest) (*rpc.ProofResponse, error) {
	return w.client.BtcDepositProve(req)
}

func (w *Worker) BtcChangeProve(req *rpc.BtcChangeRequest) (*rpc.ProofResponse, error) {
	return w.client.BtcChangeProve(req)
}

func (w *Worker) BtcBaseProve(req *rpc.BtcBaseRequest) (*rpc.ProofResponse, error) {
	return w.client.BtcBaseProve(req)
}

func (w *Worker) BtcMiddleProve(req *rpc.BtcMiddleRequest) (*rpc.ProofResponse, error) {
	return w.client.BtcMiddleProve(req)
}

func (w *Worker) BtcUpperProve(req *rpc.BtcUpperRequest) (*rpc.ProofResponse, error) {
	return w.client.BtcUpperProve(req)
}

func (w *Worker) SupportProofType() []common.ZkProofType {
	return nil
}

func (w *Worker) Close() error {
	if w.client != nil {
		return w.client.Close()
	}
	return nil
}

func (w *Worker) BtcBulkProve(req *rpc.BtcBulkRequest) (*rpc.BtcBulkResponse, error) {
	return w.client.BtcBulkProve(req)
}

func (w *Worker) BtcPackedRequest(req *rpc.BtcPackedRequest) (*rpc.BtcPackResponse, error) {
	return w.client.BtcPackedRequest(req)
}

func (w *Worker) TxInEth2Prove(req *rpc.TxInEth2ProveRequest) (*rpc.TxInEth2ProveResponse, error) {
	return w.client.TxInEth2Prove(req)
}

func (w *Worker) BlockHeaderProve(req *rpc.BlockHeaderRequest) (*rpc.BlockHeaderResponse, error) {
	return w.client.BlockHeaderProve(req)
}

func (w *Worker) BlockHeaderFinalityProve(req *rpc.BlockHeaderFinalityRequest) (*rpc.BlockHeaderFinalityResponse, error) {
	return w.client.BlockHeaderFinalityProve(req)
}

func (w *Worker) ProofInfo(proofId string) (rpc.ProofInfo, error) {
	return w.client.ProofInfo(proofId)
}

func (w *Worker) GenRedeemProof(req *rpc.RedeemRequest) (*rpc.RedeemResponse, error) {
	return w.GenRedeemProof(req)
}

func (w *Worker) GenSyncCommGenesisProof(req rpc.SyncCommGenesisRequest) (*rpc.SyncCommGenesisResponse, error) {
	return w.client.GenSyncCommGenesisProof(req)
}

func (w *Worker) GenSyncCommitUnitProof(req rpc.SyncCommUnitsRequest) (*rpc.SyncCommUnitsResponse, error) {
	return w.client.GenSyncCommitUnitProof(req)
}

func (w *Worker) GenSyncCommRecursiveProof(req rpc.SyncCommRecursiveRequest) (*rpc.SyncCommRecursiveResponse, error) {
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
