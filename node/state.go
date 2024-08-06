package node

import (
	"context"
	"fmt"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/bitcoin"
	"github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/store"
)

type IState interface {
	CheckBtcState() error
	CheckEthState() error
	CheckBeaconState() error
}

type State struct {
	proofQueue       *ArrayQueue
	fileStore        *FileStorage
	btcClient        *bitcoin.Client
	ethClient        *ethereum.Client
	store            store.IStore
	cache            *Cache
	preparedData     *PreparedData
	genesisPeriod    uint64
	genesisSlot      uint64
	btcGenesisHeight uint64 // start index
	debug            bool
}

func (s *State) CheckBtcState() error {
	logger.Debug("start check btc state ....")
	blockCount, err := s.btcClient.GetBlockCount()
	if err != nil {
		logger.Error("get block count error:%v", err)
		return err
	}
	// todo
	if s.debug {
		blockCount = 2871700 + 8*4
	}

	// btc genesis proof
	exists, err := CheckProof(s.fileStore, common.BtcGenesisType, s.btcGenesisHeight, s.btcGenesisHeight+common.BtcUpperDistance*2, "")
	if err != nil {
		logger.Error("check btc genesis proof error:%v", err)
		return err
	}
	if !exists {
		err := s.tryProofRequest(common.BtcGenesisType, s.btcGenesisHeight, s.btcGenesisHeight+common.BtcUpperDistance*2, "")
		if err != nil {
			logger.Error("try btc genesis proof error:%v", err)
			return err
		}
	}

	// btc upper proof
	btcUpperEndIndexes, err := s.fileStore.NeedBtcUpEndIndexes(uint64(blockCount))
	if err != nil {
		logger.Error("get need btc up index error:%v", err)
		return err
	}
	for _, endIndex := range btcUpperEndIndexes {
		startIndex := endIndex - common.BtcUpperDistance
		logger.Debug("start check btc upper: %v %v", startIndex, endIndex)
		ok, err := s.checkBtcUpper(startIndex, endIndex)
		if err != nil {
			logger.Error("check btc update error:%v", err)
			return err
		}
		if ok {
			err := s.tryProofRequest(common.BtcUpperType, startIndex, endIndex, "")
			if err != nil {
				logger.Error("try btc upper proof error:%v", err)
				return err
			}
		}
	}

	// btc recursive proof
	btcRecursiveEndIndexes, err := s.fileStore.NeedBtcRecursiveEndIndex(uint64(blockCount))
	if err != nil {
		logger.Error("get need btc recursive index error:%v", err)
		return err
	}
	for _, endIndex := range btcRecursiveEndIndexes {
		startIndex := endIndex - common.BtcUpperDistance
		err = s.tryProofRequest(common.BtcRecursiveType, startIndex, endIndex, "")
		if err != nil {
			logger.Error("try btc recursive proof error:%v %v %v", startIndex, endIndex, err)
			return err
		}
	}

	return nil
	// btc tx indexes
	unGenProofs, err := ReadAllUnGenProofs(s.store, Bitcoin)
	if err != nil {
		logger.Error("read unGen proof error:%v", err)
		return err
	}
	for _, tx := range unGenProofs {
		logger.Debug("bitcoin check ungen proof: %v %v", tx.ProofType.String(), tx.TxHash)
		if tx.ProofType == 0 || tx.TxHash == "" {
			logger.Warn("unGenProof error:%v %v", tx.ProofType.String(), tx.TxHash)
			err := DeleteUnGenProof(s.store, Bitcoin, tx.TxHash)
			if err != nil {
				logger.Error("delete ungen proof error:%v %v", tx.TxHash, err)
			}
			continue
		}
		switch tx.ProofType {
		case common.DepositTxType:
			err := s.checkDepositRequest(tx)
			if err != nil {
				logger.Error("check deposit request error:%v %v", tx.TxHash, err)
				continue
			}
		case common.VerifyTxType:
			err := s.tryProofRequest(common.VerifyTxType, 0, 0, tx.TxHash)
			if err != nil {
				logger.Error("try proof request error:%v %v", tx.TxHash, err)
				continue
			}
		default:
			logger.Error("unknown proof type: %v", tx.ProofType.String())
		}
	}
	return nil
}

