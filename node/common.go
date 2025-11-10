package node

import (
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	blockdepthUtil "github.com/lightec-xyz/btc_provers/utils/blockdepth"
	btcproverClient "github.com/lightec-xyz/btc_provers/utils/client"
	"github.com/lightec-xyz/daemon/node/p2p"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/ethereum/zkbridge"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
)

// todo
func getProofParams(txId, miner, network string, chainStore *ChainStore, btcClient *bitcoin.Client, proverClient btcproverClient.IClient) (*zkbridge.IBtcTxVerifierPublicWitnessParams, error) {
	dbTx, ok, err := chainStore.ReadBtcTx(txId)
	if err != nil {
		logger.Error("read btc tx error: %v %v", txId, err)
		return nil, err
	}
	if !ok {
		logger.Warn("no find btc tx: %v", txId)
		return nil, fmt.Errorf("no find btc tx:%v", txId)
	}
	if dbTx.LatestHeight == 0 {
		//logger.Warn("dbTx %v latest height is 0", txId)
		return nil, fmt.Errorf("dbTx %v latest height is 0,wait for update", txId)
	}

	cpDepth := dbTx.LatestHeight - dbTx.CheckPointHeight
	txDepth := dbTx.LatestHeight - dbTx.Height
	cpHash, ok, err := chainStore.ReadCheckpoint(dbTx.CheckPointHeight)
	if err != nil {
		logger.Error("%v", err.Error())
		return nil, err
	}
	btcTx, err := btcClient.GetRawTransaction(dbTx.Hash)
	if err != nil {
		logger.Error("get transaction error: %v %v", dbTx.Hash, err)
		return nil, err
	}
	blockHash := common.ReverseBytes(ethcommon.FromHex(btcTx.Blockhash))

	icpSignature, ok, err := chainStore.ReadIcpSignature(dbTx.LatestHeight)
	if err != nil {
		logger.Error("read dfinity sign error: %v", err)
		return nil, err
	}
	if !ok {
		logger.Warn("no find: %v icp %v signature", dbTx.Hash, dbTx.LatestHeight)
		// no work,just placeholder
		icpSignature.Hash = "6aeb6ec6f0fbc707b91a3bec690ae6536fe0abaa1994ef24c3463eb20494785d"
		icpSignature.Signature = "3f8e02c743e76a4bd655873a428db4fa2c46ac658854ba38f8be0fbbf9af9b2b6b377aaaaf231b6b890a5ee3c15a558f1ccc18dae0c844b6f06343b88a8d12e3"
	} else {
		//logger.Debug("%v icp signature: %v %v %v", txId, icpSignature.Height, icpSignature.Hash, icpSignature.Signature)
	}
	smoothedTimestamp, err := blockdepthUtil.GetSmoothedTimestampProofData(proverClient, uint32(dbTx.LatestHeight))
	if err != nil {
		logger.Error("%v", err.Error())
		return nil, err
	}
	cptimeData, err := blockdepthUtil.GetCpTimestampProofData(proverClient, uint32(dbTx.Height))
	if err != nil {
		logger.Error("%v", err)
		return nil, err
	}
	sigVerif, err := blockdepthUtil.GetSigVerifProofData(
		common.ReverseBytes(ethcommon.FromHex(icpSignature.Hash)),
		ethcommon.FromHex(icpSignature.Signature),
		ethcommon.FromHex(getIcpPublicKey(network)))
	if err != nil {
		logger.Error("%v", err.Error())
		return nil, err
	}
	flag := cptimeData.Flag<<1 | sigVerif.Flag
	params := &zkbridge.IBtcTxVerifierPublicWitnessParams{
		Checkpoint:        [32]byte(ethcommon.FromHex(cpHash)),
		CpDepth:           uint32(cpDepth),
		TxDepth:           uint32(txDepth),
		TxBlockHash:       [32]byte(blockHash),
		TxTimestamp:       uint32(btcTx.Blocktime),
		ZkpMiner:          ethcommon.HexToAddress(miner),
		Flag:              big.NewInt(int64(flag)),
		SmoothedTimestamp: smoothedTimestamp.Timestamp,
	}
	return params, nil
}

