package node

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	btctx "github.com/lightec-xyz/daemon/rpc/bitcoin/common"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/rpc/oasis"
	"github.com/lightec-xyz/daemon/store"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"strconv"
	"strings"
	"sync"
)

type EthereumAgent struct {
	btcClient        *bitcoin.Client
	ethClient        *ethrpc.Client
	oasisClient      *oasis.Client
	apiClient        *apiclient.Client // todo temporary use
	beaconClient     *beacon.Client
	store            store.IStore
	memoryStore      store.IStore
	stateCache       *Cache
	fileStore        *FileStorage
	taskManager      *TxManager
	proofRequest     chan []*common.ZkProofRequest
	multiAddressInfo MultiAddressInfo
	logAddrFilter    EthAddrFilter
	btcLockScript    string
	initHeight       int64
	txManager        *TxManager
	genesisPeriod    uint64
	lock             sync.Mutex
	genesisSlot      uint64
	debug            bool
	force            bool
	prepareData      *PreparedData
}

func NewEthereumAgent(cfg Config, genesisSlot uint64, fileStore *FileStorage, store, memoryStore store.IStore, beaClient *apiclient.Client,
	btcClient *bitcoin.Client, ethClient *ethrpc.Client, beaconClient *beacon.Client, oasisClient *oasis.Client,
	proofRequest chan []*common.ZkProofRequest, task *TxManager, state *Cache) (IAgent, error) {
	return &EthereumAgent{
		apiClient:        beaClient, // todo
		btcClient:        btcClient,
		ethClient:        ethClient,
		beaconClient:     beaconClient,
		oasisClient:      oasisClient,
		store:            store,
		fileStore:        fileStore,
		memoryStore:      memoryStore,
		proofRequest:     proofRequest,
		multiAddressInfo: cfg.MultiAddressInfo,
		initHeight:       cfg.EthInitHeight,
		logAddrFilter:    cfg.EthAddrFilter,
		btcLockScript:    cfg.BtcLockScript,
		txManager:        task,
		genesisPeriod:    genesisSlot / 8192,
		genesisSlot:      genesisSlot,
		stateCache:       state,
		debug:            common.GetEnvDebugMode(),
	}, nil
}

func (e *EthereumAgent) Init() error {
	logger.Info("init ethereum agent")
	if e.force {
		err := WriteEthereumHeight(e.store, e.initHeight)
		if err != nil {
			logger.Error("write eth height error: %v %v", e.initHeight, err)
			return err
		}
	} else {
		height, exists, err := ReadEthereumHeight(e.store)
		if err != nil {
			logger.Error("get eth current height error:%v", err)
			return err
		}
		if !exists || height < e.initHeight {
			logger.Debug("init eth current height: %v", e.initHeight)
			err := WriteEthereumHeight(e.store, e.initHeight)
			if err != nil {
				logger.Error("put eth current height error:%v", err)
				return err
			}
		}
	}
	// test rpc
	_, err := e.ethClient.GetChainId()
	if err != nil {
		logger.Error("ethClient json rpc error:%v", err)
		return err
	}
	// todo just test
	if e.debug {
		err = WriteEthereumHeight(e.store, 1584946)
		if err != nil {
			logger.Error("%v", err)
			return err
		}
	}
	return nil
}

func (e *EthereumAgent) ScanBlock() error {
	logger.Debug("ethereum scan block ...")
	ethHeight, ok, err := ReadEthereumHeight(e.store)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return err
	}
	if !ok {
		return fmt.Errorf("never should happen")
	}
	blockNumber, err := e.ethClient.BlockNumber(context.Background())
	if err != nil {
		logger.Error("get eth block number error:%v", err)
		return err
	}
	forked, err := e.CheckChainFork(blockNumber)
	if err != nil {
		logger.Error("ethereum chain fork error:%v %v", blockNumber, err)
		return err
	}
	if forked {
		logger.Error("ethereum chain forked,need to rollback %v %v", blockNumber, err)
		// todo
		//return nil
	}
	// todo
	if e.debug {
		if ethHeight > 1585877 {
			return nil
		}
	}
	// todo
	blockNumber = blockNumber - 1
	if ethHeight >= int64(blockNumber) {
		logger.Debug("eth current height:%d,latest block number :%d", ethHeight, blockNumber)
		return nil
	}
	for index := ethHeight + 1; index <= int64(blockNumber); index++ {
		logger.Debug("ethereum parse block:%d", index)
		redeemTxes, depositTxes, err := e.parseBlock(index)
		if err != nil {
			logger.Error("eth parse block error: %v %v", index, err)
			return err
		}
		err = e.updateDepositData(index, depositTxes)
		if err != nil {
			logger.Error("ethereum update deposit info error: %v %v", index, err)
			return err
		}

		allTxes := append(redeemTxes, depositTxes...)
		err = e.saveTransaction(index, allTxes)
		if err != nil {
			logger.Error("ethereum save transaction error: %v %v", index, err)
			return err
		}
		err = e.saveData(redeemTxes)
		if err != nil {
			logger.Error("ethereum save Data error: %v %v", index, err)
			return err
		}
		err = WriteEthereumHeight(e.store, index)
		if err != nil {
			logger.Error("batch write error: %v %v", index, err)
			return err
		}

	}
	return nil
}

