package node

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"sync/atomic"
	"time"
)

var _ IBeaconAgent = (*BeaconAgent)(nil)

type BeaconAgent struct {
	beaconClient      *beacon.Client
	fileStore         *FileStore
	zkProofRequest    chan []ZkProofRequest
	name              string
	genesisSyncPeriod uint64
	beaconFetch       *BeaconFetch
	stateCache        *BeaconCache
	cacheQueue        *Queue
	currentPeriod     *atomic.Uint64
}

func NewBeaconAgent(cfg NodeConfig, beaconClient *beacon.Client, zkProofReq chan []ZkProofRequest, fetchDataResp chan FetchDataResponse) (IBeaconAgent, error) {
	fileStore, err := NewFileStore(cfg.DataDir)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	// todo
	genesisPeriod := uint64(cfg.BtcInitHeight / 8192)

	beaconFetch, err := NewBeaconFetch(beaconClient, fileStore, genesisPeriod, fetchDataResp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	currentPeriod := &atomic.Uint64{}
	currentPeriod.Store(genesisPeriod)
	beaconAgent := &BeaconAgent{
		fileStore:         fileStore,
		beaconClient:      beaconClient,
		name:              "beaconAgent",
		beaconFetch:       beaconFetch,
		stateCache:        NewBeaconCache(),
		zkProofRequest:    zkProofReq,
		genesisSyncPeriod: genesisPeriod,
		cacheQueue:        NewQueue(),
		currentPeriod:     currentPeriod,
	}
	return beaconAgent, nil
}

func (b *BeaconAgent) Init() error {
	logger.Info("beacon agent init")
	go b.beaconFetch.Fetch()
	existsLatestPeriod, err := b.fileStore.CheckLatestPeriod()
	if err != nil {
		logger.Error("check latest period error: %v", err)
		return err
	}
	if !existsLatestPeriod {
		err := b.fileStore.StoreLatestPeriod(b.genesisSyncPeriod)
		if err != nil {
			logger.Error("store latest period error: %v", err)
			return err
		}
	}
	existsGenesisProof, err := b.fileStore.CheckGenesisProof()
	if err != nil {
		logger.Error("check genesis proof error: %v", err)
		return err
	}
	if existsGenesisProof {
		log.Info("genesis proof exists")
		return nil
	}
	err = b.tryProofRequest(b.genesisSyncPeriod, SyncComGenesisType)
	if err != nil {
		logger.Error("check genesis proof error: %v", err)
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
	logger.Debug("beacon scan sync period")
	currentPeriod, err := b.fileStore.GetLatestPeriod()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	latestSyncPeriod, err := b.beaconClient.GetLatestSyncPeriod()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	latestSyncPeriod = latestSyncPeriod - 1
	logger.Info("current period: %d, latest sync period: %d", currentPeriod, latestSyncPeriod)
	if currentPeriod >= latestSyncPeriod {
		return nil
	}
	for index := currentPeriod; index <= latestSyncPeriod; index++ {
		logger.Debug("beacon scan period: %d", index)
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
		logger.Warn("wait for new period,current: %v, reqPeriod: %v", currentPeriod, period)
		return nil
	}
	logger.Debug("beacon check and new proof request: %v %v", period, reqType.String())
	proofExists, err := b.CheckProofExists(period, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if proofExists {
		logger.Warn("%v %v proof exists", period, reqType.String())
		return nil
	}
	existsRequest, err := b.CheckProofRequestStatus(period, reqType)
	if err != nil {
		logger.Info("can`t send proof: %v %v", period, reqType.String())
		return err
	}
	if existsRequest {
		logger.Warn("proof request exists，%v %v  skip it now", period, reqType.String())
		return nil
	}
	data, preData, prepareDataOk, err := b.prepareProofRequestData(period, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if !prepareDataOk {
		logger.Warn("proof request data haven`t prepared now ,%v %v  can`t generate proof", period, reqType.String())
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
			preData: preData,
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

func (b *BeaconAgent) prepareProofRequestData(period uint64, reqType ZkProofType) (data, preData []byte, prepared bool, err error) {
	switch reqType {
	case SyncComGenesisType:
		data, prepared, err = b.GetGenesisRaw()
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		return data, nil, prepared, nil
	case SyncComUnitType:
		data, preData, prepared, err = b.GetUnitData(period)
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		return data, preData, prepared, nil
	case SyncComRecursiveType:
		data, preData, prepared, err = b.GetRecursiveData(period)
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		return data, preData, prepared, nil
	default:
		logger.Error(" prepare request data never should happen : %v %v", period, reqType)
		return nil, nil, false, fmt.Errorf("never should happen : %v %v", period, reqType)
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
		logger.Error("check proof exists never should happen : %v %v", period, reqType)
		return false, fmt.Errorf("check proof exists never should happen : %v %v", period, reqType)
	}
}

func (b *BeaconAgent) CheckData() error {
	return nil
}

func (b *BeaconAgent) FetchDataResponse(req FetchDataResponse) error {
	logger.Debug("beacon fetch response: %v %v", req.period, req.UpdateType.String())
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

func (b *BeaconAgent) GetGenesisRaw() ([]byte, bool, error) {
	exists, err := b.fileStore.CheckGenesisUpdate()
	if err != nil {
		logger.Error(err.Error())
		return nil, false, err
	}
	if exists {
		genesisData, err := b.fileStore.GetGenesisUpdateData()
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return genesisData, true, nil
	} else {
		logger.Warn("no find genesis update data, send new genesis request")
		b.beaconFetch.GenesisUpdateRequest()
		return nil, false, nil
	}
}

func (b *BeaconAgent) GetUnitData(period uint64) ([]byte, []byte, bool, error) {
	exists, err := b.fileStore.CheckUpdate(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, false, err
	}
	if !exists {
		logger.Warn("no find %v update data, send new update request", period)
		b.beaconFetch.NewUpdateRequest(period)
		return nil, nil, false, nil
	}
	updateData, err := b.fileStore.GetUpdateData(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, false, err
	}
	if b.genesisSyncPeriod == period {
		existsGenesisUpdate, err := b.fileStore.CheckGenesisUpdate()
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		if !existsGenesisUpdate {
			logger.Warn("get unit data,no find genesis update data, send new genesis request")
			b.beaconFetch.GenesisUpdateRequest()
			return nil, nil, false, nil
		}
		genesisData, err := b.fileStore.GetGenesisUpdateData()
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		return updateData, genesisData, true, nil
	} else {
		prePeriod := period - 1
		if prePeriod < b.genesisSyncPeriod {
			logger.Error("should never happen: %v", prePeriod)
			return nil, nil, false, nil
		}
		preUpdateExists, err := b.fileStore.CheckUpdate(prePeriod)
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		if !preUpdateExists {
			logger.Warn("get unit data,no find %v update data, send new update request", prePeriod)
			b.beaconFetch.NewUpdateRequest(prePeriod)
			return nil, nil, false, nil
		}
		preUpdateData, err := b.fileStore.GetUpdateData(prePeriod)
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		return updateData, preUpdateData, true, nil
	}

}

func (b *BeaconAgent) GetRecursiveData(period uint64) ([]byte, []byte, bool, error) {
	existUnitProof, err := b.fileStore.CheckUnitProof(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, false, err
	}
	if !existUnitProof {
		err := b.tryProofRequest(period, SyncComUnitType)
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		return nil, nil, false, nil
	}
	unitProofData, err := b.fileStore.GetUnitProofData(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, false, err
	}

	if b.genesisSyncPeriod == period {
		existsGenesisProof, err := b.fileStore.CheckGenesisProof()
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		if !existsGenesisProof {
			logger.Warn("get recursive data,no find genesis proof data, send new genesis request")
			err := b.tryProofRequest(period, SyncComGenesisType)
			if err != nil {
				logger.Error(err.Error())
				return nil, nil, false, err
			}
			return nil, nil, false, nil
		}
		genesisProofData, err := b.fileStore.GetGenesisProofData()
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		return unitProofData, genesisProofData, true, nil
	} else {
		prePeriod := period - 1
		if prePeriod < b.genesisSyncPeriod {
			logger.Error("should never happen: %v", prePeriod)
			return nil, nil, false, nil
		}
		existsRecursiveProof, err := b.fileStore.CheckRecursiveProof(prePeriod)
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		if !existsRecursiveProof {
			logger.Warn("get recursive data,no find %v proof data, send new proof request", prePeriod)
			err := b.tryProofRequest(prePeriod, SyncComRecursiveType)
			if err != nil {
				logger.Error(err.Error())
				return nil, nil, false, err
			}
			return nil, nil, false, nil
		}
		recursiveData, err := b.fileStore.GetRecursiveProofData(prePeriod)
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		return unitProofData, recursiveData, true, nil
	}

}

func (b *BeaconAgent) ProofResponse(resp ZkProofResponse) error {
	logger.Info("beacon proof response type: %v, period: %v", resp.zkProofType.String(), resp.period)
	currentPeriod := resp.period
	b.deleteCacheProofReqStatus(resp.zkProofType, resp.period)
	switch resp.zkProofType {
	case SyncComGenesisType:
		// next  recursive proof
		nextPeriod := currentPeriod + 1
		err := b.fileStore.StoreGenesisProof(resp.proof)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		err = b.tryProofRequest(nextPeriod, SyncComRecursiveType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	case SyncComUnitType:
		// current recursive proof
		// todo
		err := b.fileStore.StoreUnitProof(currentPeriod, resp.proof)
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
		// next recursive proof
		nextPeriod := currentPeriod + 1
		err := b.fileStore.StoreRecursiveProof(currentPeriod, resp.proof)
		if err != nil {
			log.Error(err.Error())
			return err
		}
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
