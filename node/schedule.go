package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/rpc/dfinity"
	"sync"
	"time"

	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
)

type IScheduler interface {
	CheckBtcState() error
	CheckEthState() error
	CheckPreBtcState() error
	CheckBeaconState() error
	UpdateBtcCp() error
}

type Scheduler struct {
	queueManager *QueueManager
	fileStore    *FileStorage
	btcClient    *bitcoin.Client
	ethClient    *ethereum.Client
	icpClient    *dfinity.Client
	chainStore   *ChainStore
	preparedData *Prepared
	lock         sync.Mutex
}

func (s *Scheduler) init() error {
	err := s.updateBtcCp()
	if err != nil {
		logger.Warn("update btc cp error:%v", err)
	}
	return nil
}

func (s *Scheduler) updateBtcCp() error {
	cpTx, ok := s.chainStore.ReadUpdateCpTx()
	if !ok {
		logger.Warn("no find cpTx")
		return nil
	}
	hash, exists, err := s.chainStore.ReadBitcoinHash(cpTx.Height)
	if err != nil {
		logger.Error("get btc hash error: %v %v", cpTx.Height, err)
		return err
	}
	if !exists {
		return nil
	}
	exist, err := s.ethClient.IsCandidateExist(hash)
	if err != nil {
		logger.Error("check candidate exist error: %v %v", hash, err)
		return err
	}
	if exist {
		logger.Warn("update cp tx is exist,skip now blockHash:%v", hash)
		return nil
	}
	// update cp when the tx time more than 24h
	latestAddedTime, err := s.ethClient.GetCpLatestAddedTime()
	if err != nil {
		logger.Error("get cp latest added time error:%v", err)
		return err
	}
	if time.Now().Add(-24*time.Hour).Unix() <= int64(latestAddedTime) {
		logger.Warn("update cp tx is too new,skip now:%v", cpTx.Hash)
		return nil
	}
	exists, err = s.fileStore.CheckProof(NewHashStoreKey(common.BtcDepositType, cpTx.Hash))
	if err != nil {
		logger.Error("%v %v", cpTx.Hash, err)
		return err
	}
	if !exists {
		return nil
	}
	btcCpKey := NewHashStoreKey(common.BtcUpdateCpType, cpTx.Hash)
	if s.queueManager.CheckId(btcCpKey.ProofId()) {
		return nil
	}
	exists, err = s.fileStore.CheckProof(btcCpKey)
	if err != nil {
		logger.Error("check proof error:%v %v", btcCpKey.ProofId(), err)
		return err
	}
	if exists {
		err = s.fileStore.DelProof(btcCpKey)
		if err != nil {
			logger.Error("del proof error: %v %v", btcCpKey.ProofId(), err)
			return err
		}
	}
	tx, ok, err := s.chainStore.ReadBtcTx(cpTx.Hash)
	if err != nil {
		logger.Error("read btc tx error:%v %v", cpTx.Hash, err)
		return err
	}
	if !ok {
		logger.Error("no find btc tx:%v", cpTx.Hash)
		return fmt.Errorf("no find btc tx:%v", cpTx.Hash)
	}
	tx.LatestHeight = 0
	tx.CheckPointHeight = 0
	tx.GenProofNums = 0
	tx.SigSigned = false
	err = s.chainStore.WriteDbTxes(tx)
	if err != nil {
		logger.Error("write db tx error: %v", err)
		return err
	}
	err = s.chainStore.WriteUnGenProof(common.BitcoinChain, &DbUnGenProof{
		Hash:      cpTx.Hash,
		ProofType: common.BtcUpdateCpType,
		ChainType: common.BitcoinChain,
		Height:    cpTx.Height,
		TxIndex:   cpTx.TxIndex,
		Amount:    uint64(cpTx.Amount),
	})
	logger.Debug("add btcUpdateCp txHash: %v", cpTx.Hash)
	if err != nil {
		logger.Error("write unGen proof error: %v", err)
		return err
	}
	return nil
}

