package node

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lightec-xyz/daemon/codec"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/store"
)

type ChainStore struct {
	store store.IStore
}

func NewChainStore(store store.IStore) *ChainStore {
	return &ChainStore{store: store}
}

func (cs *ChainStore) WriteDepositCount(height, count uint64) error {
	return cs.store.PutObj(height, count)
}
func (cs *ChainStore) IncrDepositCount(height uint64) error {
	count, _, err := cs.ReadDepositCount(height)
	if err != nil {
		return err
	}
	return cs.WriteDepositCount(height, count+1)
}

func (cs *ChainStore) ReadDepositCount(height uint64) (uint64, bool, error) {
	var count uint64
	exists, err := cs.store.GetValue(height, &count)
	return count, exists, err
}

func (cs *ChainStore) WriteBtcBlock(hash string, block string) error {
	return cs.store.PutObj(dbBtcBlockKey(hash), block)
}

func (cs *ChainStore) ReadBtcBlock(hash string) (string, bool, error) {
	var block string
	exists, err := cs.store.GetValue(dbBtcBlockKey(hash), &block)
	return block, exists, err
}

func (cs *ChainStore) WriteBlockHeader(hash string, header string) error {
	return cs.store.PutObj(dbBtcHeaderKey(hash), header)
}

func (cs *ChainStore) ReadBlockHeader(hash string) (string, bool, error) {
	var header string
	exists, err := cs.store.GetValue(dbBtcHeaderKey(hash), &header)
	if err != nil {
		return "", false, err
	}
	if !exists {
		return "", false, nil
	}
	return header, true, nil
}

func (cs *ChainStore) WriteUpdateUtxoDest(hash, dest string) error {
	return cs.store.PutObj(dbUpdateUtxoDestKey(hash), dest)
}

func (cs *ChainStore) ReadUpdateUtxoDest(hash string) (string, bool, error) {
	exists, err := cs.store.HasObj(dbUpdateUtxoDestKey(hash))
	if err != nil {
		return "", false, err
	}
	if !exists {
		return "", false, nil
	}
	var dest string
	err = cs.store.GetObj(dbUpdateUtxoDestKey(hash), &dest)
	return dest, true, err
}

func (cs *ChainStore) WriteChainFork(chain string, forkInfo *ChainFork) error {
	return cs.store.PutObj(dbChainForkKey(chain, forkInfo.Timestamp), forkInfo)
}

func (cs *ChainStore) ReadChainForks(chain string) ([]*ChainFork, error) {
	var result []*ChainFork
	err := cs.store.Iter(genPrefix(chainForkPrefix, chain), nil, func(key, value []byte) error {
		info := &ChainFork{}
		err := codec.Unmarshal(value, info)
		if err != nil {
			return err
		}
		result = append(result, info)
		return nil
	})
	return result, err

}

func (cs *ChainStore) ReadUpdateCpTx() (*DbTx, bool) {
	var dbTx DbTx
	err := cs.store.GetObj(latestUpdateCpKey, &dbTx)
	if err != nil {
		return nil, false
	}
	return &dbTx, true
}

