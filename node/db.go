package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/codec"
	"github.com/lightec-xyz/daemon/common"
	"sort"
	"strings"
	"time"

	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/store"
)

func WriteLatestBeaconSlot(store store.IStore, slot uint64) error {
	return store.PutObj(beaconLatestKey, slot)
}

func ReadLatestBeaconSlot(store store.IStore) (uint64, bool, error) {
	exists, err := store.HasObj(beaconLatestKey)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}
	var slot uint64
	err = store.GetObj(beaconLatestKey, &slot)
	if err != nil {
		return 0, false, err
	}
	return slot, true, nil
}

func WriteBeaconEthNumber(store store.IStore, slot, number uint64) error {
	return store.PutObj(DbBeaconSlotId(slot), number)
}

func ReadBeaconEthNumber(store store.IStore, slot uint64) (uint64, error) {
	var number uint64
	err := store.GetObj(DbBeaconSlotId(slot), &number)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func WriteBeaconSlot(store store.IStore, number, slot uint64) error {
	return store.PutObj(DbBeaconEthNumberId(number), slot)
}

func ReadBeaconSlot(store store.IStore, number uint64) (uint64, bool, error) {
	id := DbBeaconEthNumberId(number)
	exists, err := store.HasObj(id)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}
	var slot uint64
	err = store.GetObj(id, &slot)
	if err != nil {
		return 0, false, err
	}
	return slot, true, nil
}

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

func WriteBitcoinTxIdsByHeight(store store.IStore, height int64, txes []string) error {
	err := store.PutObj(DbBtcHeightPrefix(height), txes)
	if err != nil {
		logger.Error("put bitcoin tx error:%v", err)
		return err
	}
	return nil
}

func ReadBitcoinTxIdsByHeight(store store.IStore, height int64) ([]string, error) {
	var txIds []string
	err := store.GetObj(DbBtcHeightPrefix(height), &txIds)
	if err != nil {
		return nil, err
	}
	return txIds, nil
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

func WriteEthereumTxIdsByHeight(store store.IStore, height int64, txIds []string) error {
	err := store.PutObj(DbEthHeightPrefix(height), txIds)
	if err != nil {
		logger.Error("put ethereum tx error:%v", err)
		return err
	}
	return nil
}
func ReadEthereumTxIdsByHeight(store store.IStore, height int64) ([]string, error) {
	var txIds []string
	err := store.GetObj(DbEthHeightPrefix(height), &txIds)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
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
		err := codec.Unmarshal(iterator.Value(), &tx)
		if err != nil {
			logger.Error("get unsubmit tx error:%v", err)
			return nil, err
		}
		txes = append(txes, tx)
	}
	if err := iterator.Error(); err != nil {
		return nil, err
	}
	return txes, nil
}

func DeleteUnSubmitTx(store store.IStore, hash string) error {
	return store.DeleteObj(DbUnSubmitTxId(hash))
}