func (s *Scheduler) CheckBtcState() error {
	logger.Debug("start check btc state ....")
	blockCount, err := s.btcClient.GetBlockCount()
	if err != nil {
		logger.Error("get block count error:%v", err)
		return err
	}
	latestHeight, ok, err := s.chainStore.ReadBtcHeight()
	if err != nil {
		logger.Error("read latest btc height error: %v", err)
		return err
	}
	if !ok {
		logger.Warn("no find latest btc height")
		return nil
	}
	if latestHeight < uint64(blockCount-3) {
		logger.Warn("wait btc sync complete, block count:%v latestHeight:%v,skip check btc proof now", blockCount, latestHeight)
		return nil
	}

	cpHeight, ok, err := s.chainStore.ReadLatestCheckPoint()
	if err != nil {
		logger.Error("read latest checkpoint error: %v", err)
		return err
	}
	if !ok {
		logger.Warn("no find latest check point")
		return nil
	}
	unGenTxes, err := s.chainStore.ReadUnGenProofs(common.BitcoinChain)
	if err != nil {
		logger.Error("read unGen proof error:%v", err)
		return err
	}
	unSigProtect, err := s.ethClient.EnableUnsignedProtection()
	if err != nil {
		logger.Error("enable unsigned protection error:%v", err)
		return err
	}
	icpSig, exists, err := s.chainStore.ReadLatestIcpSig()
	if err != nil {
		logger.Error("read latest icp sig error:%v", err)
		return err
	}
	if !unSigProtect && !exists {
		err := s.BlockSignature()
		if err != nil {
			logger.Error("block signature error:%v", err)
			return err
		}
		return nil
	}
	if !unSigProtect && len(unGenTxes) > 0 && exists {
		if latestHeight < icpSig.Height {
			logger.Warn("unsigned protection is true,wait sync complete, latestHeight:%v,icpSig.Height:%v,skip check btc proof now", latestHeight, icpSig.Height)
			return nil
		}
		hash, ok, err := s.chainStore.ReadBitcoinHash(icpSig.Height)
		if err != nil {
			logger.Error("read bitcoin hash error:%v", err)
			return err
		}
		if !ok {
			logger.Error("no find bitcoin hash")
			return err
		}
		if icpSig.Height < latestHeight || (icpSig.Height == latestHeight && !common.StrEqual(icpSig.Hash, hash)) {
			err := s.BlockSignature()
			if err != nil {
				logger.Error("block signature error:%v", err)
				return err
			}
		}
	}
	for _, unGenTx := range unGenTxes {
		logger.Debug("bitcoin check unGen proof: %v %v", unGenTx.ProofType.Name(), unGenTx.Hash)
		proved, err := s.checkTxProved(unGenTx.ProofType, unGenTx.Hash)
		if err != nil {
			logger.Error("check tx proved error:%v %v %v", unGenTx.ProofType.Name(), unGenTx.Hash, err)
			return err
		}
		if proved {
			logger.Debug("%v %v proof exists ,delete ungen proof now", unGenTx.ProofType.Name(), unGenTx.Hash)
			err := s.delUnGenProof(common.BitcoinChain, unGenTx.Hash)
			if err != nil {
				logger.Error("delete ungen proof error:%v %v", unGenTx.Hash, err)
				return err
			}
			continue
		}
		btcDbTx, ok, err := s.chainStore.ReadBtcTx(unGenTx.Hash)
		if err != nil {
			logger.Error("read btc tx error:%v %v", unGenTx.Hash, err)
			return err
		}
		if !ok {
			logger.Warn("no find btc tx:%v", unGenTx.Hash)
			continue
		}
		//todo
		if btcDbTx.GenProofNums >= common.GenMaxRetryNums {
			logger.Warn("btc retry nums %v tx:%v num%v >= max %v,skip it now", unGenTx.ProofType.Name(), unGenTx.Hash, btcDbTx.GenProofNums, common.GenMaxRetryNums)
			err := s.delUnGenProof(common.BitcoinChain, unGenTx.Hash)
			if err != nil {
				logger.Error("delete ungen proof error:%v %v", unGenTx.Hash, err)
				return err
			}
			continue
		}
		depthOk, err := s.checkTxDepth(latestHeight, cpHeight, btcDbTx, unSigProtect)
		if err != nil {
			logger.Error("check tx height error:%v %v", unGenTx.Hash, err)
			return err
		}
		if !depthOk {
			logger.Warn("check tx depth:%v %v ,not ok", unGenTx.Hash, unGenTx.ProofType.Name())
			continue
		}
		logger.Debug("btcTx %v hash:%v amount: %v,cpHeight:%v, txHeight:%v,latestHeight: %v,unsignedProtect:%v",
			unGenTx.ProofType.Name(), unGenTx.Hash, unGenTx.Amount, btcDbTx.CheckPointHeight, btcDbTx.Height, btcDbTx.LatestHeight, unSigProtect)
		switch unGenTx.ProofType {
		case common.BtcDepositType, common.BtcUpdateCpType:
			err := s.checkBtcDepositRequest(unGenTx.ProofType, btcDbTx)
			if err != nil {
				logger.Error("check btc unGenTx request error:%v %v", unGenTx.Hash, err)
			}
		case common.BtcChangeType:
			err := s.checkBtcChangeRequest(btcDbTx)
			if err != nil {
				logger.Error("check btc unGenTx request error:%v %v", unGenTx.Hash, err)
			}
		default:
			return fmt.Errorf("invalid proof type:%v", unGenTx.ProofType.Name())
		}
	}
	logger.Debug("check btc scheduler done")
	return nil
}

func (s *Scheduler) CheckPreBtcState() error {
	latestHeight, ok, err := s.chainStore.ReadBtcHeight()
	if err != nil {
		logger.Error("read btc latest height error: %v", err)
		return err
	}
	if !ok {
		logger.Warn("no find latest btc height")
		return nil
	}
	chainIndex, ok, err := s.fileStore.CurrentBtcChainIndex()
	if err != nil {
		logger.Error("get current btc chainIndex error:%v", err)
		return err
	}
	if !ok {
		logger.Warn("no find current btc chainIndex")
		return nil
	}
	//currentIndex := s.upperRoundStartIndex(chainIndex.End)
	indexes := BlockChainPlan(chainIndex.End, latestHeight, true)
	for _, index := range indexes {
		if index.Step == common.BtcUpperDistance {
			err := s.chainUpperIndex(index.Start, index.End)
			if err != nil {
				logger.Error("check duper chainIndex error:%v", err)
				return err
			}
		} else if index.Step == common.BtcBaseDistance {
			err := s.chainStepBaseIndex(index.Start, index.End)
			if err != nil {
				logger.Error("check base chainIndex error:%v", err)
				return err
			}
		} else {
			_, err := s.tryProofRequest(NewDoubleStoreKey(common.BtcDuperRecursiveType, index.Start, index.End))
			if err != nil {
				logger.Error("try btc chain recursive proof error:%v", err)
				return err
			}
		}
	}
	cpHeight, ok, err := s.chainStore.ReadLatestCheckPoint()
	if err != nil {
		logger.Error("read latest check point error: %v", err)
		return err
	}
	if !ok {
		logger.Warn("no find latest check point")
		return nil
	}
	logger.Debug("latestBlockHeight: %v,latestCheckPoint: %v", latestHeight, cpHeight)
	currentDepthIndex, exists, err := s.fileStore.CurrentBtcCpDepthIndex(cpHeight)
	if err != nil {
		logger.Error("get current btc depth currentDepthIndex error:%v", err)
		return err
	}
	if !exists {
		return nil
	}
	storeKey := NewDoubleStoreKey(common.BtcBulkType, cpHeight, cpHeight+common.BtcCpMinDepth)
	exists, err = s.fileStore.CheckProof(storeKey)
	if err != nil {
		logger.Error("check proof error:%v %v", storeKey.ProofId(), err)
		return err
	}
	if !exists {
		_, err = s.tryProofRequest(storeKey)
		if err != nil {
			logger.Error("try btc bulk proof error:%v %v", storeKey.ProofId(), err)
			return err
		}
	}
	depthIndexes := BlockDepthPlan(cpHeight, currentDepthIndex.End, latestHeight, true)
	if len(depthIndexes) > 0 {
		storeKey := NewPrefixStoreKey(common.BtcDepthRecursiveType, cpHeight, depthIndexes[0].Start, depthIndexes[0].End)
		_, err := s.tryProofRequest(storeKey)
		if err != nil {
			logger.Error("try btc depth proof error:%v %v", storeKey.ProofId(), err)
			return err
		}
	}
	logger.Debug("check pre btc scheduler done")
	return nil
}

