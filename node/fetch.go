package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"math/big"
	"sync"
	"sync/atomic"
)

var _ IFetch = (*Fetch)(nil)

type Fetch struct {
	client        *beacon.Client
	genesisSlot   uint64
	genesisPeriod uint64
	fileStore     *FileStorage
	lock          sync.Mutex
	maxReqs       *atomic.Int64
	state         *Cache
	ethFetchResp  chan *FetchResponse
}

func (f *Fetch) Close() error {
	return nil
}

func (f *Fetch) Init() error {
	logger.Debug("init fetch now")
	err := f.Bootstrap()
	if err != nil {
		logger.Error("bootstrap error:%v", err)
		return err
	}
	err = f.FinalityUpdate()
	if err != nil {
		logger.Error("finality update error:%v", err)
		return err
	}
	err = f.StoreLatestPeriod()
	if err != nil {
		logger.Error("store latest Index error:%v", err)
		return err
	}
	return nil
}

func NewFetch(client *beacon.Client, fileStore *FileStorage, genesisSlot uint64, ethFetchResp chan *FetchResponse) (*Fetch, error) {
	maxReqs := atomic.Int64{}
	maxReqs.Store(0)
	return &Fetch{
		client:        client,
		fileStore:     fileStore,
		genesisSlot:   genesisSlot,
		genesisPeriod: genesisSlot / 8192,
		maxReqs:       &maxReqs,
		ethFetchResp:  ethFetchResp, // todo
		state:         NewCacheState(),
	}, nil
}

func (f *Fetch) Bootstrap() error {
	err := f.GetBootStrap(f.genesisSlot)
	if err != nil {
		logger.Error("get bootstrap error:%v %v", f.genesisSlot, err)
		return err
	}
	return nil
}

func (f *Fetch) FinalityUpdate() error {
	//logger.Debug("start finality update")
	err := f.GetFinalityUpdate()
	if err != nil {
		logger.Error("get finality update error:%v", err)
		return err
	}

	return nil
}

func (f *Fetch) StoreLatestPeriod() error {
	period, err := f.client.GetFinalizedSyncPeriod()
	if err != nil {
		logger.Error("get latest Index error:%v", err)
		return err
	}
	err = f.fileStore.StorePeriod(period)
	if err != nil {
		logger.Error("store latest Index error:%v", err)
		return err
	}
	return nil
}

func (f *Fetch) canUpdateReq() bool {
	return f.maxReqs.Load() < 3 // todo
}

func (f *Fetch) LightClientUpdate() error {
	logger.Debug("start light client update")
	err := f.StoreLatestPeriod()
	if err != nil {
		logger.Error("store latest Index error:%v", err)
		return err
	}
	if !f.canUpdateReq() {
		return nil
	}
	updateIndexes, err := f.fileStore.NeedUpdateIndexes()
	if err != nil {
		logger.Error("get update indexes error:%v", err)
		return nil
	}
	for _, index := range updateIndexes {
		if !f.canUpdateReq() {
			return nil
		}
		if f.state.Check(index) {
			continue
		}
		f.maxReqs.Add(1)
		f.state.Store(index, true)
		go f.GetLightClientUpdate(index)

	}
	return nil
}

func (f *Fetch) SendFetchResp(fetchType FetchType, index uint64, data interface{}) error {
	if fetchType == FinalityUpdateType {
		if f.ethFetchResp != nil {
			fetchResp := NewFetchResponse(fetchType, index, data)
			logger.Debug("fetch send fetch resp: %v", fetchResp.Id())
			f.ethFetchResp <- fetchResp
		}
	}
	return nil
}

func (f *Fetch) GetBootStrap(slot uint64) error {
	logger.Debug("start get bootstrap: %v", slot)
	exists, err := f.fileStore.CheckBootStrapBySlot(slot)
	if err != nil {
		logger.Error("check bootstrap error:%v %v", slot, err)
		return err
	}
	if exists {
		return nil
	}
	bootstrap, err := f.client.Bootstrap(slot)
	if err != nil {
		logger.Error("get bootstrap error:%v %v", slot, err)
		return err
	}
	err = f.fileStore.StoreBootStrapBySlot(slot, bootstrap)
	if err != nil {
		logger.Error("store bootstrap error:%v %v", slot, err)
		return err
	}
	err = f.fileStore.StoreBootStrap(bootstrap)
	if err != nil {
		logger.Error("store bootstrap error:%v %v", slot, err)
		return err
	}
	logger.Debug("success store bootstrap Data:%v", slot)
	return nil
}

