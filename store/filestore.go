package store

import (
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/lightec-xyz/daemon/logger"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

type IFileStore interface {
	Get(key FileKey, value interface{}) (bool, error)
	Store(key FileKey, data interface{}) error
	CheckExists(key FileKey) (bool, error)
	GetFilePath(key FileKey) string
	RootPath() string
	ClearAll() error
	Del(StoreKey FileKey) error
	ICache
	ISubFileStore
}

type ICache interface {
	Keys() []interface{}
	GetValue(key interface{}) (interface{}, bool)
	MaxIndex() (interface{}, interface{})
}

type ISubFileStore interface {
	SubStore(table Table, key FileKey, value interface{}) error
	SubGet(table Table, key FileKey, value interface{}) (bool, error)
	SubCheck(table Table, key FileKey) (bool, error)
	SubDel(table Table, key FileKey) error
	SubFileStore(sub Table, create ...bool) (IFileStore, bool)
	SubFileStores() ([]IFileStore, error)
}

type FileStore struct {
	SubFileMaps *sync.Map // Table -> IFileStore
	Root        string
	Prefix      string
	Suffix      string
	tree        *treemap.Map                      // cache element  endIndex <-> FileKey
	fn          func(name FileKey) (uint64, bool) // parse key to index
}

func NewFileStore(rootPath string, fn func(name FileKey) (uint64, bool), subTables ...Table) (IFileStore, error) {
	err := CheckOrCreateDir(rootPath)
	if err != nil {
		logger.Error("check or create dir error:%v %v", rootPath, err)
		return nil, err
	}

	tree := treemap.NewWith(utils.UInt64Comparator)
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}
	subFilesMap := new(sync.Map)
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			path := fmt.Sprintf("%v/%v", rootPath, name)
			subStore, err := NewFileStore(path, fn)
			if err != nil {
				return nil, err
			}
			subFilesMap.Store(Table(name), subStore)
		} else {
			if fn != nil {
				index, ok := fn(FileKey(name))
				if ok {
					tree.Put(index, FileKey(name))
				}
			}
		}
	}
	for _, subTable := range subTables {
		if _, ok := subFilesMap.Load(subTable); !ok {
			subFilesMap.Store(subTable, &FileStore{Root: fmt.Sprintf("%s/%s", rootPath, subTable)})
		}
	}
	return &FileStore{
		SubFileMaps: subFilesMap,
		Root:        rootPath,
		fn:          fn,
		tree:        tree,
	}, nil
}

func (fs *FileStore) SubStore(subTable Table, key FileKey, value interface{}) error {
	fileStore, ok := fs.SubFileStore(subTable, true)
	if !ok {
		logger.Error("get file store error %v", subTable)
		return fmt.Errorf("get file store error %v", subTable)
	}
	return fileStore.Store(key, value)
}

func (fs *FileStore) SubCheck(subTable Table, key FileKey) (bool, error) {
	fileStore, ok := fs.SubFileStore(subTable)
	if !ok {
		return false, nil
	}
	return fileStore.CheckExists(key)
}

func (fs *FileStore) SubGet(subTable Table, key FileKey, value interface{}) (bool, error) {
	fileStore, ok := fs.SubFileStore(subTable)
	if !ok {
		return false, nil
	}
	return fileStore.Get(key, value)
}
func (fs *FileStore) SubDel(subTable Table, key FileKey) error {
	fileStore, ok := fs.SubFileStore(subTable)
	if !ok {
		return nil
	}
	return fileStore.Del(key)
}
func (fs *FileStore) SubFileStores() ([]IFileStore, error) {
	var fileStores []IFileStore
	fs.SubFileMaps.Range(func(key, value interface{}) bool {
		fileStores = append(fileStores, value.(IFileStore))
		return true
	})
	return fileStores, nil
}

func (fs *FileStore) SubFileStore(sub Table, create ...bool) (IFileStore, bool) {
	filestore, ok := fs.getSubFileStore(sub)
	if !ok {
		if len(create) > 0 && create[0] {
			path := fmt.Sprintf("%s/%s", fs.Root, sub)
			err := CheckOrCreateDir(path)
			if err != nil {
				return nil, false
			}
			subStore, err := NewFileStore(path, fs.fn)
			if err != nil {
				return nil, false
			}
			fs.SubFileMaps.Store(sub, subStore)
			return subStore, true
		} else {
			return nil, false
		}
	}
	return filestore, true
}

