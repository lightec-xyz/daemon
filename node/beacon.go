package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/lightec-xyz/daemon/store"
	apiclient "github.com/lightec-xyz/provers/utils/api-client"
	"strings"
	"time"
)

var _ IAgent = (*beaconAgent)(nil)

type beaconAgent struct {
	fileStore      *FileStorage
	chainStore     *ChainStore
	name           string
	apiClient      *apiclient.Client // todo merge apiClient and beaconClient
	beaconClient   *beacon.Client
	initBeaconSlot uint64
	reScan         bool
	init           bool
	curSlot        uint64
}

func NewBeaconAgent(reScan bool, store store.IStore, beaconClient *beacon.Client, apiClient *apiclient.Client, fileStore *FileStorage,
	initBeaconSlot uint64) (IAgent, error) {
	return &beaconAgent{
		fileStore:      fileStore,
		apiClient:      apiClient,
		beaconClient:   beaconClient,
		name:           BeaconAgentName,
		chainStore:     NewChainStore(store),
		initBeaconSlot: initBeaconSlot,
		reScan:         reScan,
	}, nil
}

func (b *beaconAgent) Init() error {
	logger.Info("start init beacon agent")
	slot, exists, err := b.chainStore.ReadLatestBeaconSlot()
	if err != nil {
		logger.Error("read beacon latest slot error: %v", err)
		return err
	}
	if !exists || slot < b.initBeaconSlot || b.reScan {
		err := b.chainStore.WriteLatestBeaconSlot(b.initBeaconSlot)
		if err != nil {
			logger.Error("write beacon latest slot error: %v", err)
			return err
		}
	}
	// for some api is timeout, open a goroutine to get latest finalized slot
	go b.initSotInfo()

	return err
}

func (b *beaconAgent) initSotInfo() {
	for {
		latestFinalizedSlot, err := b.beaconClient.GetLatestFinalizedSlot()
		if err != nil {
			logger.Error("get beacon latest finalized finalizedSlot error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		if latestFinalizedSlot == 0 {
			time.Sleep(2 * time.Second)
			continue
		}

		latestSlot, exists, err := b.fileStore.GetLatestFinalizedSlot()
		if err != nil {
			logger.Error("check beacon latest slot error: %v", err)
			continue
		}
		if !exists || latestSlot < latestFinalizedSlot {
			logger.Info("use beacon genesis slot: %v", latestFinalizedSlot)
			err := b.fileStore.StoreLatestFinalizedSlot(latestFinalizedSlot)
			if err != nil {
				logger.Error("store beacon latest slot error: %v", err)
				continue
			}
		}
		latestPeriod := latestFinalizedSlot / common.SlotPerPeriod
		period, exists, err := b.fileStore.GetLatestPeriod()
		if err != nil {
			logger.Error("check beacon latest period error: %v", err)
			continue
		}
		if !exists || period < latestPeriod {
			logger.Info("beacon start period: %v", latestPeriod)
			err := b.fileStore.StoreLatestPeriod(latestPeriod)
			if err != nil {
				logger.Error("store beacon latest period: %v", err)
				continue
			}
		}
		b.init = true
		logger.Debug("beacon init success")
		return
	}

}

func (b *beaconAgent) ScanBlock() error {
	if !b.init {
		return nil
	}
	slot, ok, err := b.chainStore.ReadLatestBeaconSlot()
	if err != nil {
		logger.Error("read beacon latest slot error: %v", err)
		return err
	}
	if !ok {
		logger.Warn("no find beacon latest slot")
		return fmt.Errorf("no find beacon latest slot")
	}
	headSlot, err := b.beaconClient.GetLatestFinalizedSlot()
	if err != nil {
		logger.Error("get beaconhead slot error: %v", err)
		return err
	}
	if headSlot <= slot {
		logger.Warn("found beacon head slot %v <= dbSlot: %v", headSlot, slot)
		return nil
	}
	for index := slot + 1; index <= headSlot; index++ {
		slotMapInfo, err := b.beaconClient.Eth1MapToEth2(index)
		if err != nil {
			if strings.Contains(err.Error(), "404 NotFound response") { // todo
				logger.Warn("miss beacon slot %v", index)
				continue
			}
			logger.Error("get eth1 map to eth2 error: %v %v ", index, err)
			return err
		}
		err = b.saveSlotInfo(slotMapInfo)
		if err != nil {
			logger.Error("store beacon slot error: %v %v ", index, err)
			return err
		}
		err = b.chainStore.WriteLatestBeaconSlot(index)
		if err != nil {
			logger.Error("write beacon latest slot error: %v %v ", index, err)
			return err
		}
	}
	return nil
}

func (b *beaconAgent) saveSlotInfo(slotInfo *beacon.Eth1MapToEth2) error {
	logger.Debug("beacon slot: %v <-> eth number: %v", slotInfo.BlockSlot, slotInfo.BlockNumber)
	err := b.chainStore.WriteBeaconSlot(slotInfo.BlockNumber, slotInfo.BlockSlot)
	if err != nil {
		logger.Error("write beacon map height %v <-> slot %v,error %v ", slotInfo.BlockNumber, slotInfo.BlockSlot, err)
		return err
	}
	err = b.chainStore.WriteBeaconEthNumber(slotInfo.BlockSlot, slotInfo.BlockNumber)
	if err != nil {
		logger.Error("write eth slot %v <-> height %v, error %v", slotInfo.BlockSlot, slotInfo.BlockNumber, err)
		return err
	}
	return err
}

func (b *beaconAgent) ProofResponse(resp *common.ProofResponse) error {
	logger.Info("beacon proof response type id: %v", resp.ProofId())
	return nil
}

func (b *beaconAgent) Close() error {
	return nil
}

func (b *beaconAgent) Name() string {
	return b.name
}

func (b *beaconAgent) ReScan(height uint64) error {
	return nil
}
func (b *beaconAgent) CheckState() error {
	// beacon sync slot per half hour
	slot, ok, err := b.chainStore.ReadLatestBeaconSlot()
	if err != nil {
		logger.Error("read beacon latest slot error: %v", err)
		return err
	}
	if ok {
		diff := slot - b.curSlot
		if diff < 100 { //normal 150
			logger.Error("beacon sync too slow , node maybe offline,diff:%v, prevSlot:%v curSlot:%v", diff, b.curSlot, slot)
		}

	}
	b.curSlot = slot
	return nil
}
