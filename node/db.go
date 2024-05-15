package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"sort"
	"strings"

	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/store"
)

func WriteBitcoinHeight(store store.IStore, height int64) error {
	return store.PutObj(btcCurHeightKey, height)
}

func ReadBitcoinHeight(store store.IStore) (int64, bool, error) {
	exists, err := CheckBitcoinHeight(store)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}
	var height int64
	err = store.GetObj(btcCurHeightKey, &height)
	if err != nil {
		logger.Error("get btc current height error:%v", err)
		return 0, false, err
	}
	return height, true, nil
}

func CheckBitcoinHeight(store store.IStore) (bool, error) {
	return store.HasObj(btcCurHeightKey)
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

func WriteBitcoinTx(store store.IStore, txes []DbTx) error {
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

func WriteDestHash(store store.IStore, key, value string) error {
	return store.PutObj(DbDestId(key), value)
}

func ReadDestHash(store store.IStore, key string) (string, error) {
	var value string
	err := store.GetObj(DbDestId(key), &value)
	if err != nil {
		logger.Error("get dest hash error:%v", err)
		return value, err
	}
	return value, nil
}

func ReadDbTx(store store.IStore, txId string) (DbTx, error) {
	var tx DbTx
	err := store.GetObj(DbTxId(txId), &tx)
	if err != nil {
		logger.Error("get tx error:%v", err)
		return tx, err
	}
	return tx, nil
}

func ReadDbProof(store store.IStore, txId string) (DbProof, error) {
	var proof DbProof
	err := store.GetObj(DbProofId(txId), &proof)
	if err != nil {
		logger.Error("get Proof tx error:%v %v", txId, err)
		return proof, err
	}
	return proof, nil
}

func WriteDbProof(store store.IStore, txes []DbProof) error {
	batch := store.Batch()
	for _, tx := range txes {
		err := batch.BatchPutObj(DbProofId(tx.TxHash), tx)
		if err != nil {
			logger.Error("put Proof tx error:%v", err)
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

func UpdateProofStatus(store store.IStore, txId string, proofType common.ZkProofType, status common.ProofStatus) error {
	err := UpdateProof(store, txId, "", proofType, status)
	if err != nil {
		logger.Error("put Proof tx error:%v %v", txId, err)
		return err
	}
	return err
}

func UpdateProof(store store.IStore, txId string, proof string, proofType common.ZkProofType, status common.ProofStatus) error {
	txProof := DbProof{
		TxHash:    txId,
		Proof:     proof,
		Status:    int(status),
		ProofType: proofType,
	}
	err := store.PutObj(DbProofId(txId), txProof)
	if err != nil {
		logger.Error("put Proof tx error:%v %v", txId, err)
		return err
	}
	return nil
}

func WriteEthereumHeight(store store.IStore, height int64) error {
	return store.PutObj(ethCurHeightKey, height)
}

func ReadEthereumHeight(store store.IStore) (int64, bool, error) {
	exists, err := CheckEthereumHeight(store)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}
	var height int64
	err = store.GetObj(ethCurHeightKey, &height)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return 0, false, err
	}
	return height, true, nil
}

func CheckEthereumHeight(store store.IStore) (bool, error) {
	return store.HasObj(ethCurHeightKey)
}

func WriteEthereumTxIds(store store.IStore, height int64, txHashes []string) error {
	batch := store.Batch()
	for _, hash := range txHashes {
		err := batch.BatchPutObj(DbEthHeightPrefix(height, hash), nil)
		if err != nil {
			logger.Error("put ethereum hash error:%v", err)
			return err
		}
	}
	err := batch.BatchWriteObj()
	if err != nil {
		logger.Error("put bitcoin hash batch error:%v", err)
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

func WriteTxes(store store.IStore, txes []DbTx) error {
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

func WriteUnSubmitTx(store store.IStore, txes []DbUnSubmitTx) error {
	batch := store.Batch()
	for _, tx := range txes {
		err := batch.BatchPutObj(DbUnSubmitTxId(tx.Hash), tx)
		if err != nil {
			logger.Error("put unsubmit tx error:%v", err)
			return err
		}
	}
	err := batch.BatchWriteObj()
	if err != nil {
		logger.Error("put unsubmit tx batch error:%v", err)
		return err
	}
	return nil
}

func ReadAllUnSubmitTxs(store store.IStore) ([]DbUnSubmitTx, error) {
	var txes []DbUnSubmitTx
	iterator := store.Iterator([]byte(UnSubmitTxPrefix), nil)
	defer iterator.Release()
	for iterator.Next() {
		var tx DbUnSubmitTx
		err := store.GetObj(iterator.Key(), &tx)
		if err != nil {
			logger.Error("get unsubmit tx error:%v", err)
			return nil, err
		}
		txes = append(txes, tx)
	}
	err := iterator.Error()
	if err != nil {
		return nil, err
	}
	return txes, nil
}

func DeleteUnSubmitTx(store store.IStore, hash string) error {
	return store.DeleteObj(DbUnSubmitTxId(hash))
}

func WriteUnGenProof(store store.IStore, chain ChainType, list []*DbUnGenProof) error {
	batch := store.Batch()
	for _, item := range list {
		err := batch.BatchPutObj(DbUnGenProofId(chain, item.TxHash), item)
		if err != nil {
			logger.Error("put ungen Proof error:%v", err)
			return err
		}
	}
	err := batch.BatchWriteObj()
	if err != nil {
		logger.Error("put ungen Proof batch  error:%v", err)
		return err
	}
	return nil
}

// todo

func ReadAllUnGenProofs(store store.IStore, chainType ChainType) ([]*DbUnGenProof, error) {
	var keys []string
	queryPrefix := fmt.Sprintf("%s%d", UnGenProofPrefix, chainType)
	iterator := store.Iterator([]byte(queryPrefix), nil)
	defer iterator.Release()
	for iterator.Next() {
		keys = append(keys, string(iterator.Key()))
	}
	err := iterator.Error()
	if err != nil {
		logger.Error("read ungen Proof error:%v", err)
		return nil, err
	}
	var unGenPreProofs []*DbUnGenProof
	for _, key := range keys {
		var unGenProof DbUnGenProof
		err := store.GetObj(key, &unGenProof)
		if err != nil {
			logger.Error("read ungen Proof error:%v", err)
			return nil, err
		}
		unGenPreProofs = append(unGenPreProofs, &unGenProof)
	}
	sort.Slice(unGenPreProofs, func(i, j int) bool { return unGenPreProofs[i].Height < unGenPreProofs[j].Height })
	return unGenPreProofs, nil
}

func DeleteUnGenProof(store store.IStore, chainType ChainType, txId string) error {
	err := store.DeleteObj(DbUnGenProofId(chainType, txId))
	if err != nil {
		logger.Error("delete ungen Proof error:%v", err)
		return err
	}
	return nil
}

func WriteAddrTxs(store store.IStore, addr string, txes []DbTx) error {
	batch := store.Batch()
	for _, tx := range txes {
		err := batch.BatchPutObj(DbAddrPrefixTxId(addr, tx.TxHash), tx)
		if err != nil {
			logger.Error("put addr tx error:%v", err)
			return err
		}
	}
	err := batch.BatchWriteObj()
	if err != nil {
		logger.Error("put addr tx batch error:%v", err)
		return err
	}
	return nil
}

func ReadAddrTxs(store store.IStore, addr string) ([]DbTx, error) {
	var txes []DbTx
	iterator := store.Iterator([]byte(DbAddrPrefixTxId(addr, "")), nil)
	defer iterator.Release()
	for iterator.Next() {
		var tx DbTx
		err := store.GetObj(iterator.Key(), &tx)
		if err != nil {
			logger.Error("get addr tx error:%v", err)
			return nil, err
		}
		txes = append(txes, tx)
	}
	err := iterator.Error()
	if err != nil {
		return nil, err
	}
	return txes, nil
}

func getTxId(key string) string {
	txId := key[strings.Index(key, "_")+1:]
	return txId
}