func (e *EthereumAgent) CheckChainFork(height uint64) (bool, error) {
	// todo
	return false, nil
}

func (e *EthereumAgent) ProofResponse(resp *common.ZkProofResponse) error {
	logger.Info(" ethereumAgent receive proof response: %v proof,%x", resp.Id(), resp.Proof)
	e.stateCache.Delete(resp.Id())
	switch resp.ZkProofType {
	case common.RedeemTxType:
		err := e.DeleteRedeemCacheTx(resp)
		if err != nil {
			logger.Error("delete redeem cache tx error:%v %v", resp.TxHash, err)
		}
		err = e.updateRedeemProof(resp.TxHash, hex.EncodeToString(resp.Proof), resp.Status)
		if err != nil {
			logger.Error("update Proof error:%v %v", resp.TxHash, err)
		}
		btcHash, err := e.txManager.RedeemZkbtc(resp.TxHash, hex.EncodeToString(resp.Proof))
		if err != nil {
			logger.Error("redeem btc error:%v %v,save to db", resp.TxHash, err)
			e.txManager.AddTask(resp)
			return err
		}
		logger.Debug("success redeem btc ethHash: %v,btcHash: %v", resp.TxHash, btcHash)
	default:
	}
	return nil
}

func (e *EthereumAgent) DeleteRedeemCacheTx(resp *common.ZkProofResponse) error {
	err := DeleteTxSlot(e.store, resp.Index, resp.TxHash)
	if err != nil {
		logger.Error("delete tx slot error:%v %v", resp.TxHash, err)
		return err
	}
	logger.Debug("delete %v beaconHeader  cache", resp.TxHash)
	finalizedSlot, ok, err := e.fileStore.GetNearTxSlotFinalizedSlot(resp.Index)
	if err != nil {
		logger.Error("get latest slot error: %v", err)
		return err
	}
	if !ok {
		logger.Warn("no find latest slot: %v", resp.Index)
		return nil
	}
	err = DeleteTxFinalizedSlot(e.store, finalizedSlot, resp.TxHash)
	if err != nil {
		logger.Error("delete tx slot error:%v %v", resp.TxHash, err)
		return err
	}
	logger.Debug("delete %v beaconHeaderFinality  cache", resp.TxHash)
	return nil
}

func (e *EthereumAgent) updateDepositData(height int64, depositTxes []*Transaction) error {
	for _, tx := range depositTxes {
		// btcTxId -> ethTxHash
		err := WriteDestHash(e.store, tx.BtcTxId, tx.TxHash)
		if err != nil {
			logger.Error("update deposit final status error: %v %v", height, err)
			return err
		}
		// ethTxHash -> btcTxId
		err = WriteDestHash(e.store, tx.TxHash, tx.BtcTxId)
		if err != nil {
			logger.Error("update deposit final status error: %v %v", height, err)
			return err
		}
		// delete bitcoin ungen proof
		err = DeleteUnGenProof(e.store, Bitcoin, tx.BtcTxId)
		if err != nil {
			logger.Error("delete ungen proof error: %v %v", height, err)
			return err
		}
	}
	return nil
}

func (e *EthereumAgent) saveTransaction(height int64, txes []*Transaction) error {
	err := WriteEthereumTxIds(e.store, height, txesToTxIds(txes))
	if err != nil {
		logger.Error("write ethereum tx ids error: %v %v", height, err)
		return err
	}
	err = WriteTxes(e.store, txesToDbTxes(txes))
	if err != nil {
		logger.Error("put redeem tx error: %v %v", height, err)
		return err
	}
	return nil
}