func upperRoundStartIndex(height uint64) uint64 {
	index := height / common.BtcUpperDistance * common.BtcUpperDistance
	return index
}

func DbValue(a string) string {
	return strings.ToLower(trimOx(a))
}

func stepEndIndex(start, end, step uint64) uint64 {
	stepNums := (end - start) / step
	if stepNums > 0 {
		return start + stepNums*step
	}
	return start
}

func BlockDepthPlan(prefix, start, end uint64, skip ...bool) []ChainIndex {
	var indexes []ChainIndex
	var tmpIndex uint64
	for _, step := range common.BtcBlockDepthPlan {
		nextStartIndex := stepEndIndex(start, end, step)
		for i := start; i < nextStartIndex; i += step {
			indexes = append(indexes, ChainIndex{
				Genesis: prefix,
				Start:   i,
				End:     i + step,
				Step:    step,
			})
		}
		start = nextStartIndex
		tmpIndex = nextStartIndex
		if nextStartIndex == end {
			return indexes
		}
	}
	if len(skip) > 0 && skip[0] {
		return indexes
	}
	indexes = append(indexes, ChainIndex{
		Genesis: prefix,
		Start:   tmpIndex,
		End:     end,
		Step:    end - tmpIndex,
	})
	return indexes
}

func BlockChainPlan(start, height uint64, skip ...bool) []ChainIndex {
	var indexes []ChainIndex
	startIndex := (start / common.BtcUpperDistance) * common.BtcUpperDistance
	endIndex := startIndex + common.BtcUpperDistance
	if endIndex >= height {
		//within one upper round
		indexes = append(indexes, BlockUpperIndex(start, height, skip...)...)
	} else {
		// more one upper round
		currentIndex := start
		if start%common.BtcUpperDistance != 0 {
			// current upper round
			tmpIndexes := BlockUpperIndex(currentIndex, endIndex, false)
			indexes = append(indexes, tmpIndexes...)
			currentIndex = endIndex
		}
		//how much upper round
		roundNums := (height - currentIndex) / common.BtcUpperDistance
		for i := uint64(0); i < roundNums; i++ {
			indexes = append(indexes, ChainIndex{
				Start: currentIndex,
				End:   currentIndex + common.BtcUpperDistance,
				Step:  common.BtcUpperDistance,
			})
			currentIndex = currentIndex + common.BtcUpperDistance
		}
		//finally upper round
		if currentIndex < height {
			indexes = append(indexes, BlockUpperIndex(currentIndex, height, skip...)...)
		}
	}

	return indexes
}

// BlockUpperIndex one upper round plan
func BlockUpperIndex(start, end uint64, skip ...bool) []ChainIndex {
	blockChainPlan := []uint32{common.BtcBaseDistance}
	blockChainPlan = append(blockChainPlan, common.BlockChainPlan[:]...)
	blockChainPlan = append(blockChainPlan, common.CapacityMiniLevel)
	var indexes []ChainIndex
	var tmpIndex uint64
	for _, step := range blockChainPlan {
		nextStartIndex := stepEndIndex(start, end, uint64(step))
		for i := start; i < nextStartIndex; i += uint64(step) {
			indexes = append(indexes, ChainIndex{
				Start: i,
				End:   i + uint64(step),
				Step:  uint64(step),
			})
		}
		start = nextStartIndex
		tmpIndex = nextStartIndex
		if nextStartIndex == end {
			return indexes
		}
	}
	if len(skip) > 0 && skip[0] {
		return indexes
	}
	indexes = append(indexes, ChainIndex{
		Start: tmpIndex,
		End:   end,
		Step:  end - tmpIndex,
	})
	return indexes

}