func (cs *ChainStore) WriteTxProved(txIds []string, status bool) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		for _, id := range txIds {
			err := batch.BatchPutObj(dbTxProvedKey(id), status)
			if err != nil {
				logger.Error("put tx proved error:%v", err)
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) ReadTxProved(txId string) (bool, error) {
	dbKey := dbTxProvedKey(txId)
	var result bool
	err := cs.store.GetObj(dbKey, &result)
	if err != nil {
		return false, nil
	}
	return result, nil
}

func (cs *ChainStore) ReadCheckpointHash(txId string) (string, bool, error) {
	dbTxes, err := cs.ReadDbTxes(txId)
	if err != nil {
		logger.Error("read db tx error:%v", err)
		return "", false, err
	}
	if len(dbTxes) != 1 {
		logger.Warn("read db tx error:%v", err)
		return "", false, err
	}
	cpHash, ok, err := cs.ReadCheckpoint(dbTxes[0].CheckPointHeight)
	if err != nil {
		logger.Error("read checkpoint error:%v", err)
		return "", false, err
	}
	if !ok {
		logger.Warn("read checkpoint error:%v", err)
		return "", false, nil
	}
	return cpHash, true, nil
}

func (cs *ChainStore) WriteCheckpoint(height uint64, hash string) error {
	return cs.store.PutObj(dbCheckpointKey(height), hash)
}
func (cs *ChainStore) ReadCheckpoint(height uint64) (string, bool, error) {
	exists, err := cs.store.HasObj(dbCheckpointKey(height))
	if err != nil {
		return "", false, err
	}
	if !exists {
		return "", false, nil
	}
	var hash string
	err = cs.store.GetObj(dbCheckpointKey(height), &hash)
	if err != nil {
		return "", false, err
	}
	return hash, true, nil
}

func (cs *ChainStore) WriteMiner(addr string) error {
	return cs.store.PutObj(dbMinerAddrKey(addr), addr)
}

func (cs *ChainStore) ReadAllMiners() ([]string, error) {
	var miners []string
	iterator := cs.store.Iterator(dbMinerAddrKey(""), nil)
	defer iterator.Release()
	for iterator.Next() {
		var miner string
		err := codec.Unmarshal(iterator.Value(), &miner)
		if err != nil {
			logger.Error("unmarshal tx error:%v", err)
			return nil, err
		}
		miners = append(miners, miner)
	}
	return miners, nil
}

func (cs *ChainStore) WriteMinerPower(addr string, power, timestamp uint64) error {
	return cs.store.PutObj(dbMinerPowerKey(addr), DbMiner{
		Miner:     addr,
		Power:     power,
		Timestamp: timestamp,
	})
}
func (cs *ChainStore) ReadMinerPower(addr string) (*DbMiner, error) {
	var power DbMiner
	err := cs.store.GetObj(dbMinerPowerKey(addr), &power)
	if err != nil {
		return nil, err
	}
	return &power, nil
}

func (cs *ChainStore) WriteEthTxHeight(height uint64, dbTxIds []string) error {
	batch := cs.store.Batch()
	for _, txId := range dbTxIds {
		err := batch.BatchPutObj(ethTxHeightKey(height, txId), nil)
		if err != nil {
			return err
		}
	}
	return batch.BatchWrite()
}

func (cs *ChainStore) ReadEthTxHeight(height uint64) ([]string, error) {
	var txIds []string
	iter := cs.store.Iterator(ethTxHeightKey(height, ""), nil)
	defer iter.Release()
	for iter.Next() {
		txId, err := TxHeightKeyToTxId(iter.Key())
		if err != nil {
			return nil, err
		}
		txIds = append(txIds, txId)
	}
	return txIds, nil
}
func (cs *ChainStore) WriteBtcTxHeight(height uint64, txIds []string) error {
	batch := cs.store.Batch()
	for _, txId := range txIds {
		err := batch.BatchPutObj(btcTxHeightKey(height, txId), nil)
		if err != nil {
			return err
		}
	}
	return batch.BatchWrite()
}
func (cs *ChainStore) ReadBtcTxHeight(height uint64) ([]string, error) {
	var txIds []string
	iter := cs.store.Iterator(btcTxHeightKey(height, ""), nil)
	defer iter.Release()
	for iter.Next() {
		txId, err := TxHeightKeyToTxId(iter.Key())
		if err != nil {
			return nil, err
		}
		txIds = append(txIds, txId)
	}
	return txIds, nil
}

func (cs *ChainStore) DeleteDestHash(txIds []string) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		for _, id := range txIds {
			err := batch.BatchDeleteObj(dbDestId(id))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) DeleteTxProved(txIds []string) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		for _, id := range txIds {
			err := batch.BatchDeleteObj(dbTxProvedKey(id))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) DeleteEthTxHeight(height uint64, txIds []string) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		for _, txId := range txIds {
			err := batch.BatchDeleteObj(ethTxHeightKey(height, txId))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) DeleteBtcTxHeight(height uint64, txIds []string) error {
	batch := cs.store.Batch()
	for _, txId := range txIds {
		err := batch.BatchDeleteObj(btcTxHeightKey(height, txId))
		if err != nil {
			return err
		}
	}
	return batch.BatchWriteObj()
}

func (cs *ChainStore) ReadBtcHeaderHash(start, end uint64) ([]string, error) {
	var hashes []string
	for index := start; index <= end; index++ {
		var hash string
		err := cs.store.GetObj(dbBtcBlockHashKey(index), &hash)
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, hash)
	}
	return hashes, nil
}

func (cs *ChainStore) WriteBitcoinHash(height uint64, hash string) error {
	return cs.store.PutObj(dbBtcBlockHashKey(height), hash)
}

func (cs *ChainStore) ReadBitcoinHash(height uint64) (string, bool, error) {
	exists, err := cs.store.HasObj(dbBtcBlockHashKey(height))
	if err != nil {
		return "", false, err
	}
	if !exists {
		return "", false, nil
	}
	var hash string
	err = cs.store.GetObj(dbBtcBlockHashKey(height), &hash)
	if err != nil {
		return "", false, err
	}
	return hash, true, nil
}

func (cs *ChainStore) WriteEthHash(height uint64, hash string) error {
	return cs.store.PutObj(dbEthBlockHashKey(height), hash)
}

func (cs *ChainStore) ReadEthHash(height uint64) (string, bool, error) {
	exists, err := cs.store.HasObj(dbEthBlockHashKey(height))
	if err != nil {
		return "", false, err
	}
	if !exists {
		return "", false, nil
	}
	var hash string
	err = cs.store.GetObj(dbEthBlockHashKey(height), &hash)
	if err != nil {
		return "", false, err
	}
	return hash, true, nil
}

