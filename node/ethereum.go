package node

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/wire"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	btctx "github.com/lightec-xyz/daemon/rpc/bitcoin/common"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	redeemUtils "github.com/lightec-xyz/provers/utils/redeem-tx"
	"math/big"
	"strings"
	"time"
)

type ethereumAgent struct {
	btcClient       *bitcoin.Client
	ethClient       *ethrpc.Client
	fileStore       *FileStorage
	chainStore      *ChainStore
	taskManager     *TxManager
	ethFilter       *EthFilter
	initHeight      uint64
	txManager       *TxManager
	chainForkSignal chan<- *ChainFork
	curHeight       uint64
	reScan          bool
}

func NewEthereumAgent(cfg Config, fileStore *FileStorage, store store.IStore, btcClient *bitcoin.Client,
	ethClient *ethrpc.Client, task *TxManager, chainFork chan *ChainFork) (IAgent, error) {
	return &ethereumAgent{
		btcClient:       btcClient,
		ethClient:       ethClient,
		fileStore:       fileStore,
		initHeight:      cfg.EthInitHeight,
		ethFilter:       cfg.EthAddrFilter,
		txManager:       task,
		chainStore:      NewChainStore(store),
		chainForkSignal: chainFork,
		reScan:          cfg.EthReScan,
	}, nil
}

func (e *ethereumAgent) Init() error {
	logger.Info("init ethereum agent")
	height, exists, err := e.chainStore.ReadEthereumHeight()
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return err
	}
	if !exists && e.initHeight-32 > 0 {
		e.initHeight = e.initHeight - 32
	}
	if !exists || height < e.initHeight || e.reScan {
		logger.Debug("init eth current height: %v", e.initHeight)
		err := e.chainStore.WriteEthereumHeight(e.initHeight)
		if err != nil {
			logger.Error("put eth current height error:%v", err)
			return err
		}
	}
	// test rpc
	_, err = e.ethClient.GetChainId()
	if err != nil {
		logger.Error("ethClient json rpc error:%v", err)
		return err
	}
	err = e.GetCheckpointHeight()
	if err != nil {
		logger.Error("get checkpoint height error:%v", err)
		return err
	}
	return nil
}

func (b *ethereumAgent) GetCheckpointHeight() error {
	hash, err := b.ethClient.SuggestedCP()
	if err != nil {
		logger.Error("ethClient get checkpoint hash error:%v", err)
		return err
	}
	littleHash := hex.EncodeToString(common.ReverseBytes(hash))
	header, err := b.btcClient.GetBlockHeader(littleHash)
	if err != nil {
		logger.Error("btcClient checkpoint height  error:%v %v", err, hash)
		return err
	}
	checkpointHeight := uint64(header.Height)
	err = b.chainStore.WriteCheckpoint(checkpointHeight, hex.EncodeToString(hash))
	if err != nil {
		logger.Error("write checkpoint error:%v", err)
		return err
	}
	err = b.chainStore.WriteLatestCheckpoint(checkpointHeight)
	if err != nil {
		logger.Error("write latest checkpoint error:%v", err)
		return err
	}
	logger.Debug("checkpointHeight: %v, checkpointHash: %v", checkpointHeight, littleHash)
	return nil
}