func (e *EthereumAgent) saveData(redeemTxes []*Transaction) error {
	addrTxesMap := txesByAddrGroup(redeemTxes)
	for addr, addrTxes := range addrTxesMap {
		err := WriteTxesByAddr(e.store, addr, addrTxes)
		if err != nil {
			logger.Error("write addr txes error: %v %v", addr, err)
			return err
		}
	}
	err := WriteDbProof(e.store, txesToDbProofs(redeemTxes))
	if err != nil {
		logger.Error("put eth current height error:%v", err)
		return err
	}
	for _, tx := range redeemTxes {
		logger.Debug("save redeem tx: %v", tx.TxHash)
		err := WriteDestHash(e.store, tx.BtcTxId, tx.TxHash)
		if err != nil {
			logger.Error("batch write error: %v", err)
			return err
		}
		err = WriteDestHash(e.store, tx.TxHash, tx.BtcTxId)
		if err != nil {
			logger.Error("batch write error: %v", err)
			return err
		}
	}
	// cache need to generate redeem proof
	err = WriteUnGenProof(e.store, Ethereum, txesToUnGenProofs(Ethereum, redeemTxes))
	if err != nil {
		logger.Error("write ungen Proof error: %v", err)
		return err
	}
	return nil
}

func (e *EthereumAgent) parseBlock(height int64) ([]*Transaction, []*Transaction, error) {
	block, err := e.ethClient.GetBlock(height)
	if err != nil {
		logger.Error("ethereum rpc get block error:%v", err)
		return nil, nil, err
	}
	blockHash := block.Hash().String()
	logFilters, topicFilters := e.logAddrFilter.FilterLogs()
	logs, err := e.ethClient.GetLogs(blockHash, logFilters, topicFilters)
	if err != nil {
		logger.Error("ethereum rpc get logs error:%v", err)
		return nil, nil, err
	}
	var redeemTxes []*Transaction
	var depositTxes []*Transaction
	for _, log := range logs {
		depositTx, isDeposit, err := e.isDepositTx(log)
		if err != nil {
			logger.Error("check is deposit tx error:%v", err)
			return nil, nil, err
		}
		if isDeposit {
			depositTxes = append(depositTxes, depositTx)
			continue
		}
		redeemTx, isRedeem, err := e.isRedeemTx(log)
		if err != nil {
			logger.Error("check is redeem tx error:%v", err)
			return nil, nil, err
		}
		if isRedeem {
			proofed, err := e.CheckRedeemTx(redeemTx.BtcTxId)
			if err != nil {
				logger.Error("check redeem tx proofed error:%v", err)
				return nil, nil, err
			}
			if proofed {
				redeemTx.Proofed = proofed
			}
			redeemTxes = append(redeemTxes, redeemTx)
		}
	}
	return redeemTxes, depositTxes, nil
}

func (e *EthereumAgent) CheckRedeemTx(btxTxId string) (bool, error) {
	exists, err := e.btcClient.CheckTx(btxTxId)
	if err != nil {
		logger.Error("check btc txId error: %v", btxTxId)
		return false, err
	}
	return exists, nil
}

func (e *EthereumAgent) isDepositTx(log types.Log) (*Transaction, bool, error) {
	if log.Removed {
		return nil, false, nil
	}
	if len(log.Topics) != 3 {
		return nil, false, nil
	}
	// todo
	if strings.ToLower(log.Address.Hex()) == e.logAddrFilter.LogDepositAddr && strings.ToLower(log.Topics[0].Hex()) == e.logAddrFilter.LogTopicDepositAddr {
		btcTxId := strings.ToLower(log.Topics[1].Hex())
		hexVout := strings.TrimPrefix(strings.ToLower(log.Topics[2].Hex()), "0x")
		vout, err := strconv.ParseInt(hexVout, 16, 32)
		if err != nil {
			logger.Error("parse vout error:%v", err)
			return nil, false, err
		}
		amount, err := strconv.ParseInt(fmt.Sprintf("%x", log.Data), 16, 64)
		if err != nil {
			logger.Error("parse amount error:%v", err)
			return nil, false, err
		}
		utxo := []Utxo{
			{
				TxId:  btcTxId,
				Index: uint32(vout),
			},
		}
		txHash := log.TxHash.String()
		blockNumber := log.BlockNumber
		logger.Info("ethereum agent find deposit zkbtc height:%v ethTxHash:%v,btcTxId:%v,utxo:%v",
			blockNumber, txHash, btcTxId, formatUtxo(utxo))
		depositTx := NewDepositEthTx(blockNumber, log.TxIndex, txHash, btcTxId, utxo, amount)
		return depositTx, true, nil
	} else {
		return nil, false, nil
	}
}

