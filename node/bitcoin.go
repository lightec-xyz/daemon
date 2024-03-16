package node

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	"math/big"
	"strings"
	"time"
)

type BitcoinAgent struct {
	btcClient            *bitcoin.Client
	ethClient            *ethereum.Client
	store                store.IStore
	memoryStore          store.IStore
	blockTime            time.Duration
	proofRequest         chan<- []ZkProofRequest
	checkProofHeightNums int64
	whiteList            map[string]bool // todo
	operatorAddr         string
	submitTxEthAddr      string
	keyStore             *KeyStore
	minDepositValue      float64
	initStartHeight      int64
	autoSubmit           bool
	exitSign             chan struct{}
	submitQueue          *Queue // submit proof to eth queue
}

func NewBitcoinAgent(cfg NodeConfig, submitTxEthAddr string, store, memoryStore store.IStore, btcClient *bitcoin.Client, ethClient *ethereum.Client,
	requests chan []ZkProofRequest, keyStore *KeyStore) (IAgent, error) {
	return &BitcoinAgent{
		btcClient:            btcClient,
		ethClient:            ethClient,
		store:                store,
		memoryStore:          memoryStore,
		blockTime:            cfg.BtcScanBlockTime,
		operatorAddr:         cfg.BtcOperatorAddr,
		proofRequest:         requests,
		checkProofHeightNums: 100, // todo
		minDepositValue:      0,   // todo
		keyStore:             keyStore,
		submitTxEthAddr:      submitTxEthAddr,
		exitSign:             make(chan struct{}, 1),
		initStartHeight:      cfg.BtcInitHeight,
		autoSubmit:           cfg.AutoSubmit,
		submitQueue:          NewQueue(),
	}, nil
}

func (b *BitcoinAgent) Init() error {
	logger.Info("bitcoin agent init now")
	exists, err := CheckBitcoinHeight(b.store)
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}
	if exists {
		logger.Debug("bitcoin agent check uncompleted generate Proof tx")
		err := b.checkUnGenerateProof()
		if err != nil {
			logger.Error("check uncompleted generate Proof tx error:%v", err)
			return err
		}
	} else {
		logger.Debug("init btc current height: %v", b.initStartHeight)
		err := WriteBitcoinHeight(b.store, b.initStartHeight)
		if err != nil {
			logger.Error("put init btc current height error:%v", err)
			return err
		}
	}
	// test rpc
	_, err = b.btcClient.GetBlockCount()
	if err != nil {
		logger.Error(" bitcoin json rpc get block count error:%v", err)
		return err
	}
	logger.Info("init bitcoin agent completed")
	return nil
}

// checkUnGenerateProof check uncompleted generate Proof tx,resend again
func (b *BitcoinAgent) checkUnGenerateProof() error {
	// todo
	return nil
}

func (b *BitcoinAgent) getCurrentHeight() (int64, error) {
	return ReadBitcoinHeight(b.store)

}

func (b *BitcoinAgent) ScanBlock() error {
	logger.Debug("bitcoin scan block ...")
	curHeight, err := b.getCurrentHeight()
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}
	if curHeight < b.initStartHeight {
		curHeight = b.initStartHeight
	}
	blockCount, err := b.btcClient.GetBlockCount()
	if err != nil {
		logger.Error("bitcoin client get block count error:%v", err)
		return err
	}
	//todo
	blockCount = blockCount - 0
	if curHeight >= blockCount {
		logger.Debug("btc current height:%d,node block count:%d", curHeight, blockCount)
		return nil
	}
	for index := curHeight + 1; index <= blockCount; index++ {
		logger.Debug("bitcoin parse block height:%d", index)
		depositTxes, redeemTxes, proofRequests, proofs, err := b.parseBlock(index)
		if err != nil {
			logger.Error("bitcoin agent parse block error: %v %v", index, err)
			return err
		}
		err = b.updateRedeemInfo(index, redeemTxes)
		if err != nil {
			logger.Error("bitcoin agent update redeem info error: %v %v", index, err)
			return err
		}
		allTxes := append(depositTxes, redeemTxes...)
		err = b.saveTransaction(index, allTxes)
		if err != nil {
			logger.Error("bitcoin agent save transaction error: %v %v", index, err)
			return err
		}
		err = b.saveDepositData(proofs, proofRequests)
		if err != nil {
			logger.Error("bitcoin agent save data to db error: %v %v", index, err)
			return err
		}
		err = WriteBitcoinHeight(b.store, index)
		if err != nil {
			logger.Error("write btc height error: %v %v", index, err)
			return err
		}
		zkProofRequest, err := toDepositZkProofRequest(proofRequests)
		if err != nil {
			logger.Error("to deposit zk Proof request error: %v %v", index, err)
			return err
		}
		b.proofRequest <- zkProofRequest
		if len(redeemTxes) > 0 {
			// todo
			if b.autoSubmit {
				err = b.updateContractUtxoChange(redeemTxes)
				if err != nil {
					logger.Error("update utxo error: %v %v", index, err)
					return err
				}
			}
		}
	}
	return nil
}

