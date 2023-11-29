package node

import (
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
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
}

func NewBitcoinAgent(cfg NodeConfig, store, memoryStore store.IStore, btcClient *bitcoin.Client, ethClient *ethereum.Client,
	request chan []ProofRequest, response <-chan ProofResponse) (IAgent, error) {
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
		logger.Debug("init btc current height: %v", InitBitcoinHeight)
		err := b.store.PutObj(btcCurHeightKey, InitBitcoinHeight)
		if err != nil {
			logger.Error("put init btc current height error:%v", err)
			return err
		}
	}

	//todo
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
			proofList = append(proofList, ProofRequest{
				TxId:   proof.TxId,
				PType:  BitcoinChain,
				ToAddr: proof.ToAddr,
				Amount: proof.Amount,
				Msg:    proof.Msg,
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
			ToAddr: depositTx.Addr,
			Amount: depositTx.Amount,
			Msg:    depositTx.Extra,
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
	blockWithTx, err := b.btcClient.GetBlockWithTx(blockHash)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, err
	}
	var proofRequestList []ProofRequest
	var depositTxList []DepositTx
	for _, tx := range blockWithTx.Tx {
		depositTx, check, err := b.parseTx(tx.Vout)
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, err
		}
		if check {
			logger.Info("found deposit tx: %v", depositTx)
			depositTxList = append(depositTxList, depositTx)
			proofRequestList = append(proofRequestList, ProofRequest{
				TxId:   depositTx.TxId,
				PType:  BitcoinChain,
				ToAddr: depositTx.Addr,
				Amount: depositTx.Amount,
			})
		}
	}

	return depositTxList, proofRequestList, nil
}

func (b *BitcoinAgent) Transfer() error {
	//todo whether need queue
	for {
		select {
		case response := <-b.proofResponse:
			err := b.updateProof(response)
			if err != nil {
				logger.Error("update proof error:%v", err)
				continue
			}
			err = b.MintZKBtcTx(response)
			if err != nil {
				//todo
				logger.Error("mint btc tx error:%v", err)
				return err
			}
			logger.Info("success mint btc tx:%v", response)
		}

	}
}

func (b *BitcoinAgent) MintZKBtcTx(resp ProofResponse) error {
	//todo
	panic("implement me")
}

// todo  check rule

func (b *BitcoinAgent) parseTx(outList []types.TxOut) (DepositTx, bool, error) {
	if len(outList) < 2 {
		return DepositTx{}, false, nil
	}
	//todo
	return DepositTx{}, true, nil
}

func (b *BitcoinAgent) updateProof(resp ProofResponse) error {
	pTxId := TxIdToProofId(resp.TxId)
	err := b.store.PutObj(pTxId, TxProof{
		PTxId:  pTxId,
		TxId:   resp.TxId,
		ToAddr: resp.ToAddr,
		Amount: resp.Amount,
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
