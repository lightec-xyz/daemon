package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	"github.com/onrik/ethrpc"
	"time"
)

type EthereumAgent struct {
	btcClient     *bitcoin.Client
	ethClient     *ethereum.Client
	store         store.IStore
	memoryStore   store.IStore
	blockTime     time.Duration
	whiteList     map[string]bool
	proofResponse <-chan ProofResponse
	proofRequest  chan []ProofRequest
}

func NewEthereumAgent(cfg NodeConfig, store, memoryStore store.IStore, btcClient *bitcoin.Client, ethClient *ethereum.Client,
	proofRequest chan []ProofRequest, proofResponse <-chan ProofResponse) (IAgent, error) {
	return &EthereumAgent{
		btcClient:     btcClient,
		ethClient:     ethClient,
		store:         store,
		memoryStore:   memoryStore,
		blockTime:     time.Duration(cfg.EthBlockTime) * time.Second,
		proofRequest:  proofRequest,
		proofResponse: proofResponse,
	}, nil
}

func (e *EthereumAgent) Init() error {
	has, err := e.store.Has(ethCurHeightKey)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return err
	}
	if !has {
		err := e.store.PutObj(ethCurHeightKey, InitEthereumHeight)
		if err != nil {
			logger.Error("put init eth current height error:%v", err)
			return err
		}
	}
	//todo

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
	ethHeight, err := e.getEthHeight()
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return err
	}
	blockNumber, err := e.ethClient.EthBlockNumber()
	if err != nil {
		logger.Error("get eth block number error:%v", err)
		return err
	}
	//todo
	if ethHeight >= int64(blockNumber)-6 {
		logger.Info("current height:%d,node block count:%d", ethHeight, blockNumber)
		return nil
	}
	for index := ethHeight + 1; index <= int64(blockNumber); index++ {
		logger.Info("decode block %d", index)
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

func (e *EthereumAgent) Transfer() error {
	for {
		select {
		case response := <-e.proofResponse:
			err := e.RedeemBtcTx(response)
			if err != nil {
				//todo
				logger.Error("redeem btc tx error:%v", err)
				return err
			}
			logger.Info("success redeem btc tx:%v", response)
		}
	}

}

func (e *EthereumAgent) RedeemBtcTx(resp ProofResponse) error {
	//todo
	panic("implement me")
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
	block, err := e.ethClient.EthGetBlockByNumber(int(height), true)
	if err != nil {
		logger.Error("ethereum rpc get block error:%v", err)
		return nil, nil, err
	}
	var redeemTxList []RedeemTx
	var proofRequestList []ProofRequest
	for _, tx := range block.Transactions {
		redeemTx, ok, err := e.CheckRedeemTx(tx)
		if err != nil {
			logger.Error("check redeem tx error:%v", err)
			return nil, nil, err
		}
		if ok {
			logger.Info("found redeem tx: %v", redeemTx)
			redeemTxList = append(redeemTxList, redeemTx)
			proofRequestList = append(proofRequestList, ProofRequest{
				TxId:   redeemTx.TxId,
				PType:  EthereumChain,
				ToAddr: redeemTx.BtcAddr,
				Amount: redeemTx.Amount,
			})
		}
	}
	return redeemTxList, proofRequestList, nil
}

func (e *EthereumAgent) CheckRedeemTx(tx ethrpc.Transaction) (RedeemTx, bool, error) {
	//todo
	return RedeemTx{}, true, nil
}

func (e *EthereumAgent) Close() error {
	panic(e)
}
func (e *EthereumAgent) Name() string {
	return "Ethereum Agent"
}
func (e *EthereumAgent) BlockTime() time.Duration {
	return e.blockTime
}
