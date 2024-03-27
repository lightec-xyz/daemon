package node

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"sync/atomic"
	"time"
)

var _ IBeaconAgent = (*BeaconAgent)(nil)

type BeaconAgent struct {
	beaconClient   *beacon.Client
	fileStore      *FileStore
	zkProofRequest chan []ZkProofRequest
	name           string
	genesisPeriod  uint64
	beaconFetch    *BeaconFetch
	stateCache     *BeaconCache
	cacheQueue     *Queue
	currentPeriod  *atomic.Uint64
	circuitsFp     common.CircuitsFP
}

func NewBeaconAgent(cfg NodeConfig, beaconClient *beacon.Client, zkProofReq chan []ZkProofRequest, fetchDataResp chan FetchDataResponse) (IBeaconAgent, error) {
	fileStore, err := NewFileStore(cfg.DataDir)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	// todo
	genesisPeriod := uint64(cfg.BeaconInitHeight)
	beaconFetch, err := NewBeaconFetch(beaconClient, fileStore, genesisPeriod, fetchDataResp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	currentPeriod := &atomic.Uint64{}
	currentPeriod.Store(genesisPeriod)
	beaconAgent := &BeaconAgent{
		fileStore:      fileStore,
		beaconClient:   beaconClient,
		name:           "beaconAgent",
		beaconFetch:    beaconFetch,
		stateCache:     NewBeaconCache(),
		zkProofRequest: zkProofReq,
		genesisPeriod:  genesisPeriod,
		cacheQueue:     NewQueue(),
		currentPeriod:  currentPeriod,
	}
	return beaconAgent, nil
}

func (b *BeaconAgent) Init() error {
	logger.Info("beacon agent init")
	go b.beaconFetch.Fetch()
	_, exists, err := b.fileStore.GetLatestPeriod()
	if err != nil {
		logger.Error("check latest Period error: %v", err)
		return err
	}
	if !exists {
		err := b.fileStore.StoreLatestPeriod(b.genesisPeriod)
		if err != nil {
			logger.Error("store latest Period error: %v", err)
			return err
		}
	}
	existsGenesisProof, err := b.fileStore.CheckGenesisProof()
	if err != nil {
		logger.Error("check genesis Proof error: %v", err)
		return err
	}
	if existsGenesisProof {
		log.Info("genesis Proof exists")
		return nil
	}
	err = b.tryProofRequest(b.genesisPeriod, SyncComGenesisType)
	if err != nil {
		logger.Error("check genesis Proof error: %v", err)
		return nil
	}
	// todo check data
	return err
}

func (b *BeaconAgent) CheckRecoverData() error {
	recoverUpdateFiles, err := b.fileStore.RecoverUpdateFiles()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	recoverUnitProofFiles, err := b.fileStore.RecoverUnitProofFiles()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	recursiveProofFiles, err := b.fileStore.RecoverRecursiveProofFiles()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	for _, updatePeriod := range recoverUpdateFiles {
		b.beaconFetch.NewUpdateRequest(updatePeriod)
	}
	for _, unitPeriod := range recoverUnitProofFiles {
		err := b.tryProofRequest(unitPeriod, SyncComUnitType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	for _, recursivePeriod := range recursiveProofFiles {
		err := b.tryProofRequest(recursivePeriod, SyncComRecursiveType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	return nil

}

func (b *BeaconAgent) ScanSyncPeriod() error {
	logger.Debug("beacon scan sync Period")
	currentPeriod, ok, err := b.fileStore.GetLatestPeriod()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if !ok {
		return fmt.Errorf("get latest Period error")
	}
	latestSyncPeriod, err := b.beaconClient.GetLatestSyncPeriod()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	latestSyncPeriod = latestSyncPeriod - 1
	logger.Info("current Period: %d, latest sync Period: %d", currentPeriod, latestSyncPeriod)
	if currentPeriod >= latestSyncPeriod {
		return nil
	}
	for index := currentPeriod; index <= latestSyncPeriod; index++ {
		logger.Debug("beacon scan Period: %d", index)
		for {
			canAddRequest := b.beaconFetch.canNewRequest()
			if canAddRequest {
				b.beaconFetch.NewUpdateRequest(index)
				err := b.fileStore.StoreLatestPeriod(index)
				if err != nil {
					log.Error(err.Error())
					return err
				}
				b.currentPeriod.Store(index)
				break
			} else {
				time.Sleep(10 * time.Second)
			}
		}
	}
	return nil
}

// todo maybe use for  to replace recursive
func (b *BeaconAgent) tryProofRequest(period uint64, reqType ZkProofType) error {
	currentPeriod := b.currentPeriod.Load()
	if period > currentPeriod {
		logger.Warn("wait for new Period,current: %v, reqPeriod: %v", currentPeriod, period)
		return nil
	}
	logger.Debug("beacon check and new Proof request: %v %v", period, reqType.String())
	proofExists, err := b.CheckProofExists(period, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if proofExists {
		logger.Warn("%v %v Proof exists", period, reqType.String())
		return nil
	}
	existsRequest, err := b.CheckProofRequestStatus(period, reqType)
	if err != nil {
		logger.Info("can`t send Proof: %v %v", period, reqType.String())
		return err
	}
	if existsRequest {
		logger.Warn("Proof request exists，%v %v  skip it now", period, reqType.String())
		return nil
	}
	data, prepareDataOk, err := b.prepareProofRequestData(period, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if !prepareDataOk {
		logger.Warn("Proof request data haven`t prepared now ,%v %v  can`t generate Proof", period, reqType.String())
		return nil
	}
	err = b.cacheProofRequestStatus(period, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	b.zkProofRequest <- []ZkProofRequest{
		{
			period:  period,
			reqType: reqType,
			data:    data,
		},
	}

	return nil
}

func (b *BeaconAgent) CheckProofRequestStatus(period uint64, reqType ZkProofType) (bool, error) {
	switch reqType {
	case SyncComGenesisType:
		return b.stateCache.CheckGenesis(), nil
	case SyncComUnitType:
		return b.stateCache.CheckUnit(period), nil
	case SyncComRecursiveType:
		return b.stateCache.CheckRecursive(period), nil
	default:
		return false, fmt.Errorf("check request status never should happen: %v %v", period, reqType)
	}
}

func (b *BeaconAgent) cacheProofRequestStatus(period uint64, reqType ZkProofType) error {
	logger.Debug("beacon cache request status: %v %v", period, reqType.String())
	switch reqType {
	case SyncComGenesisType:
		return b.stateCache.StoreGenesis()
	case SyncComUnitType:
		return b.stateCache.StoreUnit(period)
	case SyncComRecursiveType:
		return b.stateCache.StoreRecursive(period)
	default:
		return fmt.Errorf("cache request status never should happen: %v %v", period, reqType)
	}
}

func (b *BeaconAgent) prepareProofRequestData(period uint64, reqType ZkProofType) (data interface{}, prepared bool, err error) {
	switch reqType {
	case SyncComGenesisType:
		data, prepared, err = b.GetGenesisRaw()
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return data, prepared, nil
	case SyncComUnitType:
		data, prepared, err = b.GetUnitData(period)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return data, prepared, nil
	case SyncComRecursiveType:
		data, prepared, err = b.GetRecursiveData(period)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return data, prepared, nil
	default:
		logger.Error(" prepare request data never should happen : %v %v", period, reqType)
		return nil, false, fmt.Errorf("never should happen : %v %v", period, reqType)
	}
}

func (b *BeaconAgent) CheckProofExists(period uint64, reqType ZkProofType) (bool, error) {
	switch reqType {
	case SyncComGenesisType:
		return b.fileStore.CheckGenesisProof()
	case SyncComUnitType:
		return b.fileStore.CheckUnitProof(period)
	case SyncComRecursiveType:
		return b.fileStore.CheckRecursiveProof(period)
	default:
		logger.Error("check Proof exists never should happen : %v %v", period, reqType)
		return false, fmt.Errorf("check Proof exists never should happen : %v %v", period, reqType)
	}
}

func (b *BeaconAgent) CheckData() error {
	return nil
}

func (b *BeaconAgent) FetchDataResponse(req FetchDataResponse) error {
	logger.Debug("beacon fetch response fetchType: %v, Period: %v", req.UpdateType.String(), req.period)
	switch req.UpdateType {
	case GenesisUpdateType:
		err := b.tryProofRequest(req.period, SyncComGenesisType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	case PeriodUpdateType:
		err := b.tryProofRequest(req.period, SyncComUnitType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	default:
		return fmt.Errorf("never should happen : %v %v", req.period, req.UpdateType)
	}

	return nil
}

func (b *BeaconAgent) GetGenesisRaw() (interface{}, bool, error) {
	genesisId, ok, err := b.fileStore.GetSyncCommitRootID(b.genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v update no find,start new request", b.genesisPeriod)
		b.beaconFetch.NewUpdateRequest(b.genesisPeriod)
		return nil, false, nil
	}

	nextPeriod := b.genesisPeriod + 1
	firstId, ok, err := b.fileStore.GetSyncCommitRootID(nextPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v update no find,start new request", nextPeriod)
		b.beaconFetch.NewUpdateRequest(nextPeriod)
		return nil, false, nil
	}
	secondPeriod := nextPeriod + 1
	secondId, ok, err := b.fileStore.GetSyncCommitRootID(secondPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v update no find,start new request", secondPeriod)
		b.beaconFetch.NewUpdateRequest(secondPeriod)
		return nil, false, nil
	}

	var firstProof StoreProof
	exists, err := b.fileStore.GetUnitProof(b.genesisPeriod, &firstProof)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("get genesis data,first proof not exists: %v period", b.genesisPeriod)
		err := b.tryProofRequest(b.genesisPeriod, SyncComUnitType)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
	}

	var secondProof StoreProof
	exists, err = b.fileStore.GetUnitProof(nextPeriod, &secondProof)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}

	if !exists {
		logger.Warn("get genesis data,second proof not exists: %v period", nextPeriod)
		err := b.tryProofRequest(nextPeriod, SyncComUnitType)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return nil, false, nil
	}

	genesisProofParam := GenesisProofParam{
		FirstProof:    firstProof.Proof,
		FirstWitness:  firstProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		GenesisId:     genesisId,
		FirstId:       firstId,
		SecondId:      secondId,
		RecursiveFp:   b.circuitsFp.RecursiveFp,
	}
	return genesisProofParam, true, nil

}

func (b *BeaconAgent) GetUnitData(period uint64) (interface{}, bool, error) {
	var currentPeriodUpdate structs.LightClientUpdate
	exists, err := b.fileStore.GetUpdate(period, &currentPeriodUpdate)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v period update data, send new update request", period)
		b.beaconFetch.NewUpdateRequest(period)
		return nil, false, nil
	}
	if b.genesisPeriod == period {
		var genesisData structs.LightClientBootstrapResponse
		genesisExists, err := b.fileStore.GetGenesisUpdate(&genesisData)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		if !genesisExists {
			logger.Warn("no find genesis update data, send new update request")
			b.beaconFetch.GenesisUpdateRequest()
			return nil, false, nil
		}
		return UnitProofParam{
			AttestedHeader:          currentPeriodUpdate.AttestedHeader,
			CurrentSyncCommittee:    genesisData.Data.CurrentSyncCommittee, // todo
			SyncAggregate:           currentPeriodUpdate.SyncAggregate,
			FinalizedHeader:         currentPeriodUpdate.FinalizedHeader,
			FinalityBranch:          currentPeriodUpdate.FinalityBranch,
			NextSyncCommittee:       currentPeriodUpdate.NextSyncCommittee,
			NextSyncCommitteeBranch: currentPeriodUpdate.NextSyncCommitteeBranch,
			SignatureSlot:           currentPeriodUpdate.SignatureSlot,
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
			logger.Warn("get unit data,no find %v period update data, send new update request", prePeriod)
			b.beaconFetch.NewUpdateRequest(prePeriod)
			return nil, false, nil
		}
		return UnitProofParam{
			AttestedHeader:          currentPeriodUpdate.AttestedHeader,
			CurrentSyncCommittee:    perUpdateData.Data.NextSyncCommittee, // TODO
			SyncAggregate:           currentPeriodUpdate.SyncAggregate,
			FinalizedHeader:         currentPeriodUpdate.FinalizedHeader,
			FinalityBranch:          currentPeriodUpdate.FinalityBranch,
			NextSyncCommittee:       currentPeriodUpdate.NextSyncCommittee,
			NextSyncCommitteeBranch: currentPeriodUpdate.NextSyncCommitteeBranch,
			SignatureSlot:           currentPeriodUpdate.SignatureSlot,
		}, true, nil
	}

}

func (b *BeaconAgent) GetRecursiveData(period uint64) (interface{}, bool, error) {
	// todo
	if period == b.genesisPeriod {
		return b.getRecursiveGenesisData(b.genesisPeriod)
	} else {
		return b.getRecursiveData(period)
	}

}

func (b *BeaconAgent) getRecursiveData(period uint64) (interface{}, bool, error) {
	//todo
	genesisId, ok, err := b.fileStore.GetSyncCommitRootID(b.genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v period update no find,start new request", b.genesisPeriod)
		b.beaconFetch.NewUpdateRequest(b.genesisPeriod)
		return nil, false, nil
	}
	relayId, ok, err := b.fileStore.GetSyncCommitRootID(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v period update no find,start new request", period)
		b.beaconFetch.NewUpdateRequest(period)
		return nil, false, nil
	}
	endPeriod := period + 1
	endId, ok, err := b.fileStore.GetSyncCommitRootID(endPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v update no find,start new request", endPeriod)
		b.beaconFetch.NewUpdateRequest(endPeriod)
		return nil, false, nil
	}
	var secondProof StoreProof
	exists, err := b.fileStore.GetUnitProof(period, &secondProof)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v unit proof data, send new proof request", period)
		err := b.tryProofRequest(period, SyncComUnitType)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return nil, false, nil
	}

	var firstProof StoreProof
	prePeriod := period - 1
	exists, err = b.fileStore.GetRecursiveProof(prePeriod, &firstProof)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !exists {
		logger.Warn("no find %v period recursive data, send new proof request", prePeriod)
		err := b.tryProofRequest(prePeriod, SyncComRecursiveType)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return nil, false, nil
	}

	return RecursiveProofParam{
		Choice:        "recursive",
		FirstProof:    firstProof.Proof,
		FirstWitness:  firstProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		BeginId:       genesisId,
		RelayId:       relayId,
		EndId:         endId,
		RecursiveFp:   b.circuitsFp.RecursiveFp,
	}, true, nil
}

func (b *BeaconAgent) getRecursiveGenesisData(period uint64) (interface{}, bool, error) {
	// todo
	genesisId, ok, err := b.fileStore.GetSyncCommitRootID(b.genesisPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v period update no find , send new update request", b.genesisPeriod)
		b.beaconFetch.NewUpdateRequest(b.genesisPeriod)
		return nil, false, nil
	}
	relayPeriod := b.genesisPeriod + 2
	relayId, ok, err := b.fileStore.GetSyncCommitRootID(relayPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v  period update no find , send new update request", relayPeriod)
		b.beaconFetch.NewUpdateRequest(relayPeriod)
		return nil, false, nil
	}
	endPeriod := b.genesisPeriod + 3
	endId, ok, err := b.fileStore.GetSyncCommitRootID(endPeriod)
	if err != nil {
		log.Error(err.Error())
		return nil, false, err
	}
	if !ok {
		logger.Warn("get %v period update no find , send new update request", endPeriod)
		b.beaconFetch.NewUpdateRequest(endPeriod)
		return nil, false, nil
	}

	var fistProof StoreProof
	firstExists, err := b.fileStore.GetGenesisProof(&fistProof)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !firstExists {
		logger.Warn("no find genesis proof ,start new proof request")
		err := b.tryProofRequest(period, SyncComGenesisType)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return nil, false, nil
	}
	var secondProof StoreProof
	secondExists, err := b.fileStore.GetUnitProof(relayPeriod, &secondProof)
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if !secondExists {
		logger.Warn("no find %v period unit proof , send new proof request", relayPeriod)
		err := b.tryProofRequest(relayPeriod, SyncComUnitType)
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return nil, false, nil
	}

	return RecursiveProofParam{
		Choice:        "genesis",
		FirstProof:    fistProof.Proof,
		FirstWitness:  fistProof.Witness,
		SecondProof:   secondProof.Proof,
		SecondWitness: secondProof.Witness,
		BeginId:       genesisId,
		RelayId:       relayId,
		EndId:         endId,
		RecursiveFp:   b.circuitsFp.RecursiveFp,
	}, true, nil

}

func (b *BeaconAgent) ProofResponse(resp ZkProofResponse) error {
	logger.Info("beacon Proof response type: %v, Period: %v", resp.ZkProofType.String(), resp.Period)
	currentPeriod := resp.Period
	b.deleteCacheProofReqStatus(resp.ZkProofType, resp.Period)
	switch resp.ZkProofType {
	case SyncComGenesisType:
		// next  recursive Proof
		err := b.fileStore.StoreGenesisProof(resp.Proof, resp.Witness)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		nextPeriod := currentPeriod + 1
		err = b.tryProofRequest(nextPeriod, SyncComRecursiveType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	case SyncComUnitType:
		// current recursive Proof
		// todo
		err := b.fileStore.StoreUnitProof(currentPeriod, resp.Proof, resp.Witness)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		err = b.tryProofRequest(currentPeriod, SyncComRecursiveType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		//nextPeriod := currentPeriod + 1
		//err = b.tryProofRequest(nextPeriod, SyncComUnitType)
		//if err != nil {
		//	logger.Error(err.Error())
		//	return err
		//}
	case SyncComRecursiveType:
		// next recursive Proof
		err := b.fileStore.StoreRecursiveProof(currentPeriod, resp.Proof, resp.Witness)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		nextPeriod := currentPeriod + 1
		err = b.tryProofRequest(nextPeriod, SyncComRecursiveType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	default:
		return fmt.Errorf("never should happen")
	}
	return nil
}

func (b *BeaconAgent) deleteCacheProofReqStatus(reqType ZkProofType, period uint64) {
	switch reqType {
	case SyncComGenesisType:
		b.stateCache.DeleteGenesis()
	case SyncComUnitType:
		b.stateCache.DeleteUnit(period)
	case SyncComRecursiveType:
		b.stateCache.DeleteRecursive(period)
	default:
		logger.Error("delete cache request status never should happen: %v %v", reqType, period)
		return
	}
}

func (b *BeaconAgent) Close() error {
	err := b.beaconFetch.Close()
	if err != nil {
		log.Error("beacon fetch close error: %v", err)
		return err
	}
	return nil
}

func (b *BeaconAgent) Name() string {
	return b.name
}

func (b *BeaconAgent) Stop() {

}
