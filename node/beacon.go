package node

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"sync"
	"sync/atomic"
	"time"
)

var _ IBeaconAgent = (*BeaconAgent)(nil)

type Status int32

const (
	Status_NONE            Status = 0 //apply to ProofG, ProofU, ProofR,  init state
	Status_ProofGenerating Status = 1 //apply to ProofG, ProofU, ProofR, generate proof
	Status_ProofGenerated  Status = 2 //apply to ProofG, ProofU, ProofR, proof generated
)

type ProofG struct {
	period          uint64
	status          Status
	bootStrapStatus DownloadStatus
}

type ProofU struct {
	period           uint64
	isGenesis        bool
	status           Status
	preUpdateStatus  DownloadStatus //if isGenesis is true, preUpdateStatus indicating bootstrap_xx.json status
	currUpdateStatus DownloadStatus
}

type ProofR struct {
	period uint64
	status Status
	proofU ProofU
}

type BeaconAgent struct {
	updatesDir        string
	unitProofDir      string
	recursiveProofDir string

	beaconClient *beacon.Client
	proofClient  rpc.ISyncCommitteeProof
	datadir      string

	targetSyncPeriod atomic.Uint64

	proofG              ProofG
	proofR              []*ProofR
	lastGeneratedProofG atomic.Int64 //indicate which sync period has been generated

	newPeriodArriveCh   chan uint64
	genesisProofReqCh   chan uint64 //notify to generate genesis proof
	unitProofReqCh      chan uint64 //notify to generate unit proof
	recursiveProofReqCh chan uint64

	downloadingUpdate        *sync.Map                                  //TODO(keep), for future use,
	downloadedUpdateChan     chan *structs.LightClientUpdateWithVersion //TODO(keep), for future use,
	generatingGenesisProof   *sync.Map                                  //TODO(keep), for future use,
	genesisProofRespCh       chan *rpc.SyncCommitteeProofResponse       //TODO(keep), for future use,
	generatingUnitProof      *sync.Map                                  //TODO(keep), for future use,
	unitProofRespCh          chan *rpc.SyncCommitteeProofResponse       //TODO(keep), for future use,
	generatingRecursiveProof *sync.Map                                  //TODO(keep), for future use,
	recursiveProofRespCh     chan *rpc.SyncCommitteeProofResponse       //TODO(keep), for future use,

	exitCh chan struct{}
	wg     sync.WaitGroup

	fileStore         *FileStore
	zkProofRequest    chan ZkProofRequest
	name              string
	genesisSyncPeriod uint64
	beaconFetch       *BeaconFetch
}

func NewBeaconAgent(cfg NodeConfig, beaconClient *beacon.Client) (IBeaconAgent, error) {
	fileStore, err := NewFileStore(cfg.DataDir)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	beaconFetch, err := NewBeaconFetch(beaconClient, fileStore)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	beaconAgent := &BeaconAgent{
		fileStore:    fileStore,
		beaconClient: beaconClient,
		name:         "beaconAgent",
		beaconFetch:  beaconFetch,
	}
	return beaconAgent, nil
}

func (b *BeaconAgent) Init() error {
	//TODO implement me
	panic("implement me")
}

func (b *BeaconAgent) initGenesis() error {
	existsGenesisRaw, err := b.fileStore.CheckGenesisUpdate()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if !existsGenesisRaw {
		log.Info("genesis not exists, start get genesis bootstrap")
		bootStrap, err := b.beaconClient.GetBootstrap(uint64(b.genesisSyncPeriod) * 32)
		if err != nil {
			logger.Error("get bootstrap error:%v", err)
			return err
		}
		err = b.fileStore.StoreGenesisUpdate(bootStrap)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	existGenesisProof, err := b.fileStore.CheckGenesisProof()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if !existGenesisProof {
		var genesisUpdate structs.LightClientBootstrapResponse
		err := b.fileStore.GetGenesisUpdate(&genesisUpdate)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		genesisData, err := json.Marshal(genesisUpdate)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
		b.zkProofRequest <- ZkProofRequest{
			reqType: SyncComGenesisType,
			body:    genesisData,
		}
	}
	logger.Info("init genesis complete")
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
			canAddRequest := b.beaconFetch.CanAddRequest()
			if canAddRequest {
				b.beaconFetch.SendRequest(index)
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

func (b *BeaconAgent) P() {

}

func (b *BeaconAgent) CheckAndGenRequest(period uint64, reqType ZkProofType) error {
	proofExists, err := b.CheckProofExists(period, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if proofExists {
		logger.Error("%v %v proof exists", period, reqType)
		return nil
	}
	requestData, canGen, err := b.prepareRequestData(period, reqType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	if !canGen {
		logger.Info("can`t generate proof: %v %v", period, reqType)
		return nil
	}
	b.zkProofRequest <- ZkProofRequest{
		period:  period,
		reqType: reqType,
		body:    requestData,
	}
	return nil
}

func (b *BeaconAgent) prepareRequestData(period uint64, reqType ZkProofType) ([]byte, bool, error) {
	switch reqType {
	case SyncComGenesisType:

	case SyncComUnitType:

	case SyncComRecursiveType:

	default:
		logger.Info("never should happen : %v %v", period, reqType)
		return nil, false, fmt.Errorf("never should happen : %v %v", period, reqType)
	}

	return nil, false, nil
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

func (b *BeaconAgent) UpdateResp(resp UpdateResponse) error {
	err := b.fileStore.StoreUpdate(resp.period, string(resp.data))
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	err = b.CheckAndGenRequest(resp.period, SyncComUnitType)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func (b *BeaconAgent) ProofResp(resp ZkProofResponse) error {
	currentPeriod := resp.period
	if resp.Status != ProofSuccess {
		err := b.CheckAndGenRequest(currentPeriod, resp.zkProofType)
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
		err = b.CheckAndGenRequest(nextPeriod, SyncComRecursiveType)
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
		err = b.CheckAndGenRequest(nextPeriod, SyncComRecursiveType)
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
		err = b.CheckAndGenRequest(nextPeriod, SyncComRecursiveType)
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
	return nil
}

func (b *BeaconAgent) Name() string {
	return b.name
}

func (b *BeaconAgent) Stop() {
	close(b.exitCh)
	b.wg.Wait()
}
