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
	GenesisProofKey = "genesisProof"
)

const (
	PeriodDir    = "period"
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
	periodDataDir := fmt.Sprintf("%s/%s", dataDir, PeriodDir)
	ok, err := DirNoTExistsAndCreate(periodDataDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir %v error", periodDataDir)
	}

	updateDataDir := fmt.Sprintf("%s/%s", dataDir, UpdateDir)
	ok, err = DirNoTExistsAndCreate(updateDataDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error: %v %v", "update", err)
	}
	genesisDir := fmt.Sprintf("%s/%s", dataDir, GenesisDir)
	ok, err = DirNoTExistsAndCreate(genesisDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error:%v %v", "genesis", err)
	}
	unitDir := fmt.Sprintf("%s/%s", dataDir, UnitDir)
	ok, err = DirNoTExistsAndCreate(unitDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error:%v %v", "unit", err)
	}
	recursiveDir := fmt.Sprintf("%s/%s", dataDir, RecursiveDir)
	ok, err = DirNoTExistsAndCreate(recursiveDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error:%v %v", "recursive", err)
	}
	return &FileStore{
		dataDir:      dataDir,
		periodDir:    periodDataDir,
		updateDir:    updateDataDir,
		genesisDir:   genesisDir,
		unitDir:      unitDir,
		recursiveDir: recursiveDir,
	}, nil
}

func (f *FileStore) StoreRecursiveProof(period uint64, data interface{}) error {
	return f.InsertData(RecursiveDir, parseKey(period), data)
}

func (f *FileStore) GetRecursiveProof(period uint64, value interface{}) error {
	return f.GetData(RecursiveDir, parseKey(period), value)
}

func (f *FileStore) CheckRecursiveProof(period uint64) (bool, error) {
	return f.CheckStorageKey(RecursiveDir, parseKey(period))
}

func (f *FileStore) CheckUnitProof(period uint64) (bool, error) {
	return f.CheckStorageKey(UnitDir, parseKey(period))
}

func (f *FileStore) StoreUnitProof(period uint64, data interface{}) error {
	return f.InsertData(UnitDir, parseKey(period), data)
}
func (f *FileStore) GetUnitProof(period uint64, value interface{}) error {
	return f.GetData(UnitDir, parseKey(period), value)
}

func (f *FileStore) StoreUpdate(period uint64, data interface{}) error {
	return f.InsertData(UpdateDir, parseKey(period), data)
}

func (f *FileStore) CheckUpdate(period uint64) (bool, error) {
	return f.CheckStorageKey(UpdateDir, parseKey(period))
}

func (f *FileStore) GetUpdate(period uint64, value interface{}) error {
	return f.GetData(UpdateDir, parseKey(period), value)
}

func (f *FileStore) GetGenesisUpdate(value interface{}) error {
	return f.GetData(GenesisDir, GenesisRawData, value)
}

func (f *FileStore) StoreGenesisUpdate(data interface{}) error {
	return f.InsertData(GenesisDir, GenesisRawData, data)
}

func (f *FileStore) CheckGenesisUpdate() (bool, error) {
	return f.CheckStorageKey(GenesisDir, GenesisRawData)
}

func (f *FileStore) StoreGenesisProof(data interface{}) error {
	return f.InsertData(GenesisDir, GenesisProofKey, data)
}

func (f *FileStore) CheckGenesisProof() (bool, error) {
	return f.CheckStorageKey(GenesisDir, GenesisProofKey)
}

func (f *FileStore) GetGenesisProof(value interface{}) error {
	return f.GetData(GenesisDir, GenesisProofKey, value)
}

func (f *FileStore) StoreLatestPeriod(period uint64) error {
	return f.InsertData(PeriodDir, LatestPeriodKey, period)
}

func (f *FileStore) CheckLatestPeriod() (bool, error) {
	return f.CheckStorageKey(PeriodDir, LatestPeriodKey)
}

func (f *FileStore) GetLatestPeriod() (uint64, error) {
	var period uint64
	err := f.GetData(PeriodDir, LatestPeriodKey, &period)
	if err != nil {
		logger.Error("get latest period error:%v", err)
		return 0, err
	}
	return period, nil
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
	case PeriodDir:
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
func parseKey(key interface{}) string {
	return fmt.Sprintf("%v", key)
}