func GenRequestData(p *Prepared, reqType common.ProofType, fIndex, sIndex uint64, hash string, prefix uint64, isCp bool) (interface{}, bool, error) {
	switch reqType {
	case common.SyncComInnerType:
		data, ok, err := p.GetSyncComInnerRequest(prefix, fIndex)
		if err != nil {
			logger.Error("get syncCommittee inner data error:%v %v", fIndex, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.SyncComOuterType:
		data, ok, err := p.GetSyncOuterRequest(fIndex)
		if err != nil {
			logger.Error("get syncCommittee outer data error:%v %v", fIndex, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.SyncComUnitType:
		data, ok, err := p.GetSyncComUnitRequest(fIndex)
		if err != nil {
			logger.Error("get SyncComUnitData error:%v %v", fIndex, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.SyncComDutyType:
		data, ok, err := p.GetDutyRequest(fIndex)
		if err != nil {
			logger.Error("get SyncComRecursiveData error:%v %v", fIndex, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.TxInEth2Type:
		data, ok, err := p.GetTxInEth2Request(hash, p.getSlotByNumber)
		if err != nil {
			logger.Error("get tx in eth2 data error: %v %v", fIndex, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BeaconHeaderType:
		data, ok, err := p.GetBlockHeaderRequest(fIndex)
		if err != nil {
			logger.Error("get block header request data error:%v %v", fIndex, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BeaconHeaderFinalityType:
		data, ok, err := p.GetBhfUpdateRequest(fIndex)
		if err != nil {
			logger.Error("get bhf update data error: %v %v", fIndex, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.RedeemTxType:
		data, ok, err := p.GetRedeemRequest(hash)
		if err != nil {
			logger.Error("get Redeem request data error:%v %v", fIndex, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BackendRedeemTxType:
		data, ok, err := p.GetRedeemRequest(hash)
		if err != nil {
			logger.Error("get Redeem request data error:%v %v", fIndex, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcBulkType:
		data, err := p.GetBtcBulkRequest(fIndex, sIndex, prefix)
		if err != nil {
			logger.Error("get mid block header error:%v %v %v", fIndex, sIndex, err)
			return nil, false, err
		}
		return data, true, nil
	case common.BtcBaseType:
		data, ok, err := p.GetBtcBaseRequest(fIndex, sIndex)
		if err != nil {
			logger.Error("get btc base data error:%v %v %v", fIndex, sIndex, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcMiddleType:
		data, ok, err := p.GetBtcMiddleRequest(fIndex, sIndex)
		if err != nil {
			logger.Error("get btc middle data error:%v %v %v", fIndex, sIndex, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcUpperType:
		data, ok, err := p.GetBtcUpperRequest(fIndex, sIndex)
		if err != nil {
			logger.Error("get btc upper data error:%v %v %v", fIndex, sIndex, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.BtcDuperRecursiveType:
		data, ok, err := p.GetBtcDuperRecursiveRequest(fIndex, sIndex)
		if err != nil {
			logger.Error("get btc duper recursive data error:%v %v %v", fIndex, sIndex, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.BtcDepthRecursiveType:
		data, ok, err := p.GetBtcDepthRecursiveRequest(prefix, fIndex, sIndex, isCp)
		if err != nil {
			logger.Error("get btc depth recursive data error:%v %v %v", fIndex, sIndex, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcTimestampType:
		data, ok, err := p.GetBtcTimestampRequest(fIndex, sIndex)
		if err != nil {
			logger.Error("get btc timestamp data error:%v %v %v", fIndex, sIndex, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcDepositType, common.BtcUpdateCpType:
		data, ok, err := p.GetBtcDepositRequest(hash)
		if err != nil {
			logger.Error("get btc deposit data error:%v %v %v", fIndex, sIndex, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcChangeType:
		data, ok, err := p.GetBtcChangeRequest(hash)
		if err != nil {
			logger.Error("get btc change data error:%v %v %v", fIndex, sIndex, err)
			return nil, false, err
		}
		return data, ok, nil
	default:
		logger.Error(" prepare request Responses never should happen : %v %v", fIndex, reqType)
		return nil, false, fmt.Errorf("never should happen : %v %v", fIndex, reqType)
	}

}

func WorkerGenProof(worker rpc.IWorker, req *common.ProofRequest) ([]*common.ProofResponse, error) {
	if v, ok := req.Data.(rpc.ICheck); ok { //todo
		err := v.Check()
		if err != nil {
			logger.Error("%v %v", req.ProofId(), err)
			return nil, err
		}
	}
	var result []*common.ProofResponse
	switch req.ProofType {
	case common.SyncComInnerType:
		var innerRequest rpc.SyncCommInnerRequest
		err := common.ParseObj(req.Data, &innerRequest)
		if err != nil {
			return nil, fmt.Errorf("not sync comm inner Proof param")
		}
		response, err := worker.SyncCommInner(&innerRequest)
		if err != nil {
			logger.Error("gen sync comm inner Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix,
			req.FIndex, req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)

	case common.SyncComOuterType:
		var outerRequest rpc.SyncCommOuterRequest
		err := common.ParseObj(req.Data, &outerRequest)
		if err != nil {
			return nil, fmt.Errorf("not sync comm outer Proof param")
		}
		response, err := worker.SyncCommOuter(&outerRequest)
		if err != nil {
			logger.Error("gen sync comm outer Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix,
			req.FIndex, req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)

	case common.SyncComUnitType:
		var commUnitsRequest rpc.SyncCommUnitsRequest
		err := common.ParseObj(req.Data, &commUnitsRequest)
		if err != nil {
			return nil, fmt.Errorf("not sync comm unit Proof param")
		}
		proofResponse, err := worker.SyncCommitUnitProve(commUnitsRequest)
		if err != nil {
			logger.Error("gen sync comm unit Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := common.NewProofResponse(req.ProofType, proofResponse.Proof, proofResponse.Witness, req.Prefix,
			req.FIndex, req.SIndex, req.Hash, req.CreateTime)
		result = append(result, zkbProofResponse)
	case common.SyncComDutyType:
		var recursiveRequest rpc.SyncCommDutyRequest
		err := common.ParseObj(req.Data, &recursiveRequest)
		if err != nil {
			return nil, fmt.Errorf("not sync comm recursive Proof param")
		}
		resp, err := worker.SyncCommDutyProve(recursiveRequest)
		if err != nil {
			logger.Error("gen sync comm recursive Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, resp.Proof, resp.Witness, req.Prefix, req.FIndex, req.SIndex,
			req.Hash, req.CreateTime)
		recursiveResp := common.NewProofResponse(common.SyncComRecursiveType, resp.RecursiveProof, resp.RecursiveWitness,
			req.Prefix, req.FIndex, req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)
		result = append(result, recursiveResp)

	case common.TxInEth2Type:
		var txInEth2Req rpc.TxInEth2ProveRequest
		err := common.ParseObj(req.Data, &txInEth2Req)
		if err != nil {
			logger.Error("parse txInEth2 Proof param error:%v", err)
			return nil, fmt.Errorf("not txInEth2 Proof param")
		}
		response, err := worker.TxInEth2Prove(&txInEth2Req)
		if err != nil {
			logger.Error("gen Redeem Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)
	case common.BeaconHeaderType:
		var blockHeaderRequest rpc.BlockHeaderRequest
		err := common.ParseObj(req.Data, &blockHeaderRequest)
		if err != nil {
			logger.Error("not block header Proof param")
			return nil, fmt.Errorf("not block header Proof param")
		}
		response, err := worker.BlockHeaderProve(&blockHeaderRequest)
		if err != nil {
			logger.Error("gen block header Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)
	case common.BeaconHeaderFinalityType:
		var finalityRequest rpc.BlockHeaderFinalityRequest
		err := common.ParseObj(req.Data, &finalityRequest)
		if err != nil {
			return nil, fmt.Errorf("not block header finality Proof param")
		}
		response, err := worker.BlockHeaderFinalityProve(&finalityRequest)
		if err != nil {
			logger.Error("gen block header finality Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)

	case common.RedeemTxType:
		var redeemRpcRequest rpc.RedeemRequest
		err := common.ParseObj(req.Data, &redeemRpcRequest)
		if err != nil {
			logger.Error("parse Redeem Proof param error:%v", req.ProofId())
			return nil, fmt.Errorf("not Redeem Proof param")
		}
		response, err := worker.RedeemProof(&redeemRpcRequest)
		if err != nil {
			logger.Error("gen Redeem Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)
		sgxRedeemResp := common.NewProofResponse(common.SgxRedeemTxType, response.ProofSgxBytes, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, sgxRedeemResp)
	case common.BackendRedeemTxType:
		var redeemRpcRequest rpc.RedeemRequest
		err := common.ParseObj(req.Data, &redeemRpcRequest)
		if err != nil {
			logger.Error("parse Redeem Proof param error:%v", req.ProofId())
			return nil, fmt.Errorf("not Redeem Proof param")
		}
		response, err := worker.BackendRedeemProof(&redeemRpcRequest)
		if err != nil {
			logger.Error("gen Redeem Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)
	case common.BtcBulkType:
		var bulkRequest rpc.BtcBulkRequest
		err := common.ParseWithNumber(req.Data, &bulkRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcBulkRequest error:%v", err)
		}
		response, err := worker.BtcBulkProve(&bulkRequest)
		if err != nil {
			logger.Error("gen btcBulk Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)

	case common.BtcBaseType:
		var baseRequest rpc.BtcBaseRequest
		err := common.ParseWithNumber(req.Data, &baseRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcBaseRequest error:%v", err)
		}
		response, err := worker.BtcBaseProve(&baseRequest)
		if err != nil {
			logger.Error("gen btcBase Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)
	case common.BtcMiddleType:
		var middleRequest rpc.BtcMiddleRequest
		err := common.ParseWithNumber(req.Data, &middleRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcMiddleRequest error:%v", err)
		}
		response, err := worker.BtcMiddleProve(&middleRequest)
		if err != nil {
			logger.Error("gen btcMiddle Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)
	case common.BtcUpperType:
		var upperRequest rpc.BtcUpperRequest
		err := common.ParseWithNumber(req.Data, &upperRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcUpperRequest error:%v", err)
		}
		response, err := worker.BtcUpperProve(&upperRequest)
		if err != nil {
			logger.Error("gen btcUpper Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)

	case common.BtcDuperRecursiveType:
		var duperRequest rpc.BtcDuperRecursiveRequest
		err := common.ParseWithNumber(req.Data, &duperRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcDuperRecursiveRequest error:%v", err)
		}
		response, err := worker.BtcDuperRecursiveProve(&duperRequest)
		if err != nil {
			logger.Error("gen btcDuperRecursive Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)
	case common.BtcDepthRecursiveType:
		var depthRequest rpc.BtcDepthRecursiveRequest
		err := common.ParseWithNumber(req.Data, &depthRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcDepthRecursiveRequest error:%v", err)
		}
		response, err := worker.BtcDepthRecursiveProve(&depthRequest)
		if err != nil {
			logger.Error("gen btcDepthRecursive Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)
	case common.BtcTimestampType:
		var timestampRequest rpc.BtcTimestampRequest
		err := common.ParseWithNumber(req.Data, &timestampRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcTimestampRequest error:%v", err)
		}
		//todo
		err = timestampRequest.Check()
		if err != nil {
			logger.Error("%v", err)
			return nil, err
		}
		response, err := worker.BtcTimestamp(&timestampRequest)
		if err != nil {
			logger.Error("gen btcTimestamp Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)
	case common.BtcDepositType, common.BtcUpdateCpType:
		var depositRequest rpc.BtcDepositRequest
		err := common.ParseWithNumber(req.Data, &depositRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcDepositRequest error:%v", err)
		}
		response, err := worker.BtcDepositProve(&depositRequest)
		if err != nil {
			logger.Error("gen btcDeposit Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)
	case common.BtcChangeType:
		var changeRequest rpc.BtcChangeRequest
		err := common.ParseWithNumber(req.Data, &changeRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcChangeRequest error:%v", err)
		}
		response, err := worker.BtcChangeProve(&changeRequest)
		if err != nil {
			logger.Error("gen btcChange Proof error:%v", err)
			return nil, err
		}
		proofResponse := common.NewProofResponse(req.ProofType, response.Proof, response.Witness, req.Prefix, req.FIndex,
			req.SIndex, req.Hash, req.CreateTime)
		result = append(result, proofResponse)
	default:
		logger.Error("never should happen Proof type:%v", req.ProofType)
		return nil, fmt.Errorf("never should happen Proof type:%v", req.ProofType)

	}
	for _, item := range result {
		logger.Info("send zkProof:%v", item.ProofId())
	}
	return result, nil

}

func BeaconBlockHeaderToSlotAndRoot(header *structs.BeaconBlockHeader) (uint64, []byte, error) {
	consensus, err := header.ToConsensus()
	if err != nil {
		logger.Error("to consensus error: %v", err)
		return 0, nil, err
	}
	root, err := consensus.HashTreeRoot()
	if err != nil {
		logger.Error("hash tree root error: %v", err)
		return 0, nil, err
	}
	slotBig, ok := big.NewInt(0).SetString(header.Slot, 10)
	if !ok {
		logger.Error("parse slot error: %v", header.Slot)
		return 0, nil, fmt.Errorf("parse slot error: %v", header.Slot)
	}
	return slotBig.Uint64(), root[0:], nil

}

func DoTask(name string, fn func() error, exit chan os.Signal) {
	defer PrintPanicStack()
	logger.Info("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		default:
			err := fn()
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}

func DoTimerTask(name string, interval time.Duration, fn func() error, exit chan os.Signal, notifies ...chan *Notify) {
	defer PrintPanicStack()
	logger.Debug("%v ticker goroutine start ...", name)
	ticker := time.NewTicker(interval)
	var lock sync.Mutex
	doFn := func() {
		lock.Lock()
		err := fn()
		lock.Unlock()
		if err != nil {
			logger.Error("%v error %v", name, err.Error())
		}
	}
	notify := make(chan *Notify, 1)
	if len(notifies) > 0 {
		notify = notifies[0]
	}
	defer ticker.Stop()
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		case <-notify:
			doFn()
		case <-ticker.C:
			doFn()
		}
	}
}

func doChainForkTask(name string, req chan *ChainFork, fn func(req *ChainFork) error, exit chan os.Signal) {
	defer PrintPanicStack()
	logger.Debug("%v goroutine start ...", name)
	var lock sync.Mutex
	for {
		select {
		case <-exit:
			logger.Info("%v proof request goroutine exit now ...", name)
			return
		case request := <-req:
			lock.Lock()
			err := fn(request)
			lock.Unlock()
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}

	}
}

func doFetchRespTask(name string, resp chan *FetchResponse, fn func(resp *FetchResponse) error, exit chan os.Signal) {
	defer PrintPanicStack()
	logger.Debug("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Debug("%v fetch goroutine exit now ...", name)
			return
		case response := <-resp:
			err := fn(response)
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}

func doLibP2pMsgTask(name string, msg <-chan *p2p.Msg, fn func(msg *p2p.Msg) error, exit chan os.Signal) {
	defer PrintPanicStack()
	logger.Debug("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Debug("%v libp2p msg goroutine exit now ...", name)
			return
		case res := <-msg:
			err := fn(res)
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}

func doReScanTask(name string, resp chan *ReScnSignal, fn func(height uint64) error, exit chan os.Signal) {
	defer PrintPanicStack()
	logger.Debug("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v proof resp goroutine exit now ...", name)
			return
		case response := <-resp:
			err := fn(response.Height)
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}

func doProofResponseTask(name string, resp chan *common.ProofResponse, fn func(resp *common.ProofResponse) error, exit chan os.Signal) {
	defer PrintPanicStack()
	logger.Debug("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v proof resp goroutine exit now ...", name)
			return
		case response := <-resp:
			err := fn(response)
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}
