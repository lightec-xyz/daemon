package node

import (
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
	proofResponse        <-chan ProofResponse
	proofRequest         chan<- []ProofRequest
	nonceManager         *NonceManager
	checkProofHeightNums int64
	whiteList            map[string]bool // todo
	operatorAddr         string
	submitTxEthAddr      string
	keyStore             *KeyStore
	minDepositValue      float64
	initStartHeight      int64
	autoSubmit           bool
	exitSign             chan struct{}
}

func NewBitcoinAgent(cfg NodeConfig, store, memoryStore store.IStore, btcClient *bitcoin.Client, ethClient *ethereum.Client,
	request chan []ProofRequest, response <-chan ProofResponse, nonceManager *NonceManager, keyStore *KeyStore) (IAgent, error) {
	submitTxEthAddr, err := privateKeyToEthAddr(cfg.EthPrivateKey)
	if err != nil {
		logger.Error("privateKeyToEthAddr error:%v", err)
		return nil, err
	}
	return &BitcoinAgent{
		btcClient:            btcClient,
		ethClient:            ethClient,
		store:                store,
		memoryStore:          memoryStore,
		blockTime:            cfg.BtcScanBlockTime,
		operatorAddr:         cfg.BtcOperatorAddr,
		proofRequest:         request,
		proofResponse:        response,
		checkProofHeightNums: 100, // todo
		minDepositValue:      0,   // todo
		nonceManager:         nonceManager,
		keyStore:             keyStore,
		submitTxEthAddr:      submitTxEthAddr,
		exitSign:             make(chan struct{}, 1),
		initStartHeight:      cfg.BtcInitHeight,
		autoSubmit:           cfg.AutoSubmit,
	}, nil
}

func (b *BitcoinAgent) Init() error {
	logger.Info("bitcoin agent init now")
	has, err := b.store.Has(btcCurHeightKey)
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}
	if has {
		logger.Debug("bitcoin agent check uncompleted generate proof tx")
		err := b.checkUnCompleteGenerateProofTx()
		if err != nil {
			logger.Error("check uncompleted generate proof tx error:%v", err)
			return err
		}
	} else {
		logger.Debug("init btc current height: %v", b.initStartHeight)
		err := b.store.PutObj(btcCurHeightKey, b.initStartHeight)
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

// checkUnCompleteGenerateProofTx check uncompleted generate proof tx,resend again
func (b *BitcoinAgent) checkUnCompleteGenerateProofTx() error {
	currentHeight, err := b.getCurrentHeight()
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}
	start := currentHeight - b.checkProofHeightNums
	var proofList []ProofRequest
	for index := start; index < currentHeight; index++ {
		var txIdList []string
		hasObj, err := b.store.HasObj(index)
		if err != nil {
			logger.Error("get txIdList error:%v", err)
			return err
		}
		if !hasObj {
			continue
		}
		err = b.store.GetObj(index, &txIdList)
		if err != nil {
			logger.Error("get txIdList error:%v", err)
			return err
		}
		for _, txId := range txIdList {
			var proof TxProof
			err := b.store.GetObj(TxIdToProofId(txId), &proof)
			if err != nil {
				logger.Error("get proof error:%v", err)
				return err
			}
			//todo
			proofList = append(proofList, ProofRequest{
				TxId:      proof.TxId,
				ProofType: Deposit,
				Msg:       proof.Msg,
			})
		}
	}
	b.proofRequest <- proofList
	return nil
}

func (b *BitcoinAgent) getCurrentHeight() (int64, error) {
	var height int64
	err := b.store.GetObj(btcCurHeightKey, &height)
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return 0, err
	}
	return height, nil

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
		depositTxList, proofRequestList, updateUtxoList, err := b.parseBlock(index)
		if err != nil {
			logger.Error("bitcoin agent parse block error: %v %v", index, err)
			return err
		}
		err = b.saveDataToDb(index, depositTxList)
		if err != nil {
			logger.Error("bitcoin agent save data to db error: %v %v", index, err)
			return err
		}
		b.proofRequest <- proofRequestList
		// todo per block to update ?
		if len(updateUtxoList) > 0 {
			err := b.updateUtxoChange(updateUtxoList)
			if err != nil {
				logger.Error("update utxo error: %v %v", index, err)
				return err
			}
		}
	}
	return nil
}