func (e *ethereumAgent) ScanBlock() error {
	logger.Debug("ethereum scan block ...")
	currentHeight, ok, err := e.chainStore.ReadEthereumHeight()
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return err
	}
	if !ok {
		logger.Warn("no find eth current height")
		return fmt.Errorf("no find eth current height")
	}
	blockNumber, err := e.ethClient.BlockNumber(context.Background())
	if err != nil {
		logger.Error("get eth block number error:%v", err)
		return err
	}
	if currentHeight >= blockNumber {
		logger.Debug("eth currentHeight:%d >= blockNumber:%d", currentHeight, blockNumber)
		return nil
	}
	err = e.GetCheckpointHeight()
	if err != nil {
		logger.Error("get checkpoint height error:%v", err)
		return err
	}

	for index := uint64(currentHeight) + 1; index <= blockNumber; index++ {
		logger.Debug("ethereum parse block:%d", index)
		preHeight := index - 1
		chainForked, err := e.chainFork(preHeight)
		if err != nil {
			logger.Error("check chain fork error: %v %v", preHeight, err)
			return err
		}
		if chainForked {
			logger.Warn("ethereum chain forked: %v", preHeight)
			err = e.rollback(preHeight)
			if err != nil {
				logger.Error("ethereum rollback error: %v %v", preHeight, err)
				return err
			}
			return nil
		}
		err = e.scan(index)
		if err != nil {
			logger.Error("scan error: %v %v", index, err)
			return err
		}
		err = e.chainStore.WriteEthereumHeight(index)
		if err != nil {
			logger.Error("batch write error: %v %v", index, err)
			return err
		}

	}
	return nil
}

func (e *ethereumAgent) ReScan(index uint64) error {
	logger.Debug("reScan eth block:%d", index)
	err := e.scan(index, true)
	if err != nil {
		logger.Error("scan error: %v %v", index, err)
		return err
	}
	return nil
}

func (e *ethereumAgent) scan(index uint64, skipCheck ...bool) error {
	depositTxes, redeemTxes, updateUtxoTxes, depositRewards, redeemRewards, err := e.parseBlock(index)
	if err != nil {
		logger.Error("eth parse block error: %v %v", index, err)
		return err
	}
	if len(skipCheck) > 0 && skipCheck[0] {
		redeemTxes = txSkipCheck(redeemTxes)
	}
	err = e.chainStore.EthSaveData(index, depositTxes, redeemTxes, updateUtxoTxes, depositRewards, redeemRewards)
	if err != nil {
		logger.Error("ethereum save data error: %v %v", index, err)
		return err
	}
	return nil
}

func (e *ethereumAgent) rollback(height uint64) error {
	startForkHeight, err := e.findChainForkHeight(height)
	if err != nil {
		logger.Error("find eth chain startForkHeight error: %v %v", height, err)
		return err
	}
	logger.Debug("find eth chain startForkHeight: %v", startForkHeight)
	for index := height; index >= startForkHeight; index-- {
		logger.Debug("eth rollback data height: %v", index)
		err := e.chainStore.EthDeleteData(index)
		if err != nil {
			logger.Error("eth rollback data error: %v %v", index, err)
			return err
		}
		err = e.chainStore.WriteEthereumHeight(index - 1)
		if err != nil {
			logger.Error("write eth height error: %v %v", index-1, err)
			return err
		}
	}
	chainFork := ChainFork{
		ForkHeight: startForkHeight,
		Chain:      common.EthereumChain,
		Timestamp:  time.Now().UnixNano(),
	}
	e.chainForkSignal <- &chainFork
	return nil
}

func (e *ethereumAgent) findChainForkHeight(height uint64) (uint64, error) {
	for index := height; index >= e.initHeight; index = index - 1 {
		localBlockHash, exists, err := e.chainStore.ReadEthHash(index)
		if err != nil {
			logger.Error("get eth localHash error: %v %v", index, err)
			return 0, err
		}
		if !exists {
			logger.Error("no find eth localHash %v", index)
			return 0, fmt.Errorf("no find eth localHash %v", index)
		}
		chainBlockHash, err := e.ethClient.GetBlock(int64(index))
		if err != nil {
			logger.Error("get eth get chainHash error: %v %v", index, err)
			return 0, err
		}
		if common.StrEqual(localBlockHash, chainBlockHash.Hash().String()) {
			logger.Info("find eth startForkHeight: %v", index)
			return index + 1, nil
		}
	}
	return e.initHeight, nil
}

