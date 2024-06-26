package node

import (
	"context"
	"encoding/hex"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lightec-xyz/btc_provers/circuits/constant"
	btcprovertypes "github.com/lightec-xyz/btc_provers/circuits/types"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/rpc/oasis"
	btctx "github.com/lightec-xyz/daemon/transaction/bitcoin"
	"github.com/lightec-xyz/daemon/transaction/ethereum"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"math/big"
	"os"
	"time"
)

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

func CheckProof(fileStore *FileStorage, zkType common.ZkProofType, index, end uint64, txHash string) (bool, error) {
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
		return fileStore.CheckTxProof(txHash)
	case common.BeaconHeaderType:
		return fileStore.CheckBeaconHeaderProof(index)
	case common.RedeemTxType:
		return fileStore.CheckRedeemProof(txHash)
	case common.DepositTxType:
		return fileStore.CheckDepositProof(txHash)
	case common.VerifyTxType:
		return fileStore.CheckVerifyProof(txHash)
	case common.BtcBulkType:
		return fileStore.CheckBtcBulkProof(index, end)
	case common.BtcPackedType:
		return fileStore.CheckBtcPackedProof(index)
	case common.BtcWrapType:
		return fileStore.CheckBtcWrapProof(index)
	default:
		return false, fmt.Errorf("unSupport now  proof type: %v", zkType.String())
	}
}

func StoreZkProof(fileStore *FileStorage, zkType common.ZkProofType, index, end uint64, txHash string, proof, witness []byte) error {
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
		return fileStore.StoreTxProof(txHash, proof, witness)
	case common.BeaconHeaderType:
		return fileStore.StoreBeaconHeaderProof(index, proof, witness)
	case common.RedeemTxType:
		return fileStore.StoreRedeemProof(txHash, proof, witness)
	case common.DepositTxType:
		return fileStore.StoreDepositProof(txHash, proof, witness)
	case common.VerifyTxType:
		return fileStore.StoreVerifyProof(txHash, proof, witness)
	case common.BtcBulkType:
		return fileStore.StoreBtcBulkProof(index, end, proof, witness)
	case common.BtcPackedType:
		return fileStore.StoreBtcPackedProof(index, proof, witness)
	case common.BtcWrapType:
		return fileStore.StoreBtcWrapProof(index, proof, witness)
	default:
		return fmt.Errorf("unSupport now  proof type: %v", zkType.String())
	}
}

func GetBtcMidBlockHeader(client *bitcoin.Client, start, end uint64) (*rpc.BtcBulkRequest, error) {
	startHash, err := client.GetBlockHash(int64(start))
	if err != nil {
		logger.Error("get block header error: %v %v", start, err)
		return nil, err
	}
	beginHash, err := common.ReverseHex(startHash)
	if err != nil {
		logger.Error("reverse hex error: %v", err)
		return nil, err
	}
	endHash, err := client.GetBlockHash(int64(end))
	if err != nil {
		logger.Error("get block header error: %v %v", end, err)
		return nil, err
	}
	eHash, err := common.ReverseHex(endHash)
	if err != nil {
		logger.Error("reverse hex error: %v", err)
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
		BeginHash:          beginHash,
		EndHeight:          end,
		EndHash:            eHash,
		MiddleBlockHeaders: middleHeaders,
	}
	err = data.Verify()
	if err != nil {
		logger.Error("verify block header error: %v", err)
		return nil, err
	}
	return &rpc.BtcBulkRequest{
		Data: data,
	}, nil

}

func GetBtcWrapData(filestore *FileStorage, client *bitcoin.Client, start, end uint64) (*rpc.BtcWrapRequest, error) {
	startHash, err := client.GetBlockHash(int64(start))
	if err != nil {
		logger.Error("get block header error: %v %v", start, err)
		return nil, err
	}
	endHash, err := client.GetBlockHash(int64(end))
	if err != nil {
		logger.Error("get block header error: %v %v", end, err)
		return nil, err
	}
	nRequired := start - end
	var proof *StoreProof
	var ok bool
	var flag string
	if nRequired <= constant.MaxNbBlockPerBulk { // todo
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
