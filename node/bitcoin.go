package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
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
	proofRequest         chan []ProofRequest
	checkProofHeightNums int64           // restart check proof height
	whiteList            map[string]bool // gen tx proof whitelist
	operateAddr          string
	submitTxEthAddr      string
	ethPrivateKey        string // todo keyStore
	minDepositValue      float64
	initStartHeight      int64
	exitSign             chan struct{}
}

func NewBitcoinAgent(cfg NodeConfig, store, memoryStore store.IStore, btcClient *bitcoin.Client, ethClient *ethereum.Client,
	request chan []ProofRequest, response <-chan ProofResponse) (IAgent, error) {
	submitTxEthAddr, err := privateKeyToEthAddr(cfg.EthPrivateKey)
	if err != nil {
		return nil, err
	}
	return &BitcoinAgent{
		btcClient:            btcClient,
		ethClient:            ethClient,
		store:                store,
		memoryStore:          memoryStore,
		blockTime:            time.Duration(cfg.BTcBtcBlockTime) * time.Second,
		operateAddr:          cfg.BtcOperatorAddr,
		proofRequest:         request,
		proofResponse:        response,
		checkProofHeightNums: 100,
		minDepositValue:      0,
		ethPrivateKey:        cfg.EthPrivateKey,
		submitTxEthAddr:      submitTxEthAddr,
		exitSign:             make(chan struct{}, 1),
		initStartHeight:      cfg.BtcInitHeight,
	}, nil
}

func (b *BitcoinAgent) Init() error {

	has, err := b.store.Has(btcCurHeightKey)
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}
	if has {
		err := b.checkUnCompleteGenerateProofTx()
		if err != nil {
			logger.Error("check uncomplete generate proof tx error:%v", err)
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
	// test btc rpc
	_, err = b.btcClient.GetBlockCount()
	if err != nil {
		logger.Error("get block count error:%v", err)
		return err
	}
	logger.Info("init bitcoin agent complete")
	//todo
	return nil
}

// checkUnCompleteGenerateProofTx check uncompleted generate proof tx,resend again
func (b *BitcoinAgent) checkUnCompleteGenerateProofTx() error {
	//todo
	return nil
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
			proofList = append(proofList, ProofRequest{
				TxId:  proof.TxId,
				PType: Deposit,
				Msg:   proof.Msg,
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
		return 0, err
	}
	return height, nil

}

func (b *BitcoinAgent) ScanBlock() error {
	//logger.Info("start bitcoin scan block")
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
		logger.Error("get block count error:%v", err)
		return err
	}
	blockCount = blockCount - 0

	//todo
	if curHeight >= blockCount {
		logger.Debug("btc urrent height:%d,node block count:%d", curHeight, blockCount)
		return nil
	}
	for index := curHeight + 1; index <= blockCount; index++ {
		logger.Debug("bitcoin parse block height:%d", index)
		depositTxList, proofRequestList, err := b.parseBlock(index)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		err = b.saveDataToDb(index, depositTxList)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		b.proofRequest <- proofRequestList
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
			logger.Error(err.Error())
			return err
		}
		pTxId := TxIdToProofId(depositTx.TxId)
		err = b.store.BatchPutObj(pTxId, TxProof{
			PTxId:  pTxId,
			TxId:   depositTx.TxId,
			ToAddr: depositTx.EthAddr,
			Amount: depositTx.Amount,
			Status: ProofDefault,
		})
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	err := b.store.BatchPutObj(height, txIdList)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	err = b.store.BatchPutObj(btcCurHeightKey, height)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	err = b.store.BatchWriteObj()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (b *BitcoinAgent) parseBlock(height int64) ([]DepositTx, []ProofRequest, error) {
	blockHash, err := b.btcClient.GetBlockHash(height)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, err
	}
	blockWithTx, err := b.btcClient.GetBlock(blockHash)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, err
	}
	var proofRequestList []ProofRequest
	var depositTxList []DepositTx
	for _, tx := range blockWithTx.Tx {
		depositTx, check, err := b.checkTx(tx.Vout)
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, err
		}
		if check {
			depositTx.TxId = tx.Txid
			depositTxList = append(depositTxList, depositTx)
			request := ProofRequest{
				TxId:    depositTx.TxId,
				Vout:    depositTx.TxIndex,
				EthAddr: depositTx.EthAddr,
				Amount:  depositTx.Amount,
				PType:   Deposit,
			}
			proofRequestList = append(proofRequestList, request)
			logger.Info("found zkbtc deposit tx: %v", request.String())
		}
	}
	return depositTxList, proofRequestList, nil
}

func (b *BitcoinAgent) Transfer() {
	//todo whether need queue ?
	logger.Debug("start bitcoin transfer goroutine")
	for {
		select {
		case <-b.exitSign:
			logger.Info("exit bitcoin transfer goroutine")
			return
		case response := <-b.proofResponse:
			logger.Info("bitcoinAgent receive deposit proof response: %v", response.String())
			err := b.updateProof(response)
			if err != nil {
				logger.Error("update proof error:%v", err)
				continue
			}
			err = b.MintZKBtcTx(response)
			if err != nil {
				//todo retry ?
				logger.Error("mint btc tx error:%v", err)
				continue
			}
			//logger.Info("success mint btc tx:%v", response)
		}

	}
}

func (b *BitcoinAgent) MintZKBtcTx(resp ProofResponse) error {
	//todo
	nonce, err := b.ethClient.GetNonce(b.submitTxEthAddr)
	if err != nil {
		logger.Error("get nonce error:%v", err)
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
	amountBig, err := Str2Big(resp.Amount, 8)
	if err != nil {
		return fmt.Errorf("parse big error amount:%v", resp.Amount)
	}
	proofBytes := []byte(resp.Proof)
	txHash, err := b.ethClient.Deposit(b.ethPrivateKey, resp.TxId, uint32(resp.Vout), nonce, gasLimit, chainId, gasPrice,
		amountBig, proofBytes)
	if err != nil {
		logger.Error("mint btc tx error:%v", err)
		return err
	}
	logger.Info("success send mint zkbtctx hash:%v, amount: %v", txHash, amountBig.String())
	return nil
}

func (b *BitcoinAgent) checkTx(txOuts []types.TxVout) (DepositTx, bool, error) {
	// todo   check rule
	depositTx := DepositTx{}
	if len(txOuts) < 2 {
		return depositTx, false, nil
	}
	if txOuts[1].ScriptPubKey.Address != b.operateAddr {
		return depositTx, false, nil
	}
	if txOuts[1].Value <= b.minDepositValue {
		return depositTx, false, nil
	}
	if !(txOuts[0].ScriptPubKey.Type == "nulldata" && strings.HasPrefix(txOuts[0].ScriptPubKey.Hex, "6a")) {
		return depositTx, false, nil
	}
	ethAddr, err := getEthAddrFromScript(txOuts[0].ScriptPubKey.Hex)
	if err != nil {
		logger.Error("get eth addr from script error:%v", err)
		return depositTx, false, err
	}
	depositTx.EthAddr = ethAddr
	depositTx.TxIndex = 1
	depositTx.Amount = fmt.Sprintf("%0.8f", txOuts[1].Value) //todo
	return depositTx, true, nil
}

func (b *BitcoinAgent) updateProof(resp ProofResponse) error {
	pTxId := TxIdToProofId(resp.TxId)
	err := b.store.PutObj(pTxId, TxProof{
		PTxId:  pTxId,
		TxId:   resp.TxId,
		Msg:    resp.Msg,
		Status: ProofSuccess,
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
