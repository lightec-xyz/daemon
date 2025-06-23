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
	BlockSignature() error
}

type Scheduler struct {
	proofQueue   *ArrayQueue
	pendingQueue *PendingQueue
	fileStore    *FileStorage
	btcClient    *bitcoin.Client
	ethClient    *ethereum.Client
	icpClient    *dfinity.Client
	chainStore   *ChainStore
	cache        *cache
	preparedData *Prepared
	lock         sync.Mutex
}

func (s *Scheduler) init() error {
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
	if s.cache.Check(btcCpKey.ProofId()) {
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
	if latestHeight < uint64(blockCount-10) {
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
		if btcDbTx.GenProofNums >= common.GenMaxRetryNums {
			//todo
			logger.Warn("btc retry nums %v tx:%v num%v >= max %v,skip it now", unGenTx.ProofType.Name(), unGenTx.Hash, btcDbTx.GenProofNums, common.GenMaxRetryNums)
			err := s.delUnGenProof(common.BitcoinChain, unGenTx.Hash)
			if err != nil {
				logger.Error("delete ungen proof error:%v %v", unGenTx.Hash, err)
				return err
			}
			continue
		}
		depthOk, err := s.checkTxDepth(latestHeight, cpHeight, btcDbTx)
		if err != nil {
			logger.Error("check tx height error:%v %v", unGenTx.Hash, err)
			return err
		}
		if !depthOk {
			logger.Warn("check tx depth:%v %v ,not ok", unGenTx.Hash, unGenTx.ProofType.Name())
			continue
		}

		logger.Debug("btcTx %v hash:%v amount: %v,cpHeight:%v, txHeight:%v,latestHeight: %v", unGenTx.ProofType.Name(), unGenTx.Hash, unGenTx.Amount,
			btcDbTx.CheckPointHeight, btcDbTx.Height, btcDbTx.LatestHeight)
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

func (s *Scheduler) checkTxDepth(curHeight, cpHeight uint64, tx *DbTx) (bool, error) {
	sig, sigExists, err := s.chainStore.ReadLatestIcpSig()
	if err != nil {
		logger.Error("read latest icp sig error:%v", err)
		return false, err
	}
	raised, err := s.getTxRaised(tx.Height, uint64(tx.Amount))
	if err != nil {
		logger.Error("get tx raised error:%v", err)
		return false, err
	}
	// check icp block signature if is ok
	if sigExists && curHeight-sig.Height <= 3 {
		return s.updateBtcTxDepth(sig.Height, cpHeight, true, raised, tx)
	} else {
		return s.updateBtcTxDepth(curHeight, cpHeight, false, raised, tx)
	}
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
		logger.Debug("check tx depth hash:%v, cpDepth:%v,txDepth:%v,height:%v,cpHeight:%v,latestHeight:%v,", tx.Hash, common.BtcCpMinDepth, txMinDepth, tx.Height, tx.CheckPointHeight, tx.LatestHeight)
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
	exists, err := s.checkBtcRequest(dbTx.LatestHeight, dbTx.CheckPointHeight, dbTx.Height)
	if err != nil {
		logger.Error("check btc depth request error:%v %v", dbTx.Hash, err)
		return err
	}
	if !exists {
		return nil
	}
	_, err = s.tryProofRequest(NewBtcStoreKey(proofType, dbTx.Height, dbTx.LatestHeight, dbTx.Hash))
	if err != nil {
		logger.Error("try proof request error:%v %v", dbTx.Hash, err)
		return err
	}
	return nil
}

func (s *Scheduler) checkBtcChangeRequest(tx *DbTx) interface{} {
	chainOK, err := s.checkBtcRequest(tx.LatestHeight, tx.CheckPointHeight, tx.Height)
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
		_, err = s.tryProofRequest(NewBtcStoreKey(common.BtcChangeType, tx.Height, tx.LatestHeight, tx.Hash))
		if err != nil {
			logger.Error("try proof request error:%v %v", tx.Hash, err)
			return err
		}
	}
	return nil
}

func (s *Scheduler) checkBtcChainRequest(latestHeight uint64) (bool, error) {
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
		storeKey := NewDoubleStoreKey(common.BtcDuperRecursiveType, chainIndexes[0].Start, chainIndexes[0].End)
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

func (s *Scheduler) checkBtcRequest(latestHeight, cpHeight, height uint64) (bool, error) {
	chainExists, err := s.checkBtcChainRequest(latestHeight + 1) // chain proof is closed boundary
	if err != nil {
		logger.Error("check btc chain request error:%v %v", latestHeight+1, err)
		return false, err
	}

	txDepthExists, err := s.checkTxDepthRequest(height, latestHeight)
	if err != nil {
		logger.Error("check depth proof request error:%v %v %v", height, latestHeight, err)
		return false, err
	}
	txCpDepthExists, err := s.checkCpDepthProofRequest(cpHeight, latestHeight)
	if err != nil {
		logger.Error("check depth proof request error:%v %v %v", cpHeight, latestHeight, err)
		return false, err
	}
	timestampKey := NewDoubleStoreKey(common.BtcTimestampType, height, latestHeight)
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

func (s *Scheduler) checkCpDepthProofRequest(depthHeight, latestHeight uint64) (bool, error) {
	step := latestHeight - depthHeight
	if step <= common.BtcCpMinDepth {
		storeKey := NewDoubleStoreKey(common.BtcBulkType, depthHeight, latestHeight)
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
		ok, err := s.checkDepthRecursive(depthHeight, latestHeight, common.BtcCpMinDepth)
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

func (s *Scheduler) checkTxDepthRequest(depthHeight, latestHeight uint64) (bool, error) {
	step := latestHeight - depthHeight
	if step >= common.BtcTxMinDepth && step <= common.BtcTxUnitMaxDepth {
		storeKey := NewDoubleStoreKey(common.BtcBulkType, depthHeight, latestHeight)
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
		ok, err := s.checkDepthRecursive(depthHeight, latestHeight, common.BtcTxUnitMaxDepth)
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

func (s *Scheduler) checkDepthRecursive(depthHeight uint64, latestHeight uint64, maxUnitDepth uint64) (bool, error) {
	_, exists, err := s.fileStore.FindDepthProof(depthHeight, latestHeight)
	if err != nil {
		logger.Error("check depth proof error:%v %v", depthHeight, err)
		return false, err
	}
	if exists {
		return exists, nil
	}
	// for depth genesis proof(btcBulk proof)
	bulkStoreKey := NewDoubleStoreKey(common.BtcBulkType, depthHeight, depthHeight+maxUnitDepth)
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
		storeKey := NewTxDepthStoreKey(common.BtcDepthRecursiveType, depthIndex[0].Genesis, depthIndex[0].Start, depthIndex[0].End)
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
	err := s.proofQueue.Filter(func(request *common.ProofRequest) (bool, error) {
		if common.IsBtcProofType(request.ProofType) && forkHeight <= request.SIndex {
			s.cache.Delete(request.ProofId())
			logger.Warn("request queue find unmatched proof request: %v", request.ProofId())
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		logger.Error("remove queue error:%v", err)
		return err
	}
	s.pendingQueue.Iterator(func(request *common.ProofRequest) error {
		if common.IsBtcProofType(request.ProofType) && forkHeight <= request.SIndex {
			logger.Warn("pending queue find unmatched proof request: %v", request.ProofId())
			s.pendingQueue.Delete(request.ProofId())
			return nil
		}
		return nil
	})
	return nil
}

func (s *Scheduler) tryProofRequest(key StoreKey) (bool, error) {
	proofId := common.GenKey(key.PType, key.Prefix, key.FIndex, key.SIndex, key.Hash).String()
	exists := s.cache.Check(proofId)
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
	req := common.NewProofRequest(key.PType, data, key.Prefix, key.FIndex, key.SIndex, key.Hash)
	if s.cache.Check(proofId) {
		return true, nil
	}
	s.cache.Store(proofId, nil)
	s.proofQueue.Push(req)
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
		if proved {
			backendRedeemStoreKey := NewHashStoreKey(common.BackendRedeemTxType, txHash)
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
		txInEth2Key := NewHashStoreKey(common.TxInEth2Type, txHash)
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
		beaconHeaderStoreKey := NewDoubleStoreKey(common.BeaconHeaderType, txSlot, finalizedSlot)
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
		beaconHeaderFinalityStoreKey := NewHeightStoreKey(common.BeaconHeaderFinalityType, finalizedSlot)
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
		redeemStoreKey := NewHashStoreKey(common.RedeemTxType, txHash)
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
		exists, err := s.btcClient.CheckTx(hash)
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
	sig, err := s.icpClient.BlockSignature()
	if err != nil {
		logger.Error("get block sig error:%v", err)
		return err
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
	if ok && count >= 21 { //todo
		return true, nil
	}
	return raised, nil
}

// when update a latestHeight of tx ,need to remove the expired request
func (s *Scheduler) removeExpiredRequest(tx *DbTx) error {
	expiredRequests := s.proofQueue.Remove(func(value *common.ProofRequest) bool {
		switch value.ProofType {
		case common.BtcBulkType:
			step := tx.LatestHeight - tx.Height
			if step >= common.BtcTxMinDepth && step < common.BtcTxUnitMaxDepth {
				return true
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
		logger.Warn("remove expired request:%v %v", tx.Hash, req.ProofId())
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
	s.pendingQueue.Add(req.ProofId(), req)
}

func (s *Scheduler) removeRequest(proofId string) {
	s.pendingQueue.Delete(proofId)
	s.cache.Delete(proofId)
}

func (s *Scheduler) PendingProofRequest() []*common.ProofRequest {
	return s.proofQueue.List()
}

func (s *Scheduler) IterPendingRequest(fn func(request *common.ProofRequest) error) {
	s.pendingQueue.Iterator(fn)
}

func NewScheduler(queue *ArrayQueue, pQueue *PendingQueue, filestore *FileStorage, store store.IStore, cache *cache, preparedData *Prepared,
	icpClient *dfinity.Client, btcClient *bitcoin.Client, ethClient *ethereum.Client) (*Scheduler, error) {
	return &Scheduler{
		proofQueue:   queue,
		pendingQueue: pQueue,
		fileStore:    filestore,
		chainStore:   NewChainStore(store),
		cache:        cache,
		preparedData: preparedData,
		btcClient:    btcClient,
		ethClient:    ethClient,
		icpClient:    icpClient,
	}, nil
}
