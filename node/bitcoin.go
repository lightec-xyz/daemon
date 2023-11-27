package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	"strings"
	"time"
)

type BitcoinAgent struct {
	btcClient   *bitcoin.Client
	ethClient   *ethereum.Client
	store       store.IStore
	memoryStore store.IStore
	proofClient rpc.ProofAPI
	blockTime   time.Duration
	operateAddr string
}

func NewBitcoinAgent(cfg NodeConfig, store, memoryStore store.IStore,
	btcClient *bitcoin.Client, ethClient *ethereum.Client, proofClient rpc.ProofAPI) (IAgent, error) {
	return &BitcoinAgent{
		btcClient:   btcClient,
		ethClient:   ethClient,
		store:       store,
		memoryStore: memoryStore,
		proofClient: proofClient,
		blockTime:   time.Duration(cfg.BTcBtcBlockTime) * time.Second,
		operateAddr: cfg.BtcOperatorAddr,
	}, nil
}

func (b *BitcoinAgent) Init() error {
	//todo
	height, err := b.getCurrentHeight()
	if err != nil && strings.Contains(err.Error(), "not found") {
		err = b.store.PutObj(BtcCurHeight, InitBitcoinHeight)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	_, err = b.btcClient.GetBlockCount()
	if err != nil {
		logger.Error("bitcoin rpc error:%v", err)
		return err
	}
	logger.Debug("bitcoin node init ok,latest height:%d", height)
	return nil
}

func (b *BitcoinAgent) getCurrentHeight() (int64, error) {
	var height int64
	err := b.store.GetObj(BtcCurHeight, &height)
	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}
	return height, nil

}

func (b *BitcoinAgent) Run() error {
	curHeight, err := b.getCurrentHeight()
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}
	blockCount, err := b.btcClient.GetBlockCount()
	if err != nil {
		logger.Error("get block count error:%v", err)
		return err
	}
	//todo
	if curHeight >= blockCount-6 {
		logger.Info("current height:%d,node block count:%d", curHeight, blockCount)
		return nil
	}
	for index := curHeight + 1; index <= blockCount; index++ {
		logger.Info("decode block %d", index)
		depositTxList, err := b.parseBlock(index)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		err = b.persistData(index, depositTxList)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		//todo gen proof

	}
	return nil

}

func (b *BitcoinAgent) persistData(height int64, depositTxList []DepositTx) error {
	//todo
	var txIdList []string
	for _, depositTx := range depositTxList {
		txIdList = append(txIdList, depositTx.TxId)
		err := b.store.BatchPutObj(depositTx.TxId, depositTx)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		pTxId := fmt.Sprintf("%s%s", ProofPrefix, depositTx.TxId)
		err = b.store.BatchPutObj(pTxId, TxProof{
			PTxId: pTxId,
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
	err = b.store.BatchPutObj(BtcCurHeight, height)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	err = b.store.BatchWrite()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

//todo

func (b *BitcoinAgent) SendTxToEth(txList []DepositTx) error {
	for _, tx := range txList {
		logger.Info("send tx to eth: %v", tx)
		//todo
	}
	return nil
}

func (b *BitcoinAgent) parseBlock(height int64) ([]DepositTx, error) {
	blockHash, err := b.btcClient.GetBlockHash(height)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	blockWithTx, err := b.btcClient.GetBlockWithTx(blockHash)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	var depositTxList []DepositTx
	for _, tx := range blockWithTx.Tx {
		depositTx, check, err := b.parseTx(tx.Vout)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		if check {
			logger.Info("found deposit tx: %v", depositTx)
			depositTxList = append(depositTxList, depositTx)
		}
	}

	return depositTxList, nil
}

// todo  check rule

func (b *BitcoinAgent) parseTx(outList []types.TxOut) (DepositTx, bool, error) {
	if len(outList) < 2 {
		return DepositTx{}, false, nil
	}
	//todo
	return DepositTx{}, true, nil
}

func (b *BitcoinAgent) Close() error {

	return nil
}
func (b *BitcoinAgent) Name() string {
	return b.Name()
}

func (b *BitcoinAgent) BlockTime() time.Duration {
	return b.blockTime
}
