package node

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
	btctx "github.com/lightec-xyz/daemon/transaction/bitcoin"
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
}

func NewEthereumAgent(cfg NodeConfig, store, memoryStore store.IStore, btcClient *bitcoin.Client, ethClient *ethereum.Client,
	proofRequest chan []ProofRequest, proofResponse <-chan ProofResponse) (IAgent, error) {
	return &EthereumAgent{
		btcClient:            btcClient,
		ethClient:            ethClient,
		store:                store,
		memoryStore:          memoryStore,
		blockTime:            time.Duration(cfg.EthBlockTime) * time.Second,
		proofRequest:         proofRequest,
		proofResponse:        proofResponse,
		checkProofHeightNums: 100,
		exitSign:             make(chan struct{}, 1),
		whiteList:            make(map[string]bool),
		multiAddressInfo:     cfg.MultiAddressInfo,
		btcNetwork:           btctx.NetWork(cfg.BtcNetwork),
	}, nil
}

func (e *EthereumAgent) Init() error {
	has, err := e.store.Has(ethCurHeightKey)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return err
	}
	if has {
		err := e.checkUnCompleteGenerateProofTx()
		if err != nil {
			logger.Error("check uncomplete generate proof tx error:%v", err)
			return err
		}
	} else {
		logger.Debug("init eth current height: %v", InitEthereumHeight)
		err := e.store.PutObj(ethCurHeightKey, InitEthereumHeight)
		if err != nil {
			logger.Error("put init eth current height error:%v", err)
			return err
		}
	}
	//todo
	return nil
}

func (e *EthereumAgent) checkUnCompleteGenerateProofTx() error {
	currentHeight, err := e.getEthHeight()
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return err
	}
	start := currentHeight - e.checkProofHeightNums
	var proofList []ProofRequest
	for index := start; index < currentHeight; index++ {
		hasObj, err := e.store.HasObj(index)
		if err != nil {
			logger.Error("get txIdList error:%v", err)
			return err
		}
		if !hasObj {
			continue
		}
		var txIdList []string
		err = e.store.GetObj(index, &txIdList)
		if err != nil {
			logger.Error("get txIdList error:%v", err)
			return err
		}
		for _, txId := range txIdList {
			var proof TxProof
			err := e.store.GetObj(TxIdToProofId(txId), &proof)
			if err != nil {
				logger.Error("get proof error:%v", err)
				return err
			}
			proofList = append(proofList, ProofRequest{
				TxId:   proof.TxId,
				PType:  EthereumChain,
				ToAddr: proof.ToAddr,
				Amount: proof.Amount,
				Msg:    proof.Msg,
			})
		}
	}
	e.proofRequest <- proofList
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
	if ethHeight >= int64(blockNumber)-12 {
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

func (e *EthereumAgent) Transfer() {
	//todo
	for {
		select {
		case <-e.exitSign:
			logger.Info("exit ethereum transfer goroutine")
			return
		case response := <-e.proofResponse:
			err := e.updateProof(response)
			if err != nil {
				logger.Error("update proof error:%v", err)
				continue
			}
			err = e.RedeemBtcTx(response)
			if err != nil {
				//todo
				logger.Error("redeem btc tx error:%v", err)
			}
			logger.Info("success redeem btc tx:%v", response)
		}
	}

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
	block, err := e.ethClient.GetBlock(height)
	if err != nil {
		logger.Error("ethereum rpc get block error:%v", err)
		return nil, nil, err
	}
	var redeemTxList []RedeemTx
	var proofRequestList []ProofRequest
	for _, tx := range block.Transactions() {
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

func (e *EthereumAgent) RedeemBtcTx(resp ProofResponse) error {
	//todo
	builder := btctx.NewMultiTransactionBuilder()
	err := builder.NetParams(e.btcNetwork)
	if err != nil {
		logger.Error("build btc tx error:%v", err)
		return err
	}
	err = builder.AddMultiPublicKey(e.multiAddressInfo.PublicKeyList, e.multiAddressInfo.NRequired)
	if err != nil {
		logger.Error("build btc tx error:%v", err)
		return err
	}
	err = builder.AddTxIn([]btctx.TxIn{})
	if err != nil {
		logger.Error("build btc tx error:%v", err)
		return err
	}
	err = builder.AddTxOut([]btctx.TxOut{})
	if err != nil {
		logger.Error("build btc tx error:%v", err)
		return err
	}
	err = builder.Sign(func(hash []byte) ([][]byte, error) {
		//todo
		panic("implement me")
	})
	if err != nil {
		logger.Error("build btc tx error:%v", err)
		return err
	}
	txBytes, err := builder.Build()
	if err != nil {
		logger.Error("build btc tx error:%v", err)
		return err
	}
	txHash, err := e.btcClient.Sendrawtransaction(hex.EncodeToString(txBytes))
	if err != nil {
		logger.Error("send btc tx error:%v", err)
		return err
	}
	logger.Info("send redeem btc tx: %v", txHash)
	return nil
}

func (e *EthereumAgent) CheckRedeemTx(tx *types.Transaction) (RedeemTx, bool, error) {
	//todo
	return RedeemTx{}, false, nil
}

func (e *EthereumAgent) updateProof(resp ProofResponse) error {
	pTxId := TxIdToProofId(resp.TxId)
	err := e.store.PutObj(pTxId, TxProof{
		PTxId:  pTxId,
		TxId:   resp.TxId,
		ToAddr: resp.ToAddr,
		Amount: resp.Amount,
		Msg:    resp.Msg,
		Status: ProofSuccess,
	})
	return err
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
