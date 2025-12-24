package node

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/beacon"
	beaconTypes "github.com/lightec-xyz/daemon/rpc/beacon/types"
	"github.com/lightec-xyz/daemon/store"
	proverType "github.com/lightec-xyz/provers/circuits/types"
)

var _ IFetch = (*Fetch)(nil)

type Fetch struct {
	client         beacon.IMultiBeacon
	genesisSlot    uint64
	genesisPeriod  uint64
	fileStore      *FileStorage
	lock           sync.Mutex
	chainStore     *ChainStore
	updateNotify   chan *Notify
	finalityNotify chan *Notify
}

func (f *Fetch) Init() error {
	logger.Debug("init fetch now")
	return nil
}

func (f *Fetch) Bootstrap() {
	for {
		err := f.bootstrap()
		if err != nil {
			logger.Error("get bootstrap error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		return
	}
}
func (f *Fetch) bootstrap() error {
	return nil
}

func (f *Fetch) FinalityUpdate() error {
	logger.Debug("start finality update")
	err := f.GetFinalityUpdate()
	if err != nil {
		logger.Error("get finality update error:%v", err)
		return err
	}

	return nil
}

func (f *Fetch) StoreLatestPeriod() error {
	logger.Debug("store beacon finalized period")
	period, err := f.client.FinalizedPeriod()
	if err != nil {
		logger.Error("get latest FIndex error:%v", err)
		return err
	}
	err = f.fileStore.StoreLatestPeriod(period)
	if err != nil {
		logger.Error("store latest FIndex error:%v", err)
		return err
	}
	logger.Debug("current beacon finalized period: %v", period)
	return nil
}

func (f *Fetch) LightClientUpdate() error {
	logger.Debug("start light client update")
	err := f.StoreLatestPeriod()
	if err != nil {
		logger.Error("store latest FIndex error:%v", err)
		return err
	}
	updateIndexes, err := f.fileStore.NeedUpdateIndexes()
	if err != nil {
		logger.Error("get update indexes error:%v", err)
		return nil
	}
	for _, index := range updateIndexes {
		go f.GetLightClientUpdate(index)
	}
	return nil
}

func (f *Fetch) GetLightClientUpdate(period uint64) {
	logger.Debug("start get light client update: %v", period)
	updates, err := f.client.LightClientUpdates(period, 1)
	if err != nil {
		logger.Error("get light client updates error:%v %v", period, err)
		return
	}
	if len(updates) == 0 {
		logger.Error("get light client updates error:%v %v", period, err)
		return
	}
	storeLightUpdate, err := parseRpcUpdateToDbUpdate(updates[0])
	if err != nil {
		logger.Error("parse light client update error:%v %v", period, err)
		return
	}
	ok, err := f.CheckLightClientUpdate(period, storeLightUpdate)
	if err != nil {
		logger.Error("check light client update error:%v %v", period, err)
		return
	}
	if !ok {
		logger.Error("check light client update error:%v %v", period, err)
		return
	}
	err = f.fileStore.StoreUpdate(period, storeLightUpdate)
	if err != nil {
		logger.Error("store update error:%v %v", period, err)
		return
	}
	if f.updateNotify != nil {
		f.updateNotify <- &Notify{}
	}
	logger.Debug("success store update Responses:%v", period)
	return
}

func (f *Fetch) GetFinalityUpdate() error {
	finalityUpdate, err := f.client.GetFinalityUpdate()
	if err != nil {
		logger.Error("get finality update error:%v", err)
		return err
	}
	//todo
	attestedSlotBig, ok := big.NewInt(0).SetString(finalityUpdate.Data.AttestedHeader.Beacon.Slot, 10)
	if !ok {
		logger.Error("parse attestedSlot error:%v", finalityUpdate.Data.AttestedHeader.Beacon.Slot)
		return fmt.Errorf("parse attestedSlot error:%v", finalityUpdate.Data.AttestedHeader.Beacon.Slot)
	}
	attestedSlot := attestedSlotBig.Uint64()

	finalizedSlotBig, ok := big.NewInt(0).SetString(finalityUpdate.Data.FinalizedHeader.Beacon.Slot, 10)
	if !ok {
		logger.Error("parse attestedSlot error:%v", finalityUpdate.Data.FinalizedHeader.Beacon.Slot)
		return fmt.Errorf("parse attestedSlot error:%v", finalityUpdate.Data.FinalizedHeader.Beacon.Slot)
	}
	finalizedSlot := finalizedSlotBig.Uint64()

	err = f.fileStore.StoreLatestFinalizedSlot(finalizedSlot)
	if err != nil {
		logger.Error("store finality update error:%v %v", finalizedSlot, err)
		return err
	}
	exists, err := f.fileStore.CheckFinalityUpdate(finalizedSlot)
	if err != nil {
		logger.Error("check finality update error:%v %v", finalizedSlot, err)
		return err
	}
	if exists {
		return nil
	}
	storeFinalityUpdate, err := parseRpcFinalityUpdateToDbFinalityUpdate(finalityUpdate)
	if err != nil {
		logger.Error("parse finality update error:%v %v", finalizedSlot, err)
		return err
	}
	ok, err = f.CheckFinalityUpdate(attestedSlot/common.SlotPerPeriod, storeFinalityUpdate)
	if err != nil {
		logger.Error("check finality update error:%v %v", finalizedSlot, err)
		return err
	}
	if !ok {
		logger.Error("store error %v finality update", finalizedSlot)
		err := f.fileStore.StoreErrorFinalityUpdate(fmt.Sprintf("error_%v", finalizedSlot), storeFinalityUpdate)
		if err != nil {
			logger.Error("store error finality update error:%v", err)
			return err
		}
		return nil
	}

	err = f.chainStore.WriteFinalityUpdateSlot(finalizedSlot)
	if err != nil {
		logger.Error("write finality update attestedSlot error:%v %v", finalizedSlot, err)
		return err
	}
	logger.Debug("success store finality update:%v", finalizedSlot)
	err = f.fileStore.StoreFinalityUpdate(finalizedSlot, storeFinalityUpdate)
	if err != nil {
		logger.Error("store finality update error:%v", err)
		return err
	}
	logger.Debug("send fetch finality update %v", finalizedSlot)
	if f.finalityNotify != nil {
		f.finalityNotify <- &Notify{}
	}
	return nil
}

func (f *Fetch) CheckFinalityUpdate(period uint64, finalityUpdate *common.LightClientFinalityUpdateEvent) (bool, error) {
	prePeriod := period - 1
	if prePeriod < 0 {
		return false, nil
	}
	var update common.LightClientUpdateResponse
	exists, err := f.fileStore.GetUpdate(prePeriod, &update)
	if err != nil {
		logger.Error("get update error:%v %v", prePeriod, err)
		return false, err
	}
	if !exists {
		logger.Warn("not found light finality update :%v", prePeriod)
		//todo
		return false, nil
	}
	proversTypeUpdate := parseUpdateToProversUpdate(update)
	proversTypeFinalityUpdate := parseFinalityUpdateToProversFinalityUpdate(finalityUpdate)
	ok, err := proversTypeFinalityUpdate.Verify(proversTypeUpdate.NextSyncCommittee)
	if err != nil {
		logger.Error("verify finality update signature error:%v %v", period, err)
		return false, nil
	}
	return ok, nil
}

func (f *Fetch) CheckLightClientUpdate(period uint64, update *common.LightClientUpdateResponse) (bool, error) {
	prePeriod := period - 1
	if prePeriod >= 0 {
		var prePeriodUpdate common.LightClientUpdateResponse
		exists, err := f.fileStore.GetUpdate(prePeriod, &prePeriodUpdate)
		if err != nil {
			logger.Error("get update error:%v %v", prePeriod, err)
			return false, err
		}
		if exists {
			var syncUpdate *proverType.SyncCommitteeUpdate
			err = common.ParseObj(update.Data, &syncUpdate)
			if err != nil {
				logger.Error("parse sync update error:%v %v", period, err)
				return false, err
			}

			logger.Debug("assigning syncUpdate.Version from update.Version: %v", update.Version)
			syncUpdate.Version = update.Version

			var currentSyncCommittee proverType.SyncCommittee
			err = common.ParseObj(prePeriodUpdate.Data.NextSyncCommittee, &currentSyncCommittee)
			if err != nil {
				logger.Error("parse current sync committee error:%v %v", period, err)
				return false, err
			}
			syncUpdate.CurrentSyncCommittee = &currentSyncCommittee
			verify, err := syncUpdate.Verify()
			if err != nil {
				logger.Error("verify sync update error:%v %v", period, err)
				return false, err
			}
			if !verify {
				logger.Error("verify sync update error:%v %v", period, err)
				f.client.Next()
				logger.Debug("verify update error,change beacon client")
				return false, nil
			}
			return true, nil
		} else {
			// todo
			logger.Warn("not found %v FIndex update Responses", prePeriod)
			return false, nil
		}
	}
	return true, nil
}

func (f *Fetch) Close() error {
	return nil
}

func NewFetch(client beacon.IMultiBeacon, store store.IStore, fileStore *FileStorage, genesisSlot uint64, update, finalityUpate chan *Notify) (*Fetch, error) {
	return &Fetch{
		client:         client,
		fileStore:      fileStore,
		genesisSlot:    genesisSlot,
		genesisPeriod:  genesisSlot / common.SlotPerPeriod,
		updateNotify:   update,
		finalityNotify: finalityUpate,
		chainStore:     NewChainStore(store),
	}, nil
}

func parseUpdateToProversUpdate(update common.LightClientUpdateResponse) *proverType.SyncCommitteeUpdate {
	return &proverType.SyncCommitteeUpdate{
		Version: update.Version,
		AttestedHeader: &proverType.BeaconHeader{
			Slot:          update.Data.AttestedHeader.Slot,
			ProposerIndex: update.Data.AttestedHeader.ProposerIndex,
			ParentRoot:    update.Data.AttestedHeader.ParentRoot,
			StateRoot:     update.Data.AttestedHeader.StateRoot,
			BodyRoot:      update.Data.AttestedHeader.BodyRoot,
		},
		NextSyncCommittee: &proverType.SyncCommittee{
			PubKeys:         update.Data.NextSyncCommittee.Pubkeys,
			AggregatePubKey: update.Data.NextSyncCommittee.AggregatePubkey,
		},
		FinalizedHeader: &proverType.BeaconHeader{
			Slot:          update.Data.FinalizedHeader.Slot,
			ProposerIndex: update.Data.FinalizedHeader.ProposerIndex,
			ParentRoot:    update.Data.FinalizedHeader.ParentRoot,
			StateRoot:     update.Data.FinalizedHeader.StateRoot,
			BodyRoot:      update.Data.FinalizedHeader.BodyRoot,
		},
		SyncAggregate: &proverType.SyncAggregate{
			SyncCommitteeBits:      update.Data.SyncAggregate.SyncCommitteeBits,
			SyncCommitteeSignature: update.Data.SyncAggregate.SyncCommitteeSignature,
		},
		NextSyncCommitteeBranch: update.Data.NextSyncCommitteeBranch,
		FinalityBranch:          update.Data.FinalityBranch,
		SignatureSlot:           update.Data.SignatureSlot,
	}
}

func parseFinalityUpdateToProversFinalityUpdate(update *common.LightClientFinalityUpdateEvent) *proverType.FinalityUpdate {
	return &proverType.FinalityUpdate{
		Version: update.Version,
		AttestedHeader: &proverType.BeaconHeader{
			Slot:          update.Data.AttestedHeader.Slot,
			ProposerIndex: update.Data.AttestedHeader.ProposerIndex,
			ParentRoot:    update.Data.AttestedHeader.ParentRoot,
			StateRoot:     update.Data.AttestedHeader.StateRoot,
			BodyRoot:      update.Data.AttestedHeader.BodyRoot,
		},
		FinalizedHeader: &proverType.BeaconHeader{
			Slot:          update.Data.FinalizedHeader.Slot,
			ProposerIndex: update.Data.FinalizedHeader.ProposerIndex,
			ParentRoot:    update.Data.FinalizedHeader.ParentRoot,
			StateRoot:     update.Data.FinalizedHeader.StateRoot,
			BodyRoot:      update.Data.FinalizedHeader.BodyRoot,
		},
		SyncAggregate: &proverType.SyncAggregate{
			SyncCommitteeBits:      update.Data.SyncAggregate.SyncCommitteeBits,
			SyncCommitteeSignature: update.Data.SyncAggregate.SyncCommitteeSignature,
		},
		FinalityBranch: update.Data.FinalityBranch,
		SignatureSlot:  update.Data.SignatureSlot,
	}

}

func parseBootstrapToDbBootstrap(bootstrap *beaconTypes.BootstrapResp) (*common.BootstrapResponse, error) {
	dbBootstrap := common.BootstrapResponse{
		Version: bootstrap.Version,
		Data: &common.Bootstrap{
			Header: &common.BeaconBlockHeader{
				Slot:          bootstrap.Data.Header.Slot,
				ProposerIndex: bootstrap.Data.Header.ProposerIndex,
				ParentRoot:    bootstrap.Data.Header.ParentRoot,
				StateRoot:     bootstrap.Data.Header.StateRoot,
				BodyRoot:      bootstrap.Data.Header.BodyRoot,
			},
			CurrentSyncCommittee: &common.SyncCommittee{
				Pubkeys:         bootstrap.Data.CurrentSyncCommittee.Pubkeys,
				AggregatePubkey: bootstrap.Data.CurrentSyncCommittee.AggregatePubkey,
			},
			CurrentSyncCommitteeBranch: bootstrap.Data.CurrentSyncCommitteeBranch,
		},
	}
	return &dbBootstrap, nil
}

func parseRpcFinalityUpdateToDbFinalityUpdate(update beaconTypes.LightClientFinalityUpdateResp) (*common.LightClientFinalityUpdateEvent, error) {
	dest := common.LightClientFinalityUpdateEvent{
		Version: update.Version,
		Data: &common.LightClientFinalityUpdate{
			AttestedHeader: &common.BeaconBlockHeader{
				Slot:          update.Data.AttestedHeader.Beacon.Slot,
				ProposerIndex: update.Data.AttestedHeader.Beacon.ProposerIndex,
				ParentRoot:    update.Data.AttestedHeader.Beacon.ParentRoot,
				StateRoot:     update.Data.AttestedHeader.Beacon.StateRoot,
				BodyRoot:      update.Data.AttestedHeader.Beacon.BodyRoot,
			},
			FinalizedHeader: &common.BeaconBlockHeader{
				Slot:          update.Data.FinalizedHeader.Beacon.Slot,
				ProposerIndex: update.Data.FinalizedHeader.Beacon.ProposerIndex,
				ParentRoot:    update.Data.FinalizedHeader.Beacon.ParentRoot,
				StateRoot:     update.Data.FinalizedHeader.Beacon.StateRoot,
				BodyRoot:      update.Data.FinalizedHeader.Beacon.BodyRoot,
			},
			SyncAggregate: &common.SyncAggregate{
				SyncCommitteeBits:      update.Data.SyncAggregate.SyncCommitteeBits,
				SyncCommitteeSignature: update.Data.SyncAggregate.SyncCommitteeSignature,
			},
			FinalityBranch: update.Data.FinalityBranch,
			SignatureSlot:  update.Data.SignatureSlot,
		},
	}
	return &dest, nil
}

func parseRpcUpdateToDbUpdate(update beaconTypes.LightClientUpdateResp) (*common.LightClientUpdateResponse, error) {
	dest := common.LightClientUpdateResponse{
		Version: update.Version,
		Data: &common.LightClientUpdate{
			AttestedHeader: &common.BeaconBlockHeader{
				Slot:          update.Data.AttestedHeader.Beacon.Slot,
				ProposerIndex: update.Data.AttestedHeader.Beacon.ProposerIndex,
				ParentRoot:    update.Data.AttestedHeader.Beacon.ParentRoot,
				StateRoot:     update.Data.AttestedHeader.Beacon.StateRoot,
				BodyRoot:      update.Data.AttestedHeader.Beacon.BodyRoot,
			},
			NextSyncCommittee: &common.SyncCommittee{
				Pubkeys:         update.Data.NextSyncCommittee.Pubkeys,
				AggregatePubkey: update.Data.NextSyncCommittee.AggregatePubkey,
			},
			FinalizedHeader: &common.BeaconBlockHeader{
				Slot:          update.Data.FinalizedHeader.Beacon.Slot,
				ProposerIndex: update.Data.FinalizedHeader.Beacon.ProposerIndex,
				ParentRoot:    update.Data.FinalizedHeader.Beacon.ParentRoot,
				StateRoot:     update.Data.FinalizedHeader.Beacon.StateRoot,
				BodyRoot:      update.Data.FinalizedHeader.Beacon.BodyRoot,
			},
			SyncAggregate: &common.SyncAggregate{
				SyncCommitteeBits:      update.Data.SyncAggregate.SyncCommitteeBits,
				SyncCommitteeSignature: update.Data.SyncAggregate.SyncCommitteeSignature,
			},
			NextSyncCommitteeBranch: update.Data.NextSyncCommitteeBranch,
			FinalityBranch:          update.Data.FinalityBranch,
			SignatureSlot:           update.Data.SignatureSlot,
		},
	}
	return &dest, nil
}
