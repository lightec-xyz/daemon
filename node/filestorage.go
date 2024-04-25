package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/store"
)

type FileStorage struct {
	RootPath       string
	Genesis        *store.FileStore
	Period         *store.FileStore
	Update         *store.FileStore
	Unit           *store.FileStore
	Outer          *store.FileStore
	Recursive      *store.FileStore
	FinalityUpdate *store.FileStore
	BhfUpdate      *store.FileStore
	BeaconHeader   *store.FileStore
	Txes           *store.FileStore
	Redeem         *store.FileStore
}

func NewFileStorage(rootPath string) (*FileStorage, error) {
	periodStore, err := CreateFileStore(rootPath, PeriodDir)
	if err != nil {
		logger.Error("create periodStore store error")
		return nil, err
	}
	genesisStore, err := CreateFileStore(rootPath, GenesisDir)
	if err != nil {
		logger.Error("create genesisStore store error")
		return nil, err
	}
	updatesStore, err := CreateFileStore(rootPath, UpdateDir)
	if err != nil {
		logger.Error("create updatesStore store error")
		return nil, err
	}
	unitStore, err := CreateFileStore(rootPath, UnitDir)
	if err != nil {
		logger.Error("create unitStore store error")
		return nil, err
	}
	outerStore, err := CreateFileStore(rootPath, Outer)
	if err != nil {
		logger.Error("create outerStore store error")
		return nil, err
	}
	recursiveStore, err := CreateFileStore(rootPath, RecursiveDir)
	if err != nil {
		logger.Error("create recursiveStore store error")
		return nil, err
	}
	finalityUpdateStore, err := CreateFileStore(rootPath, FinalityUpdate)
	if err != nil {
		logger.Error("create finalityUpdateStore store error")
		return nil, err
	}
	bhfUpdateStore, err := CreateFileStore(rootPath, BhfUpdate)
	if err != nil {
		logger.Error("create bhfUpdateStore store error")
		return nil, err
	}
	beaconHeaderStore, err := CreateFileStore(rootPath, BlockHeader)
	if err != nil {
		logger.Error("create beaconHeaderStore store error")
		return nil, err
	}
	redeemStore, err := CreateFileStore(rootPath, Redeem)
	if err != nil {
		logger.Error("create redeemStore store error")
		return nil, err
	}
	return &FileStorage{
		RootPath:       rootPath,
		Genesis:        genesisStore,
		Period:         periodStore,
		Update:         updatesStore,
		Unit:           unitStore,
		Outer:          outerStore,
		Recursive:      recursiveStore,
		FinalityUpdate: finalityUpdateStore,
		BhfUpdate:      bhfUpdateStore,
		BeaconHeader:   beaconHeaderStore,
		Txes:           nil,
		Redeem:         redeemStore,
	}, nil
}

func CreateFileStore(root, name string) (*store.FileStore, error) {
	path := fmt.Sprintf("%s/%s", root, name)
	fileStore, err := store.NewFileStore(path)
	if err != nil {
		logger.Error("create file store error", "path", path, "err", err)
		return nil, err
	}
	return fileStore, nil
}

func (fs *FileStorage) GetRootPath() string {
	return fs.RootPath
}
