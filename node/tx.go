package node

import (
	"encoding/hex"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/rpc/oasis"
	"github.com/lightec-xyz/daemon/store"
	"math/big"
	"sync"
	"time"
)

/*
todo
*/

type TxManager struct {
	ethClient   *ethrpc.Client
	btcClient   *bitcoin.Client
	oasisClient *oasis.Client
	store       store.IStore
	keyStore    *KeyStore
	lock        sync.Mutex
}

func NewTxManager(store store.IStore, keyStore *KeyStore, ethClient *ethrpc.Client, btcClient *bitcoin.Client, oasisClient *oasis.Client) (*TxManager, error) {
	return &TxManager{
		ethClient:   ethClient,
		btcClient:   btcClient,
		oasisClient: oasisClient,
		keyStore:    keyStore,
		store:       store,
	}, nil
}

func (t *TxManager) init() error {
	//allUnSubmitTxs, err := ReadAllUnSubmitTxs(t.store)
	//if err != nil {
	//	logger.Error("get unsubmit tx error:%v", err)
	//	return err
	//}
	return nil
}

func (t *TxManager) AddTask(resp *common.ZkProofResponse) {
	logger.Info("txManager manager add retry task: %v", resp.Id())
	unSubmitTx := NewDbUnSubmitTx(resp.TxHash, hex.EncodeToString(resp.Proof), resp.ZkProofType)
	err := WriteUnSubmitTx(t.store, []DbUnSubmitTx{unSubmitTx})
	if err != nil {
		logger.Error("write unsubmit tx error: %v %v", resp.Id(), err)
		return
	}
}

func (t *TxManager) Check() error {
	unSubmitTxs, err := ReadAllUnSubmitTxs(t.store)
	if err != nil {
		logger.Error("read unsubmit tx error:%v", err)
		return err
	}
	for _, tx := range unSubmitTxs {
		switch tx.ProofType {
		case common.VerifyTxType:
			hash, err := t.UpdateUtxoChange(tx.Hash, tx.Proof)
			if err != nil {
				logger.Error("update utxo error: %v %v", tx.ProofType.String(), tx.Hash)
				return err
			}
			logger.Info("success update utxo txId: %v,hash: %v", tx.Hash, hash)
		case common.RedeemTxType:
			err := t.RedeemZkbtc(tx.Hash, tx.Proof)
			if err != nil {
				logger.Error("update utxo error: %v %v", tx.ProofType.String(), tx.Hash)
			}
		default:
			logger.Warn("never should happen: %v %v", tx.ProofType.String(), tx.Hash)
		}
	}
	return nil
}

func (t *TxManager) getEthAddrNonce(addr string) (uint64, error) {
	chainNonce, err := t.ethClient.GetNonce(addr)
	if err != nil {
		logger.Error("get nonce error: %v %v", addr, err)
		return 0, err
	}
	dbNonce, exists, err := ReadNonce(t.store, common.ETH.String(), addr)
	if err != nil {
		logger.Error("read nonce error: %v %v", addr, err)
		return 0, err
	}
	if !exists {
		return chainNonce, nil
	}
	if chainNonce <= dbNonce {
		return dbNonce + 1, nil
	}
	return chainNonce, nil
}

func (t *TxManager) RedeemZkbtc(hash, proof string) error {
	proofBytes, err := hex.DecodeString(proof)
	if err != nil {
		logger.Error("decode proof error: %v %v", hash, err)
		return err
	}
	_, err = RedeemBtcTx(t.btcClient, t.ethClient, t.oasisClient, hash, proofBytes)
	if err != nil {
		logger.Error("mint zk btc error: %v", err)
		return err
	}
	err = DeleteUnSubmitTx(t.store, hash)
	if err != nil {
		logger.Error("delete unsubmit tx error: %v %v", hash, err)
		return err
	}
	return nil
}

func (t *TxManager) UpdateUtxoChange(txId, proof string) (string, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	proofBytes, err := hex.DecodeString(proof)
	if err != nil {
		logger.Error("decode proof error: %v %v", txId, err)
		return "", err
	}
	privateKey, err := t.keyStore.GetPrivateKey()
	if err != nil {
		logger.Error("get private key error: %v", err)
		return "", err
	}
	ethAddress := t.keyStore.EthAddress()
	nonce, err := t.getEthAddrNonce(ethAddress)
	if err != nil {
		logger.Error("get nonce error: %v %v", ethAddress, err)
		return "", err
	}
	chainId, err := t.ethClient.GetChainId()
	if err != nil {
		logger.Error("get chain id error:%v", err)
		return "", err
	}
	gasPrice, err := t.ethClient.GetGasPrice()
	if err != nil {
		logger.Error("get gas price error:%v", err)
		return "", err
	}
	gasLimit := uint64(500000)
	gasPrice = big.NewInt(0).Mul(gasPrice, big.NewInt(2)) // todo
	txHash, err := t.ethClient.UpdateUtxoChange(privateKey, []string{txId}, nonce, gasLimit, chainId, gasPrice, proofBytes)
	if err != nil {
		logger.Error("update utxo change error:%v", err)
		return "", err
	}
	logger.Debug("update utxo info address: %v txHash: %v,txId: %v,nonce: %v", ethAddress, txHash, txId, nonce)
	err = WriteNonce(t.store, common.ETH.String(), ethAddress, nonce)
	if err != nil {
		logger.Error("write nonce error: %v %v", ethAddress, err)
		return "", err
	}
	err = DeleteUnSubmitTx(t.store, txId)
	if err != nil {
		logger.Error("delete unsubmit tx error: %v %v", txId, err)
		return "", err
	}
	return txHash, nil
}

func (t *TxManager) Close() error {
	return nil
}

func NewDbUnSubmitTx(hash, proof string, proofType common.ZkProofType) DbUnSubmitTx {
	return DbUnSubmitTx{
		Hash:      hash,
		Proof:     proof,
		ProofType: proofType,
		Timestamp: time.Now().UnixNano(),
	}
}