func (e *ethereumAgent) chainFork(height uint64) (bool, error) {
	if height <= e.initHeight {
		return false, nil
	}
	localHash, ok, err := e.chainStore.ReadEthHash(height)
	if err != nil {
		logger.Error("get eth localHash error: %v %v", height, err)
		return false, err
	}
	if !ok {
		logger.Warn("get eth %v localHash error", height)
		return false, nil
	}
	block, err := e.ethClient.GetBlock(int64(height))
	if err != nil {
		logger.Error("get eth block error: %v %v", height, err)
		return false, err
	}
	if !common.StrEqual(localHash, block.Hash().String()) {
		logger.Error("find eth chainForked height:%v,localHash:%v,chainHash:%v", height, localHash, block.Hash().String())
		return true, nil
	}
	return false, nil

}

func (e *ethereumAgent) ProofResponse(resp *common.ProofResponse) error {
	switch resp.ProofType {
	case common.RedeemTxType:
		logger.Debug("find Redeem proof: %v %v", resp.Hash, hex.EncodeToString(resp.Proof))
		err := e.deleteRedeemTxCache(resp)
		if err != nil {
			logger.Error("delete Redeem cache tx error:%v %v", resp.Hash, err)
		}
		err = e.chainStore.UpdateProof(resp.Hash, hex.EncodeToString(resp.Proof), common.RedeemTxType, common.ProofSuccess)
		if err != nil {
			logger.Error("update Proof error: %v %v", resp.Hash, err)
			return err
		}
		btcHash, err := e.txManager.RedeemZkbtc(resp.Hash, hex.EncodeToString(resp.Proof))
		if err != nil {
			logger.Error("Redeem btc error:%v %v,save to db", resp.Hash, err)
			e.txManager.AddTask(resp)
			return err
		}
		logger.Debug("success Redeem btc ethHash: %v,btcHash: %v", resp.Hash, btcHash)
	default:
	}
	return nil
}

func (e *ethereumAgent) deleteRedeemTxCache(resp *common.ProofResponse) error {
	finalizedSlot, ok, err := e.fileStore.GetTxFinalizedSlot(resp.FIndex)
	if err != nil {
		logger.Error("get latest slot error: %v", err)
		return err
	}
	if !ok {
		logger.Warn("no find latest slot: %v", resp.FIndex)
		return nil
	}
	err = e.chainStore.DeleteRedeemSotCache(resp.FIndex, finalizedSlot, resp.Hash)
	if err != nil {
		logger.Error("delete redeem cache error: %v", err)
		return err
	}
	return nil
}

