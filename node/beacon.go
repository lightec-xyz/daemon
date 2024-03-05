package node

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"time"
)

var _ IBeaconAgent = (*BeaconAgent)(nil)

type BeaconAgent struct {
	beaconClient      *beacon.Client
	fileStore         *FileStore
	zkProofRequest    chan ZkProofRequest
	name              string
	genesisSyncPeriod uint64
	beaconFetch       *BeaconFetch
	stateCache        *BeaconCache
}

func NewBeaconAgent(cfg NodeConfig, beaconClient *beacon.Client, fetchDaaResp chan FetchDataResponse) (IBeaconAgent, error) {
	fileStore, err := NewFileStore(cfg.DataDir)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	beaconFetch, err := NewBeaconFetch(beaconClient, fileStore, fetchDaaResp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	beaconAgent := &BeaconAgent{
		fileStore:    fileStore,
		beaconClient: beaconClient,
		name:         "beaconAgent",
		beaconFetch:  beaconFetch,
		stateCache:   NewBeaconCache(),
	}
	return beaconAgent, nil
}

func (b *BeaconAgent) Init() error {
	err := b.initGenesis()
	if err != nil {
		return err
	}
	// todo check data
	return err
}

func (b *BeaconAgent) initGenesis() error {
	existsGenesisProof, err := b.fileStore.CheckGenesisProof()
	if err != nil {
		logger.Error("check genesis proof error: %v", err)
		return err
	}
	if existsGenesisProof {
		return err
	}
	err = b.CheckAndNewProofRequest(b.genesisSyncPeriod, SyncComGenesisType)
	if err != nil {
		logger.Error("check genesis proof error: %v", err)
		return nil
	}
	return nil

}

func (b *BeaconAgent) ScanSyncPeriod() error {
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
	if currentPeriod >= latestSyncPeriod {
		return nil
	}
	for index := currentPeriod + 1; index <= latestSyncPeriod; index++ {
		log.Info("beacon parse sync period: %d", index)
		// todo
		for {
			canAddRequest := b.beaconFetch.canNewRequest()
			if canAddRequest {
				b.beaconFetch.NewUpdateRequest(index)
				break
			} else {
				time.Sleep(10 * time.Second)
			}
		}
		err := b.fileStore.StoreLatestPeriod(index)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (b *BeaconAgent) CheckAndNewProofRequest(period uint64, reqType ZkProofType) error {
	proofExists, err := b.CheckProofExists(period, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if proofExists {
		logger.Error("%v %v proof exists", period, reqType)
		return nil
	}
	canSend, err := b.CheckRequestStatus(period, reqType)
	if err != nil {
		logger.Info("can`t send proof: %v %v", period, reqType)
		return err
	}
	if !canSend {
		logger.Info("can`t send proof: %v %v", period, reqType)
		return nil
	}
	data, preData, prepareDataOk, err := b.prepareRequestData(period, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if !prepareDataOk {
		logger.Info("can`t generate proof: %v %v", period, reqType)
		return nil
	}
	err = b.cacheRequestStatus(period, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	b.zkProofRequest <- ZkProofRequest{
		period:  period,
		reqType: reqType,
		data:    data,
		preData: preData,
	}

	return nil
}

func (b *BeaconAgent) CheckRequestStatus(period uint64, reqType ZkProofType) (bool, error) {
	// todo
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
		genesisData, err := b.fileStore.GetGenesisData()
		if err != nil {
			logger.Error(err.Error())
			return nil, false, err
		}
		return genesisData, true, nil
	} else {
		// todo add new request
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
		// todo add new request
		b.beaconFetch.NewUpdateRequest(period)
		return nil, nil, false, nil
	}
	prePeriod := period - 1
	preUpdateExists, err := b.fileStore.CheckUpdate(prePeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, false, err
	}
	if !preUpdateExists {
		// todo add new request
		b.beaconFetch.NewUpdateRequest(prePeriod)
		return nil, nil, false, nil
	}
	updateData, err := b.fileStore.GetUpdateData(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, false, err
	}
	preUpdateData, err := b.fileStore.GetUpdateData(prePeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, false, err
	}
	return updateData, preUpdateData, true, nil

}

func (b *BeaconAgent) GetRecursiveData(period uint64) ([]byte, []byte, bool, error) {
	existUnitProof, err := b.fileStore.CheckUnitProof(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, false, err
	}
	if !existUnitProof {
		// todo add new request
		return nil, nil, false, nil
	}
	perPeriod := period - 1
	existsRecursiveProof, err := b.fileStore.CheckRecursiveProof(perPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, false, err
	}
	if !existsRecursiveProof {
		// todo add new request
		return nil, nil, false, nil
	}
	unitData, err := b.fileStore.GetUnitData(period)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, false, err
	}
	recursiveData, err := b.fileStore.GetRecursiveData(perPeriod)
	if err != nil {
		logger.Error(err.Error())
		return nil, nil, false, err
	}
	return unitData, recursiveData, true, nil
}

func (b *BeaconAgent) prepareRequestData(period uint64, reqType ZkProofType) ([]byte, []byte, bool, error) {
	//todo
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
	//TODO implement me
	panic("implement me")

}

func (b *BeaconAgent) FetchResponse(req FetchDataResponse) error {
	switch req.reqType {
	case SyncComGenesisType:
		err := b.CheckAndNewProofRequest(req.period, SyncComGenesisType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	case SyncComUnitType:
		err := b.CheckAndNewProofRequest(req.period, SyncComUnitType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	default:
		return fmt.Errorf("never should happen : %v %v", req.period, req.reqType)
	}

	return nil
}

func (b *BeaconAgent) ProofResponse(resp ZkProofResponse) error {
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
	nextPeriod := currentPeriod + 1
	switch resp.zkProofType {
	case SyncComGenesisType:
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
		err := b.fileStore.StoreUnitProof(currentPeriod, resp.proof)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		err = b.CheckAndNewProofRequest(nextPeriod, SyncComRecursiveType)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	case SyncComRecursiveType:
		err := b.fileStore.StoreUnitProof(currentPeriod, resp.proof)
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
