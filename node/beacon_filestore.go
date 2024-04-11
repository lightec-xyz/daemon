package node

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

// Todo refactor
// store filestore protocol

type StoreProof struct {
	ProofType common.ZkProofType `json:"type"`
	Hash      string             `json:"hash"`
	Period    uint64             `json:"period"`
	Proof     []byte             `json:"proof"`
	Witness   []byte             `json:"witness"`
}

type storeProof struct {
	ProofType common.ZkProofType `json:"type"`
	Hash      string             `json:"hash"`
	Period    uint64             `json:"period"`
	Proof     string             `json:"proof"`
	Witness   string             `json:"witness"`
}

const (
	LatestPeriodKey = "latest"
	LatestSlotKey   = "latestSlot"
	GenesisRawData  = "genesisRaw"
	GenesisProofKey = "genesisProof"
)

const (
	PeriodDir    = "period"
	GenesisDir   = "genesis"
	UpdateDir    = "update"
	UnitDir      = "unit"
	RecursiveDir = "recursive"
	BhfUpdate    = "bhfUpdate"
	Tx           = "txes"
)

// todo refactor

type FileStore struct {
	dataDir       string
	periodDir     string
	genesisDir    string
	updateDir     string
	unitDir       string
	recursiveDir  string
	txDir         string
	bhfUpdateDir  string
	genesisPeriod uint64
	genesisSlot   uint64
}

func NewFileStore(dataDir string, genesisSlot, genesisPeriod uint64) (*FileStore, error) {

	dataDir = fmt.Sprintf("%s/%s", dataDir, "proofData")
	// todo
	periodDataDir := fmt.Sprintf("%s/%s", dataDir, PeriodDir)
	ok, err := dirNotExistsAndCreate(periodDataDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir %v error", periodDataDir)
	}

	updateDataDir := fmt.Sprintf("%s/%s", dataDir, UpdateDir)
	ok, err = dirNotExistsAndCreate(updateDataDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error: %v %v", "update", err)
	}
	genesisDir := fmt.Sprintf("%s/%s", dataDir, GenesisDir)
	ok, err = dirNotExistsAndCreate(genesisDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error:%v %v", "genesis", err)
	}
	unitDir := fmt.Sprintf("%s/%s", dataDir, UnitDir)
	ok, err = dirNotExistsAndCreate(unitDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error:%v %v", "unit", err)
	}
	recursiveDir := fmt.Sprintf("%s/%s", dataDir, RecursiveDir)
	ok, err = dirNotExistsAndCreate(recursiveDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error:%v %v", "recursive", err)
	}

	bhfuDir := fmt.Sprintf("%s/%s", dataDir, BhfUpdate)
	ok, err = dirNotExistsAndCreate(bhfuDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error:%v %v", "tx", err)
	}

	txDir := fmt.Sprintf("%s/%s", dataDir, Tx)
	ok, err = dirNotExistsAndCreate(txDir)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create dir error:%v %v", "tx", err)
	}
	return &FileStore{
		dataDir:       dataDir,
		periodDir:     periodDataDir,
		updateDir:     updateDataDir,
		genesisDir:    genesisDir,
		unitDir:       unitDir,
		bhfUpdateDir:  bhfuDir,
		txDir:         txDir,
		recursiveDir:  recursiveDir,
		genesisPeriod: genesisPeriod,
		genesisSlot:   genesisSlot,
	}, nil
}

func (f *FileStore) StoreBhfUpdateProof(slot uint64, proof, witness []byte) error {
	return f.InsertData(BhfUpdate, parseKey(slot), storeProof{
		Period:    slot,
		ProofType: common.BhfUpdate,
		Proof:     hex.EncodeToString(proof),
		Witness:   hex.EncodeToString(witness),
	})
}

func (f *FileStore) CheckBhfUpdateProof(slot uint64) (bool, error) {
	return f.CheckStorageKey(BhfUpdate, parseKey(slot))
}

func (f *FileStore) GetBhfUpdateProof(slot uint64) (*StoreProof, bool, error) {
	exists, err := f.CheckBhfUpdateProof(slot)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	proof, err := f.GetObj(BhfUpdate, parseKey(slot))
	if err != nil {
		return nil, false, err
	}
	return proof, true, nil
}

func (f *FileStore) CheckLatestSlot() (bool, error) {
	return f.CheckStorageKey(PeriodDir, LatestSlotKey)
}