func (e *ethereumAgent) parseBlock(height uint64) ([]*DbTx, []*DbTx, []*DbTx, []*DbTx, []*DbTx, error) {
	block, err := e.ethClient.GetBlock(int64(height))
	if err != nil {
		logger.Error("ethereum rpc get block error:%v", err)
		return nil, nil, nil, nil, nil, err
	}
	blockTime := block.Time()
	blockHash := block.Hash().String()
	err = e.chainStore.WriteEthHash(height, blockHash)
	if err != nil {
		logger.Error("write eth hash error:%v", err)
		return nil, nil, nil, nil, nil, err
	}
	logFilters, topicFilters := e.ethFilter.FilterLogs()
	logs, err := e.ethClient.GetLogs(blockHash, logFilters, topicFilters)
	if err != nil {
		logger.Error("ethereum rpc get logs error:%v", err)
		return nil, nil, nil, nil, nil, err
	}
	var depositTxes []*DbTx
	var redeemTxes []*DbTx
	var updateUtxoTxes []*DbTx
	var depositRewards []*DbTx
	var redeemRewards []*DbTx
	for _, log := range logs {
		depositTx, isDeposit, err := e.depositTx(log, blockTime)
		if err != nil {
			logger.Error("check is deposit tx error:%v", err)
			return nil, nil, nil, nil, nil, err
		}
		if isDeposit {
			depositTxes = append(depositTxes, depositTx)
			continue
		}
		depositReward, isDepositReward, err := e.depositReward(log, blockTime)
		if err != nil {
			logger.Error("check deposit reward:%v", err)
			return nil, nil, nil, nil, nil, err
		}
		if isDepositReward {
			depositRewards = append(depositRewards, depositReward)
			continue
		}
		updateUtxoTx, isUpdateUtxo, err := e.updateUtxo(log, blockTime)
		if err != nil {
			logger.Error("check is update utxo tx error:%v", err)
			return nil, nil, nil, nil, nil, err
		}
		if isUpdateUtxo {
			updateUtxoTxes = append(updateUtxoTxes, updateUtxoTx)
			continue
		}
		redeemReward, isRedeemReward, err := e.redeemReward(log, blockTime)
		if err != nil {
			logger.Error("check is Redeem reward tx error:%v", err)
			return nil, nil, nil, nil, nil, err
		}
		if isRedeemReward {
			redeemRewards = append(redeemRewards, redeemReward)
			continue
		}
		redeemTx, isRedeem, err := e.redeemTx(log, blockTime)
		if err != nil {
			logger.Error("check is Redeem tx error:%v", err)
			return nil, nil, nil, nil, nil, err
		}
		if isRedeem {
			proved, err := e.checkRedeemTxProved(redeemTx.UtxoId)
			if err != nil {
				logger.Error("check Redeem tx proved error:%v", err)
				return nil, nil, nil, nil, nil, err
			}
			if proved {
				redeemTx.Proved = proved
			}
			redeemTxes = append(redeemTxes, redeemTx)
		}
	}
	return depositTxes, redeemTxes, updateUtxoTxes, depositRewards, redeemRewards, nil
}

func (e *ethereumAgent) checkRedeemTxProved(btxTxId string) (bool, error) {
	_, exists, err := e.chainStore.ReadUpdateUtxoDest(btxTxId)
	if err != nil {
		logger.Error("read update utxo dest error: %v", err)
		return false, err
	}
	if exists {
		return true, nil
	}
	exists, err = e.btcClient.CheckTx(btxTxId)
	if err != nil {
		logger.Error("check btc txId error: %v", btxTxId)
		return false, err
	}
	return exists, nil
}

func (e *ethereumAgent) redeemReward(log types.Log, blockTime uint64) (*DbTx, bool, error) {
	if log.Removed {
		return nil, false, nil
	}
	if len(log.Topics) != 4 {
		return nil, false, nil
	}
	if !e.ethFilter.RedeemReward(log.Address.Hex(), log.Topics[0].Hex()) {
		return nil, false, nil
	}
	minerAddr := fixEthAddr(strings.ToLower(log.Topics[1].Hex()))
	reward := log.Topics[2].Big().Int64()
	txId := fmt.Sprintf("%x", common.ReverseBytes(log.Topics[3][:]))
	logger.Info("ethereum agent find Redeem reward height:%v ethTxHash:%v,miner: %v,reward:%v,txId:%v",
		log.BlockNumber, log.TxHash.String(), minerAddr, reward, txId)
	sender, err := e.ethClient.GetTxSender(log.TxHash.String(), log.BlockHash.String(), log.TxIndex)
	if err != nil {
		logger.Error("get tx sender error:%v", err)
		return nil, false, err
	}
	rewardTx := NewRedeemRewardTx(log.BlockNumber, log.TxIndex, log.Index, log.TxHash.String(), sender, minerAddr, reward, blockTime)
	return rewardTx, true, nil

}