func (cs *ChainStore) WriteBtcHeight(height uint64) error {
	return cs.store.PutObj(btcCurHeightKey, height)
}

func (cs *ChainStore) ReadBtcHeight() (uint64, bool, error) {
	exists, err := cs.store.HasObj(btcCurHeightKey)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}
	var height uint64
	err = cs.store.GetObj(btcCurHeightKey, &height)
	if err != nil {
		return 0, false, err
	}
	return height, true, nil
}

func (cs *ChainStore) WriteLatestIcpSig(sig DbIcpSignature) error {
	return cs.store.PutObj(latestIcpSignatureKey, sig)
}

func (cs *ChainStore) ReadLatestIcpSig() (*DbIcpSignature, bool, error) {
	var sig DbIcpSignature
	exists, err := cs.store.GetValue(latestIcpSignatureKey, &sig)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	return &sig, true, nil
}

func (cs *ChainStore) WriteIcpSignature(height uint64, value DbIcpSignature) error {
	return cs.store.PutObj(dbDfinityBlockSigId(height), value)
}
func (cs *ChainStore) ReadIcpSignature(height uint64) (DbIcpSignature, bool, error) {
	var value DbIcpSignature
	exists, err := cs.store.GetValue(dbDfinityBlockSigId(height), &value)
	if err != nil {
		return value, false, err
	}
	return value, exists, nil
}

func (cs *ChainStore) WriteDbTxes(txes ...*DbTx) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		for _, tx := range txes {
			err := batch.BatchPutObj(dbTxId(tx.Hash, tx.TxType, tx.LogIndex), tx)
			if err != nil {
				logger.Error("put ethereum tx error:%v", err)
				return err
			}
		}
		return nil
	},
	)
}

func (cs *ChainStore) ReadDbTxes(txId string) ([]*DbTx, error) {
	var txes []*DbTx
	err := cs.store.Iter(genPrefix(txPrefix, txId), nil, func(key, value []byte) error {
		var tx DbTx
		err := codec.Unmarshal(value, &tx)
		if err != nil {
			logger.Error("unmarshal tx error:%v", err)
			return err
		}
		txes = append(txes, &tx)
		return nil
	})
	return txes, err
}