func (f *FileStore) StoreLatestSlot(slot uint64) error {
	return f.InsertData(PeriodDir, LatestSlotKey, slot)
}
func (f *FileStore) GetLatestSlot() (uint64, bool, error) {
	exists, err := f.CheckLatestSlot()
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}
	var latestSlot uint64
	err = f.getObj(PeriodDir, LatestSlotKey, &latestSlot)
	if err != nil {
		return 0, false, err
	}
	return latestSlot, true, nil
}

func (f *FileStore) StoreRecursiveProof(period uint64, proof []byte, witness []byte) error {
	return f.InsertData(RecursiveDir, parseKey(period), storeProof{
		Period:    period,
		ProofType: common.SyncComRecursiveType,
		Proof:     hex.EncodeToString(proof),
		Witness:   hex.EncodeToString(witness),
	})
}

func (f *FileStore) GetRecursiveProof(period uint64) (*StoreProof, bool, error) {
	exists, err := f.CheckRecursiveProof(period)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	proof, err := f.GetObj(RecursiveDir, parseKey(period))
	if err != nil {
		return nil, false, err
	}
	return proof, true, nil
}

func (f *FileStore) GetRecursiveProofData(period uint64) ([]byte, error) {
	return f.GetData(RecursiveDir, parseKey(period))
}

func (f *FileStore) CheckRecursiveProof(period uint64) (bool, error) {
	return f.CheckStorageKey(RecursiveDir, parseKey(period))
}

func (f *FileStore) CheckUnitProof(period uint64) (bool, error) {
	return f.CheckStorageKey(UnitDir, parseKey(period))
}

func (f *FileStore) StoreUnitProof(period uint64, proof, witness []byte) error {
	logger.Info("store unit proof")
	return f.InsertData(UnitDir, parseKey(period), storeProof{
		Period:    period,
		ProofType: common.SyncComUnitType,
		Proof:     hex.EncodeToString(proof),
		Witness:   hex.EncodeToString(witness),
	})
}

func (f *FileStore) GetUnitProofData(period uint64) ([]byte, error) {
	return f.GetData(UnitDir, parseKey(period))
}

func (f *FileStore) GetUnitProof(period uint64) (*StoreProof, bool, error) {
	exists, err := f.CheckUnitProof(period)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	proof, err := f.GetObj(UnitDir, parseKey(period))
	if err != nil {
		return nil, false, err
	}
	return proof, true, nil
}

func (f *FileStore) StoreUpdate(period uint64, data interface{}) error {
	return f.InsertData(UpdateDir, parseKey(period), data)
}

func (f *FileStore) GetUpdateData(period uint64) ([]byte, error) {
	return f.GetData(UpdateDir, parseKey(period))
}

func (f *FileStore) CheckUpdate(period uint64) (bool, error) {
	return f.CheckStorageKey(UpdateDir, parseKey(period))
}

func (f *FileStore) GetUpdate(period uint64, value interface{}) (bool, error) {
	exists, err := f.CheckUpdate(period)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return true, f.getObj(UpdateDir, parseKey(period), value)
}

func (f *FileStore) GetBootstrap(value interface{}) (bool, error) {
	exists, err := f.CheckBootstrap()
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return true, f.getObj(GenesisDir, GenesisRawData, value)
}

func (f *FileStore) GetBootstrapData() ([]byte, error) {
	return f.GetData(GenesisDir, GenesisRawData)
}

func (f *FileStore) StoreBootstrap(data interface{}) error {
	return f.InsertData(GenesisDir, GenesisRawData, data)
}

func (f *FileStore) CheckBootstrap() (bool, error) {
	return f.CheckStorageKey(GenesisDir, GenesisRawData)
}

func (f *FileStore) StoreGenesisProof(proof []byte, witness []byte) error {
	return f.InsertData(GenesisDir, GenesisProofKey, storeProof{
		Period:    f.genesisPeriod,
		ProofType: common.SyncComGenesisType,
		Proof:     hex.EncodeToString(proof),
		Witness:   hex.EncodeToString(witness),
	})
}

func (f *FileStore) CheckGenesisProof() (bool, error) {
	return f.CheckStorageKey(GenesisDir, GenesisProofKey)
}

func (f *FileStore) GetGenesisProof() (*StoreProof, bool, error) {
	exists, err := f.CheckGenesisProof()
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	proof, err := f.GetObj(GenesisDir, GenesisProofKey)
	if err != nil {
		return nil, false, err
	}
	return proof, true, nil
}

func (f *FileStore) GetGenesisProofData() ([]byte, error) {
	return f.GetData(GenesisDir, GenesisProofKey)
}

