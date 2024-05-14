package node

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/rpc/oasis"
	"github.com/lightec-xyz/daemon/store"
	btctx "github.com/lightec-xyz/daemon/transaction/bitcoin"
	"github.com/lightec-xyz/daemon/transaction/ethereum"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"strconv"
	"strings"
)

type EthereumAgent struct {
	btcClient        *bitcoin.Client
	ethClient        *ethrpc.Client
	oasisClient      *oasis.Client
	apiClient        *apiclient.Client // todo temporary use
	beaconClient     *beacon.Client
	store            store.IStore
	memoryStore      store.IStore
	stateCache       *CacheState
	fileStore        *FileStorage
	taskManager      *TaskManager
	proofRequest     chan []*common.ZkProofRequest
	multiAddressInfo MultiAddressInfo
	logAddrFilter    EthAddrFilter
	initHeight       int64
	task             *TaskManager
	genesisPeriod    uint64
	genesisSlot      uint64
}

func NewEthereumAgent(cfg Config, genesisSlot uint64, fileStore *FileStorage, store, memoryStore store.IStore, beaClient *apiclient.Client,
	btcClient *bitcoin.Client, ethClient *ethrpc.Client, beaconClient *beacon.Client, oasisClient *oasis.Client, proofRequest chan []*common.ZkProofRequest, task *TaskManager) (IAgent, error) {
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
		task:             task,
		genesisPeriod:    genesisSlot / 8192,
		genesisSlot:      genesisSlot,
		stateCache:       NewCacheState(),
	}, nil
}