func WriteUnGenProof(store store.IStore, chain common.ChainType, list []*DbUnGenProof) error {
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

func ReadUnGenProof(store store.IStore, chainType common.ChainType, txId string) (*DbUnGenProof, error) {
	var proof DbUnGenProof
	err := store.GetObj(DbUnGenProofId(chainType, txId), &proof)
	if err != nil {
		logger.Error("read ungen Proof error:%v", err)
		return nil, err
	}
	return &proof, nil
}

func ReadAllUnGenProofs(store store.IStore, chainType common.ChainType) ([]*DbUnGenProof, error) {
	iterator := store.Iterator([]byte(DbUnGenProofId(chainType, "")), nil)
	defer iterator.Release()
	var txes []*DbUnGenProof
	for iterator.Next() {
		var tx DbUnGenProof
		err := codec.Unmarshal(iterator.Value(), &tx)
		if err != nil {
			logger.Error("read ungen Proof error:%v", err)
			return nil, err
		}
		txes = append(txes, &tx)
	}
	if err := iterator.Error(); err != nil {
		return nil, err
	}
	sort.SliceStable(txes, func(i, j int) bool {
		if txes[i].Height == txes[j].Height {
			return txes[i].TxIndex < txes[j].TxIndex
		}
		return txes[i].Height < txes[j].Height
	})
	return txes, nil
}

func DeleteUnGenProof(store store.IStore, chainType common.ChainType, txId string) error {
	err := store.DeleteObj(DbUnGenProofId(chainType, txId))
	if err != nil {
		logger.Error("delete ungen Proof error:%v", err)
		return err
	}
	return nil
}

func WriteTxIdsByAddr(store store.IStore, txType common.TxType, addr string, txes []DbTx) error {
	batch := store.Batch()
	for _, tx := range txes {
		err := batch.BatchPutObj(DbAddrPrefixTxId(addr, txType, tx.TxHash), nil)
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

func ReadTxIdsByAddr(store store.IStore, txType common.TxType, addr string) ([]string, error) {
	var txIds []string
	iterator := store.Iterator([]byte(DbAddrPrefixTxId(addr, txType, "")), nil)
	defer iterator.Release()
	for iterator.Next() {
		elems := strings.Split(string(iterator.Value()), ProtocolSeparator)
		if len(elems) == 2 {
			txIds = append(txIds, elems[1])
		}
	}
	if err := iterator.Error(); err != nil {
		return nil, err
	}
	return txIds, nil
}

func WriteTxSlot(store store.IStore, txSlot uint64, tx *DbUnGenProof) error {
	return store.PutObj(DbTxSlotId(txSlot, tx.TxHash), tx)
}

func DeleteTxSlot(store store.IStore, txSlot uint64, txHash string) error {
	return store.DeleteObj(DbTxSlotId(txSlot, txHash))
}

func ReadAllTxBySlot(store store.IStore, txSlot uint64) ([]*DbUnGenProof, error) {
	var txes []*DbUnGenProof
	iterator := store.Iterator([]byte(DbTxSlotId(txSlot, "")), nil)
	defer iterator.Release()
	for iterator.Next() {
		var tx DbUnGenProof
		err := codec.Unmarshal(iterator.Value(), &tx)
		if err != nil {
			logger.Error("unmarshal tx error:%v", err)
			return nil, err
		}
		txes = append(txes, &tx)
	}
	if err := iterator.Error(); err != nil {
		return nil, err
	}
	sort.SliceStable(txes, func(i, j int) bool {
		if txes[i].Height == txes[j].Height {
			return txes[i].TxIndex < txes[j].TxIndex
		}
		return txes[i].Height < txes[j].Height
	})
	return txes, nil
}

func WriteTxFinalizedSlot(store store.IStore, txSlot uint64, tx *DbUnGenProof) error {
	return store.PutObj(DbTxFinalizeSlotId(txSlot, tx.TxHash), tx)
}

func DeleteTxFinalizedSlot(store store.IStore, txSlot uint64, txHash string) error {
	return store.DeleteObj(DbTxFinalizeSlotId(txSlot, txHash))
}

func ReadAllTxByFinalizedSlot(store store.IStore, finalizedSlot uint64) ([]*DbUnGenProof, error) {
	var txes []*DbUnGenProof
	iterator := store.Iterator([]byte(DbTxFinalizeSlotId(finalizedSlot, "")), nil)
	defer iterator.Release()
	for iterator.Next() {
		var tx DbUnGenProof
		err := codec.Unmarshal(iterator.Value(), &tx)
		if err != nil {
			logger.Error("unmarshal tx error:%v", err)
			return nil, err
		}
		txes = append(txes, &tx)
	}
	if err := iterator.Error(); err != nil {
		return nil, err
	}
	sort.SliceStable(txes, func(i, j int) bool {
		if txes[i].Height == txes[j].Height {
			return txes[i].TxIndex < txes[j].TxIndex
		}
		return txes[i].Height < txes[j].Height
	})
	return txes, nil
}

func WritePendingRequest(store store.IStore, id string, request *common.ZkProofRequest) error {
	return store.PutObj(DbPendingRequestId(id), request)
}

func DeletePendingRequest(store store.IStore, id string) error {
	return store.DeleteObj(DbPendingRequestId(id))
}

func ReadAllPendingRequests(store store.IStore) ([]*common.ZkProofRequest, error) {
	var txes []*common.ZkProofRequest
	iterator := store.Iterator([]byte(DbPendingRequestId("")), nil)
	defer iterator.Release()
	for iterator.Next() {
		var tx common.ZkProofRequest
		err := codec.Unmarshal(iterator.Value(), &tx)
		if err != nil {
			logger.Error("unmarshal tx error:%v", err)
			return nil, err
		}
		txes = append(txes, &tx)
	}
	if err := iterator.Error(); err != nil {
		return nil, err
	}
	return txes, nil
}

func WriteProofResponse(store store.IStore, resp *common.SubmitProof) error {
	return store.PutObj(DbProofResponseId(resp.Id), resp)
}

func ReadAllProofResponse(store store.IStore) ([]*common.SubmitProof, error) {
	var txes []*common.SubmitProof
	iterator := store.Iterator([]byte(DbProofResponseId("")), nil)
	defer iterator.Release()
	for iterator.Next() {
		var tx common.SubmitProof
		err := codec.Unmarshal(iterator.Value(), &tx)
		if err != nil {
			logger.Error("unmarshal tx error:%v", err)
			return nil, err
		}
		txes = append(txes, &tx)
	}
	if err := iterator.Error(); err != nil {
		return nil, err
	}
	return txes, nil
}
func DeleteProofResponse(store store.IStore, requestId string) error {
	return store.DeleteObj(DbProofResponseId(requestId))
}

func DbProofResponseId(requestId string) string {
	return PendingProofRespPrefix + requestId
}

func ReadWorkerId(store store.IStore) (string, bool, error) {
	exists, err := store.HasObj(workerIdKey)
	if err != nil {
		return "", false, err
	}
	if !exists {
		return "", false, nil
	}
	var id string
	err = store.GetObj(workerIdKey, &id)
	if err != nil {
		return "", false, err
	}
	return id, true, nil
}

func WriteWorkerId(store store.IStore, id string) error {
	return store.PutObj(workerIdKey, id)
}

func ReadZkParamVerify(store store.IStore) (bool, error) {
	exists, err := store.HasObj(zkVerifyKey)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	var verify bool
	err = store.GetObj(zkVerifyKey, &verify)
	if err != nil {
		return false, err
	}
	return verify, nil
}

func WriteZkParamVerify(store store.IStore, verify bool) error {
	return store.PutObj(zkVerifyKey, verify)
}

func WriteNonce(store store.IStore, network, addr string, nonce uint64) error {
	return store.PutObj(DbAddrNonceId(network, addr), nonce)

}

func ReadNonce(store store.IStore, network, addr string) (uint64, bool, error) {
	id := DbAddrNonceId(network, addr)
	exists, err := store.HasObj(id)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}
	var nonce uint64
	err = store.GetObj(id, &nonce)
	if err != nil {
		return 0, false, err
	}
	return nonce, true, nil
}

func WriteUnConfirmTx(store store.IStore, network, hash, proofId string) error {
	return store.PutObj(DbUnConfirmTxId(hash), &DbUnConfirmTx{
		Network: network,
		Hash:    hash,
		ProofId: proofId,
	})
}

func DeleteUnConfirmTx(store store.IStore, hash string) error {
	return store.DeleteObj(DbUnConfirmTxId(hash))
}

func ReadAllUnConfirmTx(store store.IStore) ([]*DbUnConfirmTx, error) {
	var txes []*DbUnConfirmTx
	iterator := store.Iterator([]byte(DbUnConfirmTxId("")), nil)
	defer iterator.Release()
	for iterator.Next() {
		var tx DbUnConfirmTx
		err := codec.Unmarshal(iterator.Value(), &tx)
		if err != nil {
			logger.Error("unmarshal tx error:%v", err)
			return nil, err
		}
		txes = append(txes, &tx)
	}
	if err := iterator.Error(); err != nil {
		return nil, err
	}
	return txes, nil
}

func WriteTaskTime(store store.IStore, flag common.ProofStatus, id string, t time.Time) error {
	return store.PutObj(DbTaskTimeId(flag, id), t)
}

func ReadTaskTime(store store.IStore, flag common.ProofStatus, id string) (time.Time, bool, error) {
	exists, err := store.HasObj(DbTaskTimeId(flag, id))
	if err != nil {
		return time.Time{}, false, err
	}
	if !exists {
		return time.Time{}, false, nil
	}
	var t time.Time
	err = store.GetObj(DbTaskTimeId(flag, id), &t)
	if err != nil {
		return time.Time{}, false, err
	}
	return t, true, nil
}

func WriteFinalityUpdateSlot(store store.IStore, finalizeSlot uint64) error {
	return store.PutObj(DbFinalityUpdateSlotId(finalizeSlot), finalizeSlot)
}

func FindFinalityUpdateNearestSlot(store store.IStore, txSlot uint64) (uint64, bool, error) {
	var start []byte
	if txSlot-common.MaxDiffTxFinalitySlot > 0 {
		start = []byte(fmt.Sprintf("%d", txSlot-common.MaxDiffTxFinalitySlot))
	}
	iterator := store.Iterator([]byte(FinalityUpdateSlotPrefix), start)
	defer iterator.Release()
	for iterator.Next() {
		var slot uint64
		err := codec.Unmarshal(iterator.Value(), &slot)
		if err != nil {
			return 0, false, err
		}
		if slot >= txSlot {
			// todo
			return slot, slot-txSlot <= common.MaxDiffTxFinalitySlot, nil
		}
	}
	return 0, false, nil
}
