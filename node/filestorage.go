package node

import (
	"encoding/hex"
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/store"
	"sort"
	"sync"
)

type StoreProof struct {
	ProofType common.ZkProofType `json:"type"`
	Hash      string             `json:"hash"`
	Period    uint64             `json:"period"`
	Proof     string             `json:"proof"`
	Witness   string             `json:"witness"`
}

const (
	LatestPeriodKey = "latest"
	LatestSlotKey   = "latestFinalitySlot"
	GenesisRawData  = "genesisRaw"
	GenesisProofKey = "genesisProof"
)

type Table string

const (
	PeriodTable       Table = "period"
	GenesisTable      Table = "genesis"
	UpdateTable       Table = "update"
	OuterTable        Table = "outer"
	UnitTable         Table = "unit"
	RecursiveTable    Table = "recursive"
	FinalityTable     Table = "finalityUpdate"
	BhfTable          Table = "bhf"
	BeaconHeaderTable Table = "beaconHeader"
	TxesTable         Table = "txes"
	RedeemTable       Table = "redeem"
	RequestTable      Table = "request"
	DepositTable      Table = "deposit"
	VerifyTable       Table = "verify"
)

var InitStoreTables = []Table{PeriodTable, GenesisTable, UpdateTable, OuterTable, UnitTable, RecursiveTable,
	FinalityTable, BhfTable, BeaconHeaderTable, TxesTable, RedeemTable, RequestTable, DepositTable, VerifyTable}

type FileStorage struct {
	RootPath      string
	FileStoreMap  map[Table]*store.FileStore
	lock          sync.Mutex
	genesisSlot   uint64
	genesisPeriod uint64
}

func NewFileStorage(rootPath string, genesisSlot uint64, tables []Table) (*FileStorage, error) {
	fileStoreMap := make(map[Table]*store.FileStore)
	path := fmt.Sprintf("%s/proofData", rootPath) // todo
	logger.Info("fileStorage path: %v", path)
	for _, key := range tables {
		fileStore, err := CreateFileStore(path, string(key))
		if err != nil {
			logger.Error("create file store error")
			return nil, err
		}
		fileStoreMap[key] = fileStore
	}
	return &FileStorage{
		RootPath:      path,
		FileStoreMap:  fileStoreMap,
		genesisSlot:   genesisSlot,
		genesisPeriod: genesisSlot / 8192,
	}, nil
}