func (e *EthereumAgent) Init() error {
	logger.Info("init ethereum agent")
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
	// test rpc
	_, err = e.ethClient.GetChainId()
	if err != nil {
		logger.Error("ethClient json rpc error:%v", err)
		return err
	}
	// todo just test
	//err = WriteUnGenProof(e.store, Ethereum, []string{"622af9392653f10797297e2fa72c6236db55d28234fad5a12c098349a8c5bd3f"})
	//if err != nil {
	//	logger.Error("write ungen proof error: %v", err)
	//	return err
	//}
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
	blockNumber = blockNumber - 3
	//todo
	if ethHeight >= int64(blockNumber) {
		logger.Debug("eth current height:%d,latest block number :%d", ethHeight, blockNumber)
		return nil
	}
	for index := ethHeight + 1; index <= int64(blockNumber); index++ {
		if index%15 == 0 {
			logger.Debug("ethereum parse block:%d", index)
		}
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

func (e *EthereumAgent) ProofResponse(resp *common.ZkProofResponse) error {
	logger.Info(" ethereumAgent receive proof response: %v %v %v %x", resp.ZkProofType.String(),
		resp.Period, resp.TxHash, resp.Proof)
	err := StoreZkProof(e.fileStore, resp.ZkProofType, resp.Period, resp.TxHash, resp.Proof, resp.Witness)
	if err != nil {
		logger.Error("store zk proof error:%v", err)
		return err
	}
	proofId := common.NewProofId(resp.ZkProofType, resp.Period, resp.TxHash)
	e.stateCache.Delete(proofId)
	if resp.ZkProofType == common.RedeemTxType {
		err = e.updateRedeemProof(resp.TxHash, hex.EncodeToString(resp.Proof), resp.Status)
		if err != nil {
			logger.Error("update Proof error:%v %v", resp.TxHash, err)
			return err
		}
		_, err = RedeemBtcTx(e.btcClient, e.ethClient, e.oasisClient, resp.TxHash, resp.Proof)
		if err != nil {
			logger.Error("redeem btc tx error:%v %v", resp.TxHash, err)
			e.task.AddTask(resp)
			return err
		}
	}
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
		err := WriteAddrTxs(e.store, addr, addrTxes)
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
	// todo
	return false, nil
	//submitted, err := e.btcClient.CheckTx(btxTxId)
	//if err != nil {
	//	logger.Error("check btc tx error:%v", err)
	//	return false, err
	//}
	//return submitted, nil
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
		var outputs []TxOut
		for _, out := range transaction.TxOut {
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
		logger.Info("ethereum agent find redeem zkbtc height:%v,ethTxHash:%v,sender:%v,btcTxId:%v,input:%v,output:%v",
			blockNumber, txHash, txSender, btcTxId, formatUtxo(inputs), formatOut(outputs))
		redeemTx := NewRedeemEthTx(blockNumber, log.TxIndex, txHash, txSender, btcTxId, inputs, outputs)
		return redeemTx, true, nil
	} else {
		return nil, false, nil
	}

}

func (e *EthereumAgent) CheckState() error {
	logger.Debug("ethereum check state ...")
	unGenProofs, err := ReadAllUnGenProofs(e.store, Ethereum)
	if err != nil {
		logger.Error("read all ungen proof ids error: %v", err)
		return err
	}
	for _, unGenProof := range unGenProofs {
		txHash := unGenProof.TxHash
		logger.Debug("start check redeem proof tx: %v", txHash)
		exists, err := CheckProof(e.fileStore, common.RedeemTxType, 0, txHash)
		if err != nil {
			logger.Error("check tx proof error: %v", err)
			return err
		}
		if exists {
			logger.Debug("redeem proof exist now,delete cache: %v", txHash)
			err := DeleteUnGenProof(e.store, Ethereum, txHash)
			if err != nil {
				logger.Error("delete ungen proof error: %v", err)
				return err
			}
			logger.Debug("delete ungen proof tx: %v", txHash)
			continue
		}
		txSlot, err := e.GetSlotByHash(txHash)
		if err != nil {
			logger.Error("get txSlot error: %v", err)
			return err
		}
		finalizedSlot, ok, err := e.fileStore.GetNearTxSlotFinalizedSlot(txSlot)
		if err != nil {
			logger.Error("get near tx slot finalized slot error: %v", err)
			return err
		}
		if !ok {
			logger.Warn("no find near %v tx slot finalized slot", txSlot)
			continue
		}
		err = e.updateRedeemProofStatus(txHash, txSlot, common.ProofFinalized)
		if err != nil {
			logger.Error("update proof status error: %v %v", txHash, err)
			return err
		}
		exists, err = CheckProof(e.fileStore, common.TxInEth2, 0, txHash)
		if err != nil {
			logger.Error("check tx proof error: %v", err)
			return err
		}
		if !exists {
			err := e.tryProofRequest(common.TxInEth2, 0, txHash)
			if err != nil {
				logger.Error("try proof request error: %v", err)
				return err
			}
		}
		exists, err = CheckProof(e.fileStore, common.BeaconHeaderType, txSlot, "")
		if err != nil {
			logger.Error("check block header proof error: %v", err)
			return err
		}
		if !exists {
			err := e.tryProofRequest(common.BeaconHeaderType, txSlot, txHash)
			if err != nil {
				logger.Error("try proof request error: %v", err)
				return err
			}
		}
		logger.Debug("%v find near %v tx slot finalized slot %v", txHash, txSlot, finalizedSlot)
		exists, err = CheckProof(e.fileStore, common.BeaconHeaderFinalityType, finalizedSlot, "")
		if err != nil {
			logger.Error("check block header finality proof error: %v %v", finalizedSlot, err)
			return err
		}
		if !exists {
			err := e.tryProofRequest(common.BeaconHeaderFinalityType, finalizedSlot, "")
			if err != nil {
				logger.Error("try proof request error: %v", err)
				return err
			}
			continue
		}
		err = e.tryProofRequest(common.RedeemTxType, txSlot, txHash)
		if err != nil {
			logger.Error("try proof request error: %v", err)
			return err
		}
	}
	return nil

}

func (e *EthereumAgent) updateRedeemProofStatus(txHash string, index uint64, status common.ProofStatus) error {
	id := common.NewProofId(common.RedeemTxType, index, txHash)
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

func (e *EthereumAgent) tryProofRequest(zkType common.ZkProofType, index uint64, txHash string) error {
	proofId := common.NewProofId(zkType, index, txHash)
	exists := e.stateCache.Check(proofId)
	if exists {
		logger.Debug("proof request exists: %v", proofId)
		return nil
	}
	exists, err := CheckProof(e.fileStore, zkType, index, txHash)
	if err != nil {
		logger.Error("check proof error: %v", err)
		return err
	}
	if exists {
		logger.Error("proof exists: %v", proofId)
		return nil
	}
	match, err := e.checkRequest(zkType, index, txHash)
	if err != nil {
		logger.Error("check request error: %v", err)
		return err
	}
	if !match {
		logger.Debug("proof request not match: %v", proofId)
		return nil
	}

	data, ok, err := e.getRequestProofData(zkType, index, txHash)
	if err != nil {
		logger.Error("get request proof data error: %v", err)
		return err
	}
	if !ok {
		logger.Debug("proof request data not prepared: %v", proofId)
		return nil
	}
	proofRequest := common.NewZkProofRequest(zkType, data, index, txHash)
	err = e.sendZkProofRequest(proofRequest)
	if err != nil {
		logger.Error("send zk proof request error: %v", err)
		return err
	}
	logger.Info("success send zk proof request: %v", proofId)
	return nil
}

func (e *EthereumAgent) getTxInEth2Data(txHash string) (*rpc.TxInEth2ProveRequest, bool, error) {
	txData, err := ethblock.GenerateTxInEth2Proof(e.ethClient.Client, e.apiClient, txHash)
	if err != nil {
		logger.Error("get tx data error: %v", err)
		return nil, false, err
	}
	return &rpc.TxInEth2ProveRequest{
		TxHash: txHash,
		TxData: txData,
	}, true, nil
}

func (e *EthereumAgent) getBlockHeaderRequestData(index uint64) (*rpc.BlockHeaderRequest, bool, error) {
	// todo
	finalizedSlot, ok, err := e.fileStore.GetNearTxSlotFinalizedSlot(index)
	if err != nil {
		logger.Error("get finalized slot error: %v", err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}

	logger.Debug("get beaconHeader %v ~ %v", index, finalizedSlot)
	beaconBlockHeaders, err := e.beaconClient.RetrieveBeaconHeaders(index, finalizedSlot)
	if err != nil {
		logger.Error("get beacon block headers error: %v", err)
		return nil, false, err
	}
	if len(beaconBlockHeaders) == 0 {
		return nil, false, fmt.Errorf("never should happen %v", index)
	}
	beginSlot, beginRoot, err := BeaconBlockHeaderToSlotAndRoot(beaconBlockHeaders[0])
	if err != nil {
		logger.Error("get beacon block headers error: %v", err)
		return nil, false, err
	}
	endSlot, endRoot, err := BeaconBlockHeaderToSlotAndRoot(beaconBlockHeaders[len(beaconBlockHeaders)-1])
	if err != nil {
		logger.Error("get beacon block headers error: %v", err)
		return nil, false, err
	}
	return &rpc.BlockHeaderRequest{
		Index:     index,
		BeginSlot: beginSlot,
		EndSlot:   endSlot,
		BeginRoot: hex.EncodeToString(beginRoot),
		EndRoot:   hex.EncodeToString(endRoot),
		Headers:   beaconBlockHeaders[1:],
	}, true, nil
}

func (e *EthereumAgent) GetBeaconHeaderId(start, end uint64) ([]byte, []byte, error) {
	beaconBlockHeaders, err := e.beaconClient.RetrieveBeaconHeaders(start, end)
	if err != nil {
		logger.Error("get beacon block headers error: %v", err)
		return nil, nil, err
	}
	if len(beaconBlockHeaders) == 0 {
		return nil, nil, fmt.Errorf("never should happen %v", start)
	}
	_, beginRoot, err := BeaconBlockHeaderToSlotAndRoot(beaconBlockHeaders[0])
	if err != nil {
		logger.Error("get beacon block headers error: %v", err)
		return nil, nil, err
	}
	_, endRoot, err := BeaconBlockHeaderToSlotAndRoot(beaconBlockHeaders[len(beaconBlockHeaders)-1])
	if err != nil {
		logger.Error("get beacon block headers error: %v", err)
		return nil, nil, err
	}
	return beginRoot, endRoot, nil
}

// todo

func (e *EthereumAgent) GetSlotByHash(hash string) (uint64, error) {
	txHash := ethCommon.HexToHash(hash)
	receipt, err := e.ethClient.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return 0, err
	}
	slot, err := common.GetSlot(receipt.BlockNumber.Int64())
	if err != nil {
		return 0, err
	}
	return slot, nil
}

func (e *EthereumAgent) getRequestProofData(zkType common.ZkProofType, index uint64, txHash string) (interface{}, bool, error) {
	switch zkType {
	case common.TxInEth2:
		return e.getTxInEth2Data(txHash)
	case common.BeaconHeaderType:
		return e.getBlockHeaderRequestData(index)
	case common.RedeemTxType:
		data, ok, err := GetRedeemRequestData(e.fileStore, e.genesisPeriod, index, txHash, e.beaconClient, e.ethClient.Client)
		if err != nil {
			logger.Error("get redeem request data error: %v %v", index, err)
			return nil, false, err
		}
		return data, ok, nil
	case common.BeaconHeaderFinalityType:
		data, ok, err := GetBhfUpdateData(e.fileStore, index, e.genesisPeriod)
		if err != nil {
			logger.Error("get bhf update data error: %v %v", index, err)
			return nil, false, err
		}
		return data, ok, nil
	default:
		return nil, false, fmt.Errorf("never should happen: %v", zkType)
	}
}

func (e *EthereumAgent) checkRequest(zkType common.ZkProofType, index uint64, txHash string) (bool, error) {
	switch zkType {
	case common.TxInEth2:
		finalizedSlot, ok, err := e.fileStore.GetFinalizedSlot()
		if err != nil {
			logger.Error("get latest slot error: %v", err)
			return false, err
		}
		if !ok {
			logger.Warn("no find latest slot")
			return false, nil
		}
		receipt, err := e.ethClient.TransactionReceipt(context.Background(), ethCommon.HexToHash(txHash))
		if err != nil {
			logger.Error("get tx receipt error: %v", err)
			return false, err
		}
		txSlot, err := common.GetSlot(receipt.BlockNumber.Int64())
		if err != nil {
			logger.Error("get slot error: %v", err)
			return false, err
		}
		// todo
		if txSlot < finalizedSlot {
			return true, nil
		}
		logger.Warn("%v tx slot %v less than finalized slot %v", txHash, txSlot, finalizedSlot)
		return false, nil
	case common.BeaconHeaderType:
		// todo
		_, ok, err := e.fileStore.GetNearTxSlotFinalizedSlot(index)
		if err != nil {
			logger.Error("get latest slot error: %v", err)
			return false, err
		}
		if !ok {
			return false, nil
		}
		return true, nil
	case common.RedeemTxType:
		// todo
		return true, nil
	case common.BeaconHeaderFinalityType:
		return index >= e.genesisSlot, nil

	default:
		return false, fmt.Errorf("invalid zkType: %v", zkType)

	}
}

func (e *EthereumAgent) getBlockHeaderRoot(slot uint64) ([]byte, error) {
	response, err := e.beaconClient.BeaconHeaderBySlot(slot)
	if err != nil {
		logger.Error("get block header root error: %v", err)
		return nil, err
	}
	consensus, err := response.Data.Header.Message.ToConsensus()
	if err != nil {
		logger.Error("to consensus error: %v", err)
		return nil, err
	}
	treeRoot, err := consensus.HashTreeRoot()
	if err != nil {
		logger.Error("hash tree root error: %v", err)
		return nil, err
	}
	return treeRoot[0:], nil
}

func (e *EthereumAgent) sendZkProofRequest(requests ...*common.ZkProofRequest) error {
	e.proofRequest <- requests
	for _, req := range requests {
		logger.Info("send request: %v", req.Id())
		e.stateCache.Store(req.Id(), req)
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
	return "ethereumAgent"
}

func NewRedeemProofParam(txId string, txData *ethblock.TxInEth2ProofData) RedeemProofParam {
	return RedeemProofParam{
		TxHash: txId,
	}
}

func NewRedeemProof(txId string, status common.ProofStatus) Proof {
	return Proof{
		TxHash:    txId,
		ProofType: common.TxInEth2,
		Status:    int(status),
	}
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
func NewRedeemEthTx(height uint64, txIndex uint, txHash, sender, btcTxId string, inputs []Utxo, outputs []TxOut) *Transaction {
	return &Transaction{
		Height:    height,
		TxIndex:   txIndex,
		TxHash:    txHash,
		ChainType: Ethereum,
		TxType:    RedeemTx,
		BtcTxId:   btcTxId,
		From:      sender,
	}
}

// todo
func (e *EthereumAgent) GetSyncCommitRootID(period uint64) ([]byte, bool, error) {
	var currentPeriodUpdate structs.LightClientUpdateWithVersion
	exists, err := e.fileStore.GetUpdate(period, &currentPeriodUpdate)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v Index update Data, send new update request", period)
		return nil, false, nil
	}
	// todo
	var update utils.SyncCommitteeUpdate
	update.Version = currentPeriodUpdate.Version
	err = ParseObj(currentPeriodUpdate.Data, &update)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if e.genesisPeriod == period {
		var genesisData structs.LightClientBootstrapResponse
		genesisExists, err := e.fileStore.GetBootstrap(&genesisData)
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
		err = ParseObj(genesisData.Data.CurrentSyncCommittee, &genesisCommittee)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		update.CurrentSyncCommittee = &genesisCommittee
	} else {
		prePeriod := period - 1
		if prePeriod < e.genesisPeriod {
			logger.Error("should never happen: %v", prePeriod)
			return nil, false, nil
		}
		var preUpdateData structs.LightClientUpdateWithVersion
		preUpdateExists, err := e.fileStore.GetUpdate(prePeriod, &preUpdateData)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		if !preUpdateExists {
			logger.Warn("get unit Data,no find %v Index update Data, send new update request", prePeriod)
			return nil, false, nil
		}
		// todo
		var currentSyncCommittee utils.SyncCommittee
		err = ParseObj(preUpdateData.Data.NextSyncCommittee, &currentSyncCommittee)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		update.CurrentSyncCommittee = &currentSyncCommittee
	}
	rootId, err := circuits.SyncCommitRoot(&update)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	return rootId, true, nil

}
