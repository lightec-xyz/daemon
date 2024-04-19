package node

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
)

type BitcoinAgent struct {
	btcClient            *bitcoin.Client
	ethClient            *ethereum.Client
	store                store.IStore
	memoryStore          store.IStore
	blockTime            time.Duration
	proofRequest         chan<- []*common.ZkProofRequest
	checkProofHeightNums int64
	taskManager          *TaskManager
	whiteList            map[string]bool // todo
	operatorAddr         string
	submitTxEthAddr      string
	keyStore             *KeyStore
	minDepositValue      float64
	initStartHeight      int64
	autoSubmit           bool
	exitSign             chan struct{}
	task                 *TaskManager
}

func NewBitcoinAgent(cfg NodeConfig, submitTxEthAddr string, store, memoryStore store.IStore, btcClient *bitcoin.Client,
	ethClient *ethereum.Client, requests chan []*common.ZkProofRequest, keyStore *KeyStore, task *TaskManager) (IAgent, error) {
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
		task:                 task,
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
			logger.Error("bitcoin agent save Data to db error: %v %v", index, err)
			return err
		}
		err = WriteBitcoinHeight(b.store, index)
		if err != nil {
			logger.Error("write btc height error: %v %v", index, err)
			return err
		}
		// todo
		updateZkProofRequests, err := toUpdateZkProofRequest(redeemTxes)
		if err != nil {
			logger.Error("to update zk Proof request error: %v %v", index, err)
			return err
		}
		zkProofRequest, err := toDepositZkProofRequest(proofRequests)
		if err != nil {
			logger.Error("to deposit zk Proof request error: %v %v", index, err)
			return err
		}
		b.proofRequest <- append(updateZkProofRequests, zkProofRequest...)
		if len(zkProofRequest) > 0 {
			logger.Info("success send btc deposit proof request: %v", len(zkProofRequest))
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
		redeemTx, isRedeem := b.isRedeemTx(tx, blockHash)
		if isRedeem {
			redeemTxes = append(redeemTxes, redeemTx)
			continue
		}
		depositTx, isDeposit, err := parseDepositTx(tx, b.operatorAddr, b.minDepositValue)
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
				depositTxProof = NewDepositTxProof(tx.Txid, common.ProofSuccess)
			} else {
				requests = append(requests, NewDepositProofParam(depositTx.TxHash, blockHash))
				depositTxProof = NewDepositTxProof(tx.Txid, common.ProofDefault)
			}
			proofs = append(proofs, depositTxProof)
			depositTxes = append(depositTxes, depositTx)
		}
	}
	return depositTxes, redeemTxes, requests, proofs, nil
}

func (b *BitcoinAgent) ProofResponse(resp *common.ZkProofResponse) error {
	logger.Info("bitcoinAgent receive  Proof resp: %v %v %v %x",
		resp.ZkProofType.String(), resp.Period, resp.TxHash, resp.Proof)
	proofId := resp.TxHash
	hexProof := hex.EncodeToString(resp.Proof)
	err := b.updateDepositProof(proofId, hexProof, resp.Status)
	if err != nil {
		logger.Error("update Proof error: %v %v", proofId, err)
		return err
	}
	switch resp.ZkProofType {
	case common.DepositTxType:
	case common.VerifyTxType:
		logger.Info("start update utxo change: %v", proofId)
		err := updateContractUtxoChange(b.ethClient, b.submitTxEthAddr, b.keyStore.GetPrivateKey(), []string{resp.TxHash}, resp.Proof)
		if err != nil {
			logger.Error("update utxo error: %v %v", proofId, err)
			b.task.AddTask(resp)
			return err
		}
	default:
	}
	//// todo
	//if b.autoSubmit {
	//	txHash, err := b.taskManager.MintZkBtcRequest(proofId, resp.Proof)
	//	if err != nil {
	//		logger.Error("mint btc tx error:%v", err)
	//		return err
	//	}
	//	logger.Info("success mint zkbtc tx: %v", txHash)
	//}
	return nil
}

func updateContractUtxoChange(ethClient *ethereum.Client, address, privateKey string, txIds []string, proof []byte) error {
	// todo
	nonce, err := ethClient.GetNonce(address)
	if err != nil {
		logger.Error("get  nonce error:%v", err)
		return err
	}
	chainId, err := ethClient.GetChainId()
	if err != nil {
		logger.Error("get chain id error:%v", err)
		return err
	}
	gasPrice, err := ethClient.GetGasPrice()
	if err != nil {
		logger.Error("get gas price error:%v", err)
		return err
	}
	gasLimit := uint64(500000)
	gasPrice = big.NewInt(0).Mul(gasPrice, big.NewInt(2))
	txHash, err := ethClient.UpdateUtxoChange(privateKey, txIds, nonce, gasLimit, chainId, gasPrice, proof)
	if err != nil {
		logger.Error("update utxo change error:%v", err)
		return err
	}
	logger.Info("success send update utxo change  hash:%v", txHash)
	return nil
}

