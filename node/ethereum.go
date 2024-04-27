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
	"github.com/lightec-xyz/daemon/store"
	btctx "github.com/lightec-xyz/daemon/transaction/bitcoin"
	"github.com/lightec-xyz/daemon/transaction/ethereum"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	txineth2 "github.com/lightec-xyz/provers/circuits/tx-in-eth2"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"strconv"
	"strings"
)

type EthereumAgent struct {
	btcClient        *bitcoin.Client
	ethClient        *ethrpc.Client
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
}

func NewEthereumAgent(cfg Config, genesisPeriod uint64, fileStore *FileStorage, store, memoryStore store.IStore, beaClient *apiclient.Client,
	btcClient *bitcoin.Client, ethClient *ethrpc.Client, beaconClient *beacon.Client, proofRequest chan []*common.ZkProofRequest, task *TaskManager) (IAgent, error) {
	return &EthereumAgent{
		apiClient:        beaClient, // todo
		btcClient:        btcClient,
		ethClient:        ethClient,
		beaconClient:     beaconClient,
		store:            store,
		fileStore:        fileStore,
		memoryStore:      memoryStore,
		genesisPeriod:    genesisPeriod,
		proofRequest:     proofRequest,
		multiAddressInfo: cfg.MultiAddressInfo,
		initHeight:       cfg.EthInitHeight,
		logAddrFilter:    cfg.EthAddrFilter,
		task:             task,
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
	err = WriteUnGenProof(e.store, Ethereum, []string{"0x554ff1fa98d6dfd812bcbca1ed40bf8022403dadd5373ceff4a1df5b7b19d484"})
	if err != nil {
		logger.Error("write ungen proof error: %v", err)
		return err
	}
	return nil
}

func (e *EthereumAgent) checkUnGenerateProof() error {
	// todo
	return nil
}

func (e *EthereumAgent) ScanBlock() error {
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
		if index%100 == 0 {
			logger.Debug("ethereum parse block:%d", index)
		}
		redeemTxes, depositTxes, requests, proofs, err := e.parseBlock(index)
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
		err = e.saveRedeemData(redeemTxes, proofs, requests)
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
	e.stateCache.DeleteZkRequest(proofId)
	if resp.ZkProofType == common.RedeemTxType {
		err = e.updateRedeemProof(resp.TxHash, hex.EncodeToString(resp.Proof), resp.Status)
		if err != nil {
			logger.Error("update Proof error:%v", err)
			return err
		}
		_, err = RedeemBtcTx(e.btcClient, resp.TxHash, resp.Proof)
		if err != nil {
			logger.Error("redeem btc tx error:%v", err)
			e.task.AddTask(resp)
			return err
		}
	}
	return nil
}

func (e *EthereumAgent) updateDepositData(height int64, depositTxes []Transaction) error {
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
		err = DeleteUnGenProof(e.store, Bitcoin, tx.BtcTxId)
		if err != nil {
			logger.Error("delete ungen proof error: %v %v", height, err)
			return err
		}
	}
	return nil
}

func (e *EthereumAgent) saveTransaction(height int64, txes []Transaction) error {
	err := WriteEthereumTxIds(e.store, height, txesToTxIds(txes))
	if err != nil {
		logger.Error("write ethereum tx ids error: %v %v", height, err)
		return err
	}
	err = WriteEthereumTx(e.store, txesToDbTxes(txes))
	if err != nil {
		logger.Error("put redeem tx error: %v %v", height, err)
		return err
	}
	return nil
}

