package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"strconv"
)

var _ IBeaconAgent = (*BeaconAgent)(nil)

type BeaconAgent struct {
	beaconClient   *beacon.Client
	fileStore      *FileStorage
	zkProofRequest chan []*common.ZkProofRequest
	name           string
	genesisPeriod  uint64
	genesisSlot    uint64
	stateCache     *BeaconCache
}

func NewBeaconAgent(cfg Config, beaconClient *beacon.Client, zkProofReq chan []*common.ZkProofRequest,
	fileStore *FileStorage, genesisSlot, genesisPeriod uint64, fetchDataResp chan FetchDataResponse) (IBeaconAgent, error) {
	beaconAgent := &BeaconAgent{
		fileStore:      fileStore,
		beaconClient:   beaconClient,
		name:           "beaconAgent",
		stateCache:     NewBeaconCache(),
		zkProofRequest: zkProofReq,
		genesisPeriod:  genesisPeriod,
		genesisSlot:    genesisSlot,
	}
	return beaconAgent, nil
}

func (b *BeaconAgent) Init() error {
	logger.Info("beacon agent init")
	latestPeriod, exists, err := b.fileStore.GetPeriod()
	if err != nil {
		logger.Error("check latest Index error: %v", err)
		return err
	}
	if !exists || latestPeriod < b.genesisPeriod {
		logger.Warn("no find latest Index, store %v Index to db", b.genesisPeriod)
		err := b.fileStore.StorePeriod(b.genesisPeriod)
		if err != nil {
			logger.Error("store latest Index error: %v", err)
			return err
		}
	}
	latestSlot, exists, err := b.fileStore.GetFinalizedSlot()
	if err != nil {
		logger.Error("check latest Slot error: %v", err)
		return err
	}
	if !exists || latestSlot < b.genesisSlot {
		logger.Warn("no find latest slot, store %v slot to db", b.genesisSlot)
		err := b.fileStore.StoreFinalizedSlot(b.genesisSlot)
		if err != nil {
			logger.Error("store latest Slot error: %v", err)
			return err
		}
	}
	genesisBootStrapExists, err := b.fileStore.CheckBootstrap()
	if err != nil {
		logger.Error("check genesis Update error: %v", err)
		return err
	}
	if !genesisBootStrapExists {
		logger.Warn("no find genesis update, send request genesis update")
	}
	genesisUpdate, err := b.fileStore.CheckUpdate(b.genesisPeriod)
	if err != nil {
		logger.Error("check update error: %v", err)
		return err
	}
	if !genesisUpdate {
		logger.Warn("no find %v first Index update, send request update", b.genesisPeriod)
	}
	// todo check Data
	return err
}

