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

func WriteBitcoinTx(store store.IStore, height int64, txes []*BitcoinTx) error {
	for _, tx := range txes {
		err := store.PutObj(TxId(height, tx.TxId), tx)
		if err != nil {
			logger.Error("put bitcoin tx error:%v", err)
			return err
		}
	}
	return nil
}

func WriteProof(store store.IStore, txes []TxProof) error {
	for _, tx := range txes {
		err := store.PutObj(TxIdToProofId(tx.TxId), tx)
		if err != nil {
			logger.Error("put proof tx error:%v", err)
			return err
		}
	}
	return nil
}

func UpdateRedeemInfo(store store.IStore, txes []*BitcoinTx) error {
	//todo
	return nil
}

func UpdateProof(store store.IStore, txId, proof, proofType string, status ProofStatus) error {
	return store.PutObj(TxIdToProofId(txId), TxProof{
		TxId:      txId,
		Proof:     proof,
		Status:    status,
		ProofType: proofType,
	})
}

func WriteDestChainHash(store store.IStore, txId, destHash string) error {
	return store.PutObj(TxIdToDestId(txId), destHash)
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

func WriteEthereumTx(store store.IStore, txes []*EthereumTx) error {
	for _, tx := range txes {
		err := store.PutObj(tx.TxId, tx)
		if err != nil {
			logger.Error("put ethereum tx error:%v", err)
			return err
		}
	}
	return nil
}

func TxIdToProofId(txId string) string {
	pTxID := fmt.Sprintf("%s%s", ProofPrefix, txId)
	return pTxID
}
func TxIdToDestId(txId string) string {
	pTxID := fmt.Sprintf("%s%s", DestTxHashPrefix, txId)
	return pTxID
}

func TxId(height int64, txId string) string {
	pTxID := fmt.Sprintf("%d_%s%s", height, TxPrefix, txId)
	return pTxID
}