func (cs *ChainStore) BtcSaveData(height uint64, depositTxs, redeemTxes []*DbTx) error {
	var redeemDestHashes []string
	for _, tx := range redeemTxes {
		destHash, exists, err := cs.GetDestHash(tx.Hash)
		if err != nil {
			return err
		}
		if !exists {
			continue
		}
		redeemDestHashes = append(redeemDestHashes, destHash)
	}
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		if len(depositTxs) > 0 {
			err := batch.BatchPutObj(latestUpdateCpKey, depositTxs[0])
			if err != nil {
				return err
			}
		}
		// delete ethereum redeem proof
		for _, hash := range redeemDestHashes {
			err := batch.BatchDeleteObj(dbUnGenProofId(common.EthereumChain, hash))
			if err != nil {
				return err
			}
			err = batch.BatchPutObj(dbTxProvedKey(hash), true)
			if err != nil {
				return err
			}
		}

		allTxes := mergeDbTxes(depositTxs, redeemTxes)
		//height-[]txIds
		txIds := txesToDbTxIds(allTxes)
		for _, txId := range txIds {
			err := batch.BatchPutObj(btcTxHeightKey(height, txId), nil)
			if err != nil {
				return err
			}
		}
		// all btc db tx
		for _, tx := range allTxes {
			err := batch.BatchPutObj(dbTxId(tx.Hash, tx.TxType, tx.LogIndex), tx)
			if err != nil {
				return err
			}
		}
		// db proof
		dbProofs := txesToDbProofs(allTxes)
		for _, tx := range dbProofs {
			err := batch.BatchPutObj(dbProofId(tx.TxHash), tx)
			if err != nil {
				return err
			}
		}
		// group txes by address,maybe replace by explorer
		addrPrefixTxes := txesByAddrGroup(allTxes)
		for addr, addrDbTxes := range addrPrefixTxes {
			for _, tx := range addrDbTxes {
				err := batch.BatchPutObj(dbAddrPrefixTxId(addr, tx.TxType, tx.Hash), nil)
				if err != nil {
					logger.Error("put addr tx error:%v", err)
					return err
				}
			}
		}
		// unGenProof
		unGenProofs := txesToUnGenProofs(allTxes)
		for _, tx := range unGenProofs {
			err := batch.BatchPutObj(dbUnGenProofId(common.BitcoinChain, tx.Hash), tx)
			if err != nil {
				logger.Error("put unGenProof error:%v", err)
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) BtcDeleteData(height uint64) error {
	txIds, err := cs.ReadBtcTxHeight(height)
	if err != nil {
		logger.Error("read btc tx height error: %v %v", height, err)
		return err
	}
	hash, ok, err := cs.ReadBitcoinHash(height)
	if err != nil {
		logger.Error("read btc hash error: %v %v", height, err)
		return err
	}
	if !ok {
		logger.Error("btc hash not exist: %v", height)
	}
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		err = batch.BatchDeleteObj(dbBtcBlockHashKey(height))
		if err != nil {
			return err
		}
		err = batch.BatchDeleteObj(dbBtcBlockKey(hash))
		if err != nil {
			return err
		}
		err = batch.BatchDeleteObj(dbBtcHeaderKey(hash))
		if err != nil {
			return err
		}
		for _, id := range txIds {
			err := batch.BatchDeleteObj(genPrefix(txPrefix, id))
			if err != nil {
				return err
			}
			err = batch.BatchDeleteObj(dbUnGenProofId(common.BitcoinChain, id))
			if err != nil {
				return err
			}
			err = batch.BatchDeleteObj(dbProofId(id))
			if err != nil {
				return err
			}
			err = batch.BatchDeleteObj(btcTxHeightKey(height, id))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) EthSaveData(height uint64, depositTxes, redeemTxes, updateUtxoTxes, depositRewards, redeemRewards []*DbTx) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		//linked dest hash
		linkedIdTxes := mergeDbTxes(depositTxes, redeemTxes, updateUtxoTxes)
		for _, tx := range linkedIdTxes {
			// btcTxId -> ethTxHash
			err := batch.BatchPutObj(dbDestId(tx.UtxoId), tx.Hash)
			if err != nil {
				logger.Error("write dest id error:%v %v", tx.Hash, err)
				return err
			}
			// ethTxHash -> btcTxId
			err = batch.BatchPutObj(dbDestId(tx.Hash), tx.UtxoId)
			if err != nil {
				logger.Error("update deposit final status error: %v %v", height, err)
				return err
			}
		}
		// update utxo dest id
		for _, tx := range updateUtxoTxes {
			err := batch.BatchPutObj(dbUpdateUtxoDestKey(tx.Hash), tx.UtxoId)
			if err != nil {
				logger.Error("update deposit final status error: %v %v", height, err)
				return err
			}
			err = batch.BatchPutObj(dbUpdateUtxoDestKey(tx.UtxoId), tx.Hash)
			if err != nil {
				logger.Error("update deposit final status error: %v %v", height, err)
				return err
			}
		}
		btcGenProofIds := ethTxesToBtcIds(mergeDbTxes(depositTxes, updateUtxoTxes))
		for _, id := range btcGenProofIds {
			//remove bitcoin gen proof
			err := batch.BatchDeleteObj(dbUnGenProofId(common.BitcoinChain, id))
			if err != nil {
				logger.Error("write dest id error:%v %v", id, err)
				return err
			}
			// record a flag to skip gen proof when bitcoin check scheduler
			err = batch.BatchPutObj(dbTxProvedKey(id), true)
			if err != nil {
				logger.Error("write dest id error:%v %v", id, err)
				return err
			}
		}

		// Redeem db proofs
		dbRedeemProofs := txesToDbProofs(redeemTxes)
		for _, tx := range dbRedeemProofs {
			err := batch.BatchPutObj(dbProofId(tx.TxHash), tx)
			if err != nil {
				logger.Error("put Proof tx error:%v", err)
				return err
			}
		}
		// cache need to generate Redeem proof
		unGenProofs := txesToUnGenProofs(redeemTxes)
		for _, item := range unGenProofs {
			err := batch.BatchPutObj(dbUnGenProofId(common.EthereumChain, item.Hash), item)
			if err != nil {
				logger.Error(":%v", err)
				return err
			}
		}

		allTxes := mergeDbTxes(depositTxes, redeemTxes, updateUtxoTxes, depositRewards, redeemRewards)
		// height ->[]TxHash
		dbTxIds := txesToDbTxIds(allTxes)
		for _, txId := range dbTxIds {
			err := batch.BatchPutObj(ethTxHeightKey(height, txId), nil)
			if err != nil {
				return err
			}
		}

		// save all txes
		for _, tx := range allTxes {
			err := batch.BatchPutObj(dbTxId(tx.Hash, tx.TxType, tx.LogIndex), tx)
			if err != nil {
				logger.Error("put ethereum tx error:%v", err)
				return err
			}
		}
		// group txes by address,maybe replace by explorer
		addrPrefixTxes := txesByAddrGroup(allTxes)
		for addr, addrDbTxes := range addrPrefixTxes {
			for _, tx := range addrDbTxes {
				err := batch.BatchPutObj(dbAddrPrefixTxId(addr, tx.TxType, tx.Hash), nil)
				if err != nil {
					logger.Error("put addr tx error:%v", err)
					return err
				}
			}
		}
		return nil
	})

}