func (s *Scheduler) checkTxDepth(curHeight, cpHeight uint64, tx *DbTx, unSigProtect bool) (bool, error) {
	// todo need more check
	if tx.LatestHeight != 0 && tx.SigSigned {
		forked, err := s.checkIcpSig(tx.LatestHeight)
		if err != nil {
			logger.Error("check icp sig error:%v", err)
			return false, err
		}
		if forked {
			logger.Warn("find icp sig forked:%v,update latest height now:%v", tx.LatestHeight, tx.Hash)
			tx.SigSigned = false
			tx.LatestHeight = 0
		}
	}
	raised, err := s.getTxRaised(tx.Height, uint64(tx.Amount))
	if err != nil {
		logger.Error("get tx raised error:%v", err)
		return false, err
	}
	latestHeight, signed, err := s.fixLatestHeight(curHeight)
	if err != nil {
		logger.Error("get latest height error:%v", err)
		return false, err
	}
	if !unSigProtect && !signed {
		logger.Debug("no find icp block sig,update latest height now:%v %v", tx.Hash, latestHeight)
		return false, nil
	}

	return s.updateBtcTxDepth(latestHeight, cpHeight, signed, raised, tx)

}

func (s *Scheduler) checkIcpSig(height uint64) (bool, error) {
	//todo
	signature, ok, err := s.chainStore.ReadIcpSignature(height)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	hash, ok, err := s.chainStore.ReadBitcoinHash(height)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	if !common.StrEqual(hash, signature.Hash) {
		return true, nil
	}
	return false, nil

}

func (s *Scheduler) fixLatestHeight(curHeight uint64) (uint64, bool, error) {
	icpBlockSig, sigExists, err := s.chainStore.ReadLatestIcpSig()
	if err != nil {
		logger.Error("read latest icp sig error:%v", err)
		return 0, false, err
	}
	if !sigExists {
		return curHeight, false, nil
	}
	if icpBlockSig.Height < curHeight {
		return curHeight, false, nil
	}
	hash, ok, err := s.chainStore.ReadBitcoinHash(icpBlockSig.Height)
	if err != nil {
		logger.Error("read db bitcoin hash error: %v %v", icpBlockSig.Height, err)
		return 0, false, err
	}
	if !ok {
		logger.Error("no find bitcoin hash:%v", icpBlockSig.Height)
		return curHeight, false, nil
	}
	if !common.StrEqual(hash, icpBlockSig.Hash) {
		return curHeight, false, nil
	}
	return icpBlockSig.Height, true, nil
}

func (s *Scheduler) updateBtcTxDepth(curHeight, cpHeight uint64, signed, raised bool, tx *DbTx) (bool, error) {
	if tx.LatestHeight > curHeight {
		return false, nil
	}
	if tx.CheckPointHeight == 0 || tx.CheckPointHeight < cpHeight {
		tx.CheckPointHeight = cpHeight
		err := s.chainStore.WriteDbTxes(tx)
		if err != nil {
			logger.Error("write db tx error:%v", err)
			return false, err
		}
	}
	// the latestHeight on 24hour maybe expired
	expired := curHeight-tx.LatestHeight > common.BtcLatestBlockMaxDiff
	if tx.LatestHeight != 0 && expired {
		logger.Warn("txId latestHeight is expired:%v %v %v", tx.Hash, tx.LatestHeight, curHeight)
		s.removeExpiredRequest(tx)
	}
	if tx.LatestHeight == 0 || expired {
		tx.LatestHeight = curHeight
		tx.SigSigned = signed
		err := s.chainStore.WriteDbTxes(tx)
		if err != nil {
			logger.Error("write db tx error:%v", err)
			return false, err
		}
	}

	cpOk := tx.LatestHeight-tx.CheckPointHeight >= common.BtcCpMinDepth
	txMinDepth, err := s.ethClient.GetDepthByAmount(uint64(tx.Amount), raised, signed)
	if err != nil {
		logger.Error("get min tx depth error:%v", err)
		return false, err
	}
	txOk := tx.LatestHeight-tx.Height >= uint64(txMinDepth)
	if cpOk && txOk {
		return true, nil
	}
	if curHeight-tx.Height >= uint64(txMinDepth) && curHeight-tx.CheckPointHeight >= uint64(common.BtcCpMinDepth) {
		tx.LatestHeight = curHeight
		tx.SigSigned = signed
		logger.Debug("assignment check tx depth hash:%v, cpDepth:%v,txDepth:%v,height:%v,cpHeight:%v,latestHeight:%v, signed:%v",
			tx.Hash, common.BtcCpMinDepth, txMinDepth, tx.Height, tx.CheckPointHeight, tx.LatestHeight, signed)
		err := s.chainStore.WriteDbTxes(tx)
		if err != nil {
			logger.Error("update btc tx error:%v", err)
			return false, err
		}
		logger.Debug("update btc tx latestHeight: %v", curHeight)
		return true, nil
	}
	return false, nil
}

func (s *Scheduler) checkBtcDepositRequest(proofType common.ProofType, dbTx *DbTx) error {
	exists, err := s.checkBtcRequest(dbTx)
	if err != nil {
		logger.Error("check btc depth request error:%v %v", dbTx.Hash, err)
		return err
	}
	if !exists {
		return nil
	}
	storeKey := StoreKey{
		PType:     proofType,
		Hash:      dbTx.Hash,
		FIndex:    dbTx.Height,
		SIndex:    dbTx.LatestHeight,
		BlockTime: dbTx.BlockTime,
		TxIndex:   uint32(dbTx.TxIndex),
	}
	_, err = s.tryProofRequest(storeKey)
	if err != nil {
		logger.Error("try proof request error:%v %v", dbTx.Hash, err)
		return err
	}
	return nil
}