// todo maybe use for  to replace recursive
func (b *BeaconAgent) tryProofRequest(index uint64, reqType common.ZkProofType) error {
	ok, err := b.checkRequest(index, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if !ok {
		return nil
	}
	//	logger.Debug("beacon check and new Proof request: %v %v", Index, reqType.String())
	proofExists, err := b.CheckProofExists(index, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if proofExists {
		logger.Warn("%v %v Proof exists", index, reqType.String())
		return nil
	}
	existsRequest, err := b.CheckProofRequestStatus(index, reqType)
	if err != nil {
		logger.Info("can`t send Proof: %v %v", index, reqType.String())
		return err
	}
	if existsRequest {
		//logger.Warn("Proof request exists，%v %v  skip it now", Index, ReqType.String())
		return nil
	}
	data, prepareDataOk, err := b.prepareProofRequestData(index, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if !prepareDataOk {
		//logger.Warn("Proof request Data haven`t prepared now ,%v %v  can`t generate Proof", Index, ReqType.String())
		return nil
	}
	err = b.cacheProofRequestStatus(index, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	request := common.NewZkProofRequest(reqType, data, index, "")
	b.zkProofRequest <- []*common.ZkProofRequest{request}
	logger.Info("success send Proof request: %v %v", index, reqType.String())

	return nil
}

func (b *BeaconAgent) checkRequest(index uint64, reqType common.ZkProofType) (bool, error) {
	switch reqType {
	case common.SyncComGenesisType:
		return index == b.genesisPeriod+1, nil
	case common.SyncComUnitType:
		return index >= b.genesisPeriod, nil
	case common.SyncComRecursiveType:
		return index >= b.genesisPeriod+2, nil
	case common.BeaconHeaderFinalityType:
		return index >= b.genesisSlot, nil
	default:
		return false, fmt.Errorf("check request status never should happen: %v %v", index, reqType)
	}
}

func (b *BeaconAgent) CheckProofRequestStatus(index uint64, reqType common.ZkProofType) (bool, error) {
	switch reqType {
	case common.SyncComGenesisType:
		return b.stateCache.CheckGenesis(), nil
	case common.SyncComUnitType:
		return b.stateCache.CheckUnit(index), nil
	case common.SyncComRecursiveType:
		return b.stateCache.CheckRecursive(index), nil
	case common.BeaconHeaderFinalityType:
		return b.stateCache.CheckBhfUpdate(index), nil
	default:
		return false, fmt.Errorf("check request status never should happen: %v %v", index, reqType)
	}
}

func (b *BeaconAgent) cacheProofRequestStatus(index uint64, reqType common.ZkProofType) error {
	logger.Debug("beacon cache request status: %v %v", index, reqType.String())
	switch reqType {
	case common.SyncComGenesisType:
		return b.stateCache.StoreGenesis()
	case common.SyncComUnitType:
		return b.stateCache.StoreUnit(index)
	case common.SyncComRecursiveType:
		return b.stateCache.StoreRecursive(index)
	case common.BeaconHeaderFinalityType:
		return b.stateCache.StoreBhfUpdate(index)
	default:
		return fmt.Errorf("cache request status never should happen: %v %v", index, reqType)
	}
}

func (b *BeaconAgent) prepareProofRequestData(index uint64, reqType common.ZkProofType) (data interface{}, prepared bool, err error) {
	switch reqType {
	case common.SyncComGenesisType:
		data, prepared, err = b.GetGenesisRaw()
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return data, prepared, nil
	case common.SyncComUnitType:
		data, prepared, err = b.GetUnitData(index)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return data, prepared, nil
	case common.SyncComRecursiveType:
		data, prepared, err = b.GetRecursiveData(index)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return data, prepared, nil
	case common.BeaconHeaderFinalityType:
		data, prepared, err = b.GetBhfUpdateData(index)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return data, prepared, nil
	default:
		logger.Error(" prepare request Data never should happen : %v %v", index, reqType)
		return nil, false, fmt.Errorf("never should happen : %v %v", index, reqType)
	}
}

func (b *BeaconAgent) CheckProofExists(index uint64, reqType common.ZkProofType) (bool, error) {
	switch reqType {
	case common.SyncComGenesisType:
		return b.fileStore.CheckGenesisProof()
	case common.SyncComUnitType:
		return b.fileStore.CheckUnitProof(index)
	case common.SyncComRecursiveType:
		return b.fileStore.CheckRecursiveProof(index)
	case common.BeaconHeaderFinalityType:
		return b.fileStore.CheckBhfProof(index)
	default:
		logger.Error("check Proof exists never should happen : %v %v", index, reqType)
		return false, fmt.Errorf("check Proof exists never should happen : %v %v", index, reqType)
	}
}

func (b *BeaconAgent) CheckState() error {
	genesisProofExists, err := b.fileStore.CheckGenesisProof()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if !genesisProofExists {
		logger.Warn("no find genesis proof, send request genesis proof")
		genesisPeriod := b.genesisPeriod + 1
		err := b.tryProofRequest(genesisPeriod, common.SyncComGenesisType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	fetchIndexes, err := b.fileStore.NeedUpdateIndexes()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	for _, index := range fetchIndexes {
		if b.stateCache.CheckFetchData(index) {
			continue
		}
	}
	unitProofIndexes, err := b.fileStore.NeedGenUnitProofIndexes()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	// todo
	for _, index := range unitProofIndexes {
		if index < b.genesisPeriod {
			continue
		}
		if b.stateCache.CheckUnit(index) {
			continue
		}
		logger.Warn("need unit proof: %v", index)
		err := b.tryProofRequest(index, common.SyncComUnitType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	genRecProofIndexes, err := b.fileStore.NeedGenRecProofIndexes()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	var skip bool
	for _, index := range genRecProofIndexes {
		if index <= b.genesisPeriod+1 {
			continue
		}
		if b.stateCache.CheckRecursive(index) {
			continue
		}
		if skip {
			break
		}
		//logger.Warn("need recursive proof: %v", Index)
		err := b.tryProofRequest(index, common.SyncComRecursiveType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		skip = true
	}
	//// todo
	//bhfUpdateIndexes, err := b.fileStore.NeedGenBhfUpdateIndex()
	//if err != nil {
	//	logger.Error(err.Error())
	//	return err
	//}
	//for _, index := range bhfUpdateIndexes {
	//	logger.Info("need to update block header finality: %v", index)
	//	err := b.tryProofRequest(index, common.BeaconHeaderFinalityType)
	//	if err != nil {
	//		logger.Error(err.Error())
	//		return err
	//	}
	//}
	return nil
}

func (b *BeaconAgent) FetchDataResponse(req FetchDataResponse) error {
	logger.Debug("beacon fetch response fetchType: %v, Index: %v", req.UpdateType.String(), req.Index)
	return nil
}

// Todo

func (b *BeaconAgent) GetSyncCommitRootID(period uint64) ([]byte, bool, error) {
	var currentPeriodUpdate structs.LightClientUpdateWithVersion
	exists, err := b.fileStore.GetUpdate(period, &currentPeriodUpdate)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v Index update Data, send new update request", period)
		return nil, false, nil
	}
	// todo
	var update utils.LightClientUpdateInfo
	err = ParseObj(currentPeriodUpdate.Data, &update)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	update.Version = currentPeriodUpdate.Version
	if b.genesisPeriod == period {
		var genesisData structs.LightClientBootstrapResponse
		genesisExists, err := b.fileStore.GetBootstrap(&genesisData)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		if !genesisExists {
			logger.Warn("no find genesis update Data, send new update request")
			return nil, false, nil
		}
		// todo
		var genesisCommittee utils.SyncCommittee
		err = ParseObj(genesisData.Data.CurrentSyncCommittee, &genesisCommittee)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		update.CurrentSyncCommittee = &genesisCommittee
	} else {
		prePeriod := period - 1
		if prePeriod < b.genesisPeriod {
			logger.Error("should never happen: %v", prePeriod)
			return nil, false, nil
		}
		var preUpdateData structs.LightClientUpdateWithVersion
		preUpdateExists, err := b.fileStore.GetUpdate(prePeriod, &preUpdateData)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		if !preUpdateExists {
			logger.Warn("get unit Data,no find %v Index update Data, send new update request", prePeriod)
			return nil, false, nil
		}
		// todo
		var currentSyncCommittee utils.SyncCommittee
		err = ParseObj(preUpdateData.Data.NextSyncCommittee, &currentSyncCommittee)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		update.CurrentSyncCommittee = &currentSyncCommittee
	}
	rootId, err := circuits.SyncCommitRoot(&update)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	return rootId, true, nil

}

func (b *BeaconAgent) GetGenesisRaw() (interface{}, bool, error) {
	genesisId, ok, err := b.GetSyncCommitRootID(b.genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index genesis commitId  no find", b.genesisPeriod)
		return nil, false, nil
	}

	nextPeriod := b.genesisPeriod + 1
	firstId, ok, err := b.GetSyncCommitRootID(nextPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index first commitId no find", nextPeriod)
		return nil, false, nil
	}
	secondPeriod := nextPeriod + 1
	secondId, ok, err := b.GetSyncCommitRootID(secondPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index second commitId no find", secondPeriod)
		return nil, false, nil
	}

	firstProof, exists, err := b.fileStore.GetUnitProof(b.genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("get genesis Data,first proof not exists: %v Index", b.genesisPeriod)
		//err := b.tryProofRequest(b.genesisPeriod, SyncComUnitType)
		//if err != nil {
		//	logger.Error(err.Error())
		//	return nil, false, err
		//}
		return nil, false, nil
	}
	logger.Info("get genesis first proof: %v", b.genesisPeriod)

	secondProof, exists, err := b.fileStore.GetUnitProof(nextPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}

	if !exists {
		logger.Warn("get genesis Data,second proof not exists: %v Index", nextPeriod)
		return nil, false, nil
	}
	logger.Info("get genesis second proof: %v", nextPeriod)

	genesisProofParam := &rpc.SyncCommGenesisRequest{
		FirstProof:    firstProof.Proof,
		FirstWitness:  firstProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		GenesisID:     genesisId,
		FirstID:       firstId,
		SecondID:      secondId,
	}
	return genesisProofParam, true, nil

}

func (b *BeaconAgent) GetUnitData(period uint64) (*rpc.SyncCommUnitsRequest, bool, error) {
	var currentPeriodUpdate structs.LightClientUpdateWithVersion
	exists, err := b.fileStore.GetUpdate(period, &currentPeriodUpdate)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v Index update Data, send new update request", period)
		return nil, false, nil
	}
	if b.genesisPeriod == period {
		var genesisData structs.LightClientBootstrapResponse
		genesisExists, err := b.fileStore.GetBootstrap(&genesisData)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		if !genesisExists {
			logger.Warn("no find genesis update Data, send new update request")
			return nil, false, nil
		}
		return &rpc.SyncCommUnitsRequest{
			Period:                  period,
			Version:                 currentPeriodUpdate.Version,
			AttestedHeader:          currentPeriodUpdate.Data.AttestedHeader,
			CurrentSyncCommittee:    genesisData.Data.CurrentSyncCommittee, // todo
			SyncAggregate:           currentPeriodUpdate.Data.SyncAggregate,
			FinalizedHeader:         currentPeriodUpdate.Data.FinalizedHeader,
			FinalityBranch:          currentPeriodUpdate.Data.FinalityBranch,
			NextSyncCommittee:       currentPeriodUpdate.Data.NextSyncCommittee,
			NextSyncCommitteeBranch: currentPeriodUpdate.Data.NextSyncCommitteeBranch,
			SignatureSlot:           currentPeriodUpdate.Data.SignatureSlot,
		}, true, nil
	} else {
		prePeriod := period - 1
		if prePeriod < b.genesisPeriod {
			logger.Error("should never happen: %v", prePeriod)
			return nil, false, nil
		}
		var perUpdateData structs.LightClientUpdateWithVersion
		preUpdateExists, err := b.fileStore.GetUpdate(prePeriod, &perUpdateData)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		if !preUpdateExists {
			logger.Warn("get unit Data,no find %v Index update Data, send new update request", prePeriod)
			return nil, false, nil
		}
		return &rpc.SyncCommUnitsRequest{
			Period:                  period,
			Version:                 currentPeriodUpdate.Version,
			AttestedHeader:          currentPeriodUpdate.Data.AttestedHeader,
			CurrentSyncCommittee:    perUpdateData.Data.NextSyncCommittee, // TODO
			SyncAggregate:           currentPeriodUpdate.Data.SyncAggregate,
			FinalizedHeader:         currentPeriodUpdate.Data.FinalizedHeader,
			FinalityBranch:          currentPeriodUpdate.Data.FinalityBranch,
			NextSyncCommittee:       currentPeriodUpdate.Data.NextSyncCommittee,
			NextSyncCommitteeBranch: currentPeriodUpdate.Data.NextSyncCommitteeBranch,
			SignatureSlot:           currentPeriodUpdate.Data.SignatureSlot,
		}, true, nil
	}

}

func (b *BeaconAgent) GetRecursiveData(period uint64) (interface{}, bool, error) {
	if period == b.genesisPeriod+2 {
		// todo should  start from  (genesis+1) Index
		return b.getRecursiveGenesisData(period)
	} else if period > b.genesisPeriod+2 {
		// todo should start from (genesis+2) Index
		return b.getRecursiveData(period)
	}
	return nil, false, nil

}

func (b *BeaconAgent) getRecursiveData(period uint64) (interface{}, bool, error) {
	//todo
	genesisId, ok, err := b.GetSyncCommitRootID(b.genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index genesis commitId no find", b.genesisPeriod)
		return nil, false, nil
	}
	relayId, ok, err := b.GetSyncCommitRootID(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index relay commitId no find", period)
		return nil, false, nil
	}
	endPeriod := period + 1
	endId, ok, err := b.GetSyncCommitRootID(endPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index end commitId no find", endPeriod)
		return nil, false, nil
	}
	secondProof, exists, err := b.fileStore.GetUnitProof(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v unit proof Data, send new proof request", period)
		return nil, false, nil
	}

	prePeriod := period - 1
	firstProof, exists, err := b.fileStore.GetRecursiveProof(prePeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v Index recursive Data, send new proof request", prePeriod)
		return nil, false, nil
	}

	return &rpc.SyncCommRecursiveRequest{
		Choice:        "recursive",
		FirstProof:    firstProof.Proof,
		FirstWitness:  firstProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		BeginId:       genesisId,
		RelayId:       relayId,
		EndId:         endId,
	}, true, nil
}

func (b *BeaconAgent) getRecursiveGenesisData(period uint64) (interface{}, bool, error) {
	// todo
	genesisId, ok, err := b.GetSyncCommitRootID(b.genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index genesis commitId no find", b.genesisPeriod)
		return nil, false, nil
	}
	relayPeriod := period
	relayId, ok, err := b.GetSyncCommitRootID(relayPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v  Index relay commitId no find ", relayPeriod)
		return nil, false, nil
	}
	endPeriod := relayPeriod + 1
	endId, ok, err := b.GetSyncCommitRootID(endPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index end commitId no find", endPeriod)
		return nil, false, nil
	}

	fistProof, firstExists, err := b.fileStore.GetGenesisProof()
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !firstExists {
		logger.Warn("no find genesis proof ,start new proof request")
		//err := b.tryProofRequest(Index, SyncComGenesisType)
		//if err != nil {
		//	logger.Error(err.Error())
		//	return nil, false, err
		//}
		return nil, false, nil
	}
	secondProof, secondExists, err := b.fileStore.GetUnitProof(relayPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !secondExists {
		logger.Warn("no find %v Index unit proof , send new proof request", relayPeriod)
		//err := b.tryProofRequest(relayPeriod, SyncComUnitType)
		//if err != nil {
		//	logger.Error(err.Error())
		//	return nil, false, err
		//}
		return nil, false, nil
	}

	return &rpc.SyncCommRecursiveRequest{
		Choice:        "genesis",
		FirstProof:    fistProof.Proof,
		FirstWitness:  fistProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		BeginId:       genesisId,
		RelayId:       relayId,
		EndId:         endId,
	}, true, nil

}

func (b *BeaconAgent) ProofResponse(resp *common.ZkProofResponse) error {
	// todo
	logger.Info("beacon Proof response type: %v, Index: %v", resp.ZkProofType.String(), resp.Period)
	index := resp.Period
	if resp.ZkProofType != common.UnitOuter {
		b.deleteCacheProofReqStatus(resp.ZkProofType, resp.Period)
	}
	err := StoreZkProof(b.fileStore, resp.ZkProofType, index, resp.TxHash, resp.Proof, resp.Witness)
	if err != nil {
		logger.Error("store proof error: %v", err)
		return err
	}
	return nil
}

func (b *BeaconAgent) deleteCacheProofReqStatus(reqType common.ZkProofType, period uint64) {
	switch reqType {
	case common.SyncComGenesisType:
		b.stateCache.DeleteGenesis()
	case common.SyncComUnitType:
		b.stateCache.DeleteUnit(period)
	case common.SyncComRecursiveType:
		b.stateCache.DeleteRecursive(period)
	case common.BeaconHeaderFinalityType:
		b.stateCache.DeleteBhfUpdate(period)
	default:
		logger.Error("delete cache request status never should happen: %v %v", reqType, period)
		return
	}
}

func (b *BeaconAgent) Close() error {

	return nil
}

func (b *BeaconAgent) Name() string {
	return b.name
}

func (b *BeaconAgent) Stop() {

}

func (b *BeaconAgent) GetBhfUpdateData(slot uint64) (interface{}, bool, error) {
	logger.Debug("get bhf update data: %v", slot)

	var currentFinalityUpdate structs.LightClientUpdateWithVersion
	exists, err := b.fileStore.GetFinalityUpdate(slot, &currentFinalityUpdate)
	if err != nil {
		logger.Error("get finality update error: %v %v", slot, err)
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find finality update: %v", slot)
		return nil, false, nil
	}

	genesisId, ok, err := b.GetSyncCommitRootID(b.genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v Index genesis commitId no find", b.genesisPeriod)
		return nil, false, nil
	}
	// todo
	attestedSlot, err := strconv.ParseUint(currentFinalityUpdate.Data.AttestedHeader.Slot, 10, 64)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	period := (attestedSlot / 8192) - 1
	recursiveProof, ok, err := b.fileStore.GetRecursiveProof(period)
	if err != nil {
		logger.Error("get recursive proof error: %v %v", period, err)
		return nil, false, err
	}
	if !ok {
		logger.Warn("no find recursive proof: %v", period)
		return nil, false, nil
	}
	outerPeriod := period + 1
	logger.Debug("get bhf update data slot: %v,recPeriod: %v,outPeriod %v", slot, period, outerPeriod)
	outerProof, ok, err := b.fileStore.GetOuterProof(outerPeriod)
	if err != nil {
		logger.Error("get outer proof error: %v %v", outerPeriod, err)
		return nil, false, err
	}
	if !ok {
		logger.Warn("no find outer proof: %v", outerPeriod)
		return nil, false, nil
	}

	var finalUpdate proverType.FinalityUpdate
	err = common.ParseObj(currentFinalityUpdate.Data, &finalUpdate)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	finalUpdate.Version = currentFinalityUpdate.Version

	currentSyncCommitUpdate, ok, err := b.GetUnitData(outerPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Error("no find sync committee update: %v", period)
		return nil, false, nil
	}

	var scUpdate proverType.SyncCommitteeUpdate
	err = common.ParseObj(currentSyncCommitUpdate, &scUpdate)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	request := rpc.BlockHeaderFinalityRequest{
		Index:            slot,
		GenesisSCSSZRoot: fmt.Sprintf("%x", genesisId),
		RecursiveProof:   recursiveProof.Proof,
		RecursiveWitness: recursiveProof.Witness,
		OuterProof:       outerProof.Proof,
		OuterWitness:     outerProof.Witness,
		FinalityUpdate:   &finalUpdate,
		ScUpdate:         &scUpdate,
	}
	return &request, true, nil

}
