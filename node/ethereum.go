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
	privateKeys          []*btcec.PrivateKey //todo just env
	initStartHeight      int64
	ethSubmitAddress     string
	autoSubmit           bool
}

func NewEthereumAgent(cfg NodeConfig, store, memoryStore store.IStore, btcClient *bitcoin.Client, ethClient *ethereum.Client,
	proofRequest chan []ProofRequest, proofResponse <-chan ProofResponse) (IAgent, error) {
	submitTxEthAddr, err := privateKeyToEthAddr(cfg.EthPrivateKey)
	if err != nil {
		logger.Error("privateKeyToEthAddr error:%v", err)
		return nil, err
	}
	logger.Info("ethereum submit address:%v", submitTxEthAddr)
	// todo
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
		blockTime:            cfg.EthScanBlockTime,
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
		initStartHeight:      cfg.EthInitHeight,
		ethSubmitAddress:     submitTxEthAddr,
		autoSubmit:           cfg.AutoSubmit,
	}, nil
}

func (e *EthereumAgent) Init() error {
	logger.Info("init ethereum agent")
	exists, err := ReadInitEthereumHeight(e.store)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return err
	}
	if exists {
		logger.Debug("ethereum agent check uncompleted generate proof tx")
		err := e.checkUnCompleteGenerateProofTx()
		if err != nil {
			logger.Error("check uncompleted generate proof tx error:%v", err)
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

func (e *EthereumAgent) checkUnCompleteGenerateProofTx() error {
	//currentHeight, err := e.getEthHeight()
	//if err != nil {
	//	logger.Error("get eth current height error:%v", err)
	//	return err
	//}
	//start := currentHeight - e.checkProofHeightNums
	//var proofList []ProofRequest
	//for index := start; index < currentHeight; index++ {
	//	hasObj, err := e.store.HasObj(index)
	//	if err != nil {
	//		logger.Error("get txIdList error:%v", err)
	//		return err
	//	}
	//	if !hasObj {
	//		continue
	//	}
	//	var txIdList []string
	//	err = e.store.GetObj(index, &txIdList)
	//	if err != nil {
	//		logger.Error("get txIdList error:%v", err)
	//		return err
	//	}
	//	for _, txId := range txIdList {
	//		var proof TxProof
	//		err := e.store.GetObj(TxIdToProofId(txId), &proof)
	//		if err != nil {
	//			logger.Error("get proof error:%v", err)
	//			return err
	//		}
	//		//todo
	//		proofList = append(proofList, ProofRequest{
	//			TxId:      proof.TxId,
	//			ProofType: Redeem,
	//			Msg:       proof.Msg,
	//		})
	//	}
	//}
	//e.proofRequest <- proofList
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
	blockNumber, err := e.ethClient.EthBlockNumber()
	if err != nil {
		logger.Error("get eth block number error:%v", err)
		return err
	}
	blockNumber = blockNumber - 5
	//todo
	if ethHeight >= int64(blockNumber) {
		logger.Debug("eth current height:%d,latest block number :%d", ethHeight, blockNumber)
		return nil
	}
	for index := ethHeight + 1; index <= int64(blockNumber); index++ {
		logger.Debug("ethereum parse block:%d", index)
		redeemTxList, proofRequestList, err := e.parseBlock(index)
		if err != nil {
			logger.Error("eth parse block error: %v %v", index, err)
			return err
		}
		err = e.saveDataToDb(index, redeemTxList)
		if err != nil {
			logger.Error("ethereum save data error: %v %v", index, err)
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
			logger.Info("ethereum transfer goroutine exit now ...")
			return
		case response := <-e.proofResponse:
			logger.Info("receive redeem proof response: %v", response.String())
			err := e.updateProof(response)
			if err != nil {
				logger.Error("update proof error:%v", err)
				continue
			}
			if e.autoSubmit && ProofStatus(response.Status) == ProofSuccess {
				exists, err := e.btcClient.CheckTx(response.BtcTxId)
				if err != nil {
					logger.Error("check btc tx error: %v %v", response.BtcTxId, err)
					continue
				}
				if exists {
					logger.Warn("redeem btc tx submitted: %v", response.BtcTxId)
					continue
				}
				txHash, err := e.RedeemBtcTx(response)
				if err != nil {
					// todo add queue or cli retry
					logger.Error("redeem btc tx error:%v", err)
					continue
				}
				err = e.updateDestChainHash(response.TxId, txHash)
				if err != nil {
					logger.Error("update dest hash error: %v %v", response.TxId, err)
					continue
				}
				logger.Info("success redeem btc tx:%v", response)
			}
		}
	}

}

func (e *EthereumAgent) updateDestChainHash(txId, ethTxHash string) error {
	err := WriteDestChainHash(e.store, txId, ethTxHash)
	if err != nil {
		logger.Error("write dest hash error: %v %v", txId, err)
		return err
	}
	return nil

}

func (e *EthereumAgent) saveDataToDb(height int64, list []*EthereumTx) error {
	err := WriteEthereumTx(e.store, list)
	if err != nil {
		logger.Error("put redeem tx error: %v %v", height, err)
		return err
	}
	var txProof []TxProof
	for _, tx := range list {
		proof := TxProof{
			Height:    height,
			BlockHash: tx.BlockHash,
			TxId:      tx.TxId,
			ProofType: Redeem,
			Proof:     "",
			Status:    ProofDefault,
		}
		txProof = append(txProof, proof)
	}
	err = WriteProof(e.store, txProof)
	if err != nil {
		logger.Error("put eth current height error:%v %v", height, err)
		return err
	}
	err = WriteEthereumHeight(e.store, height)
	if err != nil {
		logger.Error("batch write error: %v %v", height, err)
		return err
	}
	return nil

}

func (e *EthereumAgent) parseBlock(height int64) ([]*EthereumTx, []ProofRequest, error) {
	block, err := e.ethClient.GetBlock(height)
	if err != nil {
		logger.Error("ethereum rpc get block error:%v", err)
		return nil, nil, err
	}
	blockHash := block.Hash().String()
	logs, err := e.ethClient.GetLogs(blockHash, e.logAddr, e.logTopic)
	if err != nil {
		logger.Error("ethereum rpc get logs error:%v", err)
		return nil, nil, err
	}
	var redeemTxList []*EthereumTx
	var proofRequestList []ProofRequest
	for _, log := range logs {
		redeemTx, ok, err := e.CheckRedeemTx(log)
		if err != nil {
			logger.Error("check redeem tx error:%v", err)
			return nil, nil, err
		}
		if ok {
			redeemTx.Height = height
			redeemTx.BlockHash = blockHash
			redeemTxList = append(redeemTxList, redeemTx)
			proofRequestList = append(proofRequestList, ProofRequest{
				Height:    height,
				BlockHash: blockHash,
				Inputs:    redeemTx.Inputs,
				Outputs:   redeemTx.Outputs,
				TxId:      redeemTx.TxId,
				ProofType: Redeem,
			})
			logger.Info("found redeem zkbtc tx: %v", redeemTx.String())
		}
	}
	return redeemTxList, proofRequestList, nil
}

func (e *EthereumAgent) RedeemBtcTx(resp ProofResponse) (string, error) {
	//todo
	var txIns []btctx.TxIn
	logger.Debug("************************************")
	for _, input := range resp.Inputs {
		utxo, err := e.btcClient.GetUtxoByTxId(input.TxId, int(input.Index))
		if err != nil {
			logger.Error("get utxo error:%v", err)
			return "", err
		}
		logger.Debug(fmt.Sprintf("utxo:%v", utxo.Amount))
		amount := BtcToSat(utxo.Amount)
		in := btctx.TxIn{
			Hash:     input.TxId,
			VOut:     input.Index,
			PkScript: utxo.ScriptPubKey,
			Amount:   amount,
		}
		txIns = append(txIns, in)
		logger.Debug("txIn: txid:%v, index:%v, amount:%v ,scriptPubKey:%v", input.TxId, input.Index, amount, utxo.ScriptPubKey)
	}

	builder := btctx.NewMultiTransactionBuilder()
	err := builder.NetParams(e.btcNetwork)
	if err != nil {
		logger.Error("multi btc tx net params error:%v", err)
		return "", err
	}
	err = builder.AddMultiPublicKey(e.multiAddressInfo.PublicKeyList, e.multiAddressInfo.NRequired)
	if err != nil {
		logger.Error("multi btc tx add public key error:%v", err)
		return "", err
	}

	err = builder.AddTxIn(txIns)
	if err != nil {
		logger.Error("multi btc tx add txIn error:%v", err)
		return "", err
	}
	txOuts := []btctx.TxOut{}
	for _, output := range resp.Outputs {
		txOuts = append(txOuts, btctx.TxOut{
			PayScript: output.PkScript,
			Amount:    output.Value,
		})
		logger.Debug("txOut: pkScript:%x, amount:%v", output.PkScript, output.Value)
	}
	err = builder.AddTxOutScript(txOuts)
	if err != nil {
		logger.Error("multi btc tx add txOut error:%v", err)
		return "", err
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
	logger.Debug("************************************")
	if err != nil {
		logger.Error("multi tx sign error:%v", err)
		return "", err
	}
	txBytes, err := builder.Build()
	if err != nil {
		logger.Error("build btc tx error:%v", err)
		return "", err
	}
	logger.Info("redeem btc tx hash: %v", builder.TxHash())
	txHash, err := e.btcClient.Sendrawtransaction(hex.EncodeToString(txBytes))
	if err != nil {
		logger.Error("send btc tx error:%v", err)
		return "", err
	}
	logger.Info("send redeem btc tx: %v", txHash)
	return txHash, nil
}

func (e *EthereumAgent) CheckRedeemTx(log types.Log) (*EthereumTx, bool, error) {
	redeemTx := &EthereumTx{}
	//todo more check
	if len(log.Data) <= 64 {
		return redeemTx, false, nil
	}
	dataLength := log.Data[32:64]
	l, err := strconv.ParseInt(fmt.Sprintf("%x", dataLength), 16, 32)
	if err != nil {
		logger.Error("parse data length error:%v", err)
		return redeemTx, false, err
	}
	txData := log.Data[64 : 64+l]
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
	redeemTx.Inputs = inputs
	redeemTx.Outputs = outputs
	redeemTx.TxId = log.TxHash.String()
	return redeemTx, true, nil
}

func (e *EthereumAgent) updateProof(resp ProofResponse) error {
	err := UpdateProof(e.store, resp.TxId, resp.Proof, Deposit, ProofStatus(resp.Status))
	if err != nil {
		logger.Error("update proof error: %v %v", resp.TxId, err)
		return err
	}
	return nil
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
