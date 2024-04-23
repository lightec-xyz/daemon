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
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"os"
	"time"
)

func GetSyncCommitUpdate(fileStore *FileStore, period uint64) (*utils.LightClientUpdateInfo, bool, error) {
	var currentPeriodUpdate structs.LightClientUpdateWithVersion
	exists, err := fileStore.GetUpdate(period, &currentPeriodUpdate)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v period update Data, send new update request", period)
		return nil, false, nil
	}
	var update utils.LightClientUpdateInfo
	err = ParseObj(currentPeriodUpdate.Data, &update)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if fileStore.GetGenesisPeriod() == period {
		var genesisData structs.LightClientBootstrapResponse
		genesisExists, err := fileStore.GetBootstrap(&genesisData)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		if !genesisExists {
			logger.Warn("no find genesis update Data, send new update request")
			return nil, false, nil
		}
		// todo
		var genesisCommittee utils.SyncCommittee
		err = common.ParseObj(genesisData.Data.CurrentSyncCommittee, &genesisCommittee)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		update.CurrentSyncCommittee = &genesisCommittee
	} else {
		prePeriod := period - 1
		if prePeriod < fileStore.GetGenesisPeriod() {
			logger.Error("should never happen: %v", prePeriod)
			return nil, false, nil
		}
		var preUpdateData structs.LightClientUpdateWithVersion
		preUpdateExists, err := fileStore.GetUpdate(prePeriod, &preUpdateData)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		if !preUpdateExists {
			logger.Warn("get unit Data,no find %v period update Data, send new update request", prePeriod)
			return nil, false, nil
		}
		// todo
		var currentSyncCommittee utils.SyncCommittee
		err = common.ParseObj(preUpdateData.Data.NextSyncCommittee, &currentSyncCommittee)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		update.CurrentSyncCommittee = &currentSyncCommittee
	}
	return &update, true, nil

}

func CheckProof(fileStore *FileStore, zkType common.ZkProofType, index uint64, txHash string) (bool, error) {
	switch zkType {
	case common.TxInEth2:
		return fileStore.CheckTxProof(txHash)
	case common.BlockHeaderType:
		return fileStore.CheckBlockHeaderProof(index)
	case common.RedeemTxType:
		return fileStore.CheckRedeemProof(txHash)
	default:
		return false, fmt.Errorf("unSupport now  proof type: %v", zkType)
	}
}

func StoreZkProof(fileStore *FileStore, zkType common.ZkProofType, index uint64, txHash string, proof, witness []byte) error {
	switch zkType {
	case common.TxInEth2:
		return fileStore.StoreTxProof(txHash, proof, witness)
	case common.BlockHeaderType:
		return fileStore.StoreBlockHeaderProof(index, proof, witness)
	case common.RedeemTxType:
		return fileStore.StoreRedeemProof(txHash, proof, witness)
	default:
		return fmt.Errorf("unSupport now  proof type: %v", zkType)
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
	txHex := hex.EncodeToString(btxTx)
	logger.Info("btx Tx: %v\n", txHex)
	TxHash, err := btcClient.Sendrawtransaction(txHex)
	if err != nil {
		logger.Error("send btc tx error:%v", err)
		// todo  just test
		_, err = bitcoin.BroadcastTx(txHex)
		if err != nil {
			logger.Error("broadcast btc tx error:%v", err)
			return "", err
		}
	}
	logger.Info("send redeem btc tx: %v", transaction.TxHash())
	return TxHash, nil
}

func doTask(name string, fn func() error, exit chan os.Signal) {
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

func doTimerTask(name string, interval time.Duration, fn func() error, exit chan os.Signal) {
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
