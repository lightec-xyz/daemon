package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"strconv"
	"strings"

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
		logger.Error("get Proof tx error:%v", err)
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

func ReadEthereumHeight(store store.IStore) (int64, error) {
	var height int64
	err := store.GetObj(ethCurHeightKey, &height)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return 0, err
	}
	return height, nil
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

func WriteEthereumTx(store store.IStore, txes []DbTx) error {
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

func WriteUnGenProof(store store.IStore, chain ChainType, txList []string) error {
	batch := store.Batch()
	for _, txHash := range txList {
		err := batch.BatchPutObj(DbUnGenProofId(chain, txHash), nil)
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

func ReadAllUnGenProof(store store.IStore) ([]common.ZkProofRequest, error) {
	panic(store)
	//var txIds []string
	//var requests []ProofRequest
	//iterator := store.Iterator([]byte(UnGenProofPrefix), nil)
	//defer iterator.Release()
	//for iterator.Next() {
	//	txIds = append(txIds, getTxId(string(iterator.Key())))
	//}
	//err := iterator.Error()
	//if err != nil {
	//	return nil, err
	//}
	//for _, txId := range txIds {
	//	id, chainType, err := parseUnGenProofId(txId)
	//	if err != nil {
	//		return nil, err
	//	}
	//	var tx DbTx
	//	err = store.GetObj(DbTxId(id), &tx)
	//	if err != nil {
	//		logger.Error("get ungen Proof error:%v", err)
	//		return nil, err
	//	}
	//	var req ProofRequest
	//	if chainType == Bitcoin {
	//		req = NewDepositProofParam(tx.TxHash, tx.EthAddr, tx.Amount, tx.Utxo)
	//	} else if chainType == Ethereum {
	//		req = NewRedeemProofParam(tx.TxHash, tx.BtcTxId, tx.Inputs, tx.Outputs)
	//	} else {
	//		return nil, fmt.Errorf("unknown chain type:%v", chainType)
	//	}
	//	requests = append(requests, req)
	//}
	//return requests, nil
}

func DeleteUnGenProof(store store.IStore, chainType ChainType, txId string) error {
	err := store.DeleteObj(DbUnGenProofId(chainType, txId))
	if err != nil {
		logger.Error("delete ungen Proof error:%v", err)
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

func getTxId(key string) string {
	txId := key[strings.Index(key, "_")+1:]
	return txId
}

func parseUnGenProofId(key string) (string, ChainType, error) {
	splits := strings.Split(key, "_")
	if len(splits) != 3 {
		return "", 0, fmt.Errorf("parse ungen Proof id error: %v ", key)
	}
	chainType, err := strconv.ParseInt(splits[1], 10, 32)
	if err != nil {
		return "", 0, err
	}
	return splits[2], ChainType(chainType), nil

}
