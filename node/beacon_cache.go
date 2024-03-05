package node

import "sync"

const (
	GenesisStateKey = "genesisStateKey"
)

type BeaconCache struct {
	downloadProof          *sync.Map
	generatingGenesisProof *sync.Map
	unitProof              *sync.Map
	RecursiveProof         *sync.Map
}

func NewBeaconCache() *BeaconCache {
	return &BeaconCache{
		downloadProof:          new(sync.Map),
		generatingGenesisProof: new(sync.Map),
		unitProof:              new(sync.Map),
		RecursiveProof:         new(sync.Map),
	}
}

func (bc *BeaconCache) CheckDownload(period uint64) bool {
	_, ok := bc.downloadProof.Load(period)
	return ok
}

func (bc *BeaconCache) CheckGenesis() bool {
	_, ok := bc.generatingGenesisProof.Load(GenesisStateKey)
	return ok
}

func (bc *BeaconCache) CheckUnit(period uint64) bool {
	_, ok := bc.unitProof.Load(period)
	return ok
}

func (bc *BeaconCache) CheckRecursive(period uint64) bool {
	_, ok := bc.RecursiveProof.Load(period)
	return ok
}

func (bc *BeaconCache) StoreGenesis() error {
	bc.generatingGenesisProof.Store(GenesisStateKey, true)
	return nil
}

func (bc *BeaconCache) StoreUnit(period uint64) error {
	bc.unitProof.Store(period, true)
	return nil
}

func (bc *BeaconCache) StoreRecursive(period uint64) error {
	bc.RecursiveProof.Store(period, true)
	return nil
}
