package store

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type FileStore struct {
	Path   string
	Prefix string
	Suffix string
}

func NewFileStore(path string, opts ...string) (*FileStore, error) {
	if !strings.HasSuffix(path, "/") {
		path = fmt.Sprintf("%s/", path)
	}
	err := common.CheckOrCreateDir(path)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	// todo default?
	var prefix, suffix string
	if len(opts) > 0 {
		prefix = opts[0]
	}
	if len(opts) > 1 {
		suffix = opts[1]
	}

	return &FileStore{
		Path:   path,
		Prefix: prefix,
		Suffix: suffix,
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
	path := fs.Path
	if !strings.HasSuffix(fs.Path, "/") {
		path = fmt.Sprintf("%s/", fs.Path)
	}
	if fs.Prefix != "" {
		path = fmt.Sprintf("%s%s_", path, fs.Prefix)
	}
	if fs.Suffix != "" {
		return fmt.Sprintf("%s%s.%s", path, name, fs.Suffix)
	}
	return fmt.Sprintf("%s%s", path, name)
}

func (fs *FileStore) RootPath() string {
	return fs.Path
}

func (fs *FileStore) AllIndexes() (map[uint64]string, error) {
	files, err := fs.AllFiles()
	if err != nil {
		logger.Error("all files error:%v", err)
		return nil, err
	}
	indexes := make(map[uint64]string)
	for fileName, _ := range files {
		match := numberPatten.FindStringSubmatch(fileName)
		if len(match) > 1 {
			i, err := strconv.ParseUint(match[1], 10, 64)
			if err != nil {
				return nil, err
			}
			indexes[i] = fileName
		}
	}
	return indexes, nil
}

func (fs *FileStore) AllFiles() (map[string]string, error) {
	files, err := traverseFile(fs.Path)
	if err != nil {
		logger.Error("traverse file error:%v %v", fs.Path, err)
		return nil, err
	}
	return files, nil
}

var numberPatten = regexp.MustCompile(`(\d+)`)

func traverseFile(path string) (map[string]string, error) {
	files := make(map[string]string)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileName, err := getFileName(info.Name())
			if err != nil {
				return err
			}
			files[fileName] = path
		}
		return nil
	})
	return files, err
}
func getFileName(path string) (string, error) {
	arrs := strings.Split(path, "/")
	if len(arrs) == 0 {
		return "", fmt.Errorf("get file name error")
	}
	return arrs[len(arrs)-1], nil
}
