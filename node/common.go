package node

import (
	"context"
	"encoding/hex"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	btcprovercom "github.com/lightec-xyz/btc_provers/circuits/common"
	btcprovertypes "github.com/lightec-xyz/btc_provers/circuits/types"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	btctx "github.com/lightec-xyz/daemon/rpc/bitcoin/common"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/rpc/oasis"
	"github.com/lightec-xyz/daemon/store"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"math/big"
	"os"
	"time"
)

func CheckProof(fileStore *FileStorage, zkType common.ZkProofType, index, end uint64, hash string) (bool, error) {
	switch zkType {
	case common.SyncComGenesisType:
		return fileStore.CheckGenesisProof()
	case common.SyncComUnitType:
		return fileStore.CheckUnitProof(index)
	case common.UnitOuter:
		return fileStore.CheckOuterProof(index)
	case common.SyncComRecursiveType:
		return fileStore.CheckRecursiveProof(index)
	case common.BeaconHeaderFinalityType:
		return fileStore.CheckBhfProof(index)
	case common.TxInEth2:
		return fileStore.CheckTxProof(hash)
	case common.BeaconHeaderType:
		return fileStore.CheckBeaconHeaderProof(index)
	case common.RedeemTxType:
		return fileStore.CheckRedeemProof(hash)
	case common.DepositTxType:
		return fileStore.CheckDepositProof(hash)
	case common.VerifyTxType:
		return fileStore.CheckVerifyProof(hash)
	case common.BtcBulkType:
		return fileStore.CheckBtcBulkProof(index, end)
	case common.BtcPackedType:
		return fileStore.CheckBtcPackedProof(index)
	case common.BtcWrapType:
		return fileStore.CheckBtcWrapProof(index)
	case common.BtcBaseType:
		return fileStore.CheckBtcBaseProof(index)
	case common.BtcMiddleType:
		return fileStore.CheckBtcMiddleProof(index)
	case common.BtcUpperType:
		return fileStore.CheckBtcUpperProof(index)
	case common.BtcGenesisType:
		return fileStore.CheckBtcGenesisProof()
	case common.BtcRecursiveType:
		return fileStore.CheckBtcRecursiveProof(index)
	default:
		return false, fmt.Errorf("unSupport now  proof type: %v", zkType.String())
	}
}

func StoreZkProof(fileStore *FileStorage, zkType common.ZkProofType, index, end uint64, hash string, proof, witness []byte) error {
	switch zkType {
	case common.SyncComUnitType:
		return fileStore.StoreUnitProof(index, proof, witness)
	case common.UnitOuter:
		return fileStore.StoreOuterProof(index, proof, witness)
	case common.SyncComGenesisType:
		return fileStore.StoreGenesisProof(index, proof, witness)
	case common.SyncComRecursiveType:
		return fileStore.StoreRecursiveProof(index, proof, witness)
	case common.BeaconHeaderFinalityType:
		return fileStore.StoreBhfProof(index, proof, witness)
	case common.TxInEth2:
		return fileStore.StoreTxProof(hash, proof, witness)
	case common.BeaconHeaderType:
		return fileStore.StoreBeaconHeaderProof(index, proof, witness)
	case common.RedeemTxType:
		return fileStore.StoreRedeemProof(hash, proof, witness)
	case common.DepositTxType:
		return fileStore.StoreDepositProof(hash, proof, witness)
	case common.VerifyTxType:
		return fileStore.StoreVerifyProof(hash, proof, witness)
	case common.BtcBulkType:
		return fileStore.StoreBtcBulkProof(index, end, proof, witness)
	case common.BtcPackedType:
		return fileStore.StoreBtcPackedProof(index, proof, witness)
	case common.BtcWrapType:
		return fileStore.StoreBtcWrapProof(index, proof, witness)
	case common.BtcBaseType:
		return fileStore.StoreBtcBaseProof(proof, witness, index)
	case common.BtcMiddleType:
		return fileStore.StoreBtcMiddleProof(proof, witness, index)
	case common.BtcUpperType:
		return fileStore.StoreBtcUpperProof(proof, witness, index)
	case common.BtcGenesisType:
		return fileStore.StoreBtcGenesisProof(proof, witness)
	case common.BtcRecursiveType:
		return fileStore.StoreBtcRecursiveProof(proof, witness, index)
	default:
		return fmt.Errorf("unSupport now  proof type: %v", zkType.String())
	}
}