func (b *BitcoinAgent) saveDataToDb(height int64, depositTxList []DepositTx) error {
	//todo
	var txIdList []string
	for _, depositTx := range depositTxList {
		txIdList = append(txIdList, depositTx.TxId)
		err := b.store.BatchPutObj(depositTx.TxId, depositTx)
		if err != nil {
			logger.Error("batch put deposit tx error: %v %v", depositTx.TxId, err)
			return err
		}
		pTxId := TxIdToProofId(depositTx.TxId)
		err = b.store.BatchPutObj(pTxId, TxProof{
			Height:    height,
			BlockHash: depositTx.BlockHash,
			TxId:      depositTx.TxId,
			ProofType: Deposit,
			Proof:     "",
			Status:    ProofDefault,
		})
		if err != nil {
			logger.Error("batch put proof tx error: %v %v", depositTx.TxId, err)
			return err
		}
	}
	err := b.store.BatchPutObj(height, txIdList)
	if err != nil {
		logger.Error("batch put txIdList error: %v %v", height, err)
		return err
	}
	err = b.store.BatchPutObj(btcCurHeightKey, height)
	if err != nil {
		logger.Error("batch put btc current height error: %v %v", height, err)
		return err
	}
	err = b.store.BatchWriteObj()
	if err != nil {
		logger.Error("batch write error: %v %v", height, err)
		return err
	}
	return nil
}

func (b *BitcoinAgent) parseBlock(height int64) ([]DepositTx, []ProofRequest, []string, error) {
	blockHash, err := b.btcClient.GetBlockHash(height)
	if err != nil {
		logger.Error("btcClient get block hash error: %v %v", height, err)
		return nil, nil, nil, err
	}
	blockWithTx, err := b.btcClient.GetBlock(blockHash)
	if err != nil {
		logger.Error("btcClient get block error: %v %v", blockHash, err)
		return nil, nil, nil, err
	}
	var proofRequestList []ProofRequest
	var depositTxList []DepositTx
	var needUpdateUtxo []string
	for _, tx := range blockWithTx.Tx {
		if b.checkRedeemTx(tx) {
			logger.Info("find redeem tx: %v", tx.Txid)
			needUpdateUtxo = append(needUpdateUtxo, tx.Txid)
			continue
		}
		depositTx, check, err := b.checkDepositTx(tx)
		if err != nil {
			logger.Error("check tx error: %v %v", tx.Txid, err)
			return nil, nil, nil, err
		}
		if check {
			depositTx.Height = height
			depositTx.BlockHash = blockHash
			depositTxList = append(depositTxList, depositTx)
			request := ProofRequest{

				Utxos:   depositTx.Utxos,
				EthAddr: depositTx.EthAddr,
				Amount:  depositTx.Amount,

				Height:    height,
				BlockHash: blockHash,
				TxId:      depositTx.TxId,
				ProofType: Deposit,
			}
			proofRequestList = append(proofRequestList, request)
			logger.Info("found zkbtc deposit tx:%v", request.String())
		}
	}
	return depositTxList, proofRequestList, needUpdateUtxo, nil
}

func (b *BitcoinAgent) Transfer() {
	//todo whether need queue ?
	logger.Debug("start bitcoin transfer goroutine")
	for {
		select {
		case <-b.exitSign:
			logger.Info("bitcoin transfer goroutine exit ...")
			return
		case response := <-b.proofResponse:
			logger.Info("bitcoinAgent receive deposit proof response: %v", response.String())
			err := b.updateProof(response)
			if err != nil {
				logger.Error("update proof error: %v %v", response.TxId, err)
				continue
			}
			if b.autoSubmit && response.Status == ProofSuccess {
				err = b.MintZKBtcTx(response)
				if err != nil {
					//todo add queue or cli retry ?
					logger.Error("mint btc tx error:%v", err)
					continue
				}
			}
		}

	}
}

func (b *BitcoinAgent) getRealNonce(address string) (uint64, error) {
	pendingNonce, err := b.ethClient.GetPendingNonce(address)
	if err != nil {
		logger.Error("get pending nonce error: %v %v", address, err)
		return 0, err
	}
	return pendingNonce, nil

}