func (s *State) checkDepositRequest(tx *DbUnGenProof) error {
	exists, err := CheckProof(s.fileStore, common.DepositTxType, 0, 0, tx.TxHash)
	if err != nil {
		logger.Error("check proof error:%v %v", tx.TxHash, err)
		return err
	}
	if exists {
		logger.Debug("%v %v proof exists ,delete ungen proof now", tx.ProofType.String(), tx.TxHash)
		err = DeleteUnGenProof(s.store, Bitcoin, tx.TxHash)
		if err != nil {
			logger.Error("delete ungen proof error:%v %v", tx.TxHash, err)
			return err
		}
		return nil
	}

	ok, confirms, err := s.CheckTxConfirms(tx.TxHash, tx.Amount)
	if err != nil {
		logger.Error("check tx confirms error: %v %v", tx.TxHash, err)
		return err
	}
	if !ok {
		logger.Warn("wait tx %v confirm: %v %v", tx.TxHash, tx.Amount, confirms)
		return nil
	}
	endHeight := tx.Height + uint64(confirms)
	if confirms <= 48 {
		exists, err := CheckProof(s.fileStore, common.BtcBulkType, tx.Height, endHeight, "")
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		if !exists {
			err := s.tryProofRequest(common.BtcBulkType, tx.Height, endHeight, "")
			if err != nil {
				logger.Error("try proof request error:%v %v", tx.TxHash, err)
				return err
			}
			return nil
		}

	} else {
		exists, err := CheckProof(s.fileStore, common.BtcPackedType, tx.Height, endHeight, "")
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		if !exists {
			err := s.tryProofRequest(common.BtcPackedType, tx.Height, endHeight, "")
			if err != nil {
				logger.Error(err.Error())
				return err
			}
			return nil
		}
	}
	wrapExists, err := CheckProof(s.fileStore, common.BtcWrapType, tx.Height, endHeight, "")
	if err != nil {
		logger.Error("check proof error:%v %v", tx.TxHash, err)
		return err
	}
	if !wrapExists {
		err := s.tryProofRequest(common.BtcWrapType, tx.Height, endHeight, "")
		if err != nil {
			logger.Error("try proof request error:%v %v", tx.TxHash, err)
			return err
		}
		return nil
	}
	err = s.tryProofRequest(common.DepositTxType, tx.Height, endHeight, tx.TxHash)
	if err != nil {
		logger.Error("try proof request error:%v %v", tx.TxHash, err)
		return err
	}
	return nil
}

func (s *State) CheckTxConfirms(hash string, amount uint64) (bool, int, error) {
	needConfirms := 0
	if amount < 100000000 {
		needConfirms = 1
	} else if amount < 200000000 {
		needConfirms = 2
	} else {
		needConfirms = 3
	}
	tx, err := s.btcClient.GetTransaction(hash)
	if err != nil {
		logger.Error("get tx error:%v %v", hash, err)
		return false, 0, err
	}
	if tx.Confirmations >= needConfirms {
		return true, needConfirms, nil
	}
	return false, 0, nil
}
func (s *State) checkBtcUpper(start, end uint64) (bool, error) {
	next := true
	for index := start; index < end; index = index + common.BtcMiddleDistance {
		startIndex := index
		endIndex := index + common.BtcMiddleDistance
		logger.Debug("start check btc middle: %v %v", startIndex, endIndex)
		exists, err := CheckProof(s.fileStore, common.BtcMiddleType, startIndex, endIndex, "")
		if err != nil {
			logger.Error("check btc update error:%v", err)
			return false, err
		}
		if !exists {
			next = false
			ok, err := s.checkBtcMiddle(startIndex, endIndex)
			if err != nil {
				logger.Error("check btc update error:%v", err)
				return false, err
			}
			if ok {
				err := s.tryProofRequest(common.BtcMiddleType, startIndex, endIndex, "")
				if err != nil {
					logger.Error("try btc middle proof error:%v", err)
					return false, err
				}
			}
		}
		logger.Debug("end check btc upper: %v %v", startIndex, endIndex)
	}
	return next, nil
}

