package node

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"os"
	"reflect"
)

const (
	LatestPeriodKey = "latest"
	GenesisRawData  = "genesisRaw"
)

const (
	BeaconPeriod = "period"
	GenesisDir   = "genesis"
	UpdateDir    = "update"
	UnitDir      = "unit"
	RecursiveDir = "recursive"
)

type FileStore struct {
	dataDir      string
	periodDir    string
	genesisDir   string
	updateDir    string
	unitDir      string
	recursiveDir string
}

func NewFileStore(dataDir string) (*FileStore, error) {
	periodDataDir := fmt.Sprintf("%s/%s", dataDir, BeaconPeriod)
	ok, err := DirNoTExistsAndCreate(periodDataDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	updateDataDir := fmt.Sprintf("%s/%s", dataDir, BeaconPeriod)
	ok, err = DirNoTExistsAndCreate(updateDataDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error: %v %v", "update", err)
	}
	genesisDir := fmt.Sprintf("%s/%s", dataDir, GenesisDir)
	ok, err = DirNoTExistsAndCreate(updateDataDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error:%v %v", "genesis", err)
	}
	unitDir := fmt.Sprintf("%s/%s", dataDir, UnitDir)
	ok, err = DirNoTExistsAndCreate(updateDataDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error:%v %v", "unit", err)
	}
	recursiveDir := fmt.Sprintf("%s/%s", dataDir, RecursiveDir)
	ok, err = DirNoTExistsAndCreate(updateDataDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error:%v %v", "recursive", err)
	}
	return &FileStore{
		periodDir:    periodDataDir,
		dataDir:      dataDir,
		updateDir:    updateDataDir,
		genesisDir:   genesisDir,
		unitDir:      unitDir,
		recursiveDir: recursiveDir,
	}, nil
}

func (f *FileStore) StoreUnitProof(period uint64, data interface{}) error {
	return f.InsertData(UnitDir, fmt.Sprintf("%v", period), data)
}

func (f *FileStore) StoreRecursiveProof(period uint64, data interface{}) error {
	return f.InsertData(RecursiveDir, fmt.Sprintf("%v", period), data)
}

func (f *FileStore) StoreUpdate(period uint64, data interface{}) error {
	return f.InsertData(UpdateDir, fmt.Sprintf("%v", period), data)
}

func (f *FileStore) CheckUpdate(period uint64) (bool, error) {
	return f.CheckStorageKey(UpdateDir, fmt.Sprintf("%v", period))
}

func (f *FileStore) GetUpdate(period uint64, value interface{}) error {
	return f.GetData(UpdateDir, fmt.Sprintf("%v", period), value)
}

func (f *FileStore) StoreLatestPeriod(period uint64) error {
	return f.InsertData(BeaconPeriod, LatestPeriodKey, period)
}

func (f *FileStore) CheckLatestPeriod() (bool, error) {
	return f.CheckStorageKey(BeaconPeriod, LatestPeriodKey)
}

func (f *FileStore) GetLatestPeriod() (uint64, error) {
	var period uint64
	err := f.GetData(BeaconPeriod, LatestPeriodKey, &period)
	if err != nil {
		logger.Error("get latest period error:%v", err)
		return 0, err
	}
	return period, nil
}

func (f *FileStore) StoreGenesisUpdate(data interface{}) error {
	return f.InsertData(GenesisDir, GenesisRawData, data)
}

func (f *FileStore) StoreGenesisProof(data interface{}) error {
	return f.InsertData(GenesisDir, "proof", data)

}
func (f *FileStore) GetGenesisUpdate(value interface{}) error {
	return f.GetData(GenesisDir, GenesisRawData, value)
}

func (f *FileStore) CheckStorageKey(table, key string) (bool, error) {
	storeKey, err := f.generateStoreKey(table, key)
	if err != nil {
		logger.Error("generate store key error:%v", err)
		return false, err
	}
	exists, err := fileExists(storeKey)
	if err != nil {
		logger.Error("file exists error:%v", err)
		return false, err
	}
	return exists, nil
}

func (f *FileStore) InsertData(table, key string, value interface{}) error {
	storeKey, err := f.generateStoreKey(table, key)
	if err != nil {
		logger.Error("generate store key error:%v", err)
		return err
	}
	exists, err := fileExists(storeKey)
	if err != nil {
		logger.Error("file exists error:%v", err)
		return err
	}
	if exists {
		err := os.Remove(storeKey)
		if err != nil {
			logger.Error("remove file error:%v", err)
			return err
		}
	}
	dataBytes, err := json.Marshal(value)
	if err != nil {
		logger.Error("marshal file error:%v", err)
		return err
	}
	err = os.WriteFile(storeKey, dataBytes, 0644)
	if err != nil {
		logger.Error("write file error:%v", err)
		return err
	}
	return nil
}

func (f *FileStore) GetData(table, key string, value interface{}) error {
	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		return fmt.Errorf("value mutst be a pointer")
	}
	storeKey, err := f.generateStoreKey(table, key)
	if err != nil {
		logger.Error("generate store key error:%v", err)
		return err
	}
	exists, err := fileExists(storeKey)
	if err != nil {
		logger.Error("file exists error:%v", err)
		return err
	}
	if !exists {
		return fmt.Errorf("no find key: %v-%v", table, key)
	}
	dataBytes, err := os.ReadFile(storeKey)
	if err != nil {
		logger.Error("read file error:%v", err)
		return err
	}
	err = json.Unmarshal(dataBytes, value)
	if err != nil {
		logger.Error("unmarshal file error:%v-%v  %v", table, key, err)
		return err
	}
	return nil
}

func (f *FileStore) generateStoreKey(table, key string) (string, error) {
	switch table {
	case BeaconPeriod:
		return fmt.Sprintf("%s/%s", f.periodDir, key), nil
	case GenesisDir:
		return fmt.Sprintf("%s/%s", f.genesisDir, key), nil
	case UpdateDir:
		return fmt.Sprintf("%s/%s", f.updateDir, key), nil
	case UnitDir:
		return fmt.Sprintf("%s/%s", f.unitDir, key), nil
	case RecursiveDir:
		return fmt.Sprintf("%s/%s", f.recursiveDir, key), nil
	default:
		return "", fmt.Errorf("no find table: %v", table)
	}
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat error: %v", err)
}

func DirNoTExistsAndCreate(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
