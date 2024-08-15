package node

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
)

func CheckProof(fileStore *FileStorage, zkType common.ZkProofType, index, end uint64, hash string) (bool, error) {
	switch zkType {
	case common.SyncComOuterType:
		return fileStore.CheckOuterProof(index)
	case common.SyncComUnitType:
		return fileStore.CheckUnitProof(index)
	case common.SyncComGenesisType:
		return fileStore.CheckGenesisProof()
	case common.SyncComRecursiveType:
		return fileStore.CheckRecursiveProof(index)
	case common.TxInEth2:
		return fileStore.CheckTxProof(hash)
	case common.BeaconHeaderType:
		return fileStore.CheckBeaconHeaderProof(index, end)
	case common.BeaconHeaderFinalityType:
		return fileStore.CheckBhfProof(index)
	case common.RedeemTxType:
		return fileStore.CheckRedeemProof(hash)
	case common.BtcBulkType:
		return fileStore.CheckBtcBulkProof(index, end)
	case common.BtcPackedType:
		return fileStore.CheckBtcPackedProof(index)
	case common.BtcBaseType:
		return fileStore.CheckBtcBaseProof(index, end)
	case common.BtcMiddleType:
		return fileStore.CheckBtcMiddleProof(index, end)
	case common.BtcUpperType:
		return fileStore.CheckBtcUpperProof(index, end)
	case common.BtcDuperGenesisType:
		return fileStore.CheckBtcDuperGenesisProof()
	case common.BtcDuperRecursive:
		return fileStore.CheckBtcDuperRecursiveProof(index, end)
	case common.BtcDepthGenesisType:
		return fileStore.CheckBtcDepthGenesisProof()
	case common.BtcDepthRecursiveType:
		return fileStore.CheckBtcDepthRecursiveProof(index, end)
	case common.BtcChainType:
		return fileStore.CheckBtcBlockChainProof(index, end)
	case common.BtcDepositType:
		return fileStore.CheckBtcDepositProof(index, end)
	case common.BtcChangeType:
		return fileStore.CheckBtcChangeProof(index, end)
	default:
		return false, fmt.Errorf("unSupport now  proof type: %v", zkType.String())
	}
}

func StoreZkProof(fileStore *FileStorage, zkType common.ZkProofType, index, end uint64, hash string, proof, witness []byte) error {
	switch zkType {
	case common.SyncComOuterType:
		return fileStore.StoreOuterProof(index, proof, witness)
	case common.SyncComUnitType:
		return fileStore.StoreUnitProof(index, proof, witness)
	case common.SyncComGenesisType:
		return fileStore.StoreGenesisProof(proof, witness)
	case common.SyncComRecursiveType:
		return fileStore.StoreRecursiveProof(index, proof, witness)
	case common.TxInEth2:
		return fileStore.StoreTxProof(hash, proof, witness)
	case common.BeaconHeaderType:
		return fileStore.StoreBeaconHeaderProof(index, end, proof, witness)
	case common.BeaconHeaderFinalityType:
		return fileStore.StoreBhfProof(index, proof, witness)
	case common.RedeemTxType:
		return fileStore.StoreRedeemProof(hash, proof, witness)
	case common.BtcBulkType:
		return fileStore.StoreBtcBulkProof(index, end, proof, witness)
	case common.BtcPackedType:
		return fileStore.StoreBtcPackedProof(index, proof, witness)
	case common.BtcBaseType:
		return fileStore.StoreBtcBaseProof(proof, witness, index, end)
	case common.BtcMiddleType:
		return fileStore.StoreBtcMiddleProof(proof, witness, index, end)
	case common.BtcUpperType:
		return fileStore.StoreBtcUpperProof(proof, witness, index, end)
	case common.BtcDuperGenesisType:
		return fileStore.StoreBtcDuperGenesisProof(proof, witness)
	case common.BtcDuperRecursive:
		return fileStore.StoreBtcDuperRecursiveProof(proof, witness, index, end)
	case common.BtcDepthGenesisType:
		return fileStore.StoreBtcDepthGenesisProof(proof, witness)
	case common.BtcDepthRecursiveType:
		return fileStore.StoreBtcDepthRecursiveProof(proof, witness, index, end)
	case common.BtcChainType:
		return fileStore.StoreBtcBlockChainProof(proof, witness, index, end)
	case common.BtcDepositType:
		return fileStore.StoreBtcDepositProof(proof, witness, index, end)
	case common.BtcChangeType:
		return fileStore.StoreBtcChangeProof(proof, witness, index, end)
	default:
		return fmt.Errorf("unSupport now  proof type: %v", zkType.String())
	}
}