func (fs *FileStorage) GetGenesisPeriod() uint64 {
	return fs.genesisPeriod
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

func (fs *FileStorage) StoreRequest(req *common.ZkProofRequest) error {
	return fs.Store(RequestTable, req.Id(), req)
}

func (fs *FileStorage) StorePeriod(period uint64) error {
	return fs.Store(PeriodTable, LatestPeriodKey, period)
}

func (fs *FileStorage) GetPeriod() (uint64, bool, error) {
	var period uint64
	exists, err := fs.Get(PeriodTable, LatestPeriodKey, &period)
	if err != nil {
		logger.Error("get period error:%v", err)
		return 0, false, err
	}
	return period, exists, nil
}

func (fs *FileStorage) StoreFinalizedSlot(slot uint64) error {
	return fs.Store(PeriodTable, LatestSlotKey, slot)
}

func (fs *FileStorage) GetFinalizedSlot() (uint64, bool, error) {
	var slot uint64
	exists, err := fs.Get(PeriodTable, LatestSlotKey, &slot)
	if err != nil {
		logger.Error("get slot error:%v", err)
		return 0, false, err
	}
	return slot, exists, nil
}

func (fs *FileStorage) StoreUpdate(period uint64, value interface{}) error {
	return fs.Store(UpdateTable, parseKey(period), value)
}

func (fs *FileStorage) CheckUpdate(period uint64) (bool, error) {
	return fs.Check(UpdateTable, parseKey(period))
}
func (fs *FileStorage) GetUpdate(period uint64, value interface{}) (bool, error) {
	return fs.Get(UpdateTable, parseKey(period), value)
}

func (fs *FileStorage) StoreBootStrapBySlot(slot uint64, data interface{}) error {
	return fs.Store(GenesisTable, parseKey(slot), data)
}

func (fs *FileStorage) GetBootStrapBySlot(slot uint64, value interface{}) (bool, error) {
	return fs.Get(GenesisTable, parseKey(slot), value)
}

func (fs *FileStorage) CheckBootStrapBySlot(slot uint64) (bool, error) {
	return fs.Check(GenesisTable, parseKey(slot))
}

func (fs *FileStorage) StoreBootStrap(data interface{}) error {
	return fs.Store(GenesisTable, GenesisRawData, data)
}
func (fs *FileStorage) GetBootstrap(value interface{}) (bool, error) {
	return fs.Get(GenesisTable, GenesisRawData, value)
}

func (fs *FileStorage) CheckBootstrap() (bool, error) {
	return fs.Check(GenesisTable, GenesisRawData)
}

func (fs *FileStorage) StoreOuterProof(period uint64, proof, witness []byte) error {
	return fs.Store(OuterTable, parseKey(period), newStoreProof(common.UnitOuter, period, "", proof, witness))
}
func (fs *FileStorage) CheckOuterProof(period uint64) (bool, error) {
	return fs.Check(OuterTable, parseKey(period))
}

func (fs *FileStorage) GetOuterProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(OuterTable, parseKey(period), &storeProof)
	if err != nil {
		logger.Error("get outer proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreUnitProof(period uint64, proof, witness []byte) error {
	return fs.Store(UnitTable, parseKey(period), newStoreProof(common.SyncComUnitType, period, "", proof, witness))
}

func (fs *FileStorage) CheckUnitProof(period uint64) (bool, error) {
	return fs.Check(UnitTable, parseKey(period))
}

func (fs *FileStorage) GetUnitProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(UnitTable, parseKey(period), &storeProof)
	if err != nil {
		logger.Error("get unit proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreGenesisProof(period uint64, proof, witness []byte) error {
	return fs.Store(GenesisTable, GenesisProofKey, newStoreProof(common.SyncComGenesisType, period, "", proof, witness))
}

func (fs *FileStorage) CheckGenesisProof() (bool, error) {
	return fs.Check(GenesisTable, GenesisProofKey)
}

func (fs *FileStorage) GetGenesisProof() (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(GenesisTable, GenesisProofKey, &storeProof)
	if err != nil {
		logger.Error("get genesis proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreRecursiveProof(period uint64, proof, witness []byte) error {
	return fs.Store(RecursiveTable, parseKey(period), newStoreProof(common.SyncComRecursiveType, period, "", proof, witness))
}

func (fs *FileStorage) CheckRecursiveProof(period uint64) (bool, error) {
	return fs.Check(RecursiveTable, parseKey(period))
}

func (fs *FileStorage) GetRecursiveProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(RecursiveTable, parseKey(period), &storeProof)
	if err != nil {
		logger.Error("get recursive proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreBhfProof(period uint64, proof, witness []byte) error {
	return fs.Store(BhfTable, parseKey(period), newStoreProof(common.BeaconHeaderFinalityType, period, "", proof, witness))
}

func (fs *FileStorage) CheckBhfProof(period uint64) (bool, error) {
	return fs.Check(BhfTable, parseKey(period))
}

func (fs *FileStorage) GetBhfProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(BhfTable, parseKey(period), &storeProof)
	if err != nil {
		logger.Error("get recursive proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreBeaconHeaderProof(period uint64, proof, witness []byte) error {
	return fs.Store(BeaconHeaderTable, parseKey(period), newStoreProof(common.BeaconHeaderType, period, "", proof, witness))
}

func (fs *FileStorage) CheckBeaconHeaderProof(period uint64) (bool, error) {
	return fs.Check(BeaconHeaderTable, parseKey(period))
}

func (fs *FileStorage) GetBeaconHeaderProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(BeaconHeaderTable, parseKey(period), &storeProof)
	if err != nil {
		logger.Error("get recursive proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreFinalityUpdate(period uint64, data interface{}) error {
	return fs.Store(FinalityTable, parseKey(period), data)
}

func (fs *FileStorage) CheckFinalityUpdate(period uint64) (bool, error) {
	return fs.Check(FinalityTable, parseKey(period))
}

func (fs *FileStorage) GetFinalityUpdate(period uint64, value interface{}) (bool, error) {
	return fs.Get(FinalityTable, parseKey(period), value)
}

func (fs *FileStorage) StoreTxProof(txHash string, proof, witness []byte) error {
	return fs.Store(TxesTable, txHash, newStoreProof(common.TxInEth2, 0, txHash, proof, witness))
}

func (fs *FileStorage) CheckTxProof(txHash string) (bool, error) {
	return fs.Check(TxesTable, txHash)
}

func (fs *FileStorage) GetTxProof(txHash string) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(TxesTable, txHash, &storeProof)
	if err != nil {
		logger.Error("get tx proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreRedeemProof(txHash string, proof, witness []byte) error {
	return fs.Store(RedeemTable, txHash, newStoreProof(common.RedeemTxType, 0, txHash, proof, witness))
}

func (fs *FileStorage) CheckRedeemProof(txHash string) (bool, error) {
	return fs.Check(RedeemTable, txHash)
}

func (fs *FileStorage) GetRedeemProof(txHash string) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(RedeemTable, txHash, &storeProof)
	if err != nil {
		logger.Error("get redeem proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreDepositProof(txHash string, proof, witness []byte) error {
	return fs.Store(DepositTable, txHash, newStoreProof(common.DepositTxType, 0, txHash, proof, witness))
}
func (fs *FileStorage) CheckDepositProof(txHash string) (bool, error) {
	return fs.Check(DepositTable, txHash)
}

func (fs *FileStorage) GetDepositProof(txHash string) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(DepositTable, txHash, &storeProof)
	if err != nil {
		logger.Error("get deposit proof error:%v %v", txHash, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreVerifyProof(txHash string, proof, witness []byte) error {
	return fs.Store(VerifyTable, txHash, newStoreProof(common.VerifyTxType, 0, txHash, proof, witness))
}
func (fs *FileStorage) CheckVerifyProof(txHash string) (bool, error) {
	return fs.Check(VerifyTable, txHash)
}

func (fs *FileStorage) GetVerifyProof(txHash string) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(VerifyTable, txHash, &storeProof)
	if err != nil {
		logger.Error("get verify proof error:%v %v", txHash, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func newStoreProof(proofType common.ZkProofType, period uint64, txHash string, proof, witness []byte) *StoreProof {
	return &StoreProof{
		Period:    period,
		ProofType: proofType,
		Hash:      txHash,
		Proof:     hex.EncodeToString(proof),   // todo
		Witness:   hex.EncodeToString(witness), // todo
	}
}

func (fs *FileStorage) GetNearTxSlotFinalizedSlot(txSlot uint64) (uint64, bool, error) {
	finalizedStore, ok := fs.GetFileStore(FinalityTable)
	if !ok {
		logger.Error("get file store error %v", FinalityTable)
		return 0, false, fmt.Errorf("get file store error %v", FinalityTable)
	}
	indexes, err := finalizedStore.AllIndexes()
	if err != nil {
		logger.Error("get update indexes error:%v", err)
		return 0, false, err
	}
	// todo
	var tmpIndexes []uint64
	for key, _ := range indexes {
		tmpIndexes = append(tmpIndexes, key)
	}
	sort.Slice(tmpIndexes, func(i, j int) bool {
		return tmpIndexes[i] < tmpIndexes[j]
	})
	var finalizedSlot uint64
	for _, index := range tmpIndexes {
		if index >= txSlot {
			finalizedSlot = index
			break
		}
	}
	if finalizedSlot-txSlot > 33 {
		logger.Warn("get txSlot nearest finalized  error slot %v, txSlot %v", finalizedSlot, txSlot)
		return 0, false, nil
	}

	return finalizedSlot, finalizedSlot != 0, nil

}

func (fs *FileStorage) NeedUpdateIndexes() ([]uint64, error) {
	fileStore, ok := fs.GetFileStore(UpdateTable)
	if !ok {
		logger.Error("get file store error %v", UpdateTable)
		return nil, fmt.Errorf("get file store error %v", UpdateTable)
	}
	latestPeriod, ok, err := fs.GetPeriod()
	if err != nil {
		logger.Error("get latest period error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest period error")
	}
	indexes, err := fileStore.AllIndexes()
	if err != nil {
		logger.Error("get update indexes error:%v", err)
		return nil, err
	}
	var needUpdateIndex []uint64
	for index := fs.genesisPeriod; index <= latestPeriod; index++ {
		if _, ok := indexes[index]; !ok {
			needUpdateIndex = append(needUpdateIndex, index)
		}
	}
	return needUpdateIndex, nil

}

func (fs *FileStorage) NeedGenUnitProofIndexes() ([]uint64, error) {
	fileStore, ok := fs.GetFileStore(UnitTable)
	if !ok {
		logger.Error("get file store error %v", UnitTable)
		return nil, fmt.Errorf("get file store error %v", UnitTable)
	}
	latestPeriod, ok, err := fs.GetPeriod()
	if err != nil {
		logger.Error("get latest period error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest period error")
	}
	indexes, err := fileStore.AllIndexes()
	if err != nil {
		logger.Error("get update indexes error:%v", err)
		return nil, err
	}
	var needUpdateIndex []uint64
	for index := fs.genesisPeriod; index <= latestPeriod; index++ {
		if _, ok := indexes[index]; !ok {
			needUpdateIndex = append(needUpdateIndex, index)
		}
	}
	return needUpdateIndex, nil
}

func (fs *FileStorage) NeedGenRecProofIndexes() ([]uint64, error) {
	fileStore, ok := fs.GetFileStore(RecursiveTable)
	if !ok {
		logger.Error("get file store error %v", RecursiveTable)
		return nil, fmt.Errorf("get file store error %v", RecursiveTable)
	}
	latestPeriod, ok, err := fs.GetPeriod()
	if err != nil {
		logger.Error("get latest period error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest period error")
	}
	indexes, err := fileStore.AllIndexes()
	if err != nil {
		logger.Error("get update indexes error:%v", err)
		return nil, err
	}
	var needUpdateIndex []uint64
	// todo
	for index := fs.genesisPeriod + 2; index <= latestPeriod; index++ {
		if _, ok := indexes[index]; !ok {
			needUpdateIndex = append(needUpdateIndex, index)
		}
	}
	return needUpdateIndex, nil
}

func (fs *FileStorage) NeedGenBhfUpdateIndex() ([]uint64, error) {
	bhfProofStore, ok := fs.GetFileStore(BhfTable)
	if !ok {
		logger.Error("get file store error %v", BhfTable)
		return nil, fmt.Errorf("get file store error %v", BhfTable)
	}
	finalityUpdateStore, ok := fs.GetFileStore(FinalityTable)
	if !ok {
		logger.Error("get file store error %v", FinalityTable)
		return nil, fmt.Errorf("get file store error %v", FinalityTable)
	}
	allFinalityIndexes, err := finalityUpdateStore.AllIndexes()
	if err != nil {
		logger.Error("get update indexes error:%v", err)
		return nil, err
	}
	bhfIndexes, err := bhfProofStore.AllIndexes()
	if err != nil {
		logger.Error("get update indexes error:%v", err)
		return nil, err
	}
	var needUpdateIndexes []uint64
	for key, _ := range allFinalityIndexes {
		if _, ok := bhfIndexes[key]; !ok {
			needUpdateIndexes = append(needUpdateIndexes, key)
		}
	}
	sort.Slice(needUpdateIndexes, func(i, j int) bool {
		return needUpdateIndexes[i] < needUpdateIndexes[j]
	})
	return needUpdateIndexes, nil
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
func parseKey(key interface{}) string {
	return fmt.Sprintf("%v", key)
}
