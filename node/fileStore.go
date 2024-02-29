package node

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"os"
	"reflect"
)

const (
	UpdateDir    = "update"
	GenesisDir   = "genesis"
	UnitDir      = "unit"
	RecursiveDir = "recursive"
)

type FileStore struct {
	dataDir      string
	updateDir    string
	genesisDir   string
	unitDir      string
	recursiveDir string
}

func NewFileStore(dataDir string) (*FileStore, error) {
	updateDataDir := fmt.Sprintf("%s/%s", dataDir, UpdateDir)
	ok, err := DirNoTExistsAndCreate(updateDataDir)
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
		dataDir:      dataDir,
		updateDir:    updateDataDir,
		genesisDir:   genesisDir,
		unitDir:      unitDir,
		recursiveDir: recursiveDir,
	}, nil
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
	case UpdateDir:
		return fmt.Sprintf("%s/%s", f.updateDir, key), nil
	case GenesisDir:
		return fmt.Sprintf("%s/%s", f.genesisDir, key), nil
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
