package node

import (
	"context"
	"encoding/hex"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/rpc/oasis"
	btctx "github.com/lightec-xyz/daemon/transaction/bitcoin"
	"github.com/lightec-xyz/daemon/transaction/ethereum"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"os"
	"strconv"
	"time"
)

func GetBhfUpdateData(fileStore *FileStorage, slot uint64) (interface{}, bool, error) {
	logger.Debug("get bhf update data: %v", slot)
	genesisPeriod := fileStore.GetGenesisPeriod()
	var currentFinalityUpdate structs.LightClientUpdateWithVersion
	exists, err := fileStore.GetFinalityUpdate(slot, &currentFinalityUpdate)
	if err != nil {
		logger.Error("get finality update error: %v %v", slot, err)
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find finality update: %v", slot)
		return nil, false, nil
	}

	genesisId, ok, err := GetSyncCommitRootId(fileStore, genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v period genesis commitId no find", genesisPeriod)
		return nil, false, nil
	}
	// todo
	attestedSlot, err := strconv.ParseUint(currentFinalityUpdate.Data.AttestedHeader.Slot, 10, 64)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	period := (attestedSlot / 8192)
	logger.Debug("get bhf update data slot: %v,period: %v", slot, period)
	recursiveProof, ok, err := fileStore.GetRecursiveProof(period)
	if err != nil {
		logger.Error("get recursive proof error: %v %v", period, err)
		return nil, false, err
	}
	if !ok {
		logger.Warn("no find recursive proof: %v", period)
		return nil, false, nil
	}

	outerProof, ok, err := fileStore.GetOuterProof(period)
	if err != nil {
		logger.Error("get outer proof error: %v %v", period, err)
		return nil, false, err
	}
	if !ok {
		logger.Warn("no find outer proof: %v", period)
		return nil, false, nil
	}

	var finalUpdate proverType.FinalityUpdate
	err = common.ParseObj(currentFinalityUpdate.Data, &finalUpdate)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	finalUpdate.Version = currentFinalityUpdate.Version

	currentSyncCommitUpdate, ok, err := GetSyncCommitUpdate(fileStore, period)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Error("no find sync committee update: %v", period)
		return nil, false, nil
	}

	var scUpdate proverType.SyncCommitteeUpdate
	err = common.ParseObj(currentSyncCommitUpdate, &scUpdate)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	request := rpc.BlockHeaderFinalityRequest{
		GenesisSCSSZRoot: fmt.Sprintf("%x", genesisId),
		RecursiveProof:   recursiveProof.Proof,
		RecursiveWitness: recursiveProof.Witness,
		OuterProof:       outerProof.Proof,
		OuterWitness:     outerProof.Witness,
		FinalityUpdate:   &finalUpdate,
		ScUpdate:         &scUpdate,
	}
	return &request, true, nil

}

func GetRecursiveData(fileStore *FileStorage, period uint64) (interface{}, bool, error) {
	//todo
	genesisPeriod := fileStore.GetGenesisPeriod()
	genesisId, ok, err := GetSyncCommitRootId(fileStore, genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v period genesis commitId no find", genesisPeriod)
		return nil, false, nil
	}
	relayId, ok, err := GetSyncCommitRootId(fileStore, period)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v period relay commitId no find", period)
		return nil, false, nil
	}
	endPeriod := period + 1
	endId, ok, err := GetSyncCommitRootId(fileStore, endPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v period end commitId no find", endPeriod)
		return nil, false, nil
	}
	secondProof, exists, err := fileStore.GetUnitProof(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v unit proof Data, send new proof request", period)
		return nil, false, nil
	}

	prePeriod := period - 1
	firstProof, exists, err := fileStore.GetRecursiveProof(prePeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v period recursive Data, send new proof request", prePeriod)
		return nil, false, nil
	}
	return &rpc.SyncCommRecursiveRequest{
		Choice:        "recursive",
		FirstProof:    firstProof.Proof,
		FirstWitness:  firstProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		BeginId:       genesisId,
		RelayId:       relayId,
		EndId:         endId,
	}, true, nil
}