func (cs *ChainStore) EthDeleteData(height uint64) error {
	txIds, err := cs.ReadEthTxHeight(height)
	if err != nil {
		logger.Error("read eth tx height error:%v", err)
		return err
	}
	var allDbTxs []*DbTx
	for _, id := range txIds {
		txes, err := cs.ReadDbTxes(id)
		if err != nil {
			return err
		}
		allDbTxs = append(allDbTxs, txes...)
	}
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		// remove tx by addr group
		for _, tx := range allDbTxs {
			if tx.Sender != "" {
				err := batch.BatchDeleteObj(dbAddrPrefixTxId(tx.Sender, tx.TxType, tx.Hash))
				if err != nil {
					return err
				}
			}
		}

		for _, id := range txIds {
			//remove all txes
			err := batch.BatchDeleteObj(genPrefix(txPrefix, id))
			if err != nil {
				return err
			}
			//remove all un gen proof
			err = batch.BatchDeleteObj(dbUnGenProofId(common.EthereumChain, id))
			if err != nil {
				return err
			}
			//remove all proof
			err = batch.BatchDeleteObj(dbProofId(id))
			if err != nil {
				return err
			}
			//remove all dest id
			err = batch.BatchDeleteObj(dbDestId(id))
			if err != nil {
				return err
			}
			//remove proof tx proved flag
			err = batch.BatchDeleteObj(dbTxProvedKey(id))
			if err != nil {
				return err
			}
			//remove tx height
			err = batch.BatchDeleteObj(ethTxHeightKey(height, id))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) DeleteDbTxes(txId []string) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		for _, id := range txId {
			err := batch.BatchDeleteObj(genPrefix(txPrefix, id))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) DeleteDbProof(txIds []string) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		for _, txId := range txIds {
			err := batch.BatchDeleteObj(dbProofId(txId))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) DeleteDbBtcTxes(txId []string) error {
	batch := cs.store.Batch()
	for _, id := range txId {
		err := batch.BatchDeleteObj(dbBtcTxId(id))
		if err != nil {
			return err
		}
	}
	return batch.BatchWrite()
}

func (cs *ChainStore) DeleteAddrTxesPrefix(txIds []string) error {
	for _, txId := range txIds {
		dbTx, err := cs.ReadDbTxes(txId)
		if err != nil {
			logger.Error("read db tx error: %v %v", txId, err)
			return err
		}
		for _, tx := range dbTx {
			if tx.Sender != "" {
				err = cs.DeleteAddrTxPrefix(tx.Sender, tx.TxType, tx.Hash)
				if err != nil {
					logger.Error("delete addr tx error: %v %v", txId, err)
					return err
				}
			}
		}
	}
	return nil
}

func (cs *ChainStore) DeleteAddrTxPrefix(addr string, txType common.TxType, hash string) error {
	return cs.store.DeleteObj(dbAddrPrefixTxId(addr, txType, hash))
}

func (cs *ChainStore) WriteLatestBeaconSlot(slot uint64) error {
	return cs.store.PutObj(beaconLatestKey, slot)
}

func (cs *ChainStore) ReadLatestBeaconSlot() (uint64, bool, error) {
	exists, err := cs.store.HasObj(beaconLatestKey)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}
	var slot uint64
	err = cs.store.GetObj(beaconLatestKey, &slot)
	if err != nil {
		return 0, false, err
	}
	return slot, true, nil
}

func (cs *ChainStore) ReadEthNumberBySlot(slot uint64) (uint64, error) {
	var number uint64
	err := cs.store.GetObj(dbBeaconSlotId(slot), &number)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (cs *ChainStore) WriteBeaconEthNumber(slot, number uint64) error {
	return cs.store.PutObj(dbBeaconSlotId(slot), number)
}

func (cs *ChainStore) WriteBeaconSlot(number, slot uint64) error {
	return cs.store.PutObj(dbBeaconEthNumberId(number), slot)
}

func (cs *ChainStore) ReadSlotByHeight(number uint64) (uint64, bool, error) {
	id := dbBeaconEthNumberId(number)
	exists, err := cs.store.HasObj(id)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}
	var slot uint64
	err = cs.store.GetObj(id, &slot)
	if err != nil {
		return 0, false, err
	}
	return slot, true, nil
}

func (cs *ChainStore) WriteDestHash(key, value string) error {
	return cs.store.PutObj(dbDestId(key), value)
}

func (cs *ChainStore) GetDestHash(key string) (string, bool, error) {
	exists, err := cs.store.HasObj(dbDestId(key))
	if err != nil {
		return "", false, err
	}
	if !exists {
		return "", false, nil
	}
	var hash string
	err = cs.store.GetObj(dbDestId(key), &hash)
	if err != nil {
		return hash, false, err
	}
	return hash, true, nil
}

func (cs *ChainStore) ReadDestHash(key string) (string, error) {
	var value string
	err := cs.store.GetObj(dbDestId(key), &value)
	if err != nil {
		//logger.Error("get dest hash error:%v", err)
		return value, err
	}
	return value, nil
}

