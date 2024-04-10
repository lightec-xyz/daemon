package node

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	dcommon "github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/rpc/oasis"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	btctx "github.com/lightec-xyz/daemon/transaction/bitcoin"
	"github.com/lightec-xyz/daemon/transaction/ethereum"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
)

type EthereumAgent struct {
	btcClient        *bitcoin.Client
	ethClient        *ethrpc.Client
	apiClient        *apiclient.Client // todo temporary use
	oasisClient      *oasis.Client     // todo temporary use
	store            store.IStore
	memoryStore      store.IStore
	fileStore        *FileStore
	blockTime        time.Duration
	taskManager      *TaskManager
	whiteList        map[string]bool
	proofRequest     chan []*dcommon.ZkProofRequest
	exitSign         chan struct{}
	multiAddressInfo MultiAddressInfo
	btcNetwork       btctx.NetWork
	logAddrFilter    EthAddrFilter
	privateKeys      []*btcec.PrivateKey //todo just test
	initStartHeight  int64
	ethSubmitAddress string
	autoSubmit       bool
	submitQueue      *Queue
}

func NewEthereumAgent(cfg NodeConfig, submitTxEthAddr string, fileStore *FileStore, store, memoryStore store.IStore, beaClient *apiclient.Client,
	btcClient *bitcoin.Client, ethClient *ethrpc.Client, proofRequest chan []*dcommon.ZkProofRequest) (IAgent, error) {
	var privateKeys []*btcec.PrivateKey
	for _, secret := range cfg.BtcPrivateKeys {
		hexPriv, err := hex.DecodeString(secret)
		if err != nil {
			logger.Error("decode private key error:%v", err)
			return nil, err
		}
		privKey, _ := btcec.PrivKeyFromBytes(hexPriv)
		privateKeys = append(privateKeys, privKey)
	}
	return &EthereumAgent{
		apiClient:        beaClient, // todo
		btcClient:        btcClient,
		ethClient:        ethClient,
		store:            store,
		fileStore:        fileStore,
		memoryStore:      memoryStore,
		blockTime:        cfg.EthScanBlockTime,
		proofRequest:     proofRequest,
		exitSign:         make(chan struct{}, 1),
		whiteList:        make(map[string]bool),
		multiAddressInfo: cfg.MultiAddressInfo,
		btcNetwork:       btctx.NetWork(cfg.BtcNetwork),
		privateKeys:      privateKeys,
		initStartHeight:  cfg.EthInitHeight,
		ethSubmitAddress: submitTxEthAddr,
		autoSubmit:       cfg.AutoSubmit,
		logAddrFilter:    cfg.EthAddrFilter,
		submitQueue:      NewQueue(),
	}, nil
}

func (e *EthereumAgent) Init() error {
	logger.Info("init ethereum agent")
	exists, err := CheckEthereumHeight(e.store)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return err
	}
	if exists {
		logger.Debug("ethereum agent check uncompleted generate Proof tx")
		err := e.checkUnGenerateProof()
		if err != nil {
			logger.Error("check uncompleted generate Proof tx error:%v", err)
			return err
		}
	} else {
		logger.Debug("init eth current height: %v", e.initStartHeight)
		err := WriteEthereumHeight(e.store, e.initStartHeight)
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
	return nil
}

func (e *EthereumAgent) checkUnGenerateProof() error {
	// todo
	return nil
}

func (e *EthereumAgent) getEthHeight() (int64, error) {
	return ReadEthereumHeight(e.store)
}

func (e *EthereumAgent) ScanBlock() error {
	ethHeight, err := e.getEthHeight()
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return err
	}
	if ethHeight < e.initStartHeight {
		ethHeight = e.initStartHeight
	}
	blockNumber, err := e.ethClient.BlockNumber(context.Background())
	if err != nil {
		logger.Error("get eth block number error:%v", err)
		return err
	}
	blockNumber = blockNumber - 0
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
		zkProofRequests, err := toRedeemZkProofRequest(requests)
		if err != nil {
			logger.Error("to redeem zk Proof request error: %v %v", index, err)
			return err
		}
		e.proofRequest <- zkProofRequests
	}
	return nil
}

func (e *EthereumAgent) ProofResponse(resp *dcommon.ZkProofResponse) error {
	logger.Info("receive proof response: %v", resp)
	hexProof := hex.EncodeToString(resp.Proof)
	err := e.updateRedeemProof(resp.TxHash, hexProof, resp.Status)
	if err != nil {
		logger.Error("update Proof error:%v", err)
		return err
	}
	_, err = e.RedeemBtcTx(resp.TxHash, resp.Proof)
	if err != nil {
		logger.Error("redeem btc tx error:%v", err)
		return err
	}
	// todo
	//if e.autoSubmit {
	//	_, err := e.taskManager.RedeemBtcRequest(resp.TxHash, nil, nil, nil)
	//	if err != nil {
	//		logger.Error("submit redeem request error:%v", err)
	//		return err
	//	}
	//}
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
				redeemTxProof = NewRedeemProof(redeemTx.TxHash, dcommon.ProofSuccess)
			} else {
				// Todo
				txData, err := ethblock.GenerateTxInEth2Proof(e.ethClient.Client, e.apiClient, redeemTx.TxHash)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				requests = append(requests, NewRedeemProofParam(redeemTx.TxHash, txData))
				redeemTxProof = NewRedeemProof(redeemTx.TxHash, dcommon.ProofDefault)
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

// todo refactor

func (e *EthereumAgent) RedeemBtcTx(txHash string, proof []byte) (interface{}, error) {
	ethTxHash := common.HexToHash(txHash)
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
	TxHash, err := e.btcClient.Sendrawtransaction(txHex)
	if err != nil {
		logger.Error("send btc tx error:%v", err)
		return "", err
	}
	logger.Info("send redeem btc tx: %v", TxHash)
	return TxHash, nil
}

func (e *EthereumAgent) CheckState() error {
	return nil
}

func (e *EthereumAgent) updateRedeemProof(txId string, proof string, status dcommon.ProofStatus) error {
	logger.Debug("update Redeem Proof status: %v %v %v", txId, proof, status)
	err := UpdateProof(e.store, txId, proof, dcommon.RedeemTxType, status)
	if err != nil {
		logger.Error("update Proof error: %v %v", txId, err)
		return err
	}
	return nil
}

func (e *EthereumAgent) Close() error {
	close(e.exitSign)
	return nil
}
func (e *EthereumAgent) Name() string {
	return "Ethereum WrapperAgent"
}

func NewRedeemProofParam(txId string, txData *ethblock.TxInEth2ProofData) RedeemProofParam {
	return RedeemProofParam{
		TxHash: txId,
		TxData: txData,
	}
}

func NewRedeemProof(txId string, status dcommon.ProofStatus) Proof {
	return Proof{
		TxHash:    txId,
		ProofType: dcommon.TxInEth2,
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