func GenRequestData(p *PreparedData, reqType common.ZkProofType, index, end uint64, hash string) (interface{}, bool, error) {
	switch reqType {
	case common.SyncComGenesisType:
		data, ok, err := p.GetSyncComGenesisData()
		if err != nil {
			logger.Error("get SyncComGenesisData error:%v", err)
			return nil, false, err
		}
		return data, ok, nil
	case common.SyncComUnitType:
		data, ok, err := p.GetSyncComUnitData(index)
		if err != nil {
			logger.Error("get SyncComUnitData error:%v %v", index, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.SyncComRecursiveType:
		if p.genesisPeriod == index+2 {
			data, ok, err := p.GetRecursiveGenesisData(index)
			if err != nil {
				logger.Error("get SyncComRecursiveGenesisData error:%v %v", index, err)
				return nil, false, err
			}
			return data, ok, nil
		} else {
			data, ok, err := p.GetRecursiveData(index)
			if err != nil {
				logger.Error("get SyncComRecursiveData error:%v %v", index, err)
				return nil, false, err
			}
			return data, ok, nil
		}
	case common.TxInEth2:
		data, ok, err := p.GetTxInEth2Data(hash, p.getSlotByNumber)
		if err != nil {
			logger.Error("get tx in eth2 data error: %v %v", index, err)
			return nil, false, err
		}
		return data, ok, err
	case common.BeaconHeaderType:
		data, ok, err := p.GetBlockHeaderRequestData(index)
		if err != nil {
			logger.Error("get block header request data error:%v %v", index, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.BeaconHeaderFinalityType:
		data, ok, err := p.GetBhfUpdateData(index, p.genesisPeriod)
		if err != nil {
			logger.Error("get bhf update data error: %v %v", index, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.RedeemTxType:
		data, ok, err := p.GetRedeemRequestData(p.genesisPeriod, index, hash)
		if err != nil {
			logger.Error("get redeem request data error:%v %v", index, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.DepositTxType:
		data, ok, err := p.GetDepositData(hash)
		if err != nil {
			logger.Error("get deposit data error:%v %v", index, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.VerifyTxType:
		data, ok, err := p.GetVerifyData(hash)
		if err != nil {
			logger.Error("get verify data error:%v %v", index, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BtcBulkType:
		data, err := p.GetBtcMidBlockHeader(index, end)
		if err != nil {
			logger.Error("get mid block header error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return rpc.BtcBulkRequest{
			Data: data,
		}, true, nil
	case common.BtcPackedType:
		data, err := p.GetBtcMidBlockHeader(index, end)
		if err != nil {
			logger.Error("get mid block header error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return rpc.BtcPackedRequest{
			Data: data,
		}, true, nil
	case common.BtcWrapType:
		data, err := p.GetBtcWrapData(index, end)
		if err != nil {
			logger.Error("get btc wrap data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, true, nil
	case common.BtcBaseType:
		data, ok, err := p.GetBtcBaseData(index)
		if err != nil {
			logger.Error("get btc base data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.BtcMiddleType:
		data, ok, err := p.GetBtcMiddleData(index)
		if err != nil {
			logger.Error("get btc middle data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.BtcUpperType:
		data, ok, err := p.GetBtcUpperData(index)
		if err != nil {
			logger.Error("get btc upper data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.BtcGenesisType:
		data, ok, err := p.GetBtcGenesisData(index, end)
		if err != nil {
			logger.Error("get btc genesis data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil

	case common.BtcRecursiveType:
		data, ok, err := p.GetBtcRecursiveData(index, end)
		if err != nil {
			logger.Error("get btc recursive data error:%v %v %v", index, end, err)
			return nil, false, err
		}
		return data, ok, nil
	default:
		logger.Error(" prepare request Data never should happen : %v %v", index, reqType)
		return nil, false, fmt.Errorf("never should happen : %v %v", index, reqType)
	}

}

func WorkerGenProof(worker rpc.IWorker, request *common.ZkProofRequest) ([]*common.ZkProofResponse, error) {
	//defer worker.DelReqNum()
	var result []*common.ZkProofResponse
	switch request.ReqType {
	case common.DepositTxType:
		var depositRpcRequest rpc.DepositRequest
		err := common.ParseObj(request.Data, &depositRpcRequest)
		if err != nil {
			logger.Error("parse deposit Proof param error: %v", request.TxHash)
			return nil, fmt.Errorf("not deposit Proof param %v", request.TxHash)
		}
		proofResponse, err := worker.GenDepositProof(depositRpcRequest)
		if err != nil {
			logger.Error("gen deposit Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.End, request.TxHash, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)
	case common.VerifyTxType:
		var verifyRpcRequest rpc.VerifyRequest
		err := common.ParseObj(request.Data, &verifyRpcRequest)
		if err != nil {
			logger.Error("parse verify Proof param error: %v", request.TxHash)
			return nil, fmt.Errorf("not verify Proof param %v", request.TxHash)
		}
		proofResponse, err := worker.GenVerifyProof(verifyRpcRequest)
		if err != nil {
			logger.Error("gen verify Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.End, request.TxHash, proofResponse.Proof, proofResponse.Wit)
		result = append(result, zkbProofResponse)
	case common.TxInEth2:
		// todo
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
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.End, request.TxHash, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)
	case common.RedeemTxType:
		// todo
		var redeemRpcRequest rpc.RedeemRequest
		err := common.ParseObj(request.Data, &redeemRpcRequest)
		if err != nil {
			logger.Error("parse redeem Proof param error:%v", request.Id())
			return nil, fmt.Errorf("not redeem Proof param")
		}
		proofResponse, err := worker.GenRedeemProof(&redeemRpcRequest)
		if err != nil {
			logger.Error("gen redeem Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.End, request.TxHash, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)
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
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.End, request.TxHash, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)

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
		// todo
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.End, request.TxHash, proofResponse.Proof, proofResponse.Witness)
		outerProof := NewProofResp(common.UnitOuter, request.Index, request.End, request.TxHash, proofResponse.OuterProof, proofResponse.OuterWitness)
		result = append(result, zkbProofResponse)
		result = append(result, outerProof)
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
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.End, request.TxHash, proofResponse.Proof, proofResponse.Witness)
		result = append(result, zkbProofResponse)

	case common.BeaconHeaderType:
		// todo
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
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.End, request.TxHash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)
	case common.BeaconHeaderFinalityType:
		// todo
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
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.End, request.TxHash, response.Proof, response.Witness)
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
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.End, request.TxHash, response.Proof, response.Witness)
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
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.End, request.TxHash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)

	case common.BtcWrapType:
		var wrapRequest rpc.BtcWrapRequest
		err := common.ParseObj(request.Data, &wrapRequest)
		if err != nil {
			return nil, fmt.Errorf("parse btcWrapRequest error:%v", err)
		}
		response, err := worker.BtcWrapProve(&wrapRequest)
		if err != nil {
			logger.Error("gen btcWrap Proof error:%v", err)
			return nil, err
		}
		zkbProofResponse := NewProofResp(request.ReqType, request.Index, request.End, request.TxHash, response.Proof, response.Witness)
		result = append(result, zkbProofResponse)

	default:
		logger.Error("never should happen Proof type:%v", request.ReqType)
		return nil, fmt.Errorf("never should happen Proof type:%v", request.ReqType)

	}

	for _, item := range result {
		logger.Info("send zkProof:%v", item.Id())
	}
	return result, nil

}
func NewProofResp(reqType common.ZkProofType, index, end uint64, hash string, proof, witness []byte) *common.ZkProofResponse {
	return &common.ZkProofResponse{
		RespId:      common.NewProofId(reqType, index, end, hash),
		ZkProofType: reqType,
		Index:       index,
		End:         end,
		Proof:       proof,
		TxHash:      hash,
		Witness:     witness,
		Status:      common.ProofSuccess,
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

func GetBtcMidBlockHeader(client *bitcoin.Client, start, end uint64) (*btcprovertypes.BlockHeaderChain, error) {
	startHash, err := GetReverseHash(client, start)
	if err != nil {
		logger.Error("get block header error: %v %v", start, err)
		return nil, err
	}

	endHash, err := GetReverseHash(client, end)
	if err != nil {
		logger.Error("get block header error: %v %v", end, err)
		return nil, err
	}
	var middleHeaders []string
	for index := start + 1; index <= end; index++ {
		header, err := client.GetHexBlockHeader(int64(index))
		if err != nil {
			logger.Error("get block header error: %v %v", index, err)
			return nil, err
		}
		middleHeaders = append(middleHeaders, header)
	}
	data := &btcprovertypes.BlockHeaderChain{
		BeginHeight:        start,
		BeginHash:          startHash,
		EndHeight:          end,
		EndHash:            endHash,
		MiddleBlockHeaders: middleHeaders,
	}
	err = data.Verify()
	if err != nil {
		logger.Error("verify block header error: %v", err)
		return nil, err
	}
	return data, nil

}

func GetReverseHash(client *bitcoin.Client, height uint64) (string, error) {
	hash, err := client.GetBlockHash(int64(height))
	if err != nil {
		logger.Error("get block header error: %v %v", height, err)
		return "", err
	}
	reverseHash, err := common.ReverseHex(hash)
	if err != nil {
		logger.Error("reverse hex error: %v", err)
		return "", err
	}
	return reverseHash, nil
}

func GetBtcWrapData(filestore *FileStorage, client *bitcoin.Client, start, end uint64) (*rpc.BtcWrapRequest, error) {
	startHash, err := GetReverseHash(client, start)
	if err != nil {
		logger.Error("get block header error: %v %v", start, err)
		return nil, err
	}
	endHash, err := GetReverseHash(client, end)
	if err != nil {
		logger.Error("get block header error: %v %v", end, err)
		return nil, err
	}
	nRequired := end - start
	var proof *StoreProof
	var ok bool
	var flag string
	if nRequired <= btcprovercom.MaxNbBlockPerBulk { // todo
		proof, ok, err = filestore.GetBtcBulkProof(start, end)
		if err != nil {
			logger.Error("get btc bulk proof error: %v", err)
			return nil, err
		}
		if !ok {
			return nil, err
		}
		flag = circuits.BtcBulk
	} else {
		proof, ok, err = filestore.GetBtcPackedProof(start)
		if err != nil {
			logger.Error("get btc bulk proof error: %v", err)
			return nil, err
		}
		if !ok {
			return nil, err
		}
		flag = circuits.BtcPacked
	}
	data := &rpc.BtcWrapRequest{
		Flag:      flag,
		Proof:     proof.Proof,
		Witness:   proof.Witness,
		BeginHash: startHash,
		EndHash:   endHash,
		NbBlocks:  nRequired,
	}
	return data, nil
}

func GetRealEthNonce(client *ethereum.Client, store store.IStore, addr string) (uint64, error) {
	chainNonce, err := client.GetNonce(addr)
	if err != nil {
		logger.Error("get nonce error: %v %v", addr, err)
		return 0, err
	}
	dbNonce, exists, err := ReadNonce(store, "eth", addr)
	if err != nil {
		logger.Error("read nonce error: %v %v", addr, err)
		return 0, err
	}
	if !exists {
		return chainNonce, nil
	}
	if chainNonce < dbNonce {
		return dbNonce + 1, nil
	}
	return chainNonce, nil
}

// todo refactor

func RedeemBtcTx(btcClient *bitcoin.Client, ethClient *ethrpc.Client, oasisClient *oasis.Client, txHash string, proof []byte) (interface{}, error) {
	if common.GetEnvDebugMode() {
		return nil, nil
	}
	ethTxHash := ethcommon.HexToHash(txHash)
	ethTx, _, err := ethClient.TransactionByHash(context.Background(), ethTxHash)
	if err != nil {
		logger.Error("get eth tx error:%v", err)
		return nil, err
	}
	receipt, err := ethClient.TransactionReceipt(context.Background(), ethTxHash)
	if err != nil {
		logger.Error("get eth tx receipt error:%v", err)
		return nil, err
	}

	btcRawTx, _, err := ethereum.DecodeRedeemLog(receipt.Logs[3].Data)
	if err != nil {
		logger.Error("decode redeem log error:%v", err)
		return nil, err
	}
	logger.Info("btcRawTx: %v\n", hexutil.Encode(btcRawTx))
	rawTx, rawReceipt := ethereum.GetRawTxAndReceipt(ethTx, receipt)
	logger.Info("rawTx: %v\n", hexutil.Encode(rawTx))
	logger.Info("rawReceipt: %v\n", hexutil.Encode(rawReceipt))

	sigs, err := oasisClient.SignBtcTx(rawTx, rawReceipt, proof)
	if err != nil {
		logger.Error("sign btc tx error:%v", err)
		return nil, err
	}
	transaction := btctx.NewMultiTransactionBuilder()
	err = transaction.Deserialize(btcRawTx)
	if err != nil {
		logger.Error("deserialize btc tx error:%v", err)
		return nil, err
	}
	multiSigScript, err := ethClient.GetMultiSigScript()
	if err != nil {
		logger.Error("get multi sig script error:%v", err)
		return nil, err
	}

	// todo
	nTotal, nRequred := 3, 2
	transaction.AddMultiScript(multiSigScript, nRequred, nTotal)
	err = transaction.MergeSignature(sigs[:nRequred])
	if err != nil {
		logger.Error("merge signature error:%v", err)
		return nil, err
	}
	btxTx, err := transaction.Serialize()
	if err != nil {
		logger.Error("serialize btc tx error:%v", err)
		return nil, err
	}
	btcTxHash := transaction.TxHash()
	_, err = btcClient.GetTransaction(btcTxHash) // todo
	if err == nil {
		logger.Warn("btc tx already exist: %v", btcTxHash)
		return "", nil
	}
	txHex := hex.EncodeToString(btxTx)
	logger.Info("btc Tx: %v\n", txHex)
	TxHash, err := btcClient.Sendrawtransaction(txHex)
	if err != nil {
		logger.Error("send btc tx error:%v %v", btcTxHash, err)
		// todo  just test
		_, err = bitcoin.BroadcastTx(txHex)
		if err != nil {
			logger.Error("broadcast btc tx error %v:%v", btcTxHash, err)
			return "", err
		}
	}
	logger.Info("send redeem btc tx: %v", btcTxHash)
	return TxHash, nil
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