func (b *BitcoinAgent) updateRedeemInfo(height int64, txList []Transaction) error {
	//todo
	return nil
}

func (e *BitcoinAgent) saveTransaction(height int64, txes []Transaction) error {
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

func (b *BitcoinAgent) saveDepositData(proofs []Proof, requests []DepositProofParam) error {
	err := WriteDbProof(b.store, proofsToDbProofs(proofs))
	if err != nil {
		logger.Error("write Proof error: %v", err)
		return err
	}
	err = WriteUnGenProof(b.store, Bitcoin, depositToTxHash(requests))
	if err != nil {
		logger.Error("write ungen Proof error:%v", err)
		return err
	}
	return nil
}

func (b *BitcoinAgent) parseBlock(height int64) ([]Transaction, []Transaction, []DepositProofParam, []Proof, error) {
	blockHash, err := b.btcClient.GetBlockHash(height)
	if err != nil {
		logger.Error("btcClient get block hash error: %v %v", height, err)
		return nil, nil, nil, nil, err
	}
	blockWithTx, err := b.btcClient.GetBlock(blockHash)
	if err != nil {
		logger.Error("btcClient get block error: %v %v", blockHash, err)
		return nil, nil, nil, nil, err
	}
	var depositTxes []Transaction
	var redeemTxes []Transaction
	var requests []DepositProofParam
	var proofs []Proof
	for _, tx := range blockWithTx.Tx {
		redeemTx, isRedeem := b.isRedeemTx(tx)
		if isRedeem {
			redeemTxes = append(redeemTxes, redeemTx)
			continue
		}
		depositTx, isDeposit, err := b.isDepositTx(tx)
		if err != nil {
			logger.Error("check deposit tx error: %v %v", tx.Txid, err)
			return nil, nil, nil, nil, err
		}
		if isDeposit {
			submitted, err := b.ethClient.CheckDepositProof(depositTx.TxHash)
			if err != nil {
				logger.Error("check deposit Proof error: %v %v", tx.Txid, err)
				return nil, nil, nil, nil, err
			}
			var depositTxProof Proof
			if submitted {
				depositTxProof = NewDepositTxProof(tx.Txid, ProofSuccess)
			} else {
				requests = append(requests, NewDepositProofRequest(depositTx.TxHash))
				depositTxProof = NewDepositTxProof(tx.Txid, ProofDefault)
			}
			proofs = append(proofs, depositTxProof)
			depositTxes = append(depositTxes, depositTx)
		}
	}
	return depositTxes, redeemTxes, requests, proofs, nil
}

func (b *BitcoinAgent) ProofResponse(response ZkProofResponse) error {

	//logger.Info("bitcoinAgent receive deposit Proof resp: %v", resp.FilterLogs())
	//err = b.updateDepositProof(resp.TxHash, resp.Proof, resp.Status)
	//if err != nil {
	//	logger.Error("update Proof error: %v %v", resp.TxHash, err)
	//	return err
	//}
	//exists, err := CheckDepositDestHash(b.store, b.ethClient, resp.TxHash)
	//if err != nil {
	//	logger.Error("check btc tx error: %v %v", resp.TxHash, err)
	//	return err
	//}
	//if exists {
	//	err := DeleteUnGenProof(b.store, Bitcoin, resp.TxHash)
	//	if err != nil {
	//		logger.Error("delete ungen Proof error: %v %v", resp.TxHash, err)
	//	}
	//	logger.Warn("deposit btc tx submitted: %v", resp.TxHash)
	//	return nil
	//}
	//if resp.Status == ProofSuccess {
	//	err = DeleteUnGenProof(b.store, Bitcoin, resp.TxHash)
	//	if err != nil {
	//		logger.Error("delete ungen Proof error: %v %v", resp.TxHash, err)
	//	}
	//	logger.Warn("deposit btc tx submitted: %v", resp.TxHash)
	//	if b.autoSubmit {
	//		txHash, err := b.MintZKBtcTx(resp.Utxos, resp.Proof, resp.EthAddr, resp.Amount)
	//		if err != nil {
	//			//todo add queue or cli retry ?
	//			logger.Error("mint btc tx error:%v", err)
	//			return err
	//		}
	//		logger.Info("success mint zkbtc tx: %v", txHash)
	//	}
	//} else {
	//	//todo retry ?
	//	logger.Warn("Proof generate fail status: %v %v", resp.TxHash, resp.Status)
	//}
	return nil
}

func (b *BitcoinAgent) updateContractUtxoChange(utxoList []Transaction) error {
	// todo
	var txIds []string
	for _, tx := range utxoList {
		txId := tx.TxHash
		if !strings.HasPrefix(txId, "0x") {
			txId = fmt.Sprintf("0x%s", txId)
		}
		txIds = append(txIds, txId)
	}
	nonce, err := b.ethClient.GetNonce(b.submitTxEthAddr)
	if err != nil {
		logger.Error("get  nonce error:%v", err)
		return err
	}
	chainId, err := b.ethClient.GetChainId()
	if err != nil {
		logger.Error("get chain id error:%v", err)
		return err
	}
	gasPrice, err := b.ethClient.GetGasPrice()
	if err != nil {
		logger.Error("get gas price error:%v", err)
		return err
	}
	gasLimit := uint64(500000)
	proofBytes := []byte("test ok")
	txHash, err := b.ethClient.UpdateUtxoChange(b.keyStore.GetPrivateKey(), txIds, nonce, gasLimit, chainId, gasPrice, proofBytes)
	if err != nil {
		logger.Error("update utxo change error:%v", err)
		return err
	}
	logger.Info("success send update utxo change  hash:%v", txHash)
	return nil
}

func (b *BitcoinAgent) MintZKBtcTx(utxo []Utxo, proof, receiverAddr string, amount int64) (string, error) {
	//todo need assign nonce ？
	nonce, err := b.ethClient.GetNonce(b.submitTxEthAddr)
	if err != nil {
		logger.Error("get nonce error:%v", err)
		return "", err
	}
	chainId, err := b.ethClient.GetChainId()
	if err != nil {
		logger.Error("get chain id error:%v %v", b.submitTxEthAddr, err)
		return "", err
	}
	gasPrice, err := b.ethClient.GetGasPrice()
	if err != nil {
		logger.Error("get gas price error:%v", err)
		return "", err
	}
	//todo
	if len(utxo) == 0 {
		logger.Error("no utxo")
		return "", fmt.Errorf("no utxo")
	}
	txId := utxo[0].TxId
	index := utxo[0].Index
	gasLimit := uint64(500000)
	amountBig := big.NewInt(amount)
	proofBytes := []byte(proof)
	txHash, err := b.ethClient.Deposit(b.keyStore.GetPrivateKey(), txId, receiverAddr, index, nonce, gasLimit, chainId, gasPrice,
		amountBig, proofBytes)
	if err != nil {
		logger.Error("mint btc tx error:%v", err)
		return "", err
	}
	logger.Info("success send mint zkbtctx hash:%v, amount: %v", txHash, amountBig.String())
	return txHash, nil
}

func (b *BitcoinAgent) isRedeemTx(tx types.Tx) (Transaction, bool) {
	// todo more check
	var inputs []Utxo
	isRedeemTx := false
	for _, vin := range tx.Vin {
		if vin.Prevout.ScriptPubKey.Address == b.operatorAddr {
			isRedeemTx = true
		}
		inputs = append(inputs, Utxo{
			TxId:  vin.TxId,
			Index: uint32(vin.Vout),
		})
	}
	var outputs []TxOut
	for _, out := range tx.Vout {
		scriptHex, err := hex.DecodeString(out.ScriptPubKey.Hex)
		if err != nil {
			logger.Error("decode hex error:%v %v", tx.Txid, err)
			return Transaction{}, false
		}
		outputs = append(outputs, TxOut{
			Value:    BtcToSat(out.Value),
			PkScript: scriptHex,
		})
	}
	if isRedeemTx {
		logger.Info("bitcoin agent find redeem tx: %v,inputs:%v ,outputs:%v", tx.Txid, formatUtxo(inputs), formatOut(outputs))
	}
	redeemBtcTx := NewRedeemBtcTx(tx.Txid, inputs, outputs)
	return redeemBtcTx, isRedeemTx
}

func (b *BitcoinAgent) isDepositTx(tx types.Tx) (Transaction, bool, error) {
	// todo more rule
	txOuts := tx.Vout
	if len(txOuts) < 2 {
		return Transaction{}, false, nil
	}
	if txOuts[1].ScriptPubKey.Address != b.operatorAddr {
		return Transaction{}, false, nil
	}
	amount := txOuts[1].Value
	if amount <= b.minDepositValue {
		logger.Warn("deposit tx less than min value: %v %v", b.minDepositValue, tx.Txid)
		return Transaction{}, false, nil
	}
	if !(txOuts[0].ScriptPubKey.Type == "nulldata" && strings.HasPrefix(txOuts[0].ScriptPubKey.Hex, "6a")) {
		logger.Warn("find deposit tx but check rule fail: %v", tx.Txid)
		return Transaction{}, false, nil
	}
	ethAddr, err := getEthAddrFromScript(txOuts[0].ScriptPubKey.Hex)
	if err != nil {
		logger.Error("get eth addr from script error:%v %v", txOuts[0].ScriptPubKey.Hex, err)
		return Transaction{}, false, err
	}

	utxoList := []Utxo{
		{
			TxId:  tx.Txid,
			Index: 1,
		},
	}
	logger.Info("bitcoin agent find  deposit tx: %v, ethAddr:%v,amount:%v,utxo:%v", tx.Txid, ethAddr, amount, formatUtxo(utxoList))
	depositTx := NewDepositBtcTx(tx.Txid, ethAddr, utxoList, BtcToSat(amount))
	return depositTx, true, nil
}

func (b *BitcoinAgent) updateDepositProof(txId, proof string, status ProofStatus) error {
	logger.Debug("update DepositTx  Proof status: %v %v %v", txId, proof, status)
	err := UpdateProof(b.store, txId, proof, DepositTxType, status)
	if err != nil {
		logger.Error("update Proof error: %v %v", txId, err)
		return err
	}
	return nil

}

func (b *BitcoinAgent) Close() error {
	return nil
}
func (b *BitcoinAgent) Name() string {
	return "Bitcoin WrapperAgent"
}

func CheckDepositDestHash(store store.IStore, ethClient *ethereum.Client, txId string) (bool, error) {
	exists, err := CheckDestHash(store, txId)
	if err != nil {
		logger.Error("check dest hash error:%v", err)
		return false, err
	}
	if exists {
		return true, nil
	}
	submitted, err := ethClient.CheckDepositProof(txId)
	if err != nil {
		logger.Error("check deposit Proof error:%v", err)
		return false, err
	}
	return submitted, nil
}

func getEthAddrFromScript(script string) (string, error) {
	// todo
	// example https://live.blockcypher.com/btc-testnet/tx/fa1bee4165f1720b33047792e47743aeb406940f4b2527874929db9cdbb9da42/
	if len(script) < 5 {
		return "", fmt.Errorf("scritp lenght is less than 4")
	}
	if !strings.HasPrefix(script, "6a") {
		return "", fmt.Errorf("script is not start with 6a")
	}
	isHexAddress := common.IsHexAddress(script[4:])
	if !isHexAddress {
		return "", fmt.Errorf("script is not hex address")
	}
	return script[4:], nil
}

func NewDepositProofRequest(txId string) DepositProofParam {
	return DepositProofParam{
		TxHash: txId,
	}
}

func NewDepositTxProof(txId string, status ProofStatus) Proof {
	return Proof{
		TxHash:    txId,
		ProofType: DepositTxType,
		Status:    int(status),
	}
}

func NewDepositBtcTx(txId, ethAddr string, utxo []Utxo, amount int64) Transaction {
	return Transaction{
		TxHash:    txId,
		TxType:    DepositTx,
		ChainType: Bitcoin,
	}
}

func NewRedeemBtcTx(txId string, inputs []Utxo, outputs []TxOut) Transaction {
	return Transaction{
		TxHash:    txId,
		TxType:    RedeemTx,
		ChainType: Bitcoin,
	}
}