func (s *Scheduler) checkBtcChangeRequest(tx *DbTx) interface{} {
	chainOK, err := s.checkBtcRequest(tx)
	if err != nil {
		logger.Error("check btc depth request error:%v %v", tx.Hash, err)
		return err
	}
	destHash, err := s.chainStore.ReadDestHash(tx.Hash)
	if err != nil {
		logger.Error("read dest hash error:%v %v", tx.Hash, err)
		return err
	}
	exists, err := s.fileStore.CheckProof(NewHashStoreKey(common.BackendRedeemTxType, destHash))
	if err != nil {
		logger.Error("check proof error:%v %v", tx.Hash, err)
		return err
	}
	if !exists {
		_, err := s.tryProofRequest(NewHashStoreKey(common.BackendRedeemTxType, destHash))
		if err != nil {
			logger.Error("try proof request error:%v %v", tx.Hash, err)
			return err
		}
		return nil
	}
	if chainOK {
		changeKey := StoreKey{
			PType:     common.BtcChangeType,
			Hash:      tx.Hash,
			FIndex:    tx.Height,
			SIndex:    tx.LatestHeight,
			BlockTime: tx.BlockTime,
			TxIndex:   uint32(tx.TxIndex),
		}
		_, err = s.tryProofRequest(changeKey)
		if err != nil {
			logger.Error("try proof request error:%v %v", tx.Hash, err)
			return err
		}
	}
	return nil
}

func (s *Scheduler) checkBtcChainRequest(latestHeight, blockTime, txIndex uint64) (bool, error) {
	_, exists, err := s.fileStore.FindBtcChainProof(latestHeight)
	if err != nil {
		logger.Error("find btc chain proof error:%v", err)
		return false, err
	}
	if exists {
		return true, nil
	}
	chainIndex, ok, err := s.fileStore.BtcChainIndex(latestHeight)
	if err != nil {
		logger.Error("get current btc chainIndex error:%v", err)
		return false, err
	}
	if !ok {
		logger.Warn("no find current btc chainIndex")
		return false, nil
	}
	if chainIndex == latestHeight {
		return true, nil
	}
	//currentIndex := s.upperRoundStartIndex(chainIndex.End)
	chainIndexes := BlockChainPlan(chainIndex, latestHeight)
	if len(chainIndexes) == 0 {
		logger.Error("never get btc chain currentIndex")
		return false, nil
	}
	if len(chainIndexes) > 0 {
		storeKey := StoreKey{
			PType:     common.BtcDuperRecursiveType,
			FIndex:    chainIndexes[0].Start,
			SIndex:    chainIndexes[0].End,
			BlockTime: blockTime,
			TxIndex:   uint32(txIndex),
		}
		proofId := storeKey.ProofId()
		exists, err := s.fileStore.CheckProof(storeKey)
		if err != nil {
			logger.Error("check proof error:%v", proofId)
			return false, err
		}
		if !exists {
			_, err := s.tryProofRequest(storeKey)
			if err != nil {
				logger.Error("try proof request error:%v %v", proofId, err)
				return false, err
			}
			return false, nil
		}
	}
	return false, nil
}

func (s *Scheduler) checkBtcRequest(tx *DbTx) (bool, error) {
	latestHeight := tx.LatestHeight
	cpHeight := tx.CheckPointHeight
	blockTime := tx.BlockTime
	height := tx.Height
	txIndex := uint64(tx.TxIndex)
	chainExists, err := s.checkBtcChainRequest(latestHeight+1, blockTime, txIndex) // chain proof is closed boundary
	if err != nil {
		logger.Error("check btc chain request error:%v %v", latestHeight+1, err)
		return false, err
	}

	txDepthExists, err := s.checkTxDepthRequest(height, latestHeight, blockTime, txIndex)
	if err != nil {
		logger.Error("check depth proof request error:%v %v %v", blockTime, latestHeight, err)
		return false, err
	}
	txCpDepthExists, err := s.checkCpDepthProofRequest(cpHeight, latestHeight, blockTime, txIndex)
	if err != nil {
		logger.Error("check depth proof request error:%v %v %v", cpHeight, latestHeight, err)
		return false, err
	}
	timestampKey := StoreKey{
		PType:     common.BtcTimestampType,
		FIndex:    height,
		SIndex:    latestHeight,
		BlockTime: blockTime,
		TxIndex:   uint32(txIndex),
	}
	exists, err := s.fileStore.CheckProof(timestampKey)
	if err != nil {
		logger.Error("check proof error:%v", timestampKey.ProofId())
		return false, err
	}
	if !exists {
		_, err := s.tryProofRequest(timestampKey)
		if err != nil {
			logger.Error("try proof request error:%v %v", timestampKey.ProofId(), err)
			return false, err
		}

	}
	return exists && chainExists && txDepthExists && txCpDepthExists, nil
}

func (s *Scheduler) checkCpDepthProofRequest(depthHeight, latestHeight, blockTime, txIndex uint64) (bool, error) {
	step := latestHeight - depthHeight
	if step <= common.BtcCpMinDepth {
		storeKey := StoreKey{
			PType:     common.BtcBulkType,
			FIndex:    depthHeight,
			SIndex:    latestHeight,
			BlockTime: blockTime,
			TxIndex:   uint32(txIndex),
		}
		proofId := storeKey.ProofId()
		exists, err := s.fileStore.CheckProof(storeKey)
		if err != nil {
			logger.Error("check proof error:%v %v", proofId, err)
			return false, err
		}
		if !exists {
			_, err := s.tryProofRequest(storeKey)
			if err != nil {
				logger.Error("try proof request error:%v %v", proofId, err)
				return false, err
			}
		}
		return exists, nil
	} else if step > common.BtcCpMinDepth {
		ok, err := s.checkDepthRecursive(depthHeight, latestHeight, common.BtcCpMinDepth, blockTime, txIndex)
		if err != nil {
			logger.Error("check depth proof error:%v %v", depthHeight, err)
			return false, err
		}
		return ok, nil
	} else {
		logger.Warn("never should happen,check cp depth %v %v", depthHeight, latestHeight)
		return false, nil
	}

}