func (e *ethereumAgent) updateUtxo(log types.Log, blockTime uint64) (*DbTx, bool, error) {
	if log.Removed {
		return nil, false, nil
	}
	if len(log.Topics) != 3 {
		return nil, false, nil
	}
	if !e.ethFilter.UpdateUtxo(log.Address.Hex(), log.Topics[0].Hex()) {
		return nil, false, nil
	}
	utxoId := fmt.Sprintf("%x", common.ReverseBytes(log.Topics[1][:]))
	utxoIndex := log.Topics[2].Big().Int64()
	amount := big.NewInt(0).SetBytes(log.Data).Int64()
	sender, err := e.ethClient.GetTxSender(log.TxHash.String(), log.BlockHash.String(), log.TxIndex)
	if err != nil {
		logger.Error("get tx sender error:%v", err)
		return nil, false, err
	}
	logger.Info("ethereum agent find update utxo  ethHash:%v,utxoId:%v,index: %v,amount:%v,height:%v,sender:%v",
		log.TxHash.String(), utxoId, utxoIndex, amount, log.BlockNumber, sender)
	updateUtxoTx := NewUpdateUtxoTx(log.BlockNumber, log.TxIndex, log.Index, log.TxHash.String(), utxoId, utxoIndex, amount, blockTime)
	return updateUtxoTx, true, nil

}

func (e *ethereumAgent) depositReward(log types.Log, blockTime uint64) (*DbTx, bool, error) {
	if log.Removed {
		return nil, false, nil
	}
	if len(log.Topics) != 3 {
		return nil, false, nil
	}
	if !e.ethFilter.DepositReward(log.Address.Hex(), log.Topics[0].Hex()) {
		return nil, false, nil
	}
	minerAddr := fixEthAddr(strings.ToLower(log.Topics[1].Hex()))
	amount := log.Topics[2].Big().Int64()
	sender, err := e.ethClient.GetTxSender(log.TxHash.String(), log.BlockHash.String(), log.TxIndex)
	if err != nil {
		logger.Error("get tx sender error:%v", err)
		return nil, false, err
	}
	logger.Info("ethereum agent find deposit reward height:%v ethTxHash:%v,miner:%v,amount:%v,sender:%v",
		log.BlockNumber, log.TxHash, minerAddr, amount, sender)
	rewardTx := NewDepositRewardTx(log.BlockNumber, log.TxIndex, log.Index, log.TxHash.String(), sender, minerAddr, amount, blockTime)
	return rewardTx, true, nil

}

func (e *ethereumAgent) depositTx(log types.Log, blockTime uint64) (*DbTx, bool, error) {
	if log.Removed {
		return nil, false, nil
	}
	if len(log.Topics) != 3 {
		return nil, false, nil
	}
	if !e.ethFilter.DepositTx(log.Address.Hex(), log.Topics[0].Hex()) {
		return nil, false, nil
	}
	utxoId := fmt.Sprintf("%x", common.ReverseBytes(log.Topics[1][:]))
	utxoIndex := log.Topics[2].Big().Int64()
	amount := big.NewInt(0).SetBytes(log.Data).Int64()
	sender, err := e.ethClient.GetTxSender(log.TxHash.String(), log.BlockHash.String(), log.TxIndex)
	if err != nil {
		logger.Error("get tx sender error:%v", err)
		return nil, false, err
	}
	logger.Info("ethereum agent find deposit zkbtc height:%v ethHash:%v,utxoId:%v,utxoIndex:%v,logIndex:%v,amount:%v,sender:%v",
		log.BlockNumber, log.TxHash.String(), utxoId, utxoIndex, log.Index, amount, sender)
	depositTx := NewDepositEthTx(log.BlockNumber, log.TxIndex, log.Index, log.TxHash.String(), sender, utxoId, utxoIndex, amount, blockTime)
	return depositTx, true, nil

}

