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

func NewBeaconAgent(cfg NodeConfig, beaconClient *beacon.Client, zkProofReq chan []ZkProofRequest, fetchDaaResp chan FetchDataResponse) (IBeaconAgent, error) {
	fileStore, err := NewFileStore(cfg.DataDir)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	// todo
	genesisPeriod := uint64(cfg.BtcInitHeight / 8192)

	beaconFetch, err := NewBeaconFetch(beaconClient, fileStore, genesisPeriod, fetchDaaResp)
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
	err = b.CheckAndNewProofRequest(b.genesisSyncPeriod, SyncComGenesisType)
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
				break
			} else {
				time.Sleep(10 * time.Second)
			}
		}
	}
	return nil
}

func (b *BeaconAgent) CheckAndNewProofRequest(period uint64, reqType ZkProofType) error {
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
	existsRequest, err := b.CheckRequestStatus(period, reqType)
	if err != nil {
		logger.Info("can`t send proof: %v %v", period, reqType.String())
		return err
	}
	if existsRequest {
		logger.Warn("proof request exists，%v %v  skip it now", period, reqType.String())
		return nil
	}
	data, preData, prepareDataOk, err := b.prepareRequestData(period, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if !prepareDataOk {
		logger.Warn("proof request data haven`t prepared now ,%v %v  can`t generate proof", period, reqType.String())
		return nil
	}
	err = b.cacheRequestStatus(period, reqType)
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

func (b *BeaconAgent) CheckRequestStatus(period uint64, reqType ZkProofType) (bool, error) {
	//todo
	logger.Debug("beacon check request status: %v %v", period, reqType.String())
	switch reqType {
	case SyncComGenesisType:
		return b.stateCache.CheckGenesis(), nil
	case SyncComUnitType:
		return b.stateCache.CheckUnit(period), nil
	case SyncComRecursiveType:
		return b.stateCache.CheckRecursive(period), nil
	default:
		return false, fmt.Errorf("never should happen: %v %v", period, reqType)
	}
}

func (b *BeaconAgent) cacheRequestStatus(period uint64, reqType ZkProofType) error {
	// todo
	logger.Debug("beacon cache request status: %v %v", period, reqType.String())
	switch reqType {
	case SyncComGenesisType:
		return b.stateCache.StoreGenesis()
	case SyncComUnitType:
		return b.stateCache.StoreUnit(period)
	case SyncComRecursiveType:
		return b.stateCache.StoreRecursive(period)
	default:
		return fmt.Errorf("never should happen: %v %v", period, reqType)
	}
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
		// todo add new request
		logger.Debug("no find genesis update data, send new genesis request")
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
		logger.Debug("no find %v update data, send new update request", period)
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
			logger.Debug("get unit data,no find genesis update data, send new genesis request")
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
		preUpdateExists, err := b.fileStore.CheckUpdate(prePeriod)
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		if !preUpdateExists {
			logger.Debug("get unit data,no find %v update data, send new update request", prePeriod)
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
		logger.Debug("no find %v proof data, send new proof request", period)
		err := b.CheckAndNewProofRequest(period, SyncComUnitType)
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

	// todo
	if b.genesisSyncPeriod == period {
		existsGenesisProof, err := b.fileStore.CheckGenesisProof()
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		if !existsGenesisProof {
			// todo add new request
			logger.Debug("get recursive data,no find genesis proof data, send new genesis request")
			err := b.CheckAndNewProofRequest(period, SyncComGenesisType)
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
		perPeriod := period - 1
		existsRecursiveProof, err := b.fileStore.CheckRecursiveProof(perPeriod)
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		if !existsRecursiveProof {
			// todo add new request
			logger.Debug("get recursive data,no find %v proof data, send new proof request", perPeriod)
			err := b.CheckAndNewProofRequest(perPeriod, SyncComRecursiveType)
			if err != nil {
				logger.Error(err.Error())
				return nil, nil, false, err
			}
			return nil, nil, false, nil
		}
		recursiveData, err := b.fileStore.GetRecursiveProofData(perPeriod)
		if err != nil {
			logger.Error(err.Error())
			return nil, nil, false, err
		}
		return unitProofData, recursiveData, true, nil
	}

}

func (b *BeaconAgent) prepareRequestData(period uint64, reqType ZkProofType) ([]byte, []byte, bool, error) {
	//todo
	logger.Debug("beacon prepare request data period:%v,proofType:%v", period, reqType.String())
	switch reqType {
	case SyncComGenesisType:
		genesisRaw, ok, err := b.GetGenesisRaw()
		if err != nil {
			logger.Info(err.Error())
			return nil, nil, false, err
		}
		return genesisRaw, nil, ok, nil
	case SyncComUnitType:
		data, preData, ok, err := b.GetUnitData(period)
		if err != nil {
			logger.Info(err.Error())
			return nil, nil, false, err
		}
		return data, preData, ok, nil
	case SyncComRecursiveType:
		data, preRecursiveData, ok, err := b.GetRecursiveData(period)
		if err != nil {
			logger.Info(err.Error())
			return nil, nil, false, err
		}
		return data, preRecursiveData, ok, nil
	default:
		logger.Info("never should happen : %v %v", period, reqType)
		return nil, nil, false, fmt.Errorf("never should happen : %v %v", period, reqType)
	}
}

func (b *BeaconAgent) CheckProofExists(period uint64, reqType ZkProofType) (bool, error) {
	logger.Debug("beacon check proof exists: %v %v", period, reqType.String())
	switch reqType {
	case SyncComGenesisType:
		return b.fileStore.CheckGenesisProof()
	case SyncComUnitType:
		return b.fileStore.CheckUnitProof(period)
	case SyncComRecursiveType:
		return b.fileStore.CheckRecursiveProof(period)
	default:
		logger.Info("never should happen : %v %v", period, reqType)
		return false, fmt.Errorf("never should happen : %v %v", period, reqType)
	}
}

func (b *BeaconAgent) CheckData() error {
	return nil
}

func (b *BeaconAgent) FetchDataResponse(req FetchDataResponse) error {
	logger.Debug("beacon fetch response: %v %v", req.period, req.UpdateType.String())
	switch req.UpdateType {
	case GenesisUpdateType:
		err := b.CheckAndNewProofRequest(req.period, SyncComGenesisType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	case PeriodUpdateType:
		err := b.CheckAndNewProofRequest(req.period, SyncComUnitType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	default:
		return fmt.Errorf("never should happen : %v %v", req.period, req.UpdateType)
	}

	return nil
}

func (b *BeaconAgent) ProofResponse(resp ZkProofResponse) error {
	logger.Debug("beacon proof response: %v %v", resp.zkProofType.String(), resp.period)
	currentPeriod := resp.period
	if resp.Status != ProofSuccess {
		err := b.CheckAndNewProofRequest(currentPeriod, resp.zkProofType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		return nil
	}
	// next generate proof

	switch resp.zkProofType {
	case SyncComGenesisType:
		// next  recursive proof
		nextPeriod := currentPeriod + 1
		err := b.fileStore.StoreGenesisProof(resp.proof)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		err = b.CheckAndNewProofRequest(nextPeriod, SyncComRecursiveType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	case SyncComUnitType:
		// current recursive proof
		err := b.fileStore.StoreUnitProof(currentPeriod, resp.proof)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		err = b.CheckAndNewProofRequest(currentPeriod, SyncComRecursiveType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	case SyncComRecursiveType:
		// next recursive proof
		nextPeriod := currentPeriod + 1
		err := b.fileStore.StoreRecursiveProof(currentPeriod, resp.proof)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		err = b.CheckAndNewProofRequest(nextPeriod, SyncComRecursiveType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	default:
		return fmt.Errorf("never should happen")
	}
	return nil
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