func (s *Scheduler) checkTxDepthRequest(depthHeight, latestHeight, blockTime, txIndex uint64) (bool, error) {
	step := latestHeight - depthHeight
	if step >= common.BtcTxMinDepth && step <= common.BtcTxUnitMaxDepth {
		//storeKey := NewDoubleStoreKey(common.BtcBulkType, depthHeight, latestHeight)
		storeKey := StoreKey{
			PType:     common.BtcBulkType,
			FIndex:    depthHeight,
			SIndex:    latestHeight,
			BlockTime: blockTime,
			TxIndex:   uint32(txIndex),
		}
		exists, err := s.fileStore.CheckProof(storeKey)
		if err != nil {
			logger.Error("check proof error:%v %v", storeKey.ProofId(), err)
			return false, err
		}
		if !exists {
			_, err := s.tryProofRequest(storeKey)
			if err != nil {
				logger.Error("try proof request error:%v %v", storeKey.ProofId(), err)
				return false, err
			}
		}
		return exists, nil
	} else if step > common.BtcTxUnitMaxDepth {
		ok, err := s.checkDepthRecursive(depthHeight, latestHeight, common.BtcTxUnitMaxDepth, blockTime, txIndex)
		if err != nil {
			logger.Error("check depth recursive error:%v %v %v", depthHeight, latestHeight, err)
			return false, err
		}
		return ok, nil
	} else {
		logger.Warn("never should happen check tx depth: %v %v", depthHeight, latestHeight)
		return false, nil
	}
}

func (s *Scheduler) checkDepthRecursive(depthHeight uint64, latestHeight uint64, maxUnitDepth, blockTime, txIndex uint64) (bool, error) {
	_, exists, err := s.fileStore.FindDepthProof(depthHeight, latestHeight)
	if err != nil {
		logger.Error("check depth proof error:%v %v", depthHeight, err)
		return false, err
	}
	if exists {
		return exists, nil
	}
	// for depth genesis proof(btcBulk proof)
	bulkStoreKey := StoreKey{
		PType:     common.BtcBulkType,
		FIndex:    depthHeight,
		SIndex:    depthHeight + maxUnitDepth,
		BlockTime: blockTime,
		TxIndex:   uint32(txIndex),
	}
	exists, err = s.fileStore.CheckProof(bulkStoreKey)
	if err != nil {
		logger.Error("check proof error:%v %v", bulkStoreKey.ProofId(), err)
		return false, err
	}
	if !exists {
		_, err := s.tryProofRequest(bulkStoreKey)
		if err != nil {
			logger.Error("try proof request error:%v %v", bulkStoreKey.ProofId(), err)
			return false, err
		}
		return false, nil
	}
	// for middle proof
	index, ok, err := s.fileStore.BtcDepthIndex(depthHeight, maxUnitDepth, latestHeight)
	if err != nil {
		logger.Error("check index error:%v %v", bulkStoreKey, err)
		return false, err
	}
	if !ok {
		return false, nil
	}
	depthIndex := BlockDepthPlan(depthHeight, index, latestHeight)
	if len(depthIndex) > 0 {
		storeKey := StoreKey{
			PType:     common.BtcDepthRecursiveType,
			Prefix:    depthIndex[0].Genesis,
			FIndex:    depthIndex[0].Start,
			SIndex:    depthIndex[0].End,
			BlockTime: blockTime,
			TxIndex:   uint32(txIndex),
		}
		exists, err := s.fileStore.CheckProof(storeKey)
		if err != nil {
			logger.Error("check proof error:%v %v", storeKey.ProofId(), err)
			return false, err
		}
		if !exists {
			_, err := s.tryProofRequest(storeKey)
			if err != nil {
				logger.Error("try proof request error:%v %v", storeKey.ProofId(), err)
				return false, err
			}
		}
	}
	return false, nil
}

func (s *Scheduler) btcStateRollback(forkHeight uint64) error {
	logger.Debug("btc scheduler roll back to height: %v", forkHeight)
	filterReqs := s.queueManager.RemoveRequest(func(request *common.ProofRequest) bool {
		if common.IsBtcProofType(request.ProofType) && forkHeight <= request.SIndex {
			return true
		}
		return false
	})
	for _, req := range filterReqs {
		logger.Warn("requests queue find forked  proof request: %v", req.ProofId())
		s.removeRequest(req.ProofId())
	}
	pendingRequests := s.queueManager.FilterPending(func(request *common.ProofRequest) bool {
		if common.IsBtcProofType(request.ProofType) && forkHeight <= request.SIndex {
			return true
		}
		return false
	})
	for _, req := range pendingRequests {
		logger.Warn("pending proof find forked proof request: %v", req.ProofId())
		s.removeRequest(req.ProofId())
	}
	return nil
}

func (s *Scheduler) tryProofRequest(key StoreKey) (bool, error) {
	proofId := common.GenKey(key.PType, key.Prefix, key.FIndex, key.SIndex, key.Hash).String()
	exists := s.queueManager.CheckId(proofId)
	if exists {
		logger.Debug("proof request exists: %v", proofId)
		return false, nil
	}
	exists, err := s.fileStore.CheckProof(key)
	if err != nil {
		logger.Error("check proof error:%v %v", proofId, err)
		return false, err
	}
	if exists {
		return false, nil
	}
	if common.IsBtcProofType(key.PType) {
		defer s.lock.Unlock()
		s.lock.Lock()
	}
	data, ok, err := GenRequestData(s.preparedData, key.PType, key.FIndex, key.SIndex, key.Hash, key.Prefix, key.isCp)
	if err != nil {
		logger.Error("get request data error:%v %v", proofId, err)
		return false, err
	}
	if !ok {
		return false, nil
	}
	req := common.NewProofRequest(key.PType, data, key.Prefix, key.FIndex, key.SIndex, key.Hash, key.BlockTime, key.TxIndex)
	if s.queueManager.CheckId(proofId) {
		return true, nil
	}
	s.queueManager.StoreId(proofId)
	s.queueManager.PushRequest(req)
	if req.ProofType == common.BtcDepositType {
		err := s.chainStore.IncrDepositCount(req.FIndex)
		if err != nil {
			logger.Error("increment deposit count error:%v %v", proofId, err)
			//return nil, false, err
		}
	}
	logger.Info("success add request to queue :%v", proofId)
	err = s.UpdateProofStatus(req, common.ProofQueued)
	if err != nil {
		logger.Error("update proof status error:%v %v", proofId, err)
	}
	return true, nil
}