func (e *EthereumAgent) saveRedeemData(redeemTxes []Transaction, proofs []Proof, requests []RedeemProofParam) error {

	err := WriteDbProof(e.store, proofsToDbProofs(proofs))
	if err != nil {
		logger.Error("put eth current height error:%v", err)
		return err
	}
	for _, tx := range redeemTxes {
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
	err = WriteUnGenProof(e.store, Ethereum, redeemToTxHashList(requests))
	if err != nil {
		logger.Error("write ungen Proof error: %v", err)
		return err
	}
	return nil
}

func (e *EthereumAgent) parseBlock(height int64) ([]Transaction, []Transaction, []RedeemProofParam, []Proof, error) {
	block, err := e.ethClient.GetBlock(height)
	if err != nil {
		logger.Error("ethereum rpc get block error:%v", err)
		return nil, nil, nil, nil, err
	}
	blockHash := block.Hash().String()
	logFilters, topicFilters := e.logAddrFilter.FilterLogs()
	logs, err := e.ethClient.GetLogs(blockHash, logFilters, topicFilters)
	if err != nil {
		logger.Error("ethereum rpc get logs error:%v", err)
		return nil, nil, nil, nil, err
	}
	var redeemTxes []Transaction
	var depositTxes []Transaction
	var proofs []Proof
	var requests []RedeemProofParam
	for _, log := range logs {
		depositTx, isDeposit, err := e.isDepositTx(log)
		if err != nil {
			logger.Error("check is deposit tx error:%v", err)
			return nil, nil, nil, nil, err
		}
		if isDeposit {
			depositTxes = append(depositTxes, depositTx)
			continue
		}
		redeemTx, isRedeem, err := e.isRedeemTx(log)
		if err != nil {
			logger.Error("check is redeem tx error:%v", err)
			return nil, nil, nil, nil, err
		}
		if isRedeem {
			submitted, err := e.btcClient.CheckTx(redeemTx.BtcTxId)
			if err != nil {
				logger.Error("check btc tx error:%v", err)
				return nil, nil, nil, nil, err
			}
			var redeemTxProof Proof
			if submitted {
				redeemTxProof = NewRedeemProof(redeemTx.TxHash, common.ProofSuccess)
			} else {
				requests = append(requests, NewRedeemProofParam(redeemTx.TxHash, nil))
				redeemTxProof = NewRedeemProof(redeemTx.TxHash, common.ProofDefault)
			}
			proofs = append(proofs, redeemTxProof)
			redeemTxes = append(redeemTxes, redeemTx)
		}
	}
	return redeemTxes, depositTxes, requests, proofs, nil
}

func (e *EthereumAgent) isDepositTx(log types.Log) (Transaction, bool, error) {
	if log.Removed {
		return Transaction{}, false, nil
	}
	if len(log.Topics) != 3 {
		return Transaction{}, false, nil
	}
	// todo
	if strings.ToLower(log.Address.Hex()) == e.logAddrFilter.LogDepositAddr && strings.ToLower(log.Topics[0].Hex()) == e.logAddrFilter.LogTopicDepositAddr {
		btcTxId := strings.ToLower(log.Topics[1].Hex())
		hexVout := strings.TrimPrefix(strings.ToLower(log.Topics[2].Hex()), "0x")
		vout, err := strconv.ParseInt(hexVout, 16, 32)
		if err != nil {
			logger.Error("parse vout error:%v", err)
			return Transaction{}, false, err
		}
		amount, err := strconv.ParseInt(fmt.Sprintf("%x", log.Data), 16, 64)
		if err != nil {
			logger.Error("parse amount error:%v", err)
			return Transaction{}, false, err
		}
		utxo := []Utxo{
			{
				TxId:  btcTxId,
				Index: uint32(vout),
			},
		}
		txHash := log.TxHash.String()
		logger.Info("ethereum agent find deposit zkbtc ethTxHash:%v,btcTxId:%v,utxo:%v",
			txHash, btcTxId, formatUtxo(utxo))
		depositTx := NewDepositEthTx(txHash, btcTxId, utxo, amount)
		return depositTx, true, nil
	} else {
		return Transaction{}, false, nil
	}
}

func (e *EthereumAgent) isRedeemTx(log types.Log) (Transaction, bool, error) {
	redeemTx := Transaction{}
	if log.Removed {
		return redeemTx, false, nil
	}
	if len(log.Topics) != 2 {
		return redeemTx, false, nil
	}
	//todo more check
	if strings.ToLower(log.Address.Hex()) == e.logAddrFilter.LogRedeemAddr && strings.ToLower(log.Topics[0].Hex()) == e.logAddrFilter.LogTopicRedeemAddr {
		btcTxId := strings.ToLower(log.Topics[1].Hex())
		txData, _, err := ethereum.DecodeRedeemLog(log.Data)
		if err != nil {
			logger.Error("decode redeem log error:%v", err)
			return redeemTx, false, err
		}
		transaction := btctx.NewTransaction()
		err = transaction.Deserialize(bytes.NewReader(txData))
		if err != nil {
			logger.Error("deserialize btc tx error:%v", err)
			return redeemTx, false, err
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
			logger.Error("never should happen btc tx not match error: TxHash:%v, logBtcTxId:%v,decodeTxHash:%v", txHash, btcTxId, transaction.TxHash().String())
			return redeemTx, false, fmt.Errorf("tx hash not match:%v", txHash)
		}
		logger.Info("ethereum agent find redeem zkbtc  ethTxHash:%v,btcTxId:%v,input:%v,output:%v",
			redeemTx.TxHash, btcTxId, formatUtxo(inputs), formatOut(outputs))
		redeemTx = NewRedeemEthTx(txHash, btcTxId, inputs, outputs)
		return redeemTx, true, nil
	} else {
		return redeemTx, false, nil
	}

}

func (e *EthereumAgent) CheckState() error {
	logger.Debug("ethereum check state ...")
	unGenProofs, err := ReadAllUnGenProofIds(e.store, Ethereum)
	if err != nil {
		logger.Error("read all ungen proof ids error: %v", err)
		return err
	}
	for _, unGenProof := range unGenProofs {
		txHash := unGenProof.TxId
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
		// todo
		txSlot, err := e.GetSlotByHash(txHash)
		if err != nil {
			logger.Error("get txSlot error: %v", err)
			return err
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

func (e *EthereumAgent) tryProofRequest(zkType common.ZkProofType, index uint64, txHash string) error {
	proofId := common.NewProofId(zkType, index, txHash)
	exists := e.stateCache.CheckZkRequest(proofId)
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

func (e *EthereumAgent) getRedeemRequestData(txSlot uint64, txHash string) (rpc.RedeemRequest, bool, error) {
	txProof, ok, err := e.fileStore.GetTxProof(txHash)
	if err != nil {
		logger.Error("get tx proof error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	if !ok {
		logger.Debug("proof request data not prepared: %v", txHash)
		return rpc.RedeemRequest{}, false, nil
	}
	blockHeaderProof, ok, err := e.fileStore.GetBeaconHeaderProof(txSlot)
	if err != nil {
		logger.Error("get block header proof error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	if !ok {
		logger.Debug("proof request data not prepared: %v", txSlot)
		return rpc.RedeemRequest{}, false, nil
	}
	finalizedSlot, ok, err := e.fileStore.GetNearTxSlotFinalizedSlot(txSlot)
	if err != nil {
		logger.Error("get bhf update proof error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	if !ok {
		logger.Debug("proof request data not prepared: %v", txSlot)
		return rpc.RedeemRequest{}, false, nil
	}
	bhfProof, ok, err := e.fileStore.GetBhfProof(finalizedSlot)
	if err != nil {
		logger.Error("get bhf update proof error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	if !ok {
		logger.Warn("no find bhf update %v", finalizedSlot)
		return rpc.RedeemRequest{}, false, nil
	}
	genesisRoot, ok, err := e.GetSyncCommitRootID(e.genesisPeriod)
	if err != nil {
		logger.Error("get genesis root error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	if !ok {
		logger.Warn("no find genesis root %v", e.genesisPeriod)
		return rpc.RedeemRequest{}, false, nil
	}

	var finalityUpdate *structs.LightClientUpdateWithVersion
	ok, err = e.fileStore.GetFinalityUpdate(finalizedSlot, &finalityUpdate)
	if err != nil {
		logger.Error("get finality update error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	if !ok {
		logger.Warn("no find finality update %v", finalizedSlot)
		return rpc.RedeemRequest{}, false, nil
	}

	attestedSlot, err := strconv.ParseUint(finalityUpdate.Data.AttestedHeader.Slot, 10, 64)
	if err != nil {
		logger.Error("parse slot error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	period := attestedSlot / common.SlotPerPeriod
	currentRoot, ok, err := e.GetSyncCommitRootID(period)
	if err != nil {
		logger.Error("get current root error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	beginID, endId, err := e.GetBeaconHeaderId(txSlot, finalizedSlot)
	if err != nil {
		logger.Error("get begin and end id error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	// todo need cache
	txVar, receiptVar, err := txineth2.GenerateTxAndReceiptU128Padded(e.ethClient.Client, txHash)
	if err != nil {
		logger.Error("get tx and receipt error: %v", err)
		return rpc.RedeemRequest{}, false, err
	}
	redeemRequest := rpc.RedeemRequest{
		TxProof:          txProof.Proof,
		TxWitness:        txProof.Witness,
		BhProof:          blockHeaderProof.Proof,
		BhWitness:        blockHeaderProof.Witness,
		BhfProof:         bhfProof.Proof,
		BhfWitness:       bhfProof.Witness,
		GenesisScRoot:    hex.EncodeToString(genesisRoot),
		BeginId:          hex.EncodeToString(beginID),
		EndId:            hex.EncodeToString(endId),
		CurrentSCSSZRoot: hex.EncodeToString(currentRoot),
		TxVar:            txVar,
		ReceiptVar:       receiptVar,
	}
	return redeemRequest, true, nil

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
		return e.getRedeemRequestData(index, txHash)
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
		e.stateCache.StoreZkRequest(req.Id(), req)
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

func NewDepositEthTx(txHash, btcTxId string, utxo []Utxo, amount int64) Transaction {
	return Transaction{
		TxHash:    txHash,
		ChainType: Ethereum,
		TxType:    DepositTx,
		BtcTxId:   btcTxId,
	}
}
func NewRedeemEthTx(txHash string, btcTxId string, inputs []Utxo, outputs []TxOut) Transaction {
	return Transaction{
		TxHash:    txHash,
		ChainType: Ethereum,
		TxType:    RedeemTx,
		BtcTxId:   btcTxId,
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
		logger.Warn("no find %v period update Data, send new update request", period)
		return nil, false, nil
	}
	// todo
	var update utils.LightClientUpdateInfo
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
			logger.Warn("get unit Data,no find %v period update Data, send new update request", prePeriod)
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
