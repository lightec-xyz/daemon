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
	"strings"
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
	addrFilter           EthAddrFilter
	privateKeys          []*btcec.PrivateKey //todo just test
	initStartHeight      int64
	ethSubmitAddress     string
	autoSubmit           bool
}

func NewEthereumAgent(cfg NodeConfig, submitTxEthAddr string, store, memoryStore store.IStore, btcClient *bitcoin.Client, ethClient *ethereum.Client,
	proofRequest chan []ProofRequest, proofResponse <-chan ProofResponse) (IAgent, error) {
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
		privateKeys:          privateKeys,
		initStartHeight:      cfg.EthInitHeight,
		ethSubmitAddress:     submitTxEthAddr,
		autoSubmit:           cfg.AutoSubmit,
		addrFilter:           cfg.EthAddrFilter,
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
		err := e.checkUnGenerateProof()
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
	blockNumber, err := e.ethClient.EthBlockNumber()
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
		err = e.saveDataToDb(index, redeemTxes, proofs)
		if err != nil {
			logger.Error("ethereum save data error: %v %v", index, err)
			return err
		}
		e.proofRequest <- requests
		if len(depositTxes) > 0 {
			err := e.updateDepositDestChainHash(depositTxes)
			if err != nil {
				logger.Error("update deposit final status error: %v %v", index, err)
				return err
			}
		}
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
		case resp := <-e.proofResponse:
			logger.Info("receive redeem proof resp: %v", resp.String())
			err := e.updateRedeemProof(resp.TxId, resp.Proof, resp.Status)
			if err != nil {
				logger.Error("update proof error:%v", err)
				continue
			}
			exists, err := e.btcClient.CheckTx(resp.BtcTxId)
			if err != nil {
				// todo ?
				logger.Error("check btc tx error: %v %v", resp.BtcTxId, err)
				continue
			}
			if exists {
				logger.Warn("redeem btc tx submitted: %v", resp.BtcTxId)
				continue
			}
			if e.autoSubmit {
				if resp.Status == ProofSuccess {
					txHash, err := e.RedeemBtcTx(resp)
					if err != nil {
						// todo add queue or cli retry
						logger.Error("redeem btc tx error:%v", err)
						continue
					}
					logger.Info("success redeem btc tx:%v", txHash)
				} else {
					// todo
					logger.Warn("proof generate failed :%v %v", resp.TxId, resp.Status)
				}
			}
		}
	}

}

func (e *EthereumAgent) saveDataToDb(height int64, redeemTxes []Transaction, proofs []Proof) error {
	err := WriteEthereumTx(e.store, height, redeemTxes)
	if err != nil {
		logger.Error("put redeem tx error: %v %v", height, err)
		return err
	}
	err = WriteProof(e.store, proofs)
	if err != nil {
		logger.Error("put eth current height error:%v %v", height, err)
		return err
	}

	err = WriteRedeemDestChainHash(e.store, redeemTxes)
	if err != nil {
		logger.Error("batch write error: %v %v", height, err)
		return err
	}

	err = WriteEthereumHeight(e.store, height)
	if err != nil {
		logger.Error("batch write error: %v %v", height, err)
		return err
	}
	return nil
}

func (e *EthereumAgent) updateDepositDestChainHash(depositTx []Transaction) error {
	err := WriteDepositDestChainHash(e.store, depositTx)
	if err != nil {
		logger.Error("update deposit final status error: %v %v", depositTx, err)
		return err
	}
	return nil
}