func (s *Scheduler) CheckEthState() error {
	logger.Debug("check eth scheduler now  ....")
	unGenProofs, err := s.chainStore.ReadUnGenProofs(common.EthereumChain)
	if err != nil {
		logger.Error("read all ungen proof ids error: %v", err)
		return err
	}
	for _, item := range unGenProofs {
		txHash := item.Hash
		logger.Debug("start check Redeem proof tx: %v %v %v", txHash, item.Height, item.TxIndex)
		proved, err := s.checkTxProved(common.RedeemTxType, txHash)
		if err != nil {
			logger.Error("check tx proved error: %v %v", txHash, err)
			return err
		}
		dbTx, ok, err := s.chainStore.ReadRedeemTx(txHash)
		if err != nil {
			logger.Error("read redeem tx error: %v %v", txHash, err)
			return err
		}
		if !ok {
			logger.Warn("read redeem tx error: %v %v", txHash, err)
			return nil
		}
		if dbTx.GenProofNums >= common.GenMaxRetryNums {
			logger.Warn("eth retry nums %v tx:%v num%v >= max %v,skip it now", dbTx.ProofType.Name(), dbTx.Hash, dbTx.GenProofNums, common.GenMaxRetryNums)
			err := s.delUnGenProof(common.EthereumChain, dbTx.Hash)
			if err != nil {
				logger.Error("delete ungen proof error:%v %v", dbTx.Hash, err)
				return err
			}
			continue
		}

		blockTime := dbTx.BlockTime
		txIndex := uint32(dbTx.TxIndex)
		if proved {
			backendRedeemStoreKey := StoreKey{
				PType:     common.BackendRedeemTxType,
				Hash:      txHash,
				BlockTime: blockTime,
				TxIndex:   txIndex,
			}
			exists, err := s.fileStore.CheckProof(backendRedeemStoreKey)
			if err != nil {
				logger.Error("check proof error: %v", err)
				return err
			}
			if !exists {
				_, err := s.tryProofRequest(backendRedeemStoreKey)
				if err != nil {
					logger.Error("try proof request error: %v", err)
					return err
				}
			} else {
				logger.Debug("Redeem proof exist now,delete cache: %v", txHash)
				err := s.delUnGenProof(common.EthereumChain, txHash)
				if err != nil {
					logger.Error("delete ungen proof error:%v %v", txHash, err)
					return err
				}
			}
			continue
		}
		txSlot, ok, err := s.chainStore.ReadSlotByHash(txHash)
		if err != nil {
			logger.Error("get txSlot error: %v %v", err, txHash)
			return err
		}
		if !ok {
			logger.Warn("no find  tx %v beacon slot", txHash)
			continue
		}
		finalizedSlot, ok, err := s.fileStore.GetTxFinalizedSlot(txSlot)
		if err != nil {
			logger.Error("get near tx slot finalized slot error: %v", err)
			return err
		}
		if !ok {
			logger.Warn("no find near %v tx slot finalized slot", txSlot)
			continue
		}
		txInEth2Key := StoreKey{
			PType:     common.TxInEth2Type,
			Hash:      txHash,
			BlockTime: blockTime,
			TxIndex:   txIndex,
		}
		exists, err := s.fileStore.CheckProof(txInEth2Key)
		if err != nil {
			logger.Error("check tx proof error: %v", err)
			return err
		}
		if !exists {
			_, err := s.tryProofRequest(txInEth2Key)
			if err != nil {
				logger.Error("try proof request error: %v", err)
				return err
			}
		}
		beaconHeaderStoreKey := StoreKey{
			PType:     common.BeaconHeaderType,
			FIndex:    txSlot,
			SIndex:    finalizedSlot,
			BlockTime: blockTime,
			TxIndex:   txIndex,
		}

		exists, err = s.fileStore.CheckProof(beaconHeaderStoreKey)
		if err != nil {
			logger.Error("check block header proof error: %v", err)
			return err
		}
		if !exists {
			err := s.chainStore.WriteTxSlot(txSlot, item)
			if err != nil {
				logger.Error("write tx slot error: %v %v %v", txHash, txSlot, err)
				return err
			}
			_, err = s.tryProofRequest(beaconHeaderStoreKey)
			if err != nil {
				logger.Error("try proof request error: %v", err)
				return err
			}
		}
		beaconHeaderFinalityStoreKey := StoreKey{
			PType:     common.BeaconHeaderFinalityType,
			FIndex:    finalizedSlot,
			BlockTime: blockTime,
			TxIndex:   txIndex,
		}
		exists, err = s.fileStore.CheckProof(beaconHeaderFinalityStoreKey)
		if err != nil {
			logger.Error("check block header finality proof error: %v %v", finalizedSlot, err)
			return err
		}
		if !exists {
			err := s.chainStore.WriteTxFinalizedSlot(finalizedSlot, item)
			if err != nil {
				logger.Error("write tx finalized slot error: %v %v %v", finalizedSlot, txHash, err)
				return err
			}
			_, err = s.tryProofRequest(beaconHeaderFinalityStoreKey)
			if err != nil {
				logger.Error("try proof request error: %v", err)
				return err
			}
			continue
		}
		redeemStoreKey := StoreKey{
			PType:     common.RedeemTxType,
			Hash:      txHash,
			BlockTime: blockTime,
			TxIndex:   txIndex,
		}
		_, err = s.tryProofRequest(redeemStoreKey)
		if err != nil {
			logger.Error("try proof request error: %v", err)
			return err
		}
	}
	logger.Debug("check eth scheduler done")
	return nil
}