func (s *State) checkBtcMiddle(start, end uint64) (bool, error) {
	next := true
	for index := start; index < end; index = index + common.BtcBaseDistance {
		logger.Debug("start check btc middle: %v %v", start, end)
		startIndex := index
		endIndex := index + common.BtcBaseDistance
		exists, err := CheckProof(s.fileStore, common.BtcBaseType, startIndex, endIndex, "")
		if err != nil {
			logger.Error("check btc update error:%v", err)
			return false, err
		}
		if !exists {
			next = false
			err := s.tryProofRequest(common.BtcBaseType, startIndex, endIndex, "")
			if err != nil {
				logger.Error("try btc base proof error:%v_%v %v", start, endIndex, err)
				return false, err
			}
		}
		logger.Debug("end check btc middle: %v %v", start, end)
	}
	return next, nil

}

func (s *State) tryProofRequest(reqType common.ZkProofType, fIndex, sIndex uint64, hash string) error {
	proofId := common.NewProofId(reqType, fIndex, sIndex, hash)
	exists := s.cache.Check(proofId)
	if exists {
		logger.Debug("proof request exists: %v", proofId)
		return nil
	}
	exists, err := CheckProof(s.fileStore, reqType, fIndex, sIndex, hash)
	if err != nil {
		logger.Error("check proof error:%v %v", proofId, err)
		return err
	}
	if exists {
		return nil
	}
	data, ok, err := GenRequestData(s.preparedData, reqType, fIndex, sIndex, hash)
	if err != nil {
		logger.Error("get request data error:%v %v", proofId, err)
		return err
	}
	if !ok {
		return nil
	}
	zkProofRequest := common.NewZkProofRequest(reqType, data, fIndex, sIndex, hash)
	// todo
	s.cache.Store(proofId, nil)
	s.proofQueue.Push(zkProofRequest)
	logger.Info("success add request:%v", proofId)
	return nil
}

func (s *State) CheckReq(reqType common.ZkProofType, index uint64, hash string) (bool, error) {
	switch reqType {
	case common.SyncComGenesisType:
		return index == s.genesisPeriod+1, nil
	case common.SyncComUnitType:
		return index >= s.genesisPeriod, nil
	case common.SyncComRecursiveType:
		return index >= s.genesisPeriod+2, nil
	case common.BeaconHeaderFinalityType:
		return index >= s.genesisSlot, nil
	case common.TxInEth2:
		finalizedSlot, ok, err := s.fileStore.GetFinalizedSlot()
		if err != nil {
			logger.Error("get latest slot error: %v", err)
			return false, err
		}
		if !ok {
			logger.Warn("no find latest slot")
			return false, nil
		}
		receipt, err := s.ethClient.TransactionReceipt(context.Background(), ethCommon.HexToHash(hash))
		if err != nil {
			logger.Error("get tx receipt error: %v", err)
			return false, err
		}
		txSlot, ok, err := ReadBeaconSlot(s.store, receipt.BlockNumber.Uint64())
		if err != nil {
			logger.Error("get beacon slot error: %v %v", receipt.BlockNumber, err)
			return false, err
		}
		if !ok {
			return false, nil
		}
		if txSlot <= finalizedSlot {
			return true, nil
		}
		logger.Warn("%v tx slot %v less than finalized slot %v", hash, txSlot, finalizedSlot)
		return false, nil
	case common.BeaconHeaderType:
		_, ok, err := s.fileStore.GetNearTxSlotFinalizedSlot(index)
		if err != nil {
			logger.Error("get latest slot error: %v", err)
			return false, err
		}
		if !ok {
			return false, nil
		}
		return true, nil
	case common.RedeemTxType:
		return true, nil
	default:
		return false, fmt.Errorf("check request status never should happen: %v %v", index, reqType)
	}
}

