package node

import (
	"context"
	"encoding/hex"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
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

func CheckProof(fileStore *FileStorage, zkType common.ZkProofType, index uint64, txHash string) (bool, error) {
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
	default:
		return false, fmt.Errorf("unSupport now  proof type: %v", zkType.String())
	}
}

func StoreZkProof(fileStore *FileStorage, zkType common.ZkProofType, index uint64, txHash string, proof, witness []byte) error {
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
	default:
		return fmt.Errorf("unSupport now  proof type: %v", zkType.String())
	}
}

// todo refactor

func RedeemBtcTx(btcClient *bitcoin.Client, txHash string, proof []byte) (interface{}, error) {
	ethTxHash := ethcommon.HexToHash(txHash)
	zkBridgeAddr, zkBtcAddr := "0x8e4f5a8f3e24a279d8ed39e868f698130777fded", "0xbf3041e37be70a58920a6fd776662b50323021c9"
	ec, err := ethrpc.NewClient("https://1rpc.io/holesky", zkBridgeAddr, zkBtcAddr)
	if err != nil {
		logger.Error("new eth client error:%v", err)
		return nil, err
	}
	ethTx, _, err := ec.TransactionByHash(context.Background(), ethTxHash)
	if err != nil {
		logger.Error("get eth tx error:%v", err)
		return nil, err
	}
	receipt, err := ec.TransactionReceipt(context.Background(), ethTxHash)
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

	btcSignerContract := "0x99e514Dc90f4Dd36850C893bec2AdC9521caF8BB"
	oasisClient, err := oasis.NewClient("https://testnet.sapphire.oasis.io", btcSignerContract)
	if err != nil {
		logger.Error("new client error:%v", err)
		return nil, err
	}

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

	multiSigScript, err := ec.GetMultiSigScript()
	if err != nil {
		logger.Error("get multi sig script error:%v", err)
		return nil, err
	}

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
	_, err = btcClient.GetTransaction(btcTxHash)
	if err != nil {
		logger.Error("get btc tx error:%v %v", btcTxHash, err)
		return "", nil
	}
	txHex := hex.EncodeToString(btxTx)
	logger.Info("btx Tx: %v\n", txHex)
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
	logger.Info("%v ticker goroutine start ...", name)
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
	logger.Info("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		case request := <-req:
			err := fn(request)
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}

	}
}

func doFetchRespTask(name string, resp chan FetchDataResponse, fn func(resp FetchDataResponse) error, exit chan os.Signal) {
	logger.Info("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
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
	logger.Info("%v goroutine start ...", name)
	for {
		select {
		case <-exit:
			logger.Info("%v goroutine exit now ...", name)
			return
		case response := <-resp:
			err := fn(response)
			if err != nil {
				logger.Error("%v error %v", name, err.Error())
			}
		}
	}
}