func (cs *ChainStore) ReadDbProof(txId string) (DbProof, error) {
	var proof DbProof
	err := cs.store.GetObj(dbProofId(txId), &proof)
	if err != nil {
		//logger.Error("get Proof tx error:%v %v", txId, err)
		return proof, err
	}
	return proof, nil
}

func (cs *ChainStore) WriteDbProof(txes ...DbProof) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		for _, tx := range txes {
			err := batch.BatchPutObj(dbProofId(tx.TxHash), tx)
			if err != nil {
				logger.Error("put Proof tx error:%v", err)
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) UpdateProofStatus(txId string, proofType common.ProofType, status common.ProofStatus) error {
	err := cs.UpdateProof(txId, "", proofType, status)
	if err != nil {
		logger.Error("put Proof tx error:%v %v", txId, err)
		return err
	}
	return err
}

func (cs *ChainStore) UpdateProof(txId string, proof string, proofType common.ProofType, status common.ProofStatus) error {
	txProof := DbProof{
		TxHash:    txId,
		Proof:     proof,
		Status:    int(status),
		ProofType: proofType,
	}
	err := cs.store.PutObj(dbProofId(txId), txProof)
	if err != nil {
		logger.Error("put Proof tx error:%v %v", txId, err)
		return err
	}
	return nil
}

func (cs *ChainStore) WriteEthereumHeight(height uint64) error {
	return cs.store.PutObj(ethCurHeightKey, height)
}

func (cs *ChainStore) ReadEthereumHeight() (uint64, bool, error) {
	exists, err := cs.store.HasObj(ethCurHeightKey)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}
	var height uint64
	err = cs.store.GetObj(ethCurHeightKey, &height)
	if err != nil {
		logger.Error("get eth current height error:%v", err)
		return 0, false, err
	}
	return height, true, nil
}

