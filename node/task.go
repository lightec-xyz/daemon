package node

import (
	"encoding/hex"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/rpc/oasis"
	"github.com/lightec-xyz/daemon/store"
	"time"
)

/*
todo
*/

type TaskManager struct {
	ethClient   *ethrpc.Client
	btcClient   *bitcoin.Client
	oasisClient *oasis.Client
	store       store.IStore
	address     string
	keyStore    *KeyStore
}

func NewTaskManager(store store.IStore, keyStore *KeyStore, ethClient *ethrpc.Client, btcClient *bitcoin.Client, oasisClient *oasis.Client) (*TaskManager, error) {
	address, err := keyStore.Address()
	if err != nil {
		logger.Error("get address error:%v", err)
		return nil, err
	}
	return &TaskManager{
		ethClient:   ethClient,
		btcClient:   btcClient,
		oasisClient: oasisClient,
		keyStore:    keyStore,
		store:       store,
		address:     address,
	}, nil
}

func (t *TaskManager) init() error {
	//allUnSubmitTxs, err := ReadAllUnSubmitTxs(t.store)
	//if err != nil {
	//	logger.Error("get unsubmit tx error:%v", err)
	//	return err
	//}
	return nil
}

func (t *TaskManager) AddTask(resp *common.ZkProofResponse) {
	logger.Info("task manager add retry task: %v", resp.Id())
	unSubmitTx := NewDbUnSubmitTx(resp.TxHash, hex.EncodeToString(resp.Proof), resp.ZkProofType)
	err := WriteUnSubmitTx(t.store, []DbUnSubmitTx{unSubmitTx})
	if err != nil {
		logger.Error("write unsubmit tx error: %v %v", resp.Id(), err)
		return
	}
}

func (t *TaskManager) Check() error {
	unSubmitTxs, err := ReadAllUnSubmitTxs(t.store)
	if err != nil {
		logger.Error("read unsubmit tx error:%v", err)
		return err
	}
	for _, tx := range unSubmitTxs {
		switch tx.ProofType {
		case common.VerifyTxType:
			err := t.UpdateUtxoChange(tx.Hash, tx.Proof)
			if err != nil {
				logger.Error("update utxo error: %v %v", tx.ProofType.String(), tx.Hash)
			}
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

func (t *TaskManager) RedeemZkbtc(hash, proof string) error {
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

func (t *TaskManager) UpdateUtxoChange(hash, proof string) error {
	proofBytes, err := hex.DecodeString(proof)
	if err != nil {
		logger.Error("decode proof error: %v %v", hash, err)
		return err
	}
	err = updateContractUtxoChange(t.ethClient, t.address, t.keyStore.GetPrivateKey(), []string{hash}, proofBytes)
	if err != nil {
		logger.Error("update utxo error: %v", err)
		return err
	}
	err = DeleteUnSubmitTx(t.store, hash)
	if err != nil {
		logger.Error("delete unsubmit tx error: %v %v", hash, err)
		return err
	}
	return nil
}

func (t *TaskManager) Close() error {
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