func GenRequestData(p *Prepared, reqType common.ZkProofType, index, end uint64, hash string) (interface{}, bool, error) {
	switch reqType {
	case common.SyncComUnitType:
		data, ok, err := p.GetSyncComUnitRequest(index)
		if err != nil {
			logger.Error("get SyncComUnitData error:%v %v", index, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.SyncComGenesisType:
		data, ok, err := p.GetSyncComGenesisRequest()
		if err != nil {
			logger.Error("get SyncComGenesisData error:%v", err)
			return nil, false, err
		}
		return data, ok, nil

	case common.SyncComRecursiveType:
		if index == p.genesisPeriod+2 { // todo
			data, ok, err := p.GetRecursiveGenesisRequest(index)
			if err != nil {
				logger.Error("get SyncComRecursiveGenesisData error:%v %v", index, err)
				return nil, false, err
			}
			return data, ok, nil
		} else {
			data, ok, err := p.GetRecursiveRequest(index)
			if err != nil {
				logger.Error("get SyncComRecursiveData error:%v %v", index, err)
				return nil, false, err
			}
			return data, ok, nil
		}
	case common.TxInEth2:
		data, ok, err := p.GetTxInEth2Request(hash, p.getSlotByNumber)
		if err != nil {
			logger.Error("get tx in eth2 data error: %v %v", index, err)
			return nil, false, err
		}
		return data, ok, err
	case common.BeaconHeaderType:
		data, ok, err := p.GetBlockHeaderRequest(index)
		if err != nil {
			logger.Error("get block header request data error:%v %v", index, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.BeaconHeaderFinalityType:
		data, ok, err := p.GetBhfUpdateRequest(index, p.genesisPeriod)
		if err != nil {
			logger.Error("get bhf update data error: %v %v", index, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.RedeemTxType:
		data, ok, err := p.GetRedeemRequest(p.genesisPeriod, index, hash)
		if err != nil {
			logger.Error("get redeem request data error:%v %v", index, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcBulkType:
		data, err := p.GetBtcBulkRequest(index, end)
		if err != nil {
			logger.Error("get mid block header error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, true, nil
	case common.BtcPackedType:
		data, ok, err := p.GetBtcPackRequest(index, end)
		if err != nil {
			logger.Error("get mid block header error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcBaseType:
		data, ok, err := p.GetBtcBaseRequest(end)
		if err != nil {
			logger.Error("get btc base data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.BtcMiddleType:
		data, ok, err := p.GetBtcMiddleRequest(end)
		if err != nil {
			logger.Error("get btc middle data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.BtcUpperType:
		data, ok, err := p.GetBtcUpperRequest(end)
		if err != nil {
			logger.Error("get btc upper data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcDuperGenesisType:
		data, ok, err := p.GetBtcDuperGenesisRequest()
		if err != nil {
			logger.Error("get btc genesis data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcDuperRecursive:
		data, ok, err := p.GetBtcDuperRecursiveRequest(index)
		if err != nil {
			logger.Error("get btc duper recursive data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcDepthGenesisType:
		data, ok, err := p.BtcDepthGenesisRequest()
		if err != nil {
			logger.Error("get btc depth genesis data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcDepthRecursiveType:
		data, ok, err := p.GetBtcDepthRecursiveRequest(index, end)
		if err != nil {
			logger.Error("get btc depth recursive data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcChainType:
		data, ok, err := p.GetBtcChainRequest(index, end)
		if err != nil {
			logger.Error("get btc chain data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcDepositType:
		data, ok, err := p.GetBtcDepositRequest(hash)
		if err != nil {
			logger.Error("get btc deposit data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcChangeType:
		data, ok, err := p.GetBtcChangeRequest(hash)
		if err != nil {
			logger.Error("get btc change data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil
	default:
		logger.Error(" prepare request Data never should happen : %v %v", index, reqType)
		return nil, false, fmt.Errorf("never should happen : %v %v", index, reqType)
	}

}

func WorkerGenProof(worker rpc.IWorker, request *common.ZkProofRequest) ([]*common.ZkProofResponse, error) {
	var result []*common.ZkProofResponse
	switch request.ProofType {
	case common.SyncComUnitType:
		var commUnitsRequest rpc.SyncCommUnitsRequest
		err := common.ParseObj(request.Data, &commUnitsRequest)
		if err != nil {
			return nil, fmt.Errorf("not sync comm unit Proof param")
		}
		proofResponse, err := worker.GenSyncCommitUnitProof(commUnitsRequest)
		if err != nil {
			logger.Error("gen sync comm unit Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, proofResponse.Proof, proofResponse.Witness)
		outerProof := NewProofResp(common.SyncComOuterType, request.Index, request.SIndex, request.Hash, proofResponse.OuterProof, proofResponse.OuterWitness)
		result = append(result, zkbProofResponse)
		result = append(result, outerProof)
	case common.SyncComGenesisType:
		var genesisRpcRequest rpc.SyncCommGenesisRequest
		err := common.ParseObj(request.Data, &genesisRpcRequest)
		if err != nil {
			return nil, fmt.Errorf("not genesis Proof param")
		}
		proofResponse, err := worker.GenSyncCommGenesisProof(genesisRpcRequest)
		if err != nil {
			logger.Error("gen sync comm genesis Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)
	case common.SyncComRecursiveType:
		var recursiveRequest rpc.SyncCommRecursiveRequest
		err := common.ParseObj(request.Data, &recursiveRequest)
		if err != nil {
			return nil, fmt.Errorf("not sync comm recursive Proof param")
		}
		proofResponse, err := worker.GenSyncCommRecursiveProof(recursiveRequest)
		if err != nil {
			logger.Error("gen sync comm recursive Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)

	case common.TxInEth2:
		var txInEth2Req rpc.TxInEth2ProveRequest
		err := common.ParseObj(request.Data, &txInEth2Req)
		if err != nil {
			logger.Error("parse txInEth2 Proof param error:%v", err)
			return nil, fmt.Errorf("not txInEth2 Proof param")
		}
		proofResponse, err := worker.TxInEth2Prove(&txInEth2Req)
		if err != nil {
			logger.Error("gen redeem Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)
	case common.BeaconHeaderType:
		var blockHeaderRequest rpc.BlockHeaderRequest
		err := common.ParseObj(request.Data, &blockHeaderRequest)
		if err != nil {
			logger.Error("not block header Proof param")
			return nil, fmt.Errorf("not block header Proof param")
		}
		response, err := worker.BlockHeaderProve(&blockHeaderRequest)
		if err != nil {
			logger.Error("gen block header Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)
	case common.BeaconHeaderFinalityType:
		var finalityRequest rpc.BlockHeaderFinalityRequest
		err := common.ParseObj(request.Data, &finalityRequest)
		if err != nil {
			return nil, fmt.Errorf("not block header finality Proof param")
		}
		response, err := worker.BlockHeaderFinalityProve(&finalityRequest)
		if err != nil {
			logger.Error("gen block header finality Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)

	case common.RedeemTxType:
		var redeemRpcRequest rpc.RedeemRequest
		err := common.ParseObj(request.Data, &redeemRpcRequest)
		if err != nil {
			logger.Error("parse redeem Proof param error:%v", request.RequestId())
			return nil, fmt.Errorf("not redeem Proof param")
		}
		proofResponse, err := worker.GenRedeemProof(&redeemRpcRequest)
		if err != nil {
			logger.Error("gen redeem Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)

	case common.BtcPackedType:
		var packedRequest rpc.BtcPackedRequest
		err := common.ParseObj(request.Data, &packedRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcPackedRequest error:%v", err)
		}
		response, err := worker.BtcPackedRequest(&packedRequest)
		if err != nil {
			logger.Error("gen btcPacked Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)
	case common.BtcBulkType:
		var bulkRequest rpc.BtcBulkRequest
		err := common.ParseObj(request.Data, &bulkRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcBulkRequest error:%v", err)
		}
		response, err := worker.BtcBulkProve(&bulkRequest)
		if err != nil {
			logger.Error("gen btcBulk Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)

	case common.BtcBaseType:
		var baseRequest rpc.BtcBaseRequest
		err := common.ParseObj(request.Data, &baseRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcBaseRequest error:%v", err)
		}
		response, err := worker.BtcBaseProve(&baseRequest)
		if err != nil {
			logger.Error("gen btcBase Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)
	case common.BtcMiddleType:
		var middleRequest rpc.BtcMiddleRequest
		err := common.ParseObj(request.Data, &middleRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcMiddleRequest error:%v", err)
		}
		response, err := worker.BtcMiddleProve(&middleRequest)
		if err != nil {
			logger.Error("gen btcMiddle Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)
	case common.BtcUpperType:
		var upperRequest rpc.BtcUpperRequest
		err := common.ParseObj(request.Data, &upperRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcUpperRequest error:%v", err)
		}
		response, err := worker.BtcUpperProve(&upperRequest)
		if err != nil {
			logger.Error("gen btcUpper Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)
	case common.BtcDuperGenesisType:
		var req rpc.BtcDuperRecursiveRequest
		err := common.ParseObj(req.Data, &req)
		if err != nil {
			return nil, fmt.Errorf("parse btcDuperRecursiveRequest error:%v", err)
		}
		response, err := worker.BtcDuperRecursiveProve(&req)
		if err != nil {
			logger.Error("gen btcDuperRecursive Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)
	case common.BtcDuperRecursive:
		var req rpc.BtcDuperRecursiveRequest
		err := common.ParseObj(req.Data, &req)
		if err != nil {
			return nil, fmt.Errorf("parse btcDuperRecursiveRequest error:%v", err)
		}
		response, err := worker.BtcDuperRecursiveProve(&req)
		if err != nil {
			logger.Error("gen btcDuperRecursive Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)

	case common.BtcDepthGenesisType:
		var req rpc.BtcDepthRecursiveRequest
		err := common.ParseObj(req.Data, &req)
		if err != nil {
			return nil, fmt.Errorf("parse btcDepthRecursiveRequest error:%v", err)
		}
		response, err := worker.BtcDepthRecursiveProve(&req)
		if err != nil {
			logger.Error("gen btcDepthRecursive Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)

	case common.BtcDepthRecursiveType:
		var req rpc.BtcDepthRecursiveRequest
		err := common.ParseObj(req.Data, &req)
		if err != nil {
			return nil, fmt.Errorf("parse btcDepthRecursiveRequest error:%v", err)
		}
		response, err := worker.BtcDepthRecursiveProve(&req)
		if err != nil {
			logger.Error("gen btcDepthRecursive Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)

	case common.BtcChainType:
		var req rpc.BtcChainRequest
		err := common.ParseObj(req.Data, &req)
		if err != nil {
			return nil, fmt.Errorf("parse btcChainRequest error:%v", err)
		}
		response, err := worker.BtcChainProve(&req)
		if err != nil {
			logger.Error("gen btcChain Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)
	case common.BtcDepositType:
		var req rpc.BtcDepositRequest
		err := common.ParseObj(req.Data, &req)
		if err != nil {
			return nil, fmt.Errorf("parse btcDepositRequest error:%v", err)
		}
		response, err := worker.BtcDepositProve(&req)
		if err != nil {
			logger.Error("gen btcDeposit Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)

	case common.BtcChangeType:
		var req rpc.BtcChangeRequest
		err := common.ParseObj(req.Data, &req)
		if err != nil {
			return nil, fmt.Errorf("parse btcChangeRequest error:%v", err)
		}
		response, err := worker.BtcChangeProve(&req)
		if err != nil {
			logger.Error("gen btcChange Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ProofType, request.Index, request.SIndex, request.Hash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)
	default:
		logger.Error("never should happen Proof type:%v", request.ProofType)
		return nil, fmt.Errorf("never should happen Proof type:%v", request.ProofType)

	}
	for _, item := range result {
		logger.Info("send zkProof:%v", item.RespId())
	}
	return result, nil

}
func NewProofResp(reqType common.ZkProofType, index, end uint64, hash string, proof, witness []byte) *common.ZkProofResponse {
	return &common.ZkProofResponse{
		Id:        common.NewProofId(reqType, index, end, hash),
		ProofType: reqType,
		Index:     index,
		SIndex:    end,
		Proof:     proof,
		Hash:      hash,
		Witness:   witness,
		Status:    common.ProofSuccess,
	}
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

func DoTimerTask(name string, interval time.Duration, fn func() error, exit chan os.Signal) {
	logger.Debug("%v ticker goroutine start ...", name)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		case <-ticker.C:
			err := fn()
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}

func doProofRequestTask(name string, req chan []*common.ZkProofRequest, fn func(req []*common.ZkProofRequest) error, exit chan os.Signal) {
	logger.Debug("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v proof request goroutine exit now ...", name)
			return
		case request := <-req:
			err := fn(request)
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}

	}
}

func doFetchRespTask(name string, resp chan *FetchResponse, fn func(resp *FetchResponse) error, exit chan os.Signal) {
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

func doProofResponseTask(name string, resp chan *common.ZkProofResponse, fn func(resp *common.ZkProofResponse) error, exit chan os.Signal) {
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