func (f *Fetch) GetLightClientUpdate(period uint64) {
	defer func() {
		f.state.Delete(period)
		f.maxReqs.Add(-1)
	}()
	logger.Debug("start get light client update: %v", period)
	exists, err := f.fileStore.CheckUpdate(period)
	if err != nil {
		logger.Error("check update error:%v %v", period, err)
		return
	}
	if exists {
		return
	}
	updates, err := f.client.GetLightClientUpdates(period, 1)
	if err != nil {
		logger.Error("get light client updates error:%v %v", period, err)
		return
	}
	if len(updates) == 0 {
		logger.Error("get light client updates error:%v %v", period, err)
		return
	}
	ok, err := f.CheckLightClientUpdate(period, &updates[0])
	if err != nil {
		logger.Error("check light client update error:%v %v", period, err)
		return
	}
	if !ok {
		logger.Error("check light client update error:%v %v", period, err)
		return
	}
	err = f.fileStore.StoreUpdate(period, updates[0])
	if err != nil {
		logger.Error("store update error:%v %v", period, err)
		return
	}
	logger.Debug("success store update Data:%v", period)
	return
}

func (f *Fetch) GetFinalityUpdate() error {
	finalityUpdate, err := f.client.GetFinalityUpdate()
	if err != nil {
		logger.Error("get finality update error:%v", err)
		return err
	}
	slotBig, ok := big.NewInt(0).SetString(finalityUpdate.Data.FinalizedHeader.Slot, 10)
	if !ok {
		logger.Error("parse slot error:%v", finalityUpdate.Data.FinalizedHeader.Slot)
		return fmt.Errorf("parse slot error:%v", finalityUpdate.Data.FinalizedHeader.Slot)
	}
	slot := slotBig.Uint64()

	err = f.fileStore.StoreFinalizedSlot(slot)
	if err != nil {
		logger.Error("store finality update error:%v %v", slot, err)
		return err
	}
	exists, err := f.fileStore.CheckFinalityUpdate(slot)
	if err != nil {
		logger.Error("check finality update error:%v %v", slot, err)
		return err
	}
	if exists {
		return nil
	}
	//todo
	if f.ethFetchResp == nil {
		logger.Debug("success store finality update:%v", slot)
		err = f.fileStore.StoreFinalityUpdate(slot, finalityUpdate)
		if err != nil {
			logger.Error("store finality update error:%v", err)
			return err
		}
	} else {
		logger.Debug("send fetch finality update %v", slot)
		err = f.SendFetchResp(FinalityUpdateType, slot, finalityUpdate)
		if err != nil {
			logger.Error("send fetch resp error:%v", err)
			return err
		}
	}
	return nil
}

func (f *Fetch) CheckLightClientUpdate(period uint64, update *structs.LightClientUpdateWithVersion) (bool, error) {
	prePeriod := period - 1
	if prePeriod >= 0 {
		var prePeriodUpdate structs.LightClientUpdateWithVersion
		exists, err := f.fileStore.GetUpdate(prePeriod, &prePeriodUpdate)
		if err != nil {
			logger.Error("get update error:%v %v", prePeriod, err)
			return false, err
		}
		if exists {
			var syncUpdate *utils.SyncCommitteeUpdate
			err = common.ParseObj(update.Data, &syncUpdate)
			if err != nil {
				logger.Error("parse sync update error:%v %v", period, err)
				return false, err
			}
			syncUpdate.Version = prePeriodUpdate.Version
			var currentSyncCommittee utils.SyncCommittee
			err = common.ParseObj(prePeriodUpdate.Data.NextSyncCommittee, &currentSyncCommittee)
			if err != nil {
				logger.Error(err.Error())
				return false, err
			}
			syncUpdate.CurrentSyncCommittee = &currentSyncCommittee
			ok, err := common.VerifyLightClientUpdate(syncUpdate)
			if err != nil {
				logger.Error(err.Error())
				return false, err
			}
			if !ok {
				logger.Error("verify sync update error:%v %v", period, err)
				return false, nil
			}
			return true, nil
		} else {
			// todo
			logger.Warn("no find %v Index update Data", prePeriod)
			return true, nil
		}
	}
	return true, nil
}
