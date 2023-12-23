package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/store"
)

func WriteBitcoinHeight(store store.IStore, height int64) error {
	return store.PutObj(btcCurHeightKey, height)
}

func ReadBitcoinHeight(store store.IStore) (int64, error) {
	var height int64
	err := store.GetObj(btcCurHeightKey, &height)
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return 0, err
	}
	return height, nil
}
func ReadInitBitcoinHeight(store store.IStore) (bool, error) {
	return store.HasObj(btcCurHeightKey)
}

func WriteBitcoinTx(store store.IStore, height int64, txes []BitcoinTx) error {
	batch := store.Batch()
	for _, tx := range txes {
		err := batch.BatchPutObj(DbTxIdKey(height, tx.TxId), tx)
		if err != nil {
			logger.Error("put bitcoin tx error:%v", err)
			return err
		}
	}
	err := batch.BatchWriteObj()
	if err != nil {
		logger.Error("put bitcoin tx batch error:%v", err)
		return err
	}
	return nil
}

func WriteDepositDestChainHash(store store.IStore, txList []EthereumTx) error {
	// todo
	batch := store.Batch()
	for _, tx := range txList {
		err := batch.BatchPutObj(DbTxIdToDestId(tx.BtcTxId), tx.TxHash)
		if err != nil {
			logger.Error("put deposit dest chain error:%v", err)
			return err
		}
	}
	err := batch.BatchWriteObj()
	if err != nil {
		logger.Error("put deposit dest chain batch  error:%v", err)
		return err
	}
	return nil
}

func WriteProof(store store.IStore, txes []Proof) error {
	batch := store.Batch()
	for _, tx := range txes {
		err := batch.BatchPutObj(DbProofId(tx.TxId), tx)
		if err != nil {
			logger.Error("put proof tx error:%v", err)
			return err
		}
	}
	err := batch.BatchWriteObj()
	if err != nil {
		logger.Error("put bitcoin tx batch error:%v", err)
		return err
	}
	return nil
}

func UpdateProof(store store.IStore, txId, proof string, proofType ProofType, status ProofStatus) error {
	txProof := Proof{
		TxId:      txId,
		Proof:     proof,
		Status:    status,
		ProofType: proofType,
	}
	err := store.PutObj(DbProofId(txId), txProof)
	if err != nil {
		logger.Error("put proof tx error:%v %v", txId, err)
		return err
	}
	return nil
}

func WriteEthereumHeight(store store.IStore, height int64) error {
	return store.PutObj(ethCurHeightKey, height)
}

func ReadEthereumHeight(store store.IStore) (int64, error) {
	var height int64
	err := store.GetObj(ethCurHeightKey, &height)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return 0, err
	}
	return height, nil
}

func ReadInitEthereumHeight(store store.IStore) (bool, error) {
	return store.HasObj(ethCurHeightKey)
}

func WriteEthereumTx(store store.IStore, height int64, txes []EthereumTx) error {
	batch := store.Batch()
	for _, tx := range txes {
		err := batch.BatchPutObj(DbTxIdKey(height, tx.BtcTxId), tx)
		if err != nil {
			logger.Error("put ethereum tx error:%v", err)
			return err
		}
	}
	err := batch.BatchWriteObj()
	if err != nil {
		logger.Error("put bitcoin tx batch error:%v", err)
		return err
	}
	return nil
}

func WriteRedeemDestChainHash(store store.IStore, txList []EthereumTx) error {
	batch := store.Batch()
	for _, tx := range txList {
		err := batch.BatchPutObj(DbTxIdToDestId(tx.TxHash), tx.BtcTxId)
		if err != nil {
			logger.Error("put deposit dest chain error:%v", err)
			return err
		}
	}
	err := batch.BatchWriteObj()
	if err != nil {
		logger.Error("put deposit dest chain batch  error:%v", err)
		return err
	}
	return nil
}

func DbProofId(txId string) string {
	pTxID := fmt.Sprintf("%s%s", ProofPrefix, txId)
	return pTxID
}

func DbTxIdKey(height int64, txId string) string {
	pTxID := fmt.Sprintf("%d_%s%s", height, TxPrefix, txId)
	return pTxID
}

func DbTxId(txId string) string {
	pTxID := fmt.Sprintf("%s%s", TxPrefix, txId)
	return pTxID
}

func DbTxIdToDestId(txId string) string {
	pTxID := fmt.Sprintf("%s%s", DestChainHashPrefix, txId)
	return pTxID
}
