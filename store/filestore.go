package store

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"os"
	"reflect"
)

type FileStore struct {
	Path string
}

func NewFileStore(path string) (*FileStore, error) {
	err := common.CheckOrCreateDir(path)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	return &FileStore{
		Path: path,
	}, nil
}

func (fs *FileStore) Clear() error {
	return os.RemoveAll(fs.Path)
}

func (fs *FileStore) Get(name string, value interface{}) (bool, error) {
	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		return false, fmt.Errorf("value mutst be a pointer")
	}
	exists, err := fs.CheckExists(name)
	if err != nil {
		logger.Error("check file exists error:%v", err)
		return false, err
	}
	if !exists {
		return false, nil
	}
	data, err := fs.GetData(name)
	if err != nil {
		logger.Error("read file error:%v", err)
		return false, err
	}
	err = json.Unmarshal(data, value)
	if err != nil {
		logger.Error("json unmarshal error:%v", err)
		return false, err
	}
	return true, nil
}

func (fs *FileStore) GetObj(name string, value interface{}) error {
	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		return fmt.Errorf("value mutst be a pointer")
	}
	exists, err := fs.CheckExists(name)
	if err != nil {
		logger.Error("check file exists error:%v", err)
		return err
	}
	if !exists {
		return fmt.Errorf("%v file not exists", name)
	}
	data, err := fs.GetData(name)
	if err != nil {
		logger.Error("read file error:%v", err)
		return err
	}
	err = json.Unmarshal(data, value)
	if err != nil {
		logger.Error("json unmarshal error:%v", err)
		return err
	}
	return nil
}

func (fs *FileStore) GetData(name string) ([]byte, error) {
	exists, err := fs.CheckExists(name)
	if err != nil {
		logger.Error("check file exists error:%v", err)
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("%v file not exists", name)
	}
	content, err := os.ReadFile(fs.GetFilePath(name))
	if err != nil {
		logger.Error("read file error:%v", err)
		return nil, err
	}
	return content, nil
}

func (fs *FileStore) Store(name string, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		logger.Error("json marshal error:%v", err)
		return err
	}

	err = common.WriteFile(fs.GetFilePath(name), dataBytes)
	if err != nil {
		logger.Error("write file error:%v", err)
		return err
	}
	return nil
}

func (fs *FileStore) CheckExists(name string) (bool, error) {
	exists, err := common.FileExists(fs.GetFilePath(name))
	if err != nil {
		logger.Error("file exists error:%v", err)
		return false, err
	}
	return exists, nil
}

func (fs *FileStore) GetFilePath(name string) string {
	return fmt.Sprintf("%s/%s", fs.Path, name)
}

func (fs *FileStore) RootPath() string {
	return fs.Path
}