func (b *BitcoinAgent) MintZKBtcTx(utxo []Utxo, proof common.ZkProof, receiverAddr string, amount int64) (string, error) {
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
	txHash, err := b.ethClient.Deposit(b.keyStore.GetPrivateKey(), txId, receiverAddr, index, nonce, gasLimit, chainId, gasPrice, amountBig, proof)
	if err != nil {
		logger.Error("mint btc tx error:%v", err)
		return "", err
	}
	logger.Info("success send mint zkbtctx hash:%v, amount: %v", txHash, amountBig.String())
	return txHash, nil
}

func (b *BitcoinAgent) isRedeemTx(tx types.Tx, blockHash string) (Transaction, bool) {
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
	redeemBtcTx := NewRedeemBtcTx(tx.Txid, blockHash, inputs, outputs)
	return redeemBtcTx, isRedeemTx
}

func (b *BitcoinAgent) updateDepositProof(txId string, proof string, status common.ProofStatus) error {
	logger.Debug("update DepositTx  Proof status: %v %v %v", txId, proof, status)
	err := UpdateProof(b.store, txId, proof, common.DepositTxType, status)
	if err != nil {
		logger.Error("update Proof error: %v %v", txId, err)
		return err
	}
	return nil

}

func (b *BitcoinAgent) CheckState() error {

	//TODO implement me
	return nil
}

func (b *BitcoinAgent) Close() error {
	return nil
}
func (b *BitcoinAgent) Name() string {
	return "Bitcoin Agent"
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

func parseDepositTx(tx types.Tx, operatorAddr string, minDepositValue float64) (Transaction, bool, error) {
	// todo more rule
	txOuts := tx.Vout
	if len(txOuts) < 2 {
		return Transaction{}, false, nil
	}
	amount, isDeposit, err := isContainOperator(tx.Vout, operatorAddr)
	if err != nil {
		return Transaction{}, false, err
	}
	if !isDeposit {
		return Transaction{}, false, nil
	}

	ethAddr, ok, err := getOPReturn(tx.Vout)
	if !ok {
		return Transaction{}, false, nil
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

func isContainOperator(txOuts []types.TxVout, operatorAddr string) (float64, bool, error) {
	var isDeposit bool
	var total float64
	for _, out := range txOuts {
		if out.ScriptPubKey.Address == operatorAddr {
			isDeposit = true
			total = total + out.Value
		}
	}
	return total, isDeposit, nil
}

func getOPReturn(txOuts []types.TxVout) (string, bool, error) {
	for _, out := range txOuts {
		if out.ScriptPubKey.Type == "nulldata" && strings.HasPrefix(out.ScriptPubKey.Hex, "6a") {
			ethAddr, err := getEthAddrFromScript(out.ScriptPubKey.Hex)
			if err != nil {
				return "", false, err
			}
			return ethAddr, true, nil
		}
	}
	return "", false, nil
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
	isHexAddress := ethcommon.IsHexAddress(script[4:])
	if !isHexAddress {
		return "", fmt.Errorf("script is not hex address")
	}
	return script[4:], nil
}

func NewDepositProofParam(txId, blockHash string) DepositProofParam {
	return DepositProofParam{
		TxHash:    txId,
		BlockHash: blockHash,
	}
}

func NewDepositTxProof(txId string, status common.ProofStatus) Proof {
	return Proof{
		TxHash:    txId,
		ProofType: common.DepositTxType,
		Status:    int(status),
	}
}

func NewDepositBtcTx(txId, ethAddr string, utxo []Utxo, amount int64) Transaction {
	return Transaction{
		TxHash:    txId,
		TxType:    DepositTx,
		ChainType: Bitcoin,
		EthAddr:   ethAddr,
		Utxo:      utxo,
		Amount:    amount,
	}
}

func NewRedeemBtcTx(txId, blockHash string, inputs []Utxo, outputs []TxOut) Transaction {
	return Transaction{
		TxHash:    txId,
		TxType:    RedeemTx,
		ChainType: Bitcoin,
		BlockHash: blockHash,
	}
}