func (s *State) CheckEthState() error {

	logger.Debug("check eth state now  ....")
	// todo
	if s.debug {
		// ethereum tx: RedeemTxType_2195577_0x291ee31eb6b8cef1ebc571fd090a1e7c96ddac5a1552dae47501581ed7d66641
		txHash := "0x291ee31eb6b8cef1ebc571fd090a1e7c96ddac5a1552dae47501581ed7d66641"
		has, err := s.store.HasObj(DbUnGenProofId(Ethereum, txHash))
		if err != nil {
			logger.Error("write ungen proof error: %v", err)
			return err
		}
		if !has {
			err := WriteBeaconSlot(s.store, 2025122, 2195577)
			if err != nil {
				logger.Error("write beacon slot error: %v", err)
				return err
			}
			err = WriteTxes(s.store, []DbTx{
				{
					Height:    2025122,
					TxHash:    txHash,
					TxIndex:   16,
					Amount:    0,
					TxType:    RedeemTx,
					ChainType: Ethereum,
				},
			})
			if err != nil {
				logger.Error("write tx error: %v", err)
				return err
			}
			err = WriteUnGenProof(s.store, Ethereum, []*DbUnGenProof{
				{
					TxHash:    txHash,
					ProofType: common.RedeemTxType,
					ChainType: Ethereum,
					Height:    2025122,
					TxIndex:   16,
					Amount:    0,
				},
			})
			if err != nil {
				logger.Error("write ungen proof error: %v", err)
				return err
			}
		}

	}

	unGenProofs, err := ReadAllUnGenProofs(s.store, Ethereum)
	if err != nil {
		logger.Error("read all ungen proof ids error: %v", err)
		return err
	}
	for _, item := range unGenProofs {
		txHash := item.TxHash
		logger.Debug("start check redeem proof tx: %v %v %v", txHash, item.Height, item.TxIndex)
		exists, err := CheckProof(s.fileStore, common.RedeemTxType, 0, 0, txHash)
		if err != nil {
			logger.Error("check tx proof error: %v", err)
			return err
		}
		if exists {
			logger.Debug("redeem proof exist now,delete cache: %v", txHash)
			err := DeleteUnGenProof(s.store, Ethereum, txHash)
			if err != nil {
				logger.Error("delete ungen proof error: %v", err)
				return err
			}
			logger.Debug("delete ungen proof tx: %v", txHash)
			continue
		}
		txSlot, ok, err := s.GetSlotByHash(txHash)
		if err != nil {
			logger.Error("get txSlot error: %v %v", err, txHash)
			return err
		}
		if !ok {
			logger.Warn("no find  tx %v beacon slot", txHash)
			continue
		}
		// todo
		finalizedSlot, ok, err := s.fileStore.GetNearTxSlotFinalizedSlot(txSlot)
		if err != nil {
			logger.Error("get near tx slot finalized slot error: %v", err)
			return err
		}
		if !ok {
			logger.Warn("no find near %v tx slot finalized slot", txSlot)
			continue
		}
		err = s.updateRedeemProofStatus(txHash, txSlot, common.ProofFinalized)
		if err != nil {
			logger.Error("update proof status error: %v %v", txHash, err)
			return err
		}
		exists, err = CheckProof(s.fileStore, common.TxInEth2, txSlot, finalizedSlot, txHash)
		if err != nil {
			logger.Error("check tx proof error: %v", err)
			return err
		}
		if !exists {
			err := s.tryProofRequest(common.TxInEth2, txSlot, finalizedSlot, txHash)
			if err != nil {
				logger.Error("try proof request error: %v", err)
				return err
			}
		}
		exists, err = CheckProof(s.fileStore, common.BeaconHeaderType, txSlot, finalizedSlot, "")
		if err != nil {
			logger.Error("check block header proof error: %v", err)
			return err
		}
		if !exists {
			err := WriteTxSlot(s.store, txSlot, item)
			if err != nil {
				logger.Error("write tx slot error: %v %v %v", txHash, txSlot, err)
				return err
			}
			err = s.tryProofRequest(common.BeaconHeaderType, txSlot, finalizedSlot, "")
			if err != nil {
				logger.Error("try proof request error: %v", err)
				return err
			}
		}
		//logger.Debug("%v find near %v tx slot finalized slot %v", txHash, txSlot, finalizedSlot)
		exists, err = CheckProof(s.fileStore, common.BeaconHeaderFinalityType, finalizedSlot, 0, "")
		if err != nil {
			logger.Error("check block header finality proof error: %v %v", finalizedSlot, err)
			return err
		}
		if !exists {
			err := WriteTxFinalizedSlot(s.store, finalizedSlot, item)
			if err != nil {
				logger.Error("write tx finalized slot error: %v %v %v", finalizedSlot, txHash, err)
				return err
			}
			err = s.tryProofRequest(common.BeaconHeaderFinalityType, finalizedSlot, 0, "")
			if err != nil {
				logger.Error("try proof request error: %v", err)
				return err
			}
			continue
		}
		err = s.tryProofRequest(common.RedeemTxType, txSlot, finalizedSlot, txHash)
		if err != nil {
			logger.Error("try proof request error: %v", err)
			return err
		}
	}
	return nil
}

