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

func WriteProof(store store.IStore, txes []Proof) error {
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

func UpdateProof(store store.IStore, txId, proof string, proofType ProofType, status ProofStatus) error {
	txProof := Proof{
		TxId:      txId,
		Proof:     proof,
		Status:    status,
		ProofType: proofType,
	}
	return store.PutObj(TxIdToProofId(txId), txProof)
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

func WriteEthereumTx(store store.IStore, txes []EthereumTx) error {
	for _, tx := range txes {
		err := store.PutObj(tx.TxHash, tx)
		if err != nil {
			logger.Error("put ethereum tx error:%v", err)
			return err
		}
	}
	return nil
}

func UpdateDepositTxFinal(store store.IStore, depositTxes []EthereumTx) error {
	// todo
	return nil
}

func ReadAddressNonce(store store.IStore, address string) (uint64, error) {
	key := fmt.Sprintf("%s%s", NoncePrefix, address)
	var nonce uint64
	err := store.GetObj(key, &nonce)
	if err != nil {
		logger.Error("nonce manager get nonce error: %v %v", address, err)
		return 0, err
	}
	return nonce, nil
}

func WriteAddressNonce(store store.IStore, address string, nonce uint64) error {
	key := fmt.Sprintf("%s%s", NoncePrefix, address)
	err := store.PutObj(key, nonce)
	if err != nil {
		logger.Error("write address nonce error: %v %v", address, err)
		return err
	}
	return nil
}
func CheckAddressNonce(store store.IStore, address string) (bool, error) {
	key := fmt.Sprintf("%s%s", NoncePrefix, address)
	ok, err := store.HasObj(key)
	if err != nil {
		logger.Error("nonce manager get nonce error: %v %v", address, err)
		return false, err
	}
	return ok, nil
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
