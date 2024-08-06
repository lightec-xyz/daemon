package node

import (
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/lightec-xyz/daemon/store"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"strings"
)

var _ IAgent = (*BeaconAgent)(nil)

type BeaconAgent struct {
	fileStore      *FileStorage
	zkProofRequest chan []*common.ZkProofRequest
	name           string
	apiClient      *apiclient.Client
	store          store.IStore
	genesisPeriod  uint64
	genesisSlot    uint64
	force          bool
}

func NewBeaconAgent(store store.IStore, beaconClient *beacon.Client, apiClient *apiclient.Client, zkProofReq chan []*common.ZkProofRequest,
	fileStore *FileStorage, state *Cache, genesisSlot, genesisPeriod uint64) (IAgent, error) {
	beaconAgent := &BeaconAgent{
		fileStore:      fileStore,
		apiClient:      apiClient,
		name:           BeaconAgentName,
		store:          store,
		zkProofRequest: zkProofReq,
		genesisPeriod:  genesisPeriod,
		genesisSlot:    genesisSlot,
	}
	return beaconAgent, nil
}

func (b *BeaconAgent) Init() error {
	logger.Info("beacon agent init")
	if b.force {
		err := WriteLatestBeaconSlot(b.store, b.genesisSlot)
		if err != nil {
			logger.Error("write latest slot error: %v", err)
			return err
		}
	} else {
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
	}

	latestPeriod, exists, err := b.fileStore.GetLatestPeriod()
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
		err := b.fileStore.StoreLatestFinalizedSlot(b.genesisSlot)
		if err != nil {
			logger.Error("store latest Slot error: %v", err)
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
		//logger.Debug("beacon parse index: %v", index)
		slotMapInfo, err := beacon.GetEth1MapToEth2(b.apiClient, index)
		if err != nil {
			if strings.Contains(err.Error(), "404 NotFound response") { // todo
				logger.Warn("miss beacon slot %v info", index)
				continue
			}
			logger.Error("get eth1 map to eth2 error: %v %v ", index, err)
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

func (b *BeaconAgent) CheckState() error {
	return nil
}

func (b *BeaconAgent) FetchDataResponse(req *FetchResponse) error {
	logger.Debug("beacon fetch response fetchType: %v, Index: %v", req.UpdateType.String(), req.Index)
	return nil
}

func (b *BeaconAgent) ProofResponse(resp *common.ZkProofResponse) error {
	logger.Info("beacon Proof response type: %v, Index: %v", resp.Id(), resp.Index)
	return nil
}

func (b *BeaconAgent) Close() error {
	return nil
}

func (b *BeaconAgent) Name() string {
	return b.name
}