func (s *Scheduler) CheckBeaconState() error {
	logger.Debug("check beacon scheduler ...")
	//beacon recursive index
	dutyIndexes, err := s.fileStore.NeedDutyIndexes()
	if err != nil {
		logger.Error("get duty indexes error:%v", err)
		return err
	}
	if len(dutyIndexes) > 0 {
		logger.Debug("beacon recursive proof: %v", dutyIndexes[0])
		syncComRecursiveKey := NewHeightStoreKey(common.SyncComDutyType, dutyIndexes[0])
		_, err = s.tryProofRequest(syncComRecursiveKey)
		if err != nil {
			logger.Error("try sync committee recursive proof error:%v", err)
			return err
		}
	}
	// beacon unit proof
	unitIndexes, err := s.fileStore.NeedGenUnitProofIndexes()
	if err != nil {
		logger.Error("get unit indexes error:%v", err)
		return err
	}
	for _, index := range unitIndexes {
		unitStorageKey := NewHeightStoreKey(common.SyncComUnitType, index)
		_, err := s.tryProofRequest(unitStorageKey)
		if err != nil {
			logger.Error("try sync committee unit proof error:%v", err)
			return err
		}
	}

	//beacon outer
	outerIndexes, err := s.fileStore.GenOuterIndexes()
	if err != nil {
		logger.Error("get outer indexes error:%v", err)
		return err
	}
	for _, index := range outerIndexes {
		syncComOuterKey := NewHeightStoreKey(common.SyncComOuterType, index)
		_, err := s.tryProofRequest(syncComOuterKey)
		if err != nil {
			logger.Error("try sync committee outer proof error:%v", err)
			return err
		}
	}
	// beacon syncCommittee inner proof
	syncComInnerIndexes, err := s.fileStore.SyncComInnerIndexes()
	if err != nil {
		logger.Error("get sync committee inner indexes error:%v", err)
		return err
	}
	for _, index := range syncComInnerIndexes {
		syncComInnerKey := NewPrefixStoreKey(common.SyncComInnerType, index.Prefix, index.Start, 0)
		_, err := s.tryProofRequest(syncComInnerKey)
		if err != nil {
			logger.Error("try sync committee inner proof error:%v", err)
			return err
		}
	}
	logger.Debug("check beacon scheduler done")
	return nil
}

func (s *Scheduler) chainStepBaseIndex(start, end uint64) error {
	_, err := s.tryProofRequest(NewDoubleStoreKey(common.BtcDuperRecursiveType, start, end))
	if err != nil {
		logger.Error("try btc duper proof error:%v", err)
		return err
	}
	_, err = s.tryProofRequest(NewDoubleStoreKey(common.BtcBaseType, start, end))
	if err != nil {
		logger.Error("try btc base proof error:%v", err)
		return err
	}
	return nil
}

func (s *Scheduler) chainUpperIndex(start, end uint64) error {
	if start%common.BtcUpperDistance != 0 || end%common.BtcUpperDistance != 0 {
		return fmt.Errorf("start or end is not multiple of 2016")
	}
	_, err := s.tryProofRequest(NewDoubleStoreKey(common.BtcDuperRecursiveType, start, end))
	if err != nil {
		logger.Error("try btc duper proof error:%v", err)
		return err
	}
	_, err = s.tryProofRequest(NewDoubleStoreKey(common.BtcUpperType, start, end))
	if err != nil {
		logger.Error("try btc upper proof error:%v", err)
		return err
	}

	middleIndexes, err := s.fileStore.BtcMiddleIndexes(start, end)
	if err != nil {
		logger.Error("get need btc middle index error:%v", err)
		return err
	}
	mCount := 0
	for _, index := range middleIndexes {
		_, err := s.tryProofRequest(NewDoubleStoreKey(common.BtcMiddleType, index.Start, index.End))
		if err != nil {
			logger.Error("try btc middle proof error:%v", err)
			return err
		}
		mCount++
		if mCount*common.BtcMiddleDistance >= common.BtcUpperDistance {
			break
		}
	}

	baseIndexes, err := s.fileStore.BtcBaseIndexes(start, end)
	if err != nil {
		logger.Error("get need btc base index error:%v", err)
		return err
	}
	bCount := 0
	for _, index := range baseIndexes {
		_, err := s.tryProofRequest(NewDoubleStoreKey(common.BtcBaseType, index, index+common.BtcBaseDistance))
		if err != nil {
			logger.Error("try btc base proof error:%v %v", index, err)
			return err
		}
		bCount++
		if bCount*common.BtcBaseDistance >= common.BtcUpperDistance {
			break
		}
	}
	logger.Debug("check pre btc state done")
	return nil
}

func (s *Scheduler) UpdateProofStatus(req *common.ProofRequest, status common.ProofStatus) error {
	if req.ProofType == common.BtcDepositType || req.ProofType == common.BtcChangeType || req.ProofType == common.RedeemTxType {
		err := s.chainStore.UpdateProof(req.Hash, "", req.ProofType, status)
		if err != nil {
			logger.Error("update Proof status error:%v %v", req.ProofId(), err)
			return err
		}
	}
	return nil
}

func (s *Scheduler) checkTxProved(proofType common.ProofType, hash string) (bool, error) {
	ok, err := s.fileStore.CheckProof(NewHashStoreKey(proofType, hash))
	if err != nil {
		logger.Error("check proof error:%v %v", hash, err)
		return false, err
	}
	if ok {
		return true, nil
	}
	switch proofType {
	case common.BtcChangeType:
		_, exists, err := s.chainStore.ReadUpdateUtxoDest(hash)
		if err != nil {
			logger.Error("check utxo error: %v %v", hash, err)
			return false, err
		}
		if exists {
			return true, nil
		}
		utxo, err := s.ethClient.GetUtxo(hash)
		if err != nil {
			logger.Error("check utxo error: %v %v", hash, err)
			return false, nil
		}
		return utxo.IsChangeConfirmed, nil
	case common.BtcDepositType:
		_, exists, err := s.chainStore.GetDestHash(hash)
		if err != nil {
			logger.Error("check deposit tx utxo error: %v %v", hash, err)
			return false, err
		}
		if exists {
			return true, nil
		}
		utxo, err := s.ethClient.GetUtxo(hash)
		if err != nil {
			logger.Warn("check deposit tx utxo error: %v %v", hash, err)
			return false, err
		}
		if TxIdIsEmpty(utxo.Txid) {
			return false, nil
		}
		return true, nil
	case common.RedeemTxType:
		exists, err := s.btcClient.CheckTxOnChain(hash)
		if err != nil {
			logger.Error("check btc tx error:%v %v", hash, err)
			return false, err
		}
		return exists, nil

	default:
		return false, nil

	}
}