func (e *ethereumAgent) redeemTx(log types.Log, blockTime uint64) (*DbTx, bool, error) {
	if log.Removed {
		return nil, false, nil
	}
	if len(log.Topics) != 3 {
		return nil, false, nil
	}
	if !e.ethFilter.RedeemTx(log.Address.Hex(), log.Topics[0].Hex()) {
		return nil, false, nil
	}
	tx, _, err := e.ethClient.TransactionByHash(context.TODO(), ethCommon.HexToHash(log.TxHash.String()))
	if err != nil {
		logger.Error("get eth tx error:%v %v", log.TxHash, err)
		return nil, false, err
	}
	if tx.Type() != 2 {
		return nil, false, nil
	}
	btcTxId := strings.ToLower(hex.EncodeToString(common.ReverseBytes(log.Topics[1].Bytes())))
	minerReward := log.Topics[2].Big()

	_, btcTxRaw, err := redeemUtils.DecodeRawTxLogData(log.Data)
	if err != nil {
		logger.Error("decode Redeem log error:%v", err)
		return nil, false, err
	}
	btcTx := btctx.NewTransaction()
	err = btcTx.Deserialize(bytes.NewReader(btcTxRaw))
	if err != nil {
		logger.Error("deserialize btc tx error:%v", err)
		return nil, false, err
	}
	txHash := log.TxHash.String()
	if !e.ethFilter.MigrateTx(btcTx.TxOut) {
		logger.Warn("check redeem status error,current  ethHash:%v,btcHash:%v", log.TxHash.String(), btcTx.TxHash().String())
		return nil, false, nil
	}
	if strings.TrimPrefix(btcTx.TxHash().String(), "0x") != strings.TrimPrefix(btcTxId, "0x") {
		logger.Error("never should happen btc tx not match error: Hash:%v, logBtcTxId:%v,decodeTxHash:%v", txHash, btcTxId, btcTx.TxHash().String())
		return nil, false, fmt.Errorf("tx hash not match:%v", txHash)
	}
	amount := getEthRedeemAmount(btcTx.TxOut, e.ethFilter.BtcLockScript)

	txSender, err := e.getTxSender(txHash, log.BlockHash.Hex(), log.TxIndex)
	if err != nil {
		logger.Error("get tx sender error:%v", err)
		return nil, false, err
	}
	blockNumber := log.BlockNumber
	logger.Info("ethereum agent find Redeem zkbtc height:%v, index: %v,ethTxHash:%v,sender:%v,btcTxId:%v,minerReward:%v,"+
		"amount:%v,input:%v,output:%v", blockNumber, log.TxIndex, txHash, txSender, btcTxId, minerReward.String(), amount,
		getInputString(btcTx.TxIn), getOutputString(btcTx.TxOut))
	redeemTx := NewRedeemEthTx(blockNumber, log.TxIndex, log.Index, txHash, txSender, btcTxId, amount, blockTime)
	return redeemTx, true, nil

}

func (e *ethereumAgent) getTxSender(txHash, blockHash string, index uint) (string, error) {
	tx, pending, err := e.ethClient.TransactionByHash(context.Background(), ethCommon.HexToHash(txHash))
	if err != nil {
		logger.Error("get eth tx error:%v %v", txHash, err)
		return "", err
	}
	if pending {
		return "", fmt.Errorf("tx %v is pending", txHash)
	}
	sender, err := e.ethClient.TransactionSender(context.Background(), tx, ethCommon.HexToHash(blockHash), index)
	if err != nil {
		logger.Error("get eth tx sender error:%v %v", txHash, err)
		return "", err
	}
	return sender.Hex(), nil
}

func (e *ethereumAgent) CheckState() error {
	// ethereum sync blocks per half hour
	height, ok, err := e.chainStore.ReadEthereumHeight()
	if err != nil {
		logger.Error("read beacon latest height error: %v", err)
		return err
	}
	if ok {
		diff := height - e.curHeight
		if diff < 100 { //normal 150
			logger.Error("ethereum sync too slow , node maybe offline diff:%v prevHeight:%v curHeight:%v", diff, e.curHeight, height)
		}

	}
	e.curHeight = height
	return nil

}
func (e *ethereumAgent) Close() error {
	return nil
}
func (e *ethereumAgent) Name() string {
	return EthereumAgentName
}

