package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/store"
	"strconv"
	"strings"
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

func WriteBitcoinTxIds(store store.IStore, height int64, txes []Transaction) error {
	batch := store.Batch()
	for _, tx := range txes {
		err := batch.BatchPutObj(DbBtcHeightPrefix(height, tx.TxHash), nil)
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

func ReadBitcoinTxIds(store store.IStore, height int64) ([]string, error) {
	var txIds []string
	iterator := store.Iterator([]byte(DbBtcHeightPrefix(height, "")), nil)
	defer iterator.Release()
	for iterator.Next() {
		txIds = append(txIds, getTxId(string(iterator.Key())))
	}
	err := iterator.Error()
	if err != nil {
		return nil, err
	}
	return txIds, nil
}

func WriteBitcoinTx(store store.IStore, txes []Transaction) error {
	// todo
	batch := store.Batch()
	for _, tx := range txes {
		err := batch.BatchPutObj(DbTxId(tx.TxHash), tx)
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

func WriteDepositDestChainHash(store store.IStore, txList []Transaction) error {
	// todo
	batch := store.Batch()
	for _, tx := range txList {
		err := batch.BatchPutObj(DbDestId(tx.BtcTxId), tx.TxHash)
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

func ReadTransaction(store store.IStore, txId string) (Transaction, error) {
	var tx Transaction
	err := store.GetObj(DbTxId(txId), &tx)
	if err != nil {
		logger.Error("get tx error:%v", err)
		return tx, err
	}
	return tx, nil
}

func ReadProof(store store.IStore, txId string) (Proof, error) {
	var proof Proof
	err := store.GetObj(DbProofId(txId), &proof)
	if err != nil {
		logger.Error("get proof tx error:%v", err)
		return proof, err
	}
	return proof, nil
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

func WriteEthereumTxIds(store store.IStore, height int64, txes []Transaction) error {
	batch := store.Batch()
	for _, tx := range txes {
		err := batch.BatchPutObj(DbEthHeightPrefix(height, tx.TxHash), nil)
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

func ReadEthereumTxIds(store store.IStore, height int64) ([]string, error) {
	var txIds []string
	iterator := store.Iterator([]byte(DbEthHeightPrefix(height, "")), nil)
	defer iterator.Release()
	for iterator.Next() {
		txIds = append(txIds, getTxId(string(iterator.Key())))
	}
	err := iterator.Error()
	if err != nil {
		return nil, err
	}
	return txIds, nil
}

func WriteEthereumTx(store store.IStore, txes []Transaction) error {
	// todo
	batch := store.Batch()
	for _, tx := range txes {
		err := batch.BatchPutObj(DbTxId(tx.TxHash), tx)
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

func WriteRedeemDestChainHash(store store.IStore, txList []Transaction) error {
	batch := store.Batch()
	for _, tx := range txList {
		err := batch.BatchPutObj(DbDestId(tx.TxHash), tx.BtcTxId)
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

func WriteUnGenProof(store store.IStore, chain ChainType, txList []ProofRequest) error {
	batch := store.Batch()
	for _, tx := range txList {
		err := batch.BatchPutObj(DbUnGenProofId(chain, tx.TxHash), nil)
		if err != nil {
			logger.Error("put ungen proof error:%v", err)
			return err
		}
	}
	err := batch.BatchWriteObj()
	if err != nil {
		logger.Error("put ungen proof batch  error:%v", err)
		return err
	}
	return nil
}

func ReadAllUnGenProof(store store.IStore) ([]ProofRequest, error) {
	var txIds []string
	var requests []ProofRequest
	iterator := store.Iterator([]byte(UnGenProofPrefix), nil)
	defer iterator.Release()
	for iterator.Next() {
		txIds = append(txIds, getTxId(string(iterator.Key())))
	}
	err := iterator.Error()
	if err != nil {
		return nil, err
	}
	for _, txId := range txIds {
		id, chainType, err := parseUnGenProofId(txId)
		if err != nil {
			return nil, err
		}
		var tx Transaction
		err = store.GetObj(DbTxId(id), &tx)
		if err != nil {
			logger.Error("get ungen proof error:%v", err)
			return nil, err
		}
		var req ProofRequest
		if chainType == Bitcoin {
			req = NewDepositProofRequest(tx.TxHash, tx.EthAddr, tx.Amount, tx.Utxo)
		} else if chainType == Ethereum {
			req = NewRedeemProofRequest(tx.TxHash, tx.BtcTxId, tx.Inputs, tx.Outputs)
		} else {
			return nil, fmt.Errorf("unknown chain type:%v", chainType)
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func DeleteUnGenProof(store store.IStore, chainType ChainType, txId string) error {
	err := store.DeleteObj(DbUnGenProofId(chainType, txId))
	if err != nil {
		logger.Error("delete ungen proof error:%v", err)
		return err
	}
	return nil
}

func DeleteUnGenProofs(store store.IStore, chainType ChainType, txes []Transaction) error {
	batch := store.Batch()
	for _, tx := range txes {
		err := batch.BatchDeleteObj(DbUnGenProofId(chainType, tx.TxHash))
		if err != nil {
			logger.Error("delete ungen proof error:%v", err)
			return err
		}
	}
	err := batch.BatchWriteObj()
	if err != nil {
		logger.Error("put ungen proof batch  error:%v", err)
		return err
	}
	return nil
}

func CheckDestHash(store store.IStore, txId string) (bool, error) {
	ok, err := store.HasObj(DbDestId(txId))
	if err != nil {
		return false, err
	}
	return ok, nil
}

func DbProofId(txId string) string {
	pTxID := fmt.Sprintf("%s%s", ProofPrefix, trimOx(txId))
	return pTxID
}

func DbBtcHeightPrefix(height int64, txId string) string {
	pTxID := fmt.Sprintf("%db_%s", height, trimOx(txId))
	return pTxID
}

func DbEthHeightPrefix(height int64, txId string) string {
	pTxID := fmt.Sprintf("%de_%s", height, trimOx(txId))
	return pTxID
}

func DbTxId(txId string) string {
	pTxID := fmt.Sprintf("%s%s", TxPrefix, trimOx(txId))
	return pTxID
}

func DbDestId(txId string) string {
	pTxID := fmt.Sprintf("%s%s", DestChainHashPrefix, trimOx(txId))
	return pTxID
}

func DbUnGenProofId(chain ChainType, txId string) string {
	pTxID := fmt.Sprintf("%s%d_%s", UnGenProofPrefix, chain, trimOx(txId))
	return pTxID
}

func trimOx(hash string) string {
	return strings.TrimPrefix(hash, "0x")
}

func getTxId(key string) string {
	txId := key[strings.Index(key, "_")+1:]
	return txId
}

func parseUnGenProofId(key string) (string, ChainType, error) {
	splits := strings.Split(key, "_")
	if len(splits) != 3 {
		return "", 0, fmt.Errorf("parse ungen proof id error: %v ", key)
	}
	chainType, err := strconv.ParseInt(splits[1], 10, 32)
	if err != nil {
		return "", 0, err
	}
	return splits[2], ChainType(chainType), nil

}
