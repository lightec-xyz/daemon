package node

import (
	"encoding/hex"
	"fmt"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"sync"
)

var _ rpc.IWorker = (*LocalWorker)(nil)

type LocalWorker struct {
	circuit     circuits.ICircuit
	dataDir     string
	maxNums     int
	currentNums int
	lock        sync.Mutex
	wid         string
}

func (w *LocalWorker) BtcTimestamp(req *rpc.BtcTimestampRequest) (*rpc.ProofResponse, error) {
	proof, err := w.circuit.BtcTimestamp(req.CpTime, req.SmoothData)
	if err != nil {
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

func (w *LocalWorker) BtcDuperRecursiveProve(req *rpc.BtcDuperRecursiveRequest) (*rpc.ProofResponse, error) {
	currentStep := req.End - req.Start // todo
	if currentStep == common.BtcUpperDistance || currentStep == common.BtcBaseDistance {
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
		proof, err := w.circuit.BtcChainRecursiveProve(req.FirstType, req.SecondType, req.FirstStep, req.SecondStep, req.BlockChainData, first, second)
		if err != nil {
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
	} else {
		first, err := circuits.HexToProof(req.First)
		if err != nil {
			logger.Error("btc recursive hex to proof error: %v", err)
			return nil, err
		}
		proof, err := w.circuit.BtcChainHybridProve(req.FirstType, req.FirstStep, req.HybridChainData, first)
		if err != nil {
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

}

func (w *LocalWorker) SyncCommInner(req *rpc.SyncCommInnerRequest) (*rpc.ProofResponse, error) {
	proof, err := w.circuit.SyncInnerProve(req.Index, req.Data)
	if err != nil {
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

func (w *LocalWorker) SyncCommOuter(req *rpc.SyncCommOuterRequest) (*rpc.ProofResponse, error) {
	innerProofs, err := circuits.HexToProofs(req.InnerProofs)
	if err != nil {
		logger.Error("sync inner hex to proofs error: %v", err)
		return nil, err
	}
	proof, err := w.circuit.SyncOutProve(req.Period, req.Data, innerProofs)
	if err != nil {
		logger.Error("sync out prove error: %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("sync out prove error: %v", err)
		return nil, err
	}
	return &rpc.ProofResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil
}

func (w *LocalWorker) BtcDepthRecursiveProve(req *rpc.BtcDepthRecursiveRequest) (*rpc.ProofResponse, error) {
	first, err := circuits.HexToProof(req.First)
	if err != nil {
		logger.Error("btc recursive hex to proof error: %v", err)
		return nil, err
	}
	proof, err := w.circuit.BtcDepthRecursiveProve(req.IsRecursive, req.PreStep, req.Data, first)
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
	panic("not support yet")

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
	timestamp, err := circuits.HexToProof(req.SigVerify)
	if err != nil {
		logger.Error("btc deposit hex to proof error: %v", err)
		return nil, err
	}
	proof, err := w.circuit.BtcDepositProve(req.ChainType, req.ChainStep, req.TxDepthStep, req.CpDepthStep, req.TxRecursive,
		req.CpRecursive, req.Data, blockChain, txDepth, cpDepth, timestamp, ethCommon.HexToAddress(req.ProverAddr), req.SmoothedTimestamp,
		req.CpFlag, req.SigVerifyData)
	if err != nil {
		logger.Error("btc deposit prove error: %v", err)
		return nil, err
	}
	proofSolBytes, err := circuits.ProofToSolBytes(proof.Proof)
	if err != nil {
		logger.Error("proof to sol bytes error: %v", err)
		return nil, err
	}
	witnessBytes, err := circuits.WitnessToBytes(proof.Witness)
	if err != nil {
		logger.Error("witness to bytes error: %v", err)
		return nil, err
	}
	return &rpc.ProofResponse{
		Proof:   proofSolBytes,
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
	timestamp, err := circuits.HexToProof(req.SigVerify)
	if err != nil {
		logger.Error("btc deposit hex to proof error: %v", err)
		return nil, err
	}
	minerRewardBytes, err := hex.DecodeString(req.MinerReward)
	if err != nil {
		logger.Error("btc change hex to proof error: %v", err)
		return nil, err
	}

	proof, err := w.circuit.BtcRedeemProve(req.ChainType, req.ChainStep, req.TxDepthStep, req.CpDepthStep, req.TxRecursive,
		req.CpRecursive, req.Data, blockChain, txDepth, cpDepth, redeem, timestamp, [32]byte(minerRewardBytes), ethCommon.HexToAddress(req.ProverAddr),
		req.SmoothedTimestamp, req.CpFlag, req.SigVerifyData)
	if err != nil {
		logger.Error("btc change prove error: %v", err)
		return nil, err
	}
	proofSolBytes, err := circuits.ProofToSolBytes(proof.Proof)
	if err != nil {
		logger.Error("proof to sol bytes error: %v", err)
		return nil, err
	}
	witnessBytes, err := circuits.WitnessToBytes(proof.Witness)
	if err != nil {
		logger.Error("witness to bytes error: %v", err)
		return nil, err
	}
	return &rpc.ProofResponse{
		Proof:   proofSolBytes,
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
		logger.Error("btc bulk prove error: %v", err)
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

func (w *LocalWorker) SupportProofType() []common.ProofType {
	return nil
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
	logger.Debug("start BlockHeaderProve %v", req.Data.BeginSlot)
	//logger.Debug("len: %v", len(headers.MiddleBeaconHeaders))
	proof, err := w.circuit.BeaconHeaderProve(req.Data)
	if err != nil {
		logger.Error("BlockHeaderProve error: %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	logger.Debug("complete BlockHeaderProve %v", req.Index)
	return &rpc.BlockHeaderResponse{
		Proof:   proofBytes,
		Witness: witnessBytes,
	}, nil

}

func (w *LocalWorker) BlockHeaderFinalityProve(req *rpc.BlockHeaderFinalityRequest) (*rpc.BlockHeaderFinalityResponse, error) {
	logger.Debug("start BlockHeaderFinalityProve %v", req.Index)
	proof, err := w.circuit.BeaconHeaderFinalityUpdateProve(req.FinalityUpdate, req.SyncCommittee)
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

func (w *LocalWorker) RedeemProof(req *rpc.RedeemRequest) (*rpc.RedeemResponse, error) {
	response, err := w.redeem(req, true)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (w *LocalWorker) BackendRedeemProof(req *rpc.RedeemRequest) (*rpc.RedeemResponse, error) {
	response, err := w.redeem(req, false)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (w *LocalWorker) redeem(req *rpc.RedeemRequest, front bool) (*rpc.RedeemResponse, error) {
	logger.Debug("start Redeem prove: %v %v", req.TxHash, front)
	tx, err := circuits.HexToProof(req.TxProof)
	if err != nil {
		logger.Error("gen Redeem proof error: %v", err)
		return nil, fmt.Errorf("gen Redeem proof error: %v", err)
	}
	bh, err := circuits.HexToProof(req.BhProof)
	if err != nil {
		logger.Error("gen Redeem proof error: %v", err)
		return nil, fmt.Errorf("gen Redeem proof error: %v", err)
	}
	bhf, err := circuits.HexToProof(req.BhfProof)
	if err != nil {
		logger.Error("gen Redeem proof error: %v", err)
		return nil, fmt.Errorf("gen Redeem proof error: %v", err)
	}
	duty, err := circuits.HexToProof(req.Duty)
	if err != nil {
		logger.Error("gen Redeem proof error: %v", err)
		return nil, fmt.Errorf("gen Redeem proof error: %v", err)
	}
	txId, err := hex.DecodeString(req.TxId)
	if err != nil {
		logger.Error("gen Redeem proof error: %v", err)
		return nil, fmt.Errorf("gen Redeem proof error: %v", err)
	}
	genesisScRoot, err := hex.DecodeString(req.GenesisScRoot)
	if err != nil {
		logger.Error(" genesisScRoot hex decode error: %v", err)
		return nil, fmt.Errorf(" genesisScRoot hex decode error: %v", err)
	}
	currentScRoot, err := hex.DecodeString(req.CurrentSCSSZRoot)
	if err != nil {
		logger.Error(" currentScRoot hex decode error: %v", err)
		return nil, fmt.Errorf(" currentScRoot hex decode error: %v", err)
	}

	minerRewardBytes, err := hex.DecodeString(req.MinerReward)
	if err != nil {
		logger.Error("gen Redeem proof error: %v", err)
		return nil, fmt.Errorf("gen Redeem proof error: %v", err)
	}
	var sigHashFixBytes [][32]byte
	for _, hash := range req.SigHashes {
		hashBytes, err := hex.DecodeString(hash)
		if err != nil {
			logger.Error("gen Redeem proof error: %v", err)
			return nil, fmt.Errorf("gen Redeem proof error: %v", err)
		}
		sigHashFixBytes = append(sigHashFixBytes, [32]byte(hashBytes))
	}
	logger.Debug(" worker Redeem prove genesisScSszRoot: %x, currentScSszRoot: %x,txid: %x, minerReward: %x,sigHashs: %x,nbBeaconHeaders: %v,isFront: %v",
		genesisScRoot, currentScRoot, txId, minerRewardBytes, sigHashFixBytes, req.NbBeaconHeaders, front)
	proof, err := w.circuit.RedeemProve(tx, bh, bhf, duty, genesisScRoot, currentScRoot, [32]byte(txId),
		[32]byte(minerRewardBytes), sigHashFixBytes, req.NbBeaconHeaders, front)
	if err != nil {
		logger.Error("gen Redeem proof error: %v", err)
		return nil, fmt.Errorf("gen Redeem proof error: %v", err)
	}
	var proofBytes []byte
	var proofSgxBytes []byte
	if front {
		proofBytes, err = circuits.ProofToSolBytes(proof.Proof)
		if err != nil {
			logger.Error("TxInEth2Prove error: %v", err)
			return nil, err
		}
		proofSgxBytes, err = circuits.ProofToBytes(proof.Proof)
		if err != nil {
			logger.Error("TxInEth2Prove error: %v", err)
			return nil, err
		}

	} else {
		proofBytes, err = circuits.ProofToBytes(proof.Proof)
		if err != nil {
			logger.Error("TxInEth2Prove error: %v", err)
			return nil, err
		}
	}
	witnessBytes, err := circuits.WitnessToBytes(proof.Witness)
	if err != nil {
		logger.Error("TxInEth2Prove error: %v", err)
		return nil, err
	}
	logger.Debug("complete gen Redeem proof: %v", req.TxHash)
	return &rpc.RedeemResponse{
		Proof:         proofBytes,
		Witness:       witnessBytes,
		ProofSgxBytes: proofSgxBytes,
	}, nil

}

func (w *LocalWorker) SyncCommitUnitProve(req rpc.SyncCommUnitsRequest) (*rpc.SyncCommUnitsResponse, error) {
	ok, err := req.Data.SyncCommitteeUpdate.Verify()
	if err != nil {
		logger.Error("verify light client update error %v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("verify light client update error")
	}
	outerProof, err := circuits.HexToProof(req.Outer)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	unitProof, err := w.circuit.SyncCommitteeUnitProve(req.Index, *outerProof, req.Data.SyncCommitteeUpdate)
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

func (w *LocalWorker) SyncCommDutyProve(req rpc.SyncCommDutyRequest) (*rpc.SyncCommDutyResponse, error) {
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
	outer, err := circuits.HexToProof(req.Outer)
	if err != nil {
		logger.Error("hex to proof error %v", err)
		return nil, err
	}
	genesisId, err := hex.DecodeString(req.BeginId)
	if err != nil {
		logger.Error("hex to proof error %v", err)
		return nil, err
	}
	relayId, err := hex.DecodeString(req.RelayId)
	if err != nil {
		logger.Error("hex to proof error %v", err)
		return nil, err
	}
	endId, err := hex.DecodeString(req.EndId)
	if err != nil {
		logger.Error("hex to proof error %v", err)
		return nil, err
	}
	proof, recursive, err := w.circuit.SyncCommitteeDutyProve(req.Choice, first, second, outer, genesisId, relayId, endId,
		req.ScIndex, req.Update)
	if err != nil {
		logger.Error("recursive prove error %v", err)
		return nil, err
	}
	proofBytes, witnessBytes, err := circuits.PlonkProofToBytes(proof)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	recursiveBytes, recursiveWitnessBytes, err := circuits.PlonkProofToBytes(recursive)
	if err != nil {
		logger.Error("btc genesis prove error: %v", err)
		return nil, err
	}
	return &rpc.SyncCommDutyResponse{
		Version:          req.Version,
		Period:           req.Period,
		ProofType:        common.SyncComDutyType,
		Proof:            proofBytes,
		Witness:          witnessBytes,
		RecursiveProof:   recursiveBytes,
		RecursiveWitness: recursiveWitnessBytes,
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

func (w *LocalWorker) Close() error {
	return nil
}
func (w *LocalWorker) Id() string {
	return w.wid
}

func (w *LocalWorker) ProofInfo(proofId string) (rpc.ProofInfo, error) {
	logger.Debug("Proof info")
	return rpc.ProofInfo{
		Status: 0,
		Proof:  "",
		TxId:   proofId,
	}, nil
}

func NewLocalWorker(btcSetupDir, ethSetupDir, dataDir, wid string, maxNums, cacheCap int) (rpc.IWorker, error) {
	config := circuits.CircuitConfig{
		EthSetupDir: ethSetupDir,
		BtcSetupDir: btcSetupDir,
		Debug:       common.GetEnvDebugMode(),
		CacheCap:    cacheCap,
	}
	circuit, err := circuits.NewCircuit(&config)
	if err != nil {
		return nil, err
	}
	return &LocalWorker{
		dataDir:     dataDir,
		maxNums:     maxNums,
		currentNums: 0,
		wid:         wid,
		circuit:     circuit,
	}, nil
}

var _ rpc.IWorker = (*Worker)(nil)

type Worker struct {
	client      rpc.IProof
	maxNums     int
	currentNums int
	lock        sync.Mutex
	wid         string
}

func (w *Worker) BtcTimestamp(req *rpc.BtcTimestampRequest) (*rpc.ProofResponse, error) {
	return w.client.BtcTimestamp(req)
}

func (w *Worker) SyncCommOuter(req *rpc.SyncCommOuterRequest) (*rpc.ProofResponse, error) {
	return w.client.SyncCommOuter(req)
}

func (w *Worker) BtcDuperRecursiveProve(req *rpc.BtcDuperRecursiveRequest) (*rpc.ProofResponse, error) {
	return w.client.BtcDuperRecursiveProve(req)
}

func (w *Worker) SyncCommInner(req *rpc.SyncCommInnerRequest) (*rpc.ProofResponse, error) {
	return w.client.SyncCommInner(req)
}

func (w *Worker) BackendRedeemProof(req *rpc.RedeemRequest) (*rpc.RedeemResponse, error) {
	return w.client.BackendRedeemProof(req)
}

func (w *Worker) BtcDepthRecursiveProve(req *rpc.BtcDepthRecursiveRequest) (*rpc.ProofResponse, error) {
	return w.client.BtcDepthRecursiveProve(req)
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

func (w *Worker) SupportProofType() []common.ProofType {
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

func (w *Worker) RedeemProof(req *rpc.RedeemRequest) (*rpc.RedeemResponse, error) {
	return w.RedeemProof(req)
}

func (w *Worker) SyncCommitUnitProve(req rpc.SyncCommUnitsRequest) (*rpc.SyncCommUnitsResponse, error) {
	return w.client.SyncCommitUnitProve(req)
}

func (w *Worker) SyncCommDutyProve(req rpc.SyncCommDutyRequest) (*rpc.SyncCommDutyResponse, error) {
	return w.client.SyncCommDutyProve(req)
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