func NewDepositEthTx(height uint64, txIndex, logIndex uint, txHash, sender, utxoId string, utxoIndex, amount int64, blockTime uint64) *DbTx {
	return &DbTx{
		Hash:      DbValue(txHash),
		TxIndex:   txIndex,
		LogIndex:  logIndex,
		Height:    height,
		ChainType: common.EthereumChain,
		TxType:    common.DepositTx,
		UtxoId:    DbValue(utxoId),
		UtxoIndex: utxoIndex,
		Amount:    amount,
		Sender:    sender,
		BlockTime: blockTime,
	}
}

func NewRedeemEthTx(height uint64, txIndex, logIndex uint, txHash, sender, btcTxId string, amount int64, blockTime uint64) *DbTx {
	return &DbTx{
		Height:    height,
		TxIndex:   txIndex,
		LogIndex:  logIndex,
		ProofType: common.RedeemTxType,
		Hash:      DbValue(txHash),
		ChainType: common.EthereumChain,
		TxType:    common.RedeemTx,
		UtxoId:    DbValue(btcTxId),
		Sender:    DbValue(sender),
		Amount:    amount,
		BlockTime: blockTime,
	}
}
func NewUpdateUtxoTx(height uint64, txIndex, logIndex uint, txHash, utxoId string, utxoIndex, amount int64, blockTime uint64) *DbTx {
	return &DbTx{
		Height:    height,
		TxIndex:   txIndex,
		LogIndex:  logIndex,
		Hash:      DbValue(txHash),
		ChainType: common.EthereumChain,
		TxType:    common.UpdateUtxoTx,
		UtxoId:    utxoId,
		UtxoIndex: utxoIndex,
		Amount:    amount,
		BlockTime: blockTime,
	}
}

func NewDepositRewardTx(height uint64, txIndex, logIndex uint, txHash, sender, minerAddr string, amount int64, blockTime uint64) *DbTx {
	return &DbTx{
		Height:    height,
		TxIndex:   txIndex,
		LogIndex:  logIndex,
		Hash:      DbValue(txHash),
		ChainType: common.EthereumChain,
		TxType:    common.DepositRewardTx,
		Sender:    DbValue(sender),
		Receiver:  DbValue(minerAddr),
		Amount:    amount,
		BlockTime: blockTime,
	}
}

func NewRedeemRewardTx(height uint64, txIndex, logIndex uint, txHash, sender, minerAddr string, amount int64, blockTime uint64) *DbTx {
	return &DbTx{
		Height:    height,
		TxIndex:   txIndex,
		LogIndex:  logIndex,
		Hash:      DbValue(txHash),
		ChainType: common.EthereumChain,
		TxType:    common.RedeemRewardTx,
		Sender:    DbValue(sender),
		Receiver:  DbValue(minerAddr),
		Amount:    amount,
		BlockTime: blockTime,
	}
}

func fixEthAddr(ethAddr string) string {
	if len(ethAddr) < 40 {
		return ""
	}
	return ethAddr[len(ethAddr)-40:]
}

func getEthRedeemAmount(outputs []*wire.TxOut, lockScript string) int64 {
	var amount int64
	for _, out := range outputs {
		if !common.StrEqual(hex.EncodeToString(out.PkScript), lockScript) {
			amount = amount + out.Value
		}
	}
	return amount
}

func getInputString(inputs []*wire.TxIn) string {
	var inputStr string
	for _, in := range inputs {
		inputStr = inputStr + in.PreviousOutPoint.String()
	}
	return inputStr
}
func getOutputString(outputs []*wire.TxOut) string {
	var outputStr string
	for _, out := range outputs {
		outputStr = outputStr + hex.EncodeToString(out.PkScript) + ":" + fmt.Sprintf("%v", out.Value)
	}
	return outputStr
}
