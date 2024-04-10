package node

import (
	"sync"
)

const (
	GenesisStateKey = "genesisStateKey"
)

type BeaconCache struct {
	fetchData                  *sync.Map
	genesisProof               *sync.Map
	unitProof                  *sync.Map
	recursiveProof             *sync.Map
	beaconHeaderFinalityUpdate *sync.Map
}

func NewBeaconCache() *BeaconCache {
	return &BeaconCache{
		fetchData:                  new(sync.Map),
		genesisProof:               new(sync.Map),
		unitProof:                  new(sync.Map),
		recursiveProof:             new(sync.Map),
		beaconHeaderFinalityUpdate: new(sync.Map),
	}
}

func (bc *BeaconCache) StoreBeaconHeaderFinalityUpdate(slot uint64) error {
	bc.beaconHeaderFinalityUpdate.Store(slot, true)
	return nil
}

func (bc *BeaconCache) CheckBeaconHeaderFinalityUpdate(slot uint64) bool {
	_, ok := bc.beaconHeaderFinalityUpdate.Load(slot)
	return ok
}

func (bc *BeaconCache) CheckFetchData(period uint64) bool {
	_, ok := bc.fetchData.Load(period)
	return ok
}

func (bc *BeaconCache) CheckGenesis() bool {
	_, ok := bc.genesisProof.Load(GenesisStateKey)
	return ok
}

func (bc *BeaconCache) CheckUnit(period uint64) bool {
	_, ok := bc.unitProof.Load(period)
	return ok
}

func (bc *BeaconCache) CheckRecursive(period uint64) bool {
	_, ok := bc.recursiveProof.Load(period)
	return ok
}

func (bc *BeaconCache) StoreGenesis() error {
	bc.genesisProof.Store(GenesisStateKey, true)
	return nil
}

func (bc *BeaconCache) StoreUnit(period uint64) error {
	bc.unitProof.Store(period, true)
	return nil
}

func (bc *BeaconCache) StoreRecursive(period uint64) error {
	bc.recursiveProof.Store(period, true)
	return nil
}

func (bc *BeaconCache) DeleteFetchData(period uint64) {
	bc.fetchData.Delete(period)
}
func (bc *BeaconCache) DeleteGenesis() {
	bc.genesisProof.Delete(GenesisStateKey)
}
func (bc *BeaconCache) DeleteUnit(period uint64) {
	bc.unitProof.Delete(period)
}

func (bc *BeaconCache) DeleteRecursive(period uint64) {
	bc.recursiveProof.Delete(period)
}