func (s *State) GetSlotByHash(hash string) (uint64, bool, error) {
	//txHash := ethCommon.HexToHash(hash)
	//receipt, err := s.ethClient.TransactionReceipt(context.Background(), txHash)
	//if err != nil {
	//	logger.Error("get tx receipt error: %v %v", hash, err)
	//	return 0, false, err
	//}
	// todo

	dbTx, err := ReadDbTx(s.store, hash)
	if err != nil {
		logger.Error("read db tx error: %v %v", hash, err)
		return 0, false, err
	}
	beaconSlot, ok, err := ReadBeaconSlot(s.store, dbTx.Height)
	if err != nil {
		logger.Error("get beacon slot error: %v %v", hash, err)
		return 0, false, err
	}
	if !ok {
		return 0, false, nil
	}
	return beaconSlot, true, nil
}

func (s *State) updateRedeemProofStatus(txHash string, index uint64, status common.ProofStatus) error {
	id := common.NewProofId(common.RedeemTxType, index, 0, txHash)
	if !s.cache.Check(id) {
		err := UpdateProof(s.store, txHash, "", common.RedeemTxType, status)
		if err != nil {
			logger.Error("update proof status error: %v %v", txHash, err)
			return err
		}
		return err
	}
	return nil
}

func (s *State) CheckBeaconState() error {

	// beacon genesis proof
	exists, err := s.fileStore.CheckGenesisProof()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if !exists {
		logger.Warn("no find genesis proof, send request genesis proof")
		genesisPeriod := s.genesisPeriod + 1
		err := s.tryProofRequest(common.SyncComGenesisType, genesisPeriod, 0, "")
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	// beacon unit proof
	unitProofIndexes, err := s.fileStore.NeedGenUnitProofIndexes()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	for _, index := range unitProofIndexes {
		if index < s.genesisPeriod {
			continue
		}
		err := s.tryProofRequest(common.SyncComUnitType, index, 0, "")
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	// beacon recursive index
	genRecProofIndexes, err := s.fileStore.NeedGenRecProofIndexes()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	for _, index := range genRecProofIndexes {
		if index <= s.genesisPeriod+1 {
			continue
		}
		err := s.tryProofRequest(common.SyncComRecursiveType, index, 0, "")
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		break
	}
	return nil

}

func (s *State) CheckProofRequest(resp *common.ZkProofResponse) error {
	requests, err := s.findNewRequests(resp)
	if err != nil {
		logger.Error("check redeem request error:%v %v", resp.Id(), err)
		return err
	}
	for _, req := range requests {
		if !s.cache.Check(req.Id()) {
			logger.Debug("add new request:%v to queue", req.Id())
			s.cache.Store(req.Id(), nil)
			s.proofQueue.Push(req)
			err := s.UpdateProofStatus(req, common.ProofQueued)
			if err != nil {
				logger.Error("update Proof status error:%v %v", req.Id(), err)
			}
		}
	}
	return nil

}

func (s *State) findNewRequests(resp *common.ZkProofResponse) ([]*common.ZkProofRequest, error) {
	switch resp.ZkProofType {
	case common.TxInEth2:
		request, ok, err := s.checkRedeemRequest(resp.TxHash)
		if err != nil {
			logger.Error("get redeem request error:%v %v", resp.Id(), err)
			return nil, err
		}
		if ok {
			return []*common.ZkProofRequest{request}, nil
		}
		return nil, nil
	case common.BeaconHeaderType:
		txes, err := ReadAllTxBySlot(s.store, resp.Index)
		if err != nil {
			logger.Error("get redeem request error:%v %v", resp.Id(), err)
			return nil, err
		}
		var result []*common.ZkProofRequest
		for _, tx := range txes {
			request, ok, err := s.checkRedeemRequest(tx.TxHash)
			if err != nil {
				logger.Error("get redeem request error:%v %v", resp.Id(), err)
				return nil, err
			}
			if ok {
				result = append(result, request)
			}
		}
		return result, nil
	case common.BeaconHeaderFinalityType:
		txes, err := ReadAllTxByFinalizedSlot(s.store, resp.Index)
		if err != nil {
			logger.Error("get redeem request error:%v %v", resp.Id(), err)
			return nil, err
		}
		var result []*common.ZkProofRequest
		for _, tx := range txes {
			request, ok, err := s.checkRedeemRequest(tx.TxHash)
			if err != nil {
				logger.Error("get redeem request error:%v %v", resp.Id(), err)
				return nil, err
			}
			if ok {
				result = append(result, request)
			}
		}
		return result, nil
	default:
		return nil, nil
	}
}

func (s *State) checkRedeemRequest(txHash string) (*common.ZkProofRequest, bool, error) {
	// todo
	txSlot, ok, err := s.GetSlotByHash(txHash)
	if err != nil {
		logger.Error("get slot by hash error: %v %v", txHash, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	data, ok, err := s.preparedData.GetRedeemRequestData(s.genesisPeriod, txSlot, txHash)
	if err != nil {
		logger.Error("get redeem request data error: %v %v", txHash, err)
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	request := common.NewZkProofRequest(common.RedeemTxType, data, txSlot, 0, txHash)
	return request, true, nil
}

func (s *State) UpdateProofStatus(req *common.ZkProofRequest, status common.ProofStatus) error {
	// todo
	if req.ReqType == common.DepositTxType || req.ReqType == common.RedeemTxType {
		err := UpdateProof(s.store, req.TxHash, "", req.ReqType, status)
		if err != nil {
			logger.Error("update Proof status error:%v %v", req.Id(), err)
			return err
		}
	}
	return nil
}

func NewState(queue *ArrayQueue, filestore *FileStorage, store store.IStore, cache *Cache, preparedData *PreparedData,
	btcGenesisHeight, genesisPeriod, genesisSlot uint64, btcClient *bitcoin.Client, ethClient *ethereum.Client) (*State, error) {
	return &State{
		proofQueue:       queue,
		fileStore:        filestore,
		store:            store,
		cache:            cache,
		preparedData:     preparedData,
		genesisPeriod:    genesisPeriod,
		btcGenesisHeight: btcGenesisHeight,
		genesisSlot:      genesisSlot,
		btcClient:        btcClient,
		ethClient:        ethClient,
		debug:            common.GetEnvDebugMode(),
	}, nil
}