func (s *Scheduler) BlockSignature() error {
	balance, err := s.icpClient.IcpBalance()
	if err != nil {
		logger.Error("get icp balance error:%v", err)
		//return err
	}
	if balance < 250_000_000_000 { // todo
		logger.Error("icp balance is not enough:%v,maybe need deposit %v", balance, s.icpClient.WalletInfo())
	}

	sig, err := s.icpClient.BlockSignatureWithCycle()
	if err != nil {
		logger.Error("get block sig error:%v", err)
		return err
	}
	if sig.Signature == "" {
		logger.Warn("block signature is empty:%v", sig.Height)
		return nil
	}
	hash, ok, err := s.chainStore.ReadBitcoinHash(uint64(sig.Height))
	if err != nil {
		logger.Error("read db bitcoin hash error: %v %v", sig.Height, err)
		return err
	}
	if !ok {
		logger.Error("no find bitcoin hash:%v", sig.Height)
		return nil
	}
	if !common.StrEqual(hash, sig.Hash) {
		logger.Warn("find icp hash %v != db hash %v,height:%v", sig.Hash, hash, sig.Height)
		return nil
	}
	dbIcpSignature := DbIcpSignature{Height: uint64(sig.Height), Hash: sig.Hash, Signature: sig.Signature}
	err = s.chainStore.WriteLatestIcpSig(dbIcpSignature)
	if err != nil {
		logger.Error("write latest icp sig error:%v", err)
		return err
	}
	err = s.chainStore.WriteIcpSignature(uint64(sig.Height), dbIcpSignature)
	if err != nil {
		logger.Error("write icp sig error:%v", err)
		return err
	}
	logger.Debug("success get icp block signature:%v %v %v", sig.Height, sig.Hash, sig.Signature)
	return nil
}
func (s *Scheduler) getTxRaised(height, amount uint64) (bool, error) {
	hash, ok, err := s.chainStore.ReadBitcoinHash(height)
	if err != nil {
		logger.Error("read bitcoin hash error:%v", err)
		return false, err
	}
	if !ok {
		fmt.Errorf("no find bitcoin hash")
		return false, nil
	}
	raised, err := s.ethClient.GetRaised(hash, amount)
	if err != nil {
		logger.Error("get raised error:%v", err)
		return false, err
	}
	if raised {
		return true, nil
	}
	count, ok, err := s.chainStore.ReadDepositCount(height)
	if err != nil {
		logger.Error("read deposit count error:%v", err)
		return false, err
	}
	if ok && count >= 10 { //todo
		return true, nil
	}
	return raised, nil
}

// when update a latestHeight of tx ,need to remove the expired request
func (s *Scheduler) removeExpiredRequest(tx *DbTx) error {
	expiredRequests := s.queueManager.RemoveRequest(func(value *common.ProofRequest) bool {
		switch value.ProofType {
		case common.BtcDepositType, common.BtcUpdateCpType, common.BtcChangeType:
			if common.StrEqual(value.Hash, tx.Hash) {
				return true
			}
		case common.BtcBulkType:
			step := tx.LatestHeight - tx.Height
			if step >= common.BtcTxMinDepth && step < common.BtcTxUnitMaxDepth {
				if value.FIndex == tx.Height && value.SIndex == tx.Height+step {
					return true
				}
			}
		case common.BtcTimestampType:
			if value.FIndex == tx.Height && value.SIndex == tx.LatestHeight {
				return true
			}
		case common.BtcDepthRecursiveType:
			if value.Prefix == tx.Height && value.SIndex == tx.LatestHeight { // tx depth
				return true

			} else if value.Prefix == tx.CheckPointHeight && value.SIndex == tx.LatestHeight { // cp depth
				return true
			}
		case common.BtcDuperRecursiveType:
			if value.SIndex == tx.LatestHeight+1 { // tx chain depth
				return true
			}
		default:
			return false
		}
		return false

	})
	for _, req := range expiredRequests {
		logger.Warn("remove expired request:%v %v %v", tx.Hash, tx.ProofType.Name(), req.ProofId())
		s.removeRequest(req.ProofId())
	}
	return nil
}

func (s *Scheduler) upperRoundStartIndex(height uint64) uint64 {
	index := upperRoundStartIndex(height)
	if index < s.fileStore.btcGenesisHeight {
		return s.fileStore.btcGenesisHeight
	}
	return index
}

func (s *Scheduler) Locks() func() {
	s.lock.Lock()
	return func() {
		s.lock.Unlock()
	}
}

func (s *Scheduler) delUnGenProof(chain common.ChainType, hash string) error {
	err := s.chainStore.DeleteUnGenProof(chain, hash)
	if err != nil {
		logger.Error("delete ungen proof error: %v", err)
		return err
	}
	return err
}

func (s *Scheduler) addRequestToPending(req *common.ProofRequest) {
	logger.Debug("add request to pending queue: %v", req.ProofId())
	s.queueManager.AddPending(req.ProofId(), req)
}

func (s *Scheduler) removeRequest(proofId string) {
	s.queueManager.DeletePending(proofId)
	s.queueManager.DeleteId(proofId)
}

func (s *Scheduler) PendingProofRequest() []*common.ProofRequest {
	return s.queueManager.ListRequest()
}

func (s *Scheduler) PendingRequest() []*common.ProofRequest {
	return s.queueManager.FilterPending(func(value *common.ProofRequest) bool {
		return true
	})
}

func NewScheduler(filestore *FileStorage, store store.IStore, preparedData *Prepared,
	icpClient *dfinity.Client, btcClient *bitcoin.Client, ethClient *ethereum.Client) (*Scheduler, error) {
	return &Scheduler{
		queueManager: NewQueueManager(),
		fileStore:    filestore,
		chainStore:   NewChainStore(store),
		preparedData: preparedData,
		btcClient:    btcClient,
		ethClient:    ethClient,
		icpClient:    icpClient,
	}, nil
}