func GetRecursiveGenesisData(fileStore *FileStorage, period uint64) (interface{}, bool, error) {
	genesisPeriod := fileStore.GetGenesisPeriod()
	genesisId, ok, err := GetSyncCommitRootId(fileStore, genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v period genesis commitId no find", genesisPeriod)
		return nil, false, nil
	}
	relayPeriod := period
	relayId, ok, err := GetSyncCommitRootId(fileStore, relayPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v  period relay commitId no find ", relayPeriod)
		return nil, false, nil
	}
	endPeriod := relayPeriod + 1
	endId, ok, err := GetSyncCommitRootId(fileStore, endPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v period end commitId no find", endPeriod)
		return nil, false, nil
	}

	fistProof, firstExists, err := fileStore.GetGenesisProof()
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !firstExists {
		logger.Warn("no find genesis proof ,start new proof request")
		return nil, false, nil
	}
	secondProof, secondExists, err := fileStore.GetUnitProof(relayPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !secondExists {
		logger.Warn("no find %v period unit proof , send new proof request", relayPeriod)
		return nil, false, nil
	}
	return &rpc.SyncCommRecursiveRequest{
		Choice:        "genesis",
		FirstProof:    fistProof.Proof,
		FirstWitness:  fistProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		BeginId:       genesisId,
		RelayId:       relayId,
		EndId:         endId,
	}, true, nil

}

func GetGenesisData(fileStore *FileStorage) (*rpc.SyncCommGenesisRequest, bool, error) {
	genesisPeriod := fileStore.GetGenesisPeriod()
	genesisId, ok, err := GetSyncCommitRootId(fileStore, genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v period genesis commitId  no find", genesisPeriod)
		return nil, false, nil
	}

	nextPeriod := genesisPeriod + 1
	firstId, ok, err := GetSyncCommitRootId(fileStore, nextPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v period first commitId no find", nextPeriod)
		return nil, false, nil
	}
	secondPeriod := nextPeriod + 1
	secondId, ok, err := GetSyncCommitRootId(fileStore, secondPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v period second commitId no find", secondPeriod)
		return nil, false, nil
	}

	firstProof, exists, err := fileStore.GetUnitProof(genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("get genesis Data,first proof not exists: %v period", genesisPeriod)
		return nil, false, nil
	}
	logger.Info("get genesis first proof: %v", genesisPeriod)

	secondProof, exists, err := fileStore.GetUnitProof(nextPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}

	if !exists {
		logger.Warn("get genesis Data,second proof not exists: %v period", nextPeriod)
		return nil, false, nil
	}
	logger.Info("get genesis second proof: %v", nextPeriod)
	genesisProofParam := &rpc.SyncCommGenesisRequest{
		FirstProof:    firstProof.Proof,
		FirstWitness:  firstProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		GenesisID:     genesisId,
		FirstID:       firstId,
		SecondID:      secondId,
	}
	return genesisProofParam, true, nil

}

func GetSyncCommitRootId(fileStore *FileStorage, period uint64) ([]byte, bool, error) {
	update, ok, err := GetSyncCommitUpdate(fileStore, period)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	syncCommitRoot, err := circuits.SyncCommitRoot(update)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	return syncCommitRoot, true, nil
}

func GetSyncCommitUpdate(fileStore *FileStorage, period uint64) (*utils.LightClientUpdateInfo, bool, error) {
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
	err = common.ParseObj(currentPeriodUpdate.Data, &update)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	update.Version = currentPeriodUpdate.Version
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

func CheckProof(fileStore *FileStorage, zkType common.ZkProofType, index uint64, txHash string) (bool, error) {
	switch zkType {
	case common.SyncComGenesisType:
		return fileStore.CheckGenesisProof()
	case common.SyncComUnitType:
		return fileStore.CheckUnitProof(index)
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
	default:
		return false, fmt.Errorf("unSupport now  proof type: %v", zkType)
	}
}

func StoreZkProof(fileStore *FileStorage, zkType common.ZkProofType, index uint64, txHash string, proof, witness []byte) error {
	switch zkType {
	case common.SyncComUnitType:
		return fileStore.StoreUnitProof(index, proof, witness)
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