func (f *FileStore) StoreLatestPeriod(period uint64) error {
	return f.InsertData(PeriodDir, LatestPeriodKey, period)
}

func (f *FileStore) CheckLatestPeriod() (bool, error) {
	return f.CheckStorageKey(PeriodDir, LatestPeriodKey)
}

func (f *FileStore) GetLatestPeriod() (uint64, bool, error) {
	exists, err := f.CheckLatestPeriod()
	if err != nil {
		return 0, false, err
	}
	if !exists {
		return 0, false, nil
	}
	var period uint64
	err = f.getObj(PeriodDir, LatestPeriodKey, &period)
	if err != nil {
		logger.Error("get latest Period error:%v", err)
		return 0, false, err
	}
	return period, true, nil
}

func (f *FileStore) StoreTxInEth2Proof(hash string, data interface{}) error {
	key := fmt.Sprintf("%s/txInEth2", hash)
	return f.InsertData(Tx, key, data)
}

func (f *FileStore) StoreCheckPointFinalityProve(hash string, data interface{}) error {
	key := fmt.Sprintf("%s/checkPointFinality", hash)
	return f.InsertData(Tx, key, data)
}

func (f *FileStore) TxBlockIsParentOfCheckPointProve(hash string, data interface{}) error {
	key := fmt.Sprintf("%s/txBlockIsParentOfCheckPoint", hash)
	return f.InsertData(Tx, key, data)
}

func (f *FileStore) GetTxInEth2Proof(hash string, value interface{}) (bool, error) {
	key := fmt.Sprintf("%s/txInEth2", hash)
	exists, err := f.CheckStorageKey(Tx, key)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return true, f.getObj(Tx, key, value)
}

func (f *FileStore) GetCheckPointFinalityProve(hash string, value interface{}) (bool, error) {
	key := fmt.Sprintf("%s/checkPointFinality", hash)
	exists, err := f.CheckStorageKey(Tx, key)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return true, f.getObj(Tx, key, value)
}

func (f *FileStore) GetTxBlockIsParentOfCheckPointProve(hash string, value interface{}) (bool, error) {
	key := fmt.Sprintf("%s/txBlockIsParentOfCheckPoint", hash)
	exists, err := f.CheckStorageKey(Tx, key)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return true, f.getObj(Tx, key, value)
}

func (f *FileStore) txDirCheckOrCreate(hash string) error {
	txHashDir := fmt.Sprintf("%s/%s", f.txDir, hash)
	exists, err := fileExists(txHashDir)
	if err != nil {
		logger.Error("file exists error:%v", err)
		return err
	}
	if !exists {
		err := os.MkdirAll(txHashDir, os.ModePerm)
		if err != nil {
			logger.Error("create tx dir error:%v", err)
			return err
		}
	}
	return nil
}

func (f *FileStore) Clear() error {
	return os.RemoveAll(f.dataDir)
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
	err = dirCheckOrCreate(storeKey)
	if err != nil {
		logger.Error("create dir error:%v", err)
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
	file, err := os.Create(storeKey)
	if err != nil {
		logger.Error("open file error:%v", err)
		return err
	}
	defer file.Close()
	dataBytes, err := json.Marshal(value)
	if err != nil {
		logger.Error("marshal file error:%v", err)
		return err
	}
	_, err = file.Write(dataBytes)
	if err != nil {
		logger.Error("write file error:%v", err)
		return err
	}
	err = file.Sync()
	if err != nil {
		logger.Error("sync file error:%v", err)
		return err
	}
	return nil
}

func (f *FileStore) GetData(table, key string) ([]byte, error) {
	storeKey, err := f.generateStoreKey(table, key)
	if err != nil {
		logger.Error("generate store key error:%v", err)
		return nil, err
	}
	exists, err := fileExists(storeKey)
	if err != nil {
		logger.Error("file exists error:%v", err)
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("no find key: %v-%v", table, key)
	}
	dataBytes, err := os.ReadFile(storeKey)
	if err != nil {
		logger.Error("read file error:%v", err)
		return nil, err
	}
	return dataBytes, nil
}

func (f *FileStore) GetObj(table, key string) (*StoreProof, error) {
	var obj storeProof
	err := f.getObj(table, key, &obj)
	if err != nil {
		return nil, err
	}
	proofBytes, err := hex.DecodeString(obj.Proof)
	if err != nil {
		return nil, err
	}
	witnessBytes, err := hex.DecodeString(obj.Witness)
	if err != nil {
		return nil, err
	}
	return &StoreProof{
		ProofType: obj.ProofType,
		Period:    obj.Period,
		Hash:      obj.Hash,
		Proof:     proofBytes,
		Witness:   witnessBytes,
	}, nil

}

func (f *FileStore) getObj(table, key string, value interface{}) error {
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
	case Tx:
		return fmt.Sprintf("%s/%s", f.txDir, key), nil
	default:
		return "", fmt.Errorf("no find table: %v", table)
	}
}

