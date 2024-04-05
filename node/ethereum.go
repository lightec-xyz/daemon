package node

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
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
	apiClient        *apiclient.Client // todo temp use
	store            store.IStore
	memoryStore      store.IStore
	fileStore        *FileStore
	blockTime        time.Duration
	taskManager      *TaskManager
	whiteList        map[string]bool
	proofRequest     chan []ZkProofRequest
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
	btcClient *bitcoin.Client, ethClient *ethrpc.Client, proofRequest chan []ZkProofRequest) (IAgent, error) {
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
		logger.Debug("ethereum parse block:%d", index)
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
			logger.Error("ethereum save data error: %v %v", index, err)
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

func (e *EthereumAgent) ProofResponse(resp ZkProofResponse) error {
	logger.Info("receive proof response: %v", resp.TxHash)
	err := e.updateRedeemProof(resp.TxHash, resp.ProofStr, resp.Status)
	if err != nil {
		logger.Error("update Proof error:%v", err)
		return err
	}
	// todo
	if e.autoSubmit {
		_, err := e.taskManager.RedeemBtcRequest(resp.TxHash, nil, nil, nil)
		if err != nil {
			logger.Error("submit redeem request error:%v", err)
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
				redeemTxProof = NewRedeemProof(redeemTx.TxHash, ProofSuccess)
			} else {
				// Todo
				txData, err := ethblock.GenerateTxInEth2Proof(e.ethClient.Client, e.apiClient, redeemTx.TxHash)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				requests = append(requests, NewRedeemProofParam(redeemTx.TxHash, txData))
				redeemTxProof = NewRedeemProof(redeemTx.TxHash, ProofDefault)
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

//func (e *EthereumAgent) RedeemBtcTx(resp RedeemProof) (string, error) {
////todo
//var txIns []btctx.TxIn
//logger.Debug("************************************")
//for _, input := range resp.Inputs {
//	utxo, err := e.btcClient.GetUtxoByTxId(input.TxHash, int(input.Index))
//	if err != nil {
//		logger.Error("get utxo error:%v", err)
//		return "", err
//	}
//	logger.Debug(fmt.Sprintf("utxo:%v", utxo.Amount))
//	amount := BtcToSat(utxo.Amount)
//	in := btctx.TxIn{
//		Hash:     input.TxHash,
//		VOut:     input.Index,
//		PkScript: utxo.ScriptPubKey,
//		Amount:   amount,
//	}
//	txIns = append(txIns, in)
//	logger.Debug("txIn: txid:%v, index:%v, amount:%v ,scriptPubKey:%v", input.TxHash, input.Index, amount, utxo.ScriptPubKey)
//}
//
//builder := btctx.NewMultiTransactionBuilder()
//err := builder.NetParams(e.btcNetwork)
//if err != nil {
//	logger.Error("multi btc tx net params error:%v", err)
//	return "", err
//}
//err = builder.AddMultiPublicKey(e.multiAddressInfo.PublicKeyList, e.multiAddressInfo.NRequired)
//if err != nil {
//	logger.Error("multi btc tx add public key error:%v", err)
//	return "", err
//}
//
//err = builder.AddTxIn(txIns)
//if err != nil {
//	logger.Error("multi btc tx add txIn error:%v", err)
//	return "", err
//}
//txOuts := []btctx.TxOut{}
//for _, output := range resp.Outputs {
//	txOuts = append(txOuts, btctx.TxOut{
//		PayScript: output.PkScript,
//		Amount:    output.Value,
//	})
//	logger.Debug("txOut: pkScript:%x, amount:%v", output.PkScript, output.Value)
//}
//err = builder.AddTxOutScript(txOuts)
//if err != nil {
//	logger.Error("multi btc tx add txOut error:%v", err)
//	return "", err
//}
//err = builder.Sign(func(hash []byte) ([][]byte, error) {
//	// todo
//	var sigs [][]byte
//	for _, privkey := range e.privateKeys {
//		sig := ecdsa.Sign(privkey, hash)
//		sigWithType := append(sig.Serialize(), byte(txscript.SigHashAll))
//		sigs = append(sigs, sigWithType)
//	}
//	return sigs, nil
//
//})
//logger.Debug("************************************")
//if err != nil {
//	logger.Error("multi tx sign error:%v", err)
//	return "", err
//}
//txBytes, err := builder.Build()
//if err != nil {
//	logger.Error("build btc tx error:%v", err)
//	return "", err
//}
//logger.Info("redeem btc tx hash: %v", builder.TxHash())
//TxHash, err := e.btcClient.Sendrawtransaction(hex.EncodeToString(txBytes))
//if err != nil {
//	logger.Error("send btc tx error:%v", err)
//	return "", err
//}
//logger.Info("send redeem btc tx: %v", TxHash)
//return TxHash, nil
//}

func (e *EthereumAgent) CheckState() error {
	panic(e)
}

func (e *EthereumAgent) updateRedeemProof(txId string, proof string, status ProofStatus) error {
	logger.Debug("update Redeem Proof status: %v %v %v", txId, proof, status)
	err := UpdateProof(e.store, txId, proof, RedeemTxType, status)
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

func NewRedeemProof(txId string, status ProofStatus) Proof {
	return Proof{
		TxHash:    txId,
		ProofType: RedeemTxType,
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
