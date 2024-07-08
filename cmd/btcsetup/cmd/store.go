package cmd

import (
	"encoding/hex"
	"fmt"
	native_plonk "github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/store"
	"sync"
)

type StoreProof struct {
	Type    string `json:"type"`
	Id      string `json:"id"`
	Proof   string `json:"proof"`
	Witness string `json:"witness"`
}
type Table string

const (
	baseTable   Table = "base"
	middleTable Table = "middle"
	upTable     Table = "up"
)

var InitStoreTables = []Table{baseTable, middleTable, upTable}

type FileStorage struct {
	RootPath     string
	FileStoreMap map[Table]*store.FileStore
	lock         sync.Mutex
}

func NewFileStorage(rootPath string, tables ...Table) (*FileStorage, error) {
	fileStoreMap := make(map[Table]*store.FileStore)
	path := fmt.Sprintf("%s/proofData", rootPath) // todo
	logger.Info("fileStorage path: %v", path)
	if tables == nil {
		tables = InitStoreTables
	}
	for _, key := range tables {
		fileStore, err := CreateFileStore(path, string(key))
		if err != nil {
			logger.Error("create file store error")
			return nil, err
		}
		fileStoreMap[key] = fileStore
	}
	return &FileStorage{
		RootPath:     path,
		FileStoreMap: fileStoreMap,
	}, nil
}

func (fs *FileStorage) StoreBase(key string, proof native_plonk.Proof, wit witness.Witness) error {
	storeProof, err := newStoreProof(string(baseTable), key, proof, wit)
	if err != nil {
		return err
	}
	return fs.Store(baseTable, key, storeProof)
}

func (fs *FileStorage) StoreMiddle(key string, proof native_plonk.Proof, wit witness.Witness) error {
	storeProof, err := newStoreProof(string(middleTable), key, proof, wit)
	if err != nil {
		return err
	}
	return fs.Store(middleTable, key, storeProof)
}

func (fs *FileStorage) StoreUp(key string, proof native_plonk.Proof, wit witness.Witness) error {
	storeProof, err := newStoreProof(string(upTable), key, proof, wit)
	if err != nil {
		return err
	}
	return fs.Store(upTable, key, storeProof)
}

func (fs *FileStorage) GetFileStore(table Table) (*store.FileStore, bool) {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	filestore, ok := fs.FileStoreMap[table]
	if !ok {
		return nil, false
	}
	return filestore, true
}

func (fs *FileStorage) Store(table Table, key string, value interface{}) error {
	fileStore, ok := fs.GetFileStore(table)
	if !ok {
		logger.Error("get file store error %v", table)
		return fmt.Errorf("get file store error %v", table)
	}
	return fileStore.Store(key, value)
}

func (fs *FileStorage) Get(table Table, key string, value interface{}) (bool, error) {
	fileStore, ok := fs.GetFileStore(table)
	if !ok {
		logger.Error("get file store error %v", table)
		return false, fmt.Errorf("get file store error %v", table)
	}
	return fileStore.Get(key, value)
}

func (fs *FileStorage) Check(table Table, key string) (bool, error) {
	fileStore, ok := fs.GetFileStore(table)
	if !ok {
		logger.Error("get file store error %v", table)
		return false, fmt.Errorf("get file store error %v", table)
	}
	return fileStore.CheckExists(key)
}
func CreateFileStore(root, name string) (*store.FileStore, error) {
	path := fmt.Sprintf("%s/%s/", root, name)
	fileStore, err := store.NewFileStore(path)
	if err != nil {
		logger.Error("create file store error %v %v %v", root, name, err)
		return nil, err
	}
	return fileStore, nil
}

func (fs *FileStorage) GetRootPath() string {
	return fs.RootPath
}
func parseKey(keys ...interface{}) string {
	var key string
	for index, tmp := range keys {
		if index == 0 {
			key = fmt.Sprintf("%v", tmp)
		} else {
			key = key + fmt.Sprintf("_%v", tmp)
		}
	}
	return key
}

func genKey(prefix string, args ...interface{}) string {
	name := fmt.Sprintf("%v", prefix)
	for _, arg := range args {
		name = name + fmt.Sprintf("_%v", arg)
	}
	return name
}
func newStoreProof(proofType, id string, proof native_plonk.Proof, wit witness.Witness) (*StoreProof, error) {
	proofBytes, err := circuits.ProofToBytes(proof)
	if err != nil {
		logger.Error("proof to bytes error: %v", err)
		return nil, err
	}
	witnessBytes, err := circuits.WitnessToBytes(wit)
	if err != nil {
		return nil, err
	}
	return &StoreProof{
		Type:    proofType,
		Id:      id,
		Proof:   hex.EncodeToString(proofBytes),   // todo
		Witness: hex.EncodeToString(witnessBytes), // todo
	}, nil
}