func (f *FileStore) NeedGenBhfUpdateIndex() ([]uint64, error) {
	latestSlot, ok, err := f.GetLatestSlot()
	if err != nil {
		logger.Error("get latest slot error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest slot error")
	}
	if latestSlot <= f.genesisSlot {
		return nil, nil
	}
	files, err := traverseFile(f.bhfUpdateDir)
	if err != nil {
		return nil, err
	}
	var recoverFile []uint64
	for index := f.genesisSlot + common.BeaconHeaderSlot; index <= latestSlot; index = index + common.BeaconHeaderSlot {
		if _, ok := files[fmt.Sprintf("%d", index)]; !ok {
			recoverFile = append(recoverFile, index)
		}
	}
	return recoverFile, nil

}

func (f *FileStore) NeedUpdateIndexes() ([]uint64, error) {
	latestPeriod, ok, err := f.GetLatestPeriod()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest period error")
	}
	if latestPeriod <= f.genesisPeriod {
		return nil, nil
	}
	files, err := traverseFile(f.updateDir)
	if err != nil {
		return nil, err
	}
	var recoverFile []uint64
	for index := f.genesisPeriod; index <= latestPeriod; index++ {
		if _, ok := files[fmt.Sprintf("%d", index)]; !ok {
			recoverFile = append(recoverFile, index)
		}
	}
	return recoverFile, nil
}

func (f *FileStore) NeedGenUnitProofIndexes() ([]uint64, error) {
	latestPeriod, ok, err := f.GetLatestPeriod()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest period error")
	}

	if latestPeriod <= f.genesisPeriod {
		return nil, nil
	}
	files, err := traverseFile(f.unitDir)
	if err != nil {
		return nil, err
	}
	var recoverFile []uint64
	for index := f.genesisPeriod; index <= latestPeriod; index++ {
		if _, ok := files[fmt.Sprintf("%d", index)]; !ok {
			recoverFile = append(recoverFile, index)
		}
	}
	return recoverFile, nil
}

func (f *FileStore) NeedGenRecProofIndexes() ([]uint64, error) {
	latestPeriod, ok, err := f.GetLatestPeriod()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest period error")
	}
	if latestPeriod <= f.genesisPeriod {
		return nil, nil
	}
	files, err := traverseFile(f.recursiveDir)
	if err != nil {
		return nil, err
	}
	var recoverFile []uint64
	for index := f.genesisPeriod; index <= latestPeriod; index++ {
		if _, ok := files[fmt.Sprintf("%d", index)]; !ok {
			recoverFile = append(recoverFile, index)
		}
	}
	return recoverFile, nil
}

var numberPattern = regexp.MustCompile(`^\d+$`)

func traverseFile(path string) (map[string]string, error) {
	files := make(map[string]string)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filePath, err := getFileName(info.Name())
			if err != nil {
				return err
			}
			if numberPattern.MatchString(filePath) {
				files[filePath] = path
			}
		}
		return nil
	})
	return files, err
}

func WriteFile(path string, data []byte) error {
	err := dirCheckOrCreate(path)
	if err != nil {
		logger.Error("create dir error:%v", err)
		return err
	}
	exists, err := fileExists(path)
	if err != nil {
		logger.Error("file exists error:%v", err)
		return err
	}
	if exists {
		err := os.Remove(path)
		if err != nil {
			logger.Error("remove file error:%v", err)
			return err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		logger.Error("open file error:%v", err)
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		logger.Error("write file error:%v", err)
		return err
	}
	err = file.Sync()
	if err != nil {
		logger.Error("sync file error:%v", err)
		return err
	}
	return nil
}

func dirCheckOrCreate(path string) error {
	dir := filepath.Dir(path)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
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

func getFileName(path string) (string, error) {
	arrs := strings.Split(path, "/")
	if len(arrs) == 0 {
		return "", fmt.Errorf("get file name error")
	}
	return arrs[len(arrs)-1], nil
}
func dirNotExistsAndCreate(path string) (bool, error) {
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
