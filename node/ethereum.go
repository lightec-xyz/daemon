package node

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	btctx "github.com/lightec-xyz/daemon/transaction/bitcoin"
	"strconv"
	"time"
)

type EthereumAgent struct {
	btcClient            *bitcoin.Client
	ethClient            *ethereum.Client
	store                store.IStore
	memoryStore          store.IStore
	blockTime            time.Duration
	whiteList            map[string]bool
	checkProofHeightNums int64
	proofResponse        <-chan ProofResponse
	proofRequest         chan []ProofRequest
	exitSign             chan struct{}
	multiAddressInfo     MultiAddressInfo
	btcNetwork           btctx.NetWork
	logAddr              []string
	logTopic             []string
	privateKeys          []*btcec.PrivateKey
}

func NewEthereumAgent(cfg NodeConfig, store, memoryStore store.IStore, btcClient *bitcoin.Client, ethClient *ethereum.Client,
	proofRequest chan []ProofRequest, proofResponse <-chan ProofResponse) (IAgent, error) {
	//todo
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
		btcClient:            btcClient,
		ethClient:            ethClient,
		store:                store,
		memoryStore:          memoryStore,
		blockTime:            time.Duration(cfg.EthBlockTime) * time.Second,
		proofRequest:         proofRequest,
		proofResponse:        proofResponse,
		checkProofHeightNums: 100,
		exitSign:             make(chan struct{}, 1),
		whiteList:            make(map[string]bool),
		multiAddressInfo:     cfg.MultiAddressInfo,
		btcNetwork:           btctx.NetWork(cfg.BtcNetwork),
		logAddr:              cfg.LogAddr,
		logTopic:             cfg.LogTopic,
		privateKeys:          privateKeys,
	}, nil
}

func (e *EthereumAgent) Init() error {
	logger.Info("init ethereum agent")
	has, err := e.store.Has(ethCurHeightKey)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return err
	}
	if has {
		err := e.checkUnCompleteGenerateProofTx()
		if err != nil {
			logger.Error("check uncomplete generate proof tx error:%v", err)
			return err
		}
	} else {
		logger.Debug("init eth current height: %v", InitEthereumHeight)
		err := e.store.PutObj(ethCurHeightKey, InitEthereumHeight)
		if err != nil {
			logger.Error("put init eth current height error:%v", err)
			return err
		}
	}
	// test rpc
	_, err = e.ethClient.GetChainId()
	if err != nil {
		logger.Error("get chain id error:%v", err)
		return err
	}

	//todo
	return nil
}

func (e *EthereumAgent) checkUnCompleteGenerateProofTx() error {
	currentHeight, err := e.getEthHeight()
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}
	start := currentHeight - e.checkProofHeightNums
	var proofList []ProofRequest
	for index := start; index < currentHeight; index++ {
		hasObj, err := e.store.HasObj(index)
		if err != nil {
			logger.Error("get txIdList error:%v", err)
			return err
		}
		if !hasObj {
			continue
		}
		var txIdList []string
		err = e.store.GetObj(index, &txIdList)
		if err != nil {
			logger.Error("get txIdList error:%v", err)
			return err
		}
		for _, txId := range txIdList {
			var proof TxProof
			err := e.store.GetObj(TxIdToProofId(txId), &proof)
			if err != nil {
				logger.Error("get proof error:%v", err)
				return err
			}
			proofList = append(proofList, ProofRequest{
				TxId:  proof.TxId,
				PType: Redeem,
				Msg:   proof.Msg,
			})
		}
	}
	e.proofRequest <- proofList
	return nil
}

func (e *EthereumAgent) getEthHeight() (int64, error) {
	var curHeight int64
	err := e.store.GetObj(ethCurHeightKey, &curHeight)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return 0, err
	}
	return curHeight, nil
}

