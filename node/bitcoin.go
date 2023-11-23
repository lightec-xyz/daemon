package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"github.com/lightec-xyz/daemon/store"
	"time"
)

type BitcoinAgent struct {
	client      *bitcoin.Client
	store       *store.Store
	blockTime   time.Duration
	exitSignal  chan struct{}
	operateAddr string
}

func NewBitcoinAgent(cfg BtcConfig, store *store.Store) (IAgent, error) {
	btcClient, err := bitcoin.NewClient(cfg.Url, cfg.User, cfg.Pwd, cfg.Network)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return &BitcoinAgent{
		client:      btcClient,
		store:       store,
		blockTime:   time.Duration(cfg.BlockTime) * time.Second,
		exitSignal:  make(chan struct{}, 1),
		operateAddr: cfg.OperatorAddr,
	}, nil
}

func (b *BitcoinAgent) Init() error {
	return nil
}

func (b *BitcoinAgent) Run() error {
	ticker := time.NewTicker(b.blockTime)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := b.ScanBlock()
			if err != nil {
				logger.Error(err.Error())
			}
		case <-b.exitSignal:
			logger.Info("%v node exit now", b.Name())
			return nil
		}
	}

}

func (b *BitcoinAgent) ScanBlock() error {
	var curHeight int64
	err := getCurrentHeight(b.store, BtcCurHeight, &curHeight)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	blockCount, err := b.client.GetBlockCount()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	//todo
	if curHeight >= blockCount-6 {
		return nil
	}
	for index := curHeight + 1; index <= blockCount; index++ {
		logger.Info("decode block %d", index)
		depositTxList, err := b.parseBlock(index)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		err = b.SaveTxData(index, depositTxList)
		if err != nil {
			logger.Error(err.Error())
			return err
		}

	}
	return nil

}

func (b *BitcoinAgent) SaveTxData(height int64, depositTxList []DepositTx) error {
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

func (b *BitcoinAgent) parseBlock(height int64) ([]DepositTx, error) {
	blockHash, err := b.client.GetBlockHash(height)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	blockWithTx, err := b.client.GetBlockWithTx(blockHash)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	var depositTxList []DepositTx
	for _, tx := range blockWithTx.Tx {
		depositTx, check, err := b.CheckDeposit(tx.Vout)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		if !check {
			logger.Warn("invalid tx: %v", tx.Hash)
			continue
		}
		depositTxList = append(depositTxList, depositTx)
	}

	return depositTxList, nil
}

// todo  check rule

func (b *BitcoinAgent) CheckDeposit(outList []types.TxOut) (DepositTx, bool, error) {
	if len(outList) < 2 {
		return DepositTx{}, false, nil
	}
	//todo
	return DepositTx{}, true, nil
}

func (b *BitcoinAgent) Close() error {
	b.exitSignal <- struct{}{}
	return nil
}
func (b *BitcoinAgent) Name() string {
	return b.Name()
}
