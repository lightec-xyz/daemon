package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/lightec-xyz/daemon/store"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"strconv"
	"strings"
)

var _ IBeaconAgent = (*BeaconAgent)(nil)

type BeaconAgent struct {
	beaconClient   *beacon.Client
	fileStore      *FileStorage
	zkProofRequest chan []*common.ZkProofRequest
	name           string
	apiClient      *apiclient.Client
	store          store.IStore
	genesisPeriod  uint64
	genesisSlot    uint64
	stateCache     *CacheState
}

func NewBeaconAgent(store store.IStore, beaconClient *beacon.Client, apiClient *apiclient.Client, zkProofReq chan []*common.ZkProofRequest,
	fileStore *FileStorage, genesisSlot, genesisPeriod uint64, fetchDataResp chan FetchDataResponse) (IBeaconAgent, error) {
	beaconAgent := &BeaconAgent{
		fileStore:      fileStore,
		beaconClient:   beaconClient,
		apiClient:      apiClient,
		name:           "beaconAgent",
		store:          store,
		stateCache:     NewCacheState(),
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
	slot, exists, err := ReadLatestBeaconSlot(b.store)
	if err != nil {
		logger.Error("read latest slot error: %v", err)
		return err
	}
	if !exists || slot < b.genesisSlot {
		err := WriteLatestBeaconSlot(b.store, b.genesisSlot)
		if err != nil {
			logger.Error("write latest slot error: %v", err)
			return err
		}
	}
	return err
}

func (b *BeaconAgent) ScanBlock() error {
	slot, ok, err := ReadLatestBeaconSlot(b.store)
	if err != nil {
		logger.Error("read latest slot error: %v", err)
		return err
	}
	if !ok {
		return nil
	}
	headSlot, err := beacon.GetHeadSlot(b.apiClient)
	if err != nil {
		logger.Error("get head slot error: %v", err)
		return err
	}
	if headSlot <= slot {
		logger.Warn("head slot %v, dbSlot: %v", headSlot, slot)
		return nil
	}
	for index := slot + 1; index <= headSlot; index++ {
		logger.Debug("beacon parse index: %v", index)
		slotMapInfo, err := beacon.GetEth1MapToEth2(b.apiClient, index)
		if err != nil {
			logger.Error("get eth1 map to eth2 error: %v %v ", index, err)
			if strings.Contains(err.Error(), "404 NotFound response") { // todo
				logger.Warn("no find beacon slot %v info", index)
				continue
			}
			return err
		}
		err = b.saveSlotInfo(slotMapInfo)
		if err != nil {
			logger.Error("parse slot info error: %v %v ", index, err)
			return err
		}
		err = WriteLatestBeaconSlot(b.store, index)
		if err != nil {
			logger.Error("write latest slot error: %v %v ", index, err)
			return err
		}
	}
	return nil
}

func (b *BeaconAgent) saveSlotInfo(slotInfo *beacon.Eth1MapToEth2) error {
	logger.Debug("beacon slot: %v <-> eth number: %v", slotInfo.BlockSlot, slotInfo.BlockNumber)
	err := WriteBeaconSlot(b.store, slotInfo.BlockNumber, slotInfo.BlockSlot)
	if err != nil {
		logger.Error("write slot error: %v %v %v ", slotInfo.BlockNumber, slotInfo.BlockSlot, err)
		return err
	}
	err = WriteBeaconEthNumber(b.store, slotInfo.BlockSlot, slotInfo.BlockNumber)
	if err != nil {
		logger.Error("write eth number error: %v %v %v", slotInfo.BlockSlot, slotInfo.BlockNumber, err)
		return err
	}
	return err
}

func (b *BeaconAgent) tryProofRequest(index uint64, reqType common.ZkProofType) error {
	proofId := common.NewProofId(reqType, index, "")
	exists := b.stateCache.Check(proofId)
	if exists {
		return nil
	}
	ok, err := b.checkRequest(index, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if !ok {
		return nil
	}
	proofExists, err := CheckProof(b.fileStore, reqType, index, "")
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if proofExists {
		logger.Warn("%v %v Proof exists", index, reqType.String())
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
	b.stateCache.Store(proofId, true)
	request := common.NewZkProofRequest(reqType, data, index, "")
	b.zkProofRequest <- []*common.ZkProofRequest{request}
	logger.Info("beacon success send Proof request: %v", request.Id())

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
			logger.Error("get RecursiveData %v err: %v", index, err)
			return nil, false, err
		}
		return data, prepared, nil
	default:
		logger.Error(" prepare request Data never should happen : %v %v", index, reqType)
		return nil, false, fmt.Errorf("never should happen : %v %v", index, reqType)
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
	unitProofIndexes, err := b.fileStore.NeedGenUnitProofIndexes()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	for _, index := range unitProofIndexes {
		if index < b.genesisPeriod {
			continue
		}
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
		if skip {
			break
		}
		err := b.tryProofRequest(index, common.SyncComRecursiveType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		skip = true
	}
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
		logger.Warn("no find %v Index update Data", period)
		return nil, false, nil
	}
	// todo
	var update utils.SyncCommitteeUpdate
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
			logger.Warn("get unit Data,no find %v Index update Data", prePeriod)
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
		logger.Error("get commitId %v err: %v", endPeriod, err)
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
		logger.Warn("no find genesis proof")
		return nil, false, nil
	}
	secondProof, secondExists, err := b.fileStore.GetUnitProof(relayPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !secondExists {
		logger.Warn("no find %v Index unit proof ", relayPeriod)
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
		proofId := common.NewProofId(resp.ZkProofType, index, resp.TxHash)
		b.stateCache.Delete(proofId)
	}
	err := StoreZkProof(b.fileStore, resp.ZkProofType, index, resp.TxHash, resp.Proof, resp.Witness)
	if err != nil {
		logger.Error("store proof error: %v", err)
		return err
	}
	return nil
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