func (e *EthereumAgent) ScanBlock() error {
	logger.Info("start ethereum scan block")
	ethHeight, err := e.getEthHeight()
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return err
	}

	ethHeight = 456924 - 1
	blockNumber, err := e.ethClient.EthBlockNumber()
	if err != nil {
		logger.Error("get eth block number error:%v", err)
		return err
	}
	//todo
	if ethHeight >= int64(blockNumber) {
		logger.Info("current height:%d,node block count:%d", ethHeight, blockNumber)
		return nil
	}
	for index := ethHeight + 1; index <= int64(blockNumber); index++ {
		redeemTxList, proofRequestList, err := e.parseBlock(index)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		err = e.saveDataToDb(index, redeemTxList)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		e.proofRequest <- proofRequestList

	}
	return nil
}

func (e *EthereumAgent) Transfer() {
	//todo
	logger.Info("start ethereum transfer goroutine")
	for {
		select {
		case <-e.exitSign:
			logger.Info("exit ethereum transfer goroutine")
			return
		case response := <-e.proofResponse:
			logger.Info("receive redeem proof response:%v %v", response.TxId, response)
			err := e.updateProof(response)
			if err != nil {
				logger.Error("update proof error:%v", err)
				continue
			}
			err = e.RedeemBtcTx(response)
			if err != nil {
				//todo
				logger.Error("redeem btc tx error:%v", err)
				continue
			}
			logger.Info("success redeem btc tx:%v", response)
		}
	}

}