func (cs *ChainStore) WriteUnSubmitTx(txes ...DbUnSubmitTx) error {
	batch := cs.store.Batch()
	for _, tx := range txes {
		err := batch.BatchPutObj(dbUnSubmitTxId(tx.Hash), tx)
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

func (cs *ChainStore) ReadUnSubmitTxs() ([]DbUnSubmitTx, error) {
	var txes []DbUnSubmitTx
	err := cs.store.Iter(dbUnSubmitTxId(""), nil, func(key []byte, value []byte) error {
		var tx DbUnSubmitTx
		err := codec.Unmarshal(value, &tx)
		if err != nil {
			return err
		}
		txes = append(txes, tx)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return txes, nil
}

func (cs *ChainStore) DeleteUnSubmitTx(hash string) error {
	return cs.store.DeleteObj(dbUnSubmitTxId(hash))
}

func (cs *ChainStore) ReadDbUnGenProof(chainType common.ChainType, txId string) (*DbUnGenProof, bool, error) {
	dbKey := dbUnGenProofId(chainType, txId)
	has, err := cs.store.HasObj(dbKey)
	if err != nil {
		return nil, false, err
	}
	if !has {
		return nil, false, nil
	}
	var value DbUnGenProof
	err = cs.store.GetObj(dbKey, &value)
	if err != nil {
		return nil, false, err
	}
	return &value, true, nil

}

func (cs *ChainStore) DeleteDbUnGenProofs(chainType common.ChainType, txIds []string) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		for _, id := range txIds {
			err := batch.BatchDeleteObj(dbUnGenProofId(chainType, id))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) DeleteUnGenProof(chainType common.ChainType, txIds ...string) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		for _, txId := range txIds {
			err := batch.BatchDeleteObj(dbUnGenProofId(chainType, txId))
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) WriteUnGenProof(chain common.ChainType, list ...*DbUnGenProof) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		for _, item := range list {
			err := batch.BatchPutObj(dbUnGenProofId(chain, item.Hash), item)
			if err != nil {
				logger.Error("put ungen Proof error:%v", err)
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) ReadUnGenProofs(chainType common.ChainType) ([]*DbUnGenProof, error) {
	iterator := cs.store.Iterator(dbUnGenProofId(chainType, ""), nil)
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

func (cs *ChainStore) WriteDepositAddrPrefix(txes []*Transaction) error {
	batch := cs.store.Batch()
	for _, tx := range txes {
		if tx.TxType == common.DepositTx {
			for _, addr := range tx.BtcFrom {
				err := batch.BatchPutObj(dbAddrPrefixTxId(addr, common.DepositTx, tx.Hash), nil)
				if err != nil {
					logger.Error("put addr prefix tx error:%v", err)
					return err
				}
			}

		}
	}
	err := batch.BatchWriteObj()
	if err != nil {
		logger.Error("put addr tx batch error:%v", err)
		return err
	}
	return nil

}

func (cs *ChainStore) WriteAddrPrefixTx(txes []*DbTx) error {
	addrPrefixTxes := txesByAddrGroup(txes)
	for addr, addrDbTxes := range addrPrefixTxes {
		err := cs.WriteTxIdsByAddr(addr, addrDbTxes)
		if err != nil {
			logger.Error("write addr txes error: %v %v", addr, err)
			return err
		}
	}
	return nil
}

func (cs *ChainStore) WriteTxIdsByAddr(addr string, txes []DbTx) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		for _, tx := range txes {
			err := batch.BatchPutObj(dbAddrPrefixTxId(addr, tx.TxType, tx.Hash), nil)
			if err != nil {
				logger.Error("put addr tx error:%v", err)
				return err
			}
		}
		return nil
	})
}

func (cs *ChainStore) ReadTxIdsByAddr(txType common.TxType, addr string) ([]string, error) {
	var txIds []string
	iterator := cs.store.Iterator(dbAddrPrefixTxId(addr, txType, ""), nil)
	defer iterator.Release()
	for iterator.Next() {
		elems := strings.Split(string(iterator.Key()), protocolSeparator)
		if len(elems) == 3 {
			txIds = append(txIds, elems[2])
		}
	}
	if err := iterator.Error(); err != nil {
		return nil, err
	}
	return txIds, nil
}

func (cs *ChainStore) DeleteRedeemSotCache(txSlot, finalizeSlot uint64, hash string) error {
	return cs.store.WrapBatch(func(batch store.IBatch) error {
		err := batch.BatchDeleteObj(dbTxSlotId(txSlot, hash))
		if err != nil {
			return err
		}
		err = batch.BatchDeleteObj(dbTxFinalizeSlotId(finalizeSlot, hash))
		if err != nil {
			return err
		}
		return nil
	})
}

func (cs *ChainStore) WriteTxSlot(txSlot uint64, tx *DbUnGenProof) error {
	return cs.store.PutObj(dbTxSlotId(txSlot, tx.Hash), tx)
}

func (cs *ChainStore) DeleteTxSlot(txSlot uint64, txHash string) error {
	return cs.store.DeleteObj(dbTxSlotId(txSlot, txHash))
}

func (cs *ChainStore) ReadAllTxBySlot(txSlot uint64) ([]*DbUnGenProof, error) {
	var txes []*DbUnGenProof
	iterator := cs.store.Iterator(dbTxSlotId(txSlot, ""), nil)
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

func (cs *ChainStore) WriteTxFinalizedSlot(txSlot uint64, tx *DbUnGenProof) error {
	return cs.store.PutObj(dbTxFinalizeSlotId(txSlot, tx.Hash), tx)
}

func (cs *ChainStore) DeleteTxFinalizedSlot(txSlot uint64, txHash string) error {
	return cs.store.DeleteObj(dbTxFinalizeSlotId(txSlot, txHash))
}

func (cs *ChainStore) ReadAllTxByFinalizedSlot(finalizedSlot uint64) ([]*DbUnGenProof, error) {
	var txes []*DbUnGenProof
	iterator := cs.store.Iterator(dbTxFinalizeSlotId(finalizedSlot, ""), nil)
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

func (cs *ChainStore) WritePendingRequest(proofId string, request *common.ProofRequest) error {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		logger.Error("proofId:%v marshal request error:%v", proofId, err)
		return err
	}
	return cs.store.PutObj(dbPendingRequestId(proofId), string(reqBytes))
}

func (cs *ChainStore) DeletePendingRequest(proofId string) error {
	return cs.store.DeleteObj(dbPendingRequestId(proofId))
}

func (cs *ChainStore) ReadAllPendingRequests() ([]*common.ProofRequest, error) {
	var requests []*common.ProofRequest
	iterator := cs.store.Iterator(dbPendingRequestId(""), nil)
	defer iterator.Release()
	for iterator.Next() {
		var value string
		err := codec.Unmarshal(iterator.Value(), &value)
		if err != nil {
			logger.Error("unmarshal tx error:%v", err)
			return nil, err
		}
		var req common.ProofRequest
		err = json.Unmarshal([]byte(value), &req)
		if err != nil {
			logger.Error("unmarshal tx error:%v", err)
			return nil, err
		}
		requests = append(requests, &req)
	}
	if err := iterator.Error(); err != nil {
		return nil, err
	}
	return requests, nil
}

func (cs *ChainStore) WriteProofResponse(resp *common.SubmitProof) error {
	return cs.store.PutObj(dbProofResponseId(resp.Id), resp)
}

func (cs *ChainStore) ReadAllProofResponse() ([]*common.SubmitProof, error) {
	var txes []*common.SubmitProof
	iterator := cs.store.Iterator(dbProofResponseId(""), nil)
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
func (cs *ChainStore) DeleteProofResponse(requestId string) error {
	return cs.store.DeleteObj(dbProofResponseId(requestId))
}

func (cs *ChainStore) ReadWorkerId() (string, bool, error) {
	exists, err := cs.store.HasObj(workerIdKey)
	if err != nil {
		return "", false, err
	}
	if !exists {
		return "", false, nil
	}
	var id string
	err = cs.store.GetObj(workerIdKey, &id)
	if err != nil {
		return "", false, err
	}
	return id, true, nil
}

func (cs *ChainStore) WriteWorkerId(id string) error {
	return cs.store.PutObj(workerIdKey, id)
}

func (cs *ChainStore) ReadZkParamVerify() (bool, error) {
	exists, err := cs.store.HasObj(zkVerifyKey)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	var verify bool
	err = cs.store.GetObj(zkVerifyKey, &verify)
	if err != nil {
		return false, err
	}
	return verify, nil
}

func (cs *ChainStore) WriteZkParamVerify(verify bool) error {
	return cs.store.PutObj(zkVerifyKey, verify)
}

func (cs *ChainStore) WriteNonce(network, addr string, nonce uint64) error {
	return cs.store.PutObj(dbAddrNonceId(network, addr), nonce)

}

func (cs *ChainStore) ReadNonce(network, addr string) (uint64, bool, error) {
	id := dbAddrNonceId(network, addr)
	exists, err := cs.store.HasObj(id)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}
	var nonce uint64
	err = cs.store.GetObj(id, &nonce)
	if err != nil {
		return 0, false, err
	}
	return nonce, true, nil
}

func (cs *ChainStore) WriteTaskTime(id string, status common.ProofStatus, value time.Time) error {
	return cs.store.PutObj(dbTaskTimeId(id, status), value)
}

func (cs *ChainStore) ReadTaskTime(id string, status common.ProofStatus) (time.Time, error) {
	var t time.Time
	err := cs.store.GetObj(dbTaskTimeId(id, status), &t)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

type TaskTime struct {
	QueueTime      time.Time `json:"queueTime"`
	GeneratingTime time.Time `json:"generatingTime"`
	EndTime        time.Time `json:"endTime"`
}

func (cs *ChainStore) ReadAllTaskTime(id string) (TaskTime, error) {
	task := TaskTime{}
	queueTime, _ := cs.ReadTaskTime(id, common.ProofQueued)
	task.QueueTime = queueTime
	generatingTime, _ := cs.ReadTaskTime(id, common.ProofGenerating)
	task.GeneratingTime = generatingTime
	endTime, _ := cs.ReadTaskTime(id, common.ProofSuccess)
	task.EndTime = endTime
	return task, nil
}

func (cs *ChainStore) WriteFinalityUpdateSlot(finalizeSlot uint64) error {
	return cs.store.PutObj(dbFinalityUpdateSlotId(finalizeSlot), finalizeSlot)
}

func (cs *ChainStore) FindFinalityUpdateNearestSlot(txSlot uint64) (uint64, bool, error) {
	var start []byte
	if txSlot-common.MaxDiffTxFinalitySlot > 0 {
		start = []byte(fmt.Sprintf("%d", txSlot-common.MaxDiffTxFinalitySlot))
	}
	iterator := cs.store.Iterator([]byte(finalityUpdateSlotPrefix), start)
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

func (cs *ChainStore) ReadSlotByHash(hash string) (uint64, bool, error) {
	dbTx, err := cs.ReadDbTxes(hash)
	if err != nil {
		return 0, false, err
	}
	if len(dbTx) == 0 {
		return 0, false, nil
	}
	slot, ok, err := cs.ReadSlotByHeight(dbTx[0].Height)
	if err != nil {
		return 0, false, err
	}
	if !ok {
		return 0, false, nil
	}
	return slot, true, nil
}

func (cs *ChainStore) ReadBtcTx(hash string) (*DbTx, bool, error) {
	dbTxes, err := cs.ReadDbTxes(hash)
	if err != nil {
		return nil, false, err
	}
	// for btc ,only one tx
	if len(dbTxes) != 1 {
		return nil, false, nil
	}
	return dbTxes[0], true, nil
}
func (cs *ChainStore) WriteLatestCheckpoint(height uint64) error {
	return cs.store.PutObj(latestCheckPointHeightKey, height)
}

func (cs *ChainStore) ReadLatestCheckPoint() (uint64, bool, error) {
	exists, err := cs.store.HasObj(latestCheckPointHeightKey)
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}
	var height uint64
	err = cs.store.GetObj(latestCheckPointHeightKey, &height)
	if err != nil {
		return 0, false, err
	}
	return height, true, nil
}

func TxHeightKeyToTxId(key []byte) (string, error) {
	// dbTxId
	split := strings.Split(string(key), protocolSeparator)
	if len(split) != 3 {
		return "", fmt.Errorf("invalid tx height key %s", string(key))
	}
	return split[2], nil
}

func DbSignatureKeyToHeight(key string) (uint64, error) {
	ids := strings.Split(key, protocolSeparator)
	if len(ids) != 2 {
		return 0, fmt.Errorf("invalid signature key %s", key)
	}
	return strconv.ParseUint(ids[1], 10, 64)
}