func (e *BitcoinAgent) updateUtxoChange(utxoList []string) error {
	// todo eth contract array param
	for _, utxo := range utxoList {
		nonce, err := e.ethClient.GetNonce(e.submitTxEthAddr)
		if err != nil {
			logger.Error("get nonce error:%v", err)
			return err
		}
		chainId, err := e.ethClient.GetChainId()
		if err != nil {
			logger.Error("get chain id error:%v", err)
			return err
		}
		gasPrice, err := e.ethClient.GetGasPrice()
		if err != nil {
			logger.Error("get gas price error:%v", err)
			return err
		}
		gasLimit := uint64(500000)
		proofBytes := []byte("test ok")
		txHash, err := e.ethClient.UpdateUtxoChange(e.keyStore.GetPrivateKey(), utxo, nonce, gasLimit, chainId, gasPrice,
			proofBytes)
		if err != nil {
			logger.Error("update utxo change error:%v", err)
			return err
		}
		logger.Info("success send update utxo change  hash:%v", txHash)
		return nil
	}
	return nil
}

func (b *BitcoinAgent) MintZKBtcTx(resp ProofResponse) error {
	exists, err := b.ethClient.CheckDepositProof(resp.TxId)
	if err != nil {
		logger.Error("check deposit proof error: %v %v", resp.TxId, err)
		return err
	}
	if exists {
		logger.Warn("deposit proof exists now: %v", resp.TxId)
		return nil
	}
	//todo
	nonce, err := b.getRealNonce(b.submitTxEthAddr)
	if err != nil {
		logger.Error("get nonce error: %v %v", b.submitTxEthAddr, err)
		return err
	}
	chainId, err := b.ethClient.GetChainId()
	if err != nil {
		logger.Error("get chain id error:%v %v", b.submitTxEthAddr, err)
		return err
	}
	gasPrice, err := b.ethClient.GetGasPrice()
	if err != nil {
		logger.Error("get gas price error:%v", err)
		return err
	}
	gasLimit := uint64(500000)
	amountBig := big.NewInt(resp.Amount)
	//todo
	proofBytes := []byte(resp.Proof)
	index := resp.Utxos[0].Index
	txHash, err := b.ethClient.Deposit(b.keyStore.GetPrivateKey(), resp.TxId, index, nonce, gasLimit, chainId, gasPrice,
		amountBig, proofBytes)
	if err != nil {
		logger.Error("mint btc tx error:%v", err)
		return err
	}
	logger.Info("success send mint zkbtctx hash:%v, amount: %v", txHash, amountBig.String())
	return nil
}

func (b *BitcoinAgent) checkRedeemTx(tx types.Tx) bool {
	for _, vin := range tx.Vin {
		if vin.Prevout.ScriptPubKey.Address == b.operatorAddr {
			return true
		}
	}
	return false

}

func (b *BitcoinAgent) checkDepositTx(tx types.Tx) (DepositTx, bool, error) {
	// todo   check rule
	txOuts := tx.Vout
	depositTx := DepositTx{}
	if len(txOuts) < 2 {
		return depositTx, false, nil
	}
	if txOuts[1].ScriptPubKey.Address != b.operatorAddr {
		return depositTx, false, nil
	}
	if txOuts[1].Value <= b.minDepositValue {
		logger.Warn("deposit tx less than min value: %v %v", b.minDepositValue, tx.Txid)
		return depositTx, false, nil
	}
	if !(txOuts[0].ScriptPubKey.Type == "nulldata" && strings.HasPrefix(txOuts[0].ScriptPubKey.Hex, "6a")) {
		logger.Warn("find deposit tx but check rule fail: %v", tx.Txid)
		return depositTx, false, nil
	}
	ethAddr, err := getEthAddrFromScript(txOuts[0].ScriptPubKey.Hex)
	if err != nil {
		logger.Error("get eth addr from script error:%v %v", txOuts[0].ScriptPubKey.Hex, err)
		return depositTx, false, err
	}
	//todo
	utxoList := []Utxo{
		{
			TxId:  tx.Txid,
			Index: 1,
		},
	}
	depositTx.Utxos = utxoList
	depositTx.TxId = tx.Txid
	depositTx.EthAddr = ethAddr
	depositTx.Amount = BtcToSat(txOuts[1].Value)
	return depositTx, true, nil
}

func (b *BitcoinAgent) updateProof(resp ProofResponse) error {
	pTxId := TxIdToProofId(resp.TxId)
	err := b.store.PutObj(pTxId, TxProof{
		Height:    resp.Height,
		BlockHash: resp.BlockHash,
		TxId:      resp.TxId,
		ProofType: Deposit,
		Proof:     resp.Proof,
		Status:    resp.Status,
	})
	return err
}

func (b *BitcoinAgent) Close() error {
	return nil
}
func (b *BitcoinAgent) Name() string {
	return "Bitcoin Agent"
}

func (b *BitcoinAgent) BlockTime() time.Duration {
	return b.blockTime
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