func (e *EthereumAgent) saveDataToDb(height int64, list []RedeemTx) error {
	var txIdList []string
	for _, tx := range list {
		err := e.store.BatchPutObj(tx.TxId, tx)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		pTxId := fmt.Sprintf("%s%s", ProofPrefix, tx.TxId)
		err = e.store.BatchPutObj(pTxId, TxProof{
			PTxId: pTxId,
		})
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	err := e.store.BatchPutObj(height, txIdList)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	err = e.store.BatchPutObj(ethCurHeightKey, height)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	err = e.store.BatchWriteObj()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil

}

func (e *EthereumAgent) parseBlock(height int64) ([]RedeemTx, []ProofRequest, error) {
	block, err := e.ethClient.GetBlock(height)
	if err != nil {
		logger.Error("ethereum rpc get block error:%v", err)
		return nil, nil, err
	}
	logs, err := e.ethClient.GetLogs(block.Hash().String(), e.logAddr, e.logTopic)
	if err != nil {
		logger.Error("ethereum rpc get logs error:%v", err)
		return nil, nil, err
	}
	var redeemTxList []RedeemTx
	var proofRequestList []ProofRequest
	for _, log := range logs {
		redeemTx, ok, err := e.CheckRedeemTx(log)
		if err != nil {
			logger.Error("check redeem tx error:%v", err)
			return nil, nil, err
		}
		if ok {
			logger.Info("found redeem tx: %v", redeemTx)
			redeemTxList = append(redeemTxList, redeemTx)
			proofRequestList = append(proofRequestList, ProofRequest{
				Inputs:  redeemTx.Inputs,
				Outputs: redeemTx.Outputs,
				TxId:    redeemTx.TxId,
				TxIndex: int(redeemTx.TxIndex),
				PType:   Redeem,
			})
		}
	}
	return redeemTxList, proofRequestList, nil
}

func (e *EthereumAgent) getUtxoScriptInfo(txId string, Vout int) (string, int64, error) {
	//todo
	transaction, err := e.btcClient.GetRawtransaction(txId)
	if err != nil {
		logger.Error("get btc transaction error:%v", err)
		return "", 0, err
	}
	for index, out := range transaction.Vout {
		if index == Vout {
			valueBig, err := Str2Big(fmt.Sprintf("%v", out.Value), 8)
			if err != nil {
				logger.Error("parse str to big error:%v", err)
				return "", 0, err
			}
			return out.ScriptPubKey.Hex, valueBig.Int64(), nil
		}
	}
	return "", 0, fmt.Errorf("nevel should happen: %v %v", txId, Vout)

}

func (e *EthereumAgent) RedeemBtcTx(resp ProofResponse) error {
	//todo
	txIns := []btctx.TxIn{}
	for _, input := range resp.Inputs {
		utxoScriptInfo, value, err := e.getUtxoScriptInfo(input.TxId, int(input.Index))
		if err != nil {
			logger.Error("get utxo script info error:%v", err)
			return err
		}
		txIns = append(txIns, btctx.TxIn{
			Hash:     input.TxId,
			VOut:     input.Index,
			PkScript: utxoScriptInfo,
			Amount:   value,
		})
	}

	builder := btctx.NewMultiTransactionBuilder()
	err := builder.NetParams(e.btcNetwork)
	if err != nil {
		logger.Error("build btc tx error:%v", err)
		return err
	}
	err = builder.AddMultiPublicKey(e.multiAddressInfo.PublicKeyList, e.multiAddressInfo.NRequired)
	if err != nil {
		logger.Error("build btc tx error:%v", err)
		return err
	}

	err = builder.AddTxIn(txIns)
	if err != nil {
		logger.Error("build btc tx error:%v", err)
		return err
	}
	txOuts := []btctx.TxOut{}
	for _, output := range resp.Outputs {
		txOuts = append(txOuts, btctx.TxOut{
			PayScript: output.PkScript,
			Amount:    output.Value,
		})
	}
	err = builder.AddTxOutScript(txOuts)
	if err != nil {
		logger.Error("build btc tx error:%v", err)
		return err
	}
	err = builder.Sign(func(hash []byte) ([][]byte, error) {
		// todo
		var sigs [][]byte
		for _, privkey := range e.privateKeys {
			sig := ecdsa.Sign(privkey, hash)
			sigWithType := append(sig.Serialize(), byte(txscript.SigHashAll))
			sigs = append(sigs, sigWithType)
		}
		return sigs, nil

	})
	if err != nil {
		logger.Error("build btc tx error:%v", err)
		return err
	}
	txBytes, err := builder.Build()
	if err != nil {
		logger.Error("build btc tx error:%v", err)
		return err
	}
	txHash, err := e.btcClient.Sendrawtransaction(hex.EncodeToString(txBytes))
	if err != nil {
		logger.Error("send btc tx error:%v", err)
		return err
	}
	logger.Info("send redeem btc tx: %v", txHash)
	return nil
}

func (e *EthereumAgent) CheckRedeemTx(log types.Log) (RedeemTx, bool, error) {
	redeemTx := RedeemTx{}
	//todo more check
	if len(log.Data) <= 64 {
		return redeemTx, false, nil
	}
	//version := log.Data[0:32]
	dataLength := log.Data[32:64]
	l, err := strconv.ParseInt(fmt.Sprintf("%x", dataLength), 16, 32)
	if err != nil {
		return redeemTx, false, err
	}
	txData := log.Data[64 : 64+l]
	transaction := btctx.NewTransaction()
	err = transaction.Deserialize(bytes.NewReader(txData))
	if err != nil {
		return redeemTx, true, nil
	}
	var inputs []TxIn
	for _, in := range transaction.TxIn {
		inputs = append(inputs, TxIn{
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
	redeemTx.Inputs = inputs
	redeemTx.Outputs = outputs
	redeemTx.TxId = log.TxHash.String()
	return redeemTx, true, nil
}

func (e *EthereumAgent) updateProof(resp ProofResponse) error {
	pTxId := TxIdToProofId(resp.TxId)
	err := e.store.PutObj(pTxId, TxProof{
		PTxId:  pTxId,
		TxId:   resp.TxId,
		Msg:    resp.Msg,
		Status: ProofSuccess,
	})
	return err
}

func (e *EthereumAgent) Close() error {
	//todo
	return nil
}
func (e *EthereumAgent) Name() string {
	return "Ethereum Agent"
}
func (e *EthereumAgent) BlockTime() time.Duration {
	return e.blockTime
}