func (e *EthereumAgent) getTxSender(txHash, blockHash string, index uint) (string, error) {
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

func (e *EthereumAgent) isRedeemTx(log types.Log) (*Transaction, bool, error) {

	if log.Removed {
		return nil, false, nil
	}
	if len(log.Topics) != 2 {
		return nil, false, nil
	}
	//todo more check
	if strings.ToLower(log.Address.Hex()) == e.logAddrFilter.LogRedeemAddr && strings.ToLower(log.Topics[0].Hex()) == e.logAddrFilter.LogTopicRedeemAddr {
		btcTxId := strings.ToLower(log.Topics[1].Hex())
		txData, _, err := ethereum.DecodeRedeemLog(log.Data)
		if err != nil {
			logger.Error("decode redeem log error:%v", err)
			return nil, false, err
		}
		transaction := btctx.NewTransaction()
		err = transaction.Deserialize(bytes.NewReader(txData))
		if err != nil {
			logger.Error("deserialize btc tx error:%v", err)
			return nil, false, err
		}
		var inputs []Utxo
		for _, in := range transaction.TxIn {
			inputs = append(inputs, Utxo{
				TxId:  in.PreviousOutPoint.Hash.String(),
				Index: in.PreviousOutPoint.Index,
			})
		}
		var amount int64
		var outputs []TxOut
		for _, out := range transaction.TxOut {
			if hex.EncodeToString(out.PkScript) != e.btcLockScript {
				amount = amount + out.Value
			}
			outputs = append(outputs, TxOut{
				Value:    out.Value,
				PkScript: out.PkScript,
			})
		}
		txHash := log.TxHash.String()
		if strings.TrimPrefix(transaction.TxHash().String(), "0x") != strings.TrimPrefix(btcTxId, "0x") {
			logger.Error("never should happen btc tx not match error: Hash:%v, logBtcTxId:%v,decodeTxHash:%v", txHash, btcTxId, transaction.TxHash().String())
			return nil, false, fmt.Errorf("tx hash not match:%v", txHash)
		}

		txSender, err := e.getTxSender(txHash, log.BlockHash.Hex(), log.TxIndex)
		if err != nil {
			logger.Error("get tx sender error:%v", err)
			return nil, false, err
		}
		blockNumber := log.BlockNumber
		logger.Info("ethereum agent find redeem zkbtc height:%v, index: %v,ethTxHash:%v,sender:%v,btcTxId:%v,amount:%v,input:%v,output:%v",
			blockNumber, log.TxIndex, txHash, txSender, btcTxId, amount, formatUtxo(inputs), formatOut(outputs))
		redeemTx := NewRedeemEthTx(blockNumber, log.TxIndex, txHash, txSender, btcTxId, amount, inputs, outputs)
		return redeemTx, true, nil
	} else {
		return nil, false, nil
	}

}

func (e *EthereumAgent) FetchDataResponse(resp *FetchResponse) error {
	logger.Debug("ethereum fetch response fetchType: %v", resp.Id())
	return nil
}
func (e *EthereumAgent) CheckState() error {
	return nil

}

func (e *EthereumAgent) updateRedeemProofStatus(txHash string, index uint64, status common.ProofStatus) error {
	id := common.NewProofId(common.RedeemTxType, index, 0, txHash)
	if !e.stateCache.Check(id) {
		err := UpdateProof(e.store, txHash, "", common.RedeemTxType, status)
		if err != nil {
			logger.Error("update proof status error: %v %v", txHash, err)
			return err
		}
		return err
	}
	return nil
}

func (e *EthereumAgent) updateRedeemProof(txId string, proof string, status common.ProofStatus) error {
	logger.Debug("update Redeem Proof status: %v %v %v", txId, proof, status)
	err := UpdateProof(e.store, txId, proof, common.RedeemTxType, status)
	if err != nil {
		logger.Error("update Proof error: %v %v", txId, err)
		return err
	}
	return nil
}

func (e *EthereumAgent) Close() error {
	return nil
}
func (e *EthereumAgent) Name() string {
	return EthereumAgentName
}

func NewDepositEthTx(height uint64, txIndex uint, txHash, btcTxId string, utxo []Utxo, amount int64) *Transaction {
	return &Transaction{
		TxHash:    txHash,
		TxIndex:   txIndex,
		Height:    height,
		ChainType: Ethereum,
		TxType:    DepositTx,
		BtcTxId:   btcTxId,
	}
}
func NewRedeemEthTx(height uint64, txIndex uint, txHash, sender, btcTxId string, amount int64, inputs []Utxo, outputs []TxOut) *Transaction {
	return &Transaction{
		Height:    height,
		TxIndex:   txIndex,
		TxHash:    txHash,
		ChainType: Ethereum,
		TxType:    RedeemTx,
		BtcTxId:   btcTxId,
		From:      sender,
		Amount:    amount,
	}
}