func (e *EthereumAgent) parseBlock(height int64) ([]Transaction, []Transaction, []ProofRequest, []Proof, error) {
	block, err := e.ethClient.GetBlock(height)
	if err != nil {
		logger.Error("ethereum rpc get block error:%v", err)
		return nil, nil, nil, nil, err
	}
	blockHash := block.Hash().String()
	logAddrs := []string{e.addrFilter.LogDepositAddr, e.addrFilter.LogRedeemAddr}
	logTopics := []string{e.addrFilter.LogTopicDepositAddr, e.addrFilter.LogTopicRedeemAddr}
	logs, err := e.ethClient.GetLogs(blockHash, logAddrs, logTopics)
	if err != nil {
		logger.Error("ethereum rpc get logs error:%v", err)
		return nil, nil, nil, nil, err
	}
	var redeemTxes []Transaction
	var depositTxes []Transaction
	var proofs []Proof
	var requests []ProofRequest
	for _, log := range logs {
		depositTx, isDeposit, err := e.isDepositTx(log)
		if err != nil {
			logger.Error("check is deposit tx error:%v", err)
			return nil, nil, nil, nil, err
		}
		if isDeposit {
			logger.Info("ethereum agent find deposit zkbtc ethTxHash:%v,btcTxId:%v,utxo:%v",
				depositTx.TxHash, depositTx.BtcTxId, formatUtxo(depositTx.Utxo))
			depositTxes = append(depositTxes, depositTx)
			continue
		}

		redeemTx, isRedeem, err := e.isRedeemTx(log)
		if err != nil {
			logger.Error("check is redeem tx error:%v", err)
			return nil, nil, nil, nil, err
		}
		if isRedeem {
			logger.Info("ethereum agent find redeem zkbtc  ethTxHash:%v,btcTxId:%v,input:%v,output:%v",
				redeemTx.TxHash, redeemTx.BtcTxId, formatUtxo(redeemTx.Inputs), formatOut(redeemTx.Outputs))
			proofs = append(proofs, NewRedeemProof(redeemTx.TxHash))
			requests = append(requests, NewRedeemProofRequest(redeemTx.TxHash, redeemTx.BtcTxId, redeemTx.Inputs, redeemTx.Outputs))
			redeemTxes = append(redeemTxes, redeemTx)
			continue
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
	if strings.ToLower(log.Address.Hex()) == e.addrFilter.LogDepositAddr && strings.ToLower(log.Topics[0].Hex()) == e.addrFilter.LogTopicDepositAddr {
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
		depositTx := NewDepositEthTx(log.TxHash.String(), btcTxId, utxo, amount)
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
	if strings.ToLower(log.Address.Hex()) == e.addrFilter.LogRedeemAddr && strings.ToLower(log.Topics[0].Hex()) == e.addrFilter.LogTopicRedeemAddr {
		btcTxId := strings.ToLower(log.Topics[1].Hex())
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
		if transaction.TxHash().String() != btcTxId {
			logger.Error("never should happen btc tx not match error: %v %v", log.TxHash.String())
			return redeemTx, false, fmt.Errorf("tx hash not match:%v", log.TxHash.String())
		}
		redeemTx = NewRedeemEthTx(log.TxHash.String(), btcTxId, inputs, outputs)
		return redeemTx, true, nil
	} else {
		return redeemTx, false, nil
	}

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

func (e *EthereumAgent) updateRedeemProof(txId, proof string, status ProofStatus) error {
	logger.Debug("update Redeem proof status: %v %v %v", txId, proof, status)
	err := UpdateProof(e.store, txId, proof, Redeem, status)
	if err != nil {
		logger.Error("update proof error: %v %v", txId, err)
		return err
	}
	return nil
}

func (e *EthereumAgent) Close() error {
	close(e.exitSign)
	return nil
}
func (e *EthereumAgent) Name() string {
	return "Ethereum Agent"
}
func (e *EthereumAgent) BlockTime() time.Duration {
	return e.blockTime
}

func NewRedeemProofRequest(txId, btcTxId string, inputs []Utxo, outputs []TxOut) ProofRequest {
	return ProofRequest{
		TxId:      txId,
		ProofType: Redeem,
		Inputs:    inputs,
		Outputs:   outputs,
		BtcTxId:   btcTxId,
	}
}

func NewRedeemProof(txId string) Proof {
	return Proof{
		TxId:      txId,
		ProofType: Redeem,
		Proof:     "",
		Status:    ProofDefault,
	}
}

func NewDepositEthTx(txHash, btcTxId string, utxo []Utxo, amount int64) Transaction {
	return Transaction{
		TxHash:    txHash,
		BtcTxId:   btcTxId,
		Utxo:      utxo,
		Amount:    amount,
		ChainType: Ethereum,
		TxType:    DepositTx,
	}
}
func NewRedeemEthTx(txId string, destHash string, inputs []Utxo, outputs []TxOut) Transaction {
	return Transaction{
		TxHash:    txId,
		DestHash:  destHash,
		BtcTxId:   destHash,
		Inputs:    inputs,
		Outputs:   outputs,
		ChainType: Ethereum,
		TxType:    RedeemTx,
	}
}