func (fs *FileStore) getSubFileStore(sub Table) (IFileStore, bool) {
	filestore, ok := fs.SubFileMaps.Load(sub)
	if !ok {
		return nil, false
	}
	fileStore, ok := filestore.(IFileStore)
	if !ok {
		return nil, false
	}
	return fileStore, true
}
func (fs *FileStore) Get(name FileKey, value interface{}) (bool, error) {
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
	data, err := os.ReadFile(fs.GetFilePath(name))
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

func (fs *FileStore) Store(name FileKey, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		logger.Error("json marshal error:%v", err)
		return err
	}
	if fs.fn != nil {
		key, ok := fs.fn(name)
		if ok {
			fs.tree.Put(key, name)
		}
	}
	err = OverWriteFile(fs.GetFilePath(name), dataBytes)
	if err != nil {
		logger.Error("write file error:%v", err)
		return err
	}
	return nil
}

func (fs *FileStore) CheckExists(name FileKey) (bool, error) {
	exists, err := FileExists(fs.GetFilePath(name))
	if err != nil {
		logger.Error("file exists error:%v", err)
		return false, err
	}
	return exists, nil
}

func (fs *FileStore) Del(name FileKey) error {
	if fs.fn != nil {
		key, ok := fs.fn(name)
		if ok {
			fs.tree.Remove(key)
		}
	}
	err := os.Remove(fs.GetFilePath(name))
	if err != nil {
		logger.Error("remove file error:%v", err)
		return err
	}
	return nil
}

func (fs *FileStore) GetFilePath(name FileKey) string {
	path := fs.Root
	if !strings.HasSuffix(fs.Root, "/") {
		path = fmt.Sprintf("%s/", fs.Root)
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
	return fs.Root
}

func (fs *FileStore) Files(all ...bool) (map[FileKey]string, error) {
	files, err := traverseFile(fs.Root, all...)
	if err != nil {
		logger.Error("traverse file error:%v %v", fs.Root, err)
		return nil, err
	}
	return files, nil
}

func (fs *FileStore) ClearAll() error {
	err := os.RemoveAll(fs.Root)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStore) MaxIndex() (interface{}, interface{}) {
	return fs.tree.Max()
}

func (fs *FileStore) Keys() []interface{} {
	return fs.tree.Keys()
}

func (fs *FileStore) GetValue(key interface{}) (interface{}, bool) {
	return fs.tree.Get(key)

}

type Table string

type FileKey string

func (k FileKey) String() string {
	return string(k)
}

type IndexKey struct {
	End uint64
}

func traverseFile(root string, all ...bool) (map[FileKey]string, error) {
	files := make(map[FileKey]string)
	if len(all) > 0 && all[0] {
		err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && root != path {
				fileName, err := getFileName(info.Name())
				if err != nil {
					return err
				}
				files[FileKey(fileName)] = path
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		dirEntries, err := os.ReadDir(root)
		if err != nil {
			return nil, err
		}
		for _, dirEntry := range dirEntries {
			if !dirEntry.IsDir() && root != dirEntry.Name() {
				fileName, err := getFileName(dirEntry.Name())
				if err != nil {
					return nil, err
				}
				files[FileKey(fileName)] = filepath.Join(root, dirEntry.Name())
			}
		}
	}
	return files, nil
}
func getFileName(path string) (string, error) {
	arrs := strings.Split(path, "/")
	if len(arrs) == 0 {
		return "", fmt.Errorf("get file name error")
	}
	return arrs[len(arrs)-1], nil
}

func GenFileKey(prefix Table, args ...interface{}) FileKey {
	var key string
	key = fmt.Sprintf("%v", prefix)
	for _, arg := range args {
		key = key + fmt.Sprintf("_%v", arg)
	}
	return FileKey(strings.ToLower(strings.TrimPrefix(key, "_")))
}

func CheckOrCreateDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
func OverWriteFile(path string, data []byte) error {
	err := CheckOrCreateDir(path)
	if err != nil {
		return err
	}
	exists, err := FileExists(path)
	if err != nil {
		return err
	}
	if exists {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat error: %v", err)
}

type Indexes []IndexKey
