package node

import (
	"encoding/hex"
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/store"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type StoreProof struct {
	Id        string             `json:"id"`
	ProofType common.ZkProofType `json:"type"`
	Hash      string             `json:"hash"`
	Index     uint64             `json:"Index"`
	Proof     string             `json:"proof"`
	Witness   string             `json:"witness"`
}

const (
	LatestPeriodKey   = "latestPeriod"
	LatestSlotKey     = "latestFinalitySlot"
	GenesisRawKey     = "genesisRaw"
	SyncComGenesisKey = "syncComGenesisProof"
)

type Table string

const (
	IndexTable        Table = "index"
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
	BtcBulkTable      Table = "btcBulk"
	BtcPackedTable    Table = "btcPack"
	BtcWrapTable      Table = "btcWrap"
	BtcBaseTable      Table = "btcBase"
	BtcMiddleTable    Table = "btcMiddle"
	BtcUpperTable     Table = "btcUpper"
	BtcGenesisTable   Table = "btcGenesis"
	BtcRecursiveTable Table = "btcRecursive"
)

var InitStoreTables = []Table{IndexTable, GenesisTable, UpdateTable, OuterTable, UnitTable, RecursiveTable,
	FinalityTable, BhfTable, BeaconHeaderTable, TxesTable, RedeemTable, RequestTable, DepositTable, VerifyTable,
	BtcBulkTable, BtcPackedTable, BtcWrapTable, BtcBaseTable, BtcMiddleTable, BtcUpperTable, BtcGenesisTable, BtcRecursiveTable}

type FileStorage struct {
	RootPath         string
	FileStoreMap     map[Table]*store.FileStore
	lock             sync.Mutex
	genesisSlot      uint64
	genesisPeriod    uint64
	btcGenesisHeight uint64
}

func NewFileStorage(rootPath string, genesisSlot, btcGenesisHeight uint64, tables ...Table) (*FileStorage, error) {
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
		RootPath:         path,
		FileStoreMap:     fileStoreMap,
		genesisSlot:      genesisSlot,
		genesisPeriod:    genesisSlot / common.SlotPerPeriod,
		btcGenesisHeight: btcGenesisHeight,
	}, nil
}

func (fs *FileStorage) StoreObj(table Table, key string, value interface{}) error {
	fileStore, ok := fs.GetFileStore(table)
	if !ok {
		logger.Error("get file store error %v", table)
		return fmt.Errorf("get file store error %v", table)
	}
	return fileStore.Store(key, value)
}
func (fs *FileStorage) GetObj(table Table, key string, value interface{}) (bool, error) {
	fileStore, ok := fs.GetFileStore(table)
	if !ok {
		logger.Error("get file store error %v", table)
		return false, fmt.Errorf("get file store error %v", table)
	}
	return fileStore.Get(key, value)
}

func (fs *FileStorage) CheckObj(table Table, key string) (bool, error) {
	fileStore, ok := fs.GetFileStore(table)
	if !ok {
		logger.Error("get file store error %v", table)
		return false, fmt.Errorf("get file store error %v", table)
	}
	return fileStore.CheckExists(key)
}

func (fs *FileStorage) Store(table Table, pType common.ZkProofType, proof, witness []byte, keys ...interface{}) error {
	key := genKey(table, keys...)
	obj := newStoreProof(pType, key, proof, witness)
	return fs.StoreObj(table, key, obj)
}
func (fs *FileStorage) Get(table Table, value interface{}, key ...interface{}) (bool, error) {
	return fs.GetObj(table, genKey(table, key...), value)
}

func (fs *FileStorage) Check(table Table, key ...interface{}) (bool, error) {
	return fs.CheckObj(table, genKey(table, key...))
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

func (fs *FileStorage) GetGenesisPeriod() uint64 {
	return fs.genesisPeriod
}
func (fs *FileStorage) StoreRequest(req *common.ZkProofRequest) error {
	return fs.StoreObj(RequestTable, req.Id(), req)
}

func (fs *FileStorage) StorePeriod(period uint64) error {
	return fs.StoreObj(IndexTable, LatestPeriodKey, period)
}

func (fs *FileStorage) GetLatestPeriod() (uint64, bool, error) {
	var period uint64
	exists, err := fs.GetObj(IndexTable, LatestPeriodKey, &period)
	if err != nil {
		logger.Error("get Index error:%v", err)
		return 0, false, err
	}
	return period, exists, nil
}

func (fs *FileStorage) StoreLatestFinalizedSlot(slot uint64) error {
	return fs.StoreObj(IndexTable, LatestSlotKey, slot)
}

func (fs *FileStorage) GetFinalizedSlot() (uint64, bool, error) {
	var slot uint64
	exists, err := fs.GetObj(IndexTable, LatestSlotKey, &slot)
	if err != nil {
		logger.Error("get slot error:%v", err)
		return 0, false, err
	}
	return slot, exists, nil
}

func (fs *FileStorage) StoreUpdate(period uint64, value interface{}) error {
	return fs.StoreObj(UpdateTable, genKey(UpdateTable, period), value)
}

func (fs *FileStorage) CheckUpdate(period uint64) (bool, error) {
	return fs.CheckObj(UpdateTable, genKey(UpdateTable, period))
}
func (fs *FileStorage) GetUpdate(period uint64, value interface{}) (bool, error) {
	return fs.GetObj(UpdateTable, genKey(UpdateTable, period), value)
}

func (fs *FileStorage) StoreBootStrapBySlot(slot uint64, data interface{}) error {
	return fs.StoreObj(GenesisTable, genKey(GenesisTable, slot), data)
}

func (fs *FileStorage) GetBootStrapBySlot(slot uint64, value interface{}) (bool, error) {
	return fs.GetObj(GenesisTable, genKey(GenesisTable, slot), value)
}

func (fs *FileStorage) CheckBootStrapBySlot(slot uint64) (bool, error) {
	return fs.CheckObj(GenesisTable, genKey(GenesisTable, slot))
}

func (fs *FileStorage) StoreBootStrap(data interface{}) error {
	return fs.StoreObj(GenesisTable, GenesisRawKey, data)
}
func (fs *FileStorage) GetBootstrap(value interface{}) (bool, error) {
	return fs.GetObj(GenesisTable, GenesisRawKey, value)
}

func (fs *FileStorage) CheckBootstrap() (bool, error) {
	return fs.CheckObj(GenesisTable, GenesisRawKey)
}

func (fs *FileStorage) StoreOuterProof(period uint64, proof, witness []byte) error {
	return fs.Store(OuterTable, common.UnitOuter, proof, witness, period)
}
func (fs *FileStorage) CheckOuterProof(period uint64) (bool, error) {
	return fs.Check(OuterTable, period)
}

func (fs *FileStorage) GetOuterProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(OuterTable, &storeProof, period)
	if err != nil {
		logger.Error("get outer proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreUnitProof(period uint64, proof, witness []byte) error {
	return fs.Store(UnitTable, common.SyncComUnitType, proof, witness, period)
}

func (fs *FileStorage) CheckUnitProof(period uint64) (bool, error) {
	return fs.Check(UnitTable, period)
}

func (fs *FileStorage) GetUnitProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(UnitTable, &storeProof, period)
	if err != nil {
		logger.Error("get unit proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreGenesisProof(proof, witness []byte) error {
	obj := newStoreProof(common.SyncComGenesisType, SyncComGenesisKey, proof, witness)
	return fs.StoreObj(GenesisTable, SyncComGenesisKey, obj)
}

func (fs *FileStorage) CheckGenesisProof() (bool, error) {
	return fs.CheckObj(GenesisTable, SyncComGenesisKey)
}

func (fs *FileStorage) GetGenesisProof() (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.GetObj(GenesisTable, SyncComGenesisKey, &storeProof)
	if err != nil {
		logger.Error("get genesis proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreRecursiveProof(period uint64, proof, witness []byte) error {
	return fs.Store(RecursiveTable, common.SyncComRecursiveType, proof, witness, period)
}

func (fs *FileStorage) CheckRecursiveProof(period uint64) (bool, error) {
	return fs.Check(RecursiveTable, period)
}

func (fs *FileStorage) GetRecursiveProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(RecursiveTable, &storeProof, period)
	if err != nil {
		logger.Error("get recursive proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreBhfProof(period uint64, proof, witness []byte) error {
	return fs.Store(BhfTable, common.BeaconHeaderFinalityType, proof, witness, period)
}

func (fs *FileStorage) CheckBhfProof(period uint64) (bool, error) {
	return fs.Check(BhfTable, period)
}

func (fs *FileStorage) GetBhfProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(BhfTable, &storeProof, period)
	if err != nil {
		logger.Error("get recursive proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreBeaconHeaderProof(start, end uint64, proof, witness []byte) error {
	return fs.Store(BeaconHeaderTable, common.BeaconHeaderType, proof, witness, start, end)
}

func (fs *FileStorage) CheckBeaconHeaderProof(start, end uint64) (bool, error) {
	return fs.Check(BeaconHeaderTable, start, end)
}

func (fs *FileStorage) GetBeaconHeaderProof(start, end uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(BeaconHeaderTable, &storeProof, start, end)
	if err != nil {
		logger.Error("get recursive proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreFinalityUpdate(period uint64, data interface{}) error {
	return fs.StoreObj(FinalityTable, genKey(FinalityTable, period), data)
}

func (fs *FileStorage) CheckFinalityUpdate(period uint64) (bool, error) {
	return fs.CheckObj(FinalityTable, genKey(FinalityTable, period))
}

func (fs *FileStorage) GetFinalityUpdate(period uint64, value interface{}) (bool, error) {
	return fs.GetObj(FinalityTable, genKey(FinalityTable, period), value)
}

func (fs *FileStorage) StoreTxProof(txHash string, proof, witness []byte) error {
	return fs.Store(TxesTable, common.TxInEth2, proof, witness, txHash)
}

func (fs *FileStorage) CheckTxProof(txHash string) (bool, error) {
	return fs.Check(TxesTable, txHash)
}

func (fs *FileStorage) GetTxProof(txHash string) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(TxesTable, &storeProof, txHash)
	if err != nil {
		logger.Error("get tx proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreRedeemProof(txHash string, proof, witness []byte) error {
	return fs.Store(RedeemTable, common.RedeemTxType, proof, witness, txHash)
}

func (fs *FileStorage) CheckRedeemProof(txHash string) (bool, error) {
	return fs.Check(RedeemTable, txHash)
}

func (fs *FileStorage) GetRedeemProof(txHash string) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(RedeemTable, &storeProof, txHash)
	if err != nil {
		logger.Error("get redeem proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreDepositProof(txHash string, proof, witness []byte) error {
	return fs.Store(DepositTable, common.DepositTxType, proof, witness, txHash)
}
func (fs *FileStorage) CheckDepositProof(txHash string) (bool, error) {
	return fs.Check(DepositTable, txHash)
}

func (fs *FileStorage) GetDepositProof(txHash string) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(DepositTable, &storeProof, txHash)
	if err != nil {
		logger.Error("get deposit proof error:%v %v", txHash, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreVerifyProof(txHash string, proof, witness []byte) error {
	return fs.Store(VerifyTable, common.VerifyTxType, proof, witness, txHash)
}
func (fs *FileStorage) CheckVerifyProof(txHash string) (bool, error) {
	return fs.Check(VerifyTable, txHash)
}

func (fs *FileStorage) GetVerifyProof(txHash string) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(VerifyTable, &storeProof, txHash)
	if err != nil {
		logger.Error("get verify proof error:%v %v", txHash, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreBtcBulkProof(index, end uint64, proof, witness []byte) error {
	return fs.Store(BtcBulkTable, common.BtcBulkType, proof, witness, index, end)
}

func (fs *FileStorage) CheckBtcBulkProof(index, end uint64) (bool, error) {
	return fs.Check(BtcBulkTable, index, end)
}

func (fs *FileStorage) GetBtcBulkProof(index, end uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(BtcBulkTable, &storeProof, index, end)
	if err != nil {
		logger.Error("get btc bulk proof error:%v %v", index, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreBtcPackedProof(index uint64, proof, witness []byte) error {
	return fs.Store(BtcPackedTable, common.BtcPackedType, proof, witness, index)
}

func (fs *FileStorage) CheckBtcPackedProof(index uint64) (bool, error) {
	return fs.Check(BtcPackedTable, index)
}

func (fs *FileStorage) GetBtcPackedProof(index uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(BtcPackedTable, &storeProof, index)
	if err != nil {
		logger.Error("get btc packed proof error:%v %v", index, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreBtcWrapProof(index uint64, proof, witness []byte) error {
	return fs.Store(BtcWrapTable, common.BtcWrapType, proof, witness, index)
}

func (fs *FileStorage) CheckBtcWrapProof(index uint64) (bool, error) {
	return fs.Check(BtcWrapTable, index)
}

func (fs *FileStorage) GetBtcWrapProof(index uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(BtcWrapTable, &storeProof, index)
	if err != nil {
		logger.Error("get btc wrap proof error:%v %v", index, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreBtcBaseProof(proof, witness []byte, key ...interface{}) error {
	return fs.Store(BtcBaseTable, common.BtcBaseType, proof, witness, key...)
}

func (fs *FileStorage) CheckBtcBaseProof(key ...interface{}) (bool, error) {
	return fs.Check(BtcBaseTable, key...)
}
func (fs *FileStorage) GetBtcBaseProof(key ...interface{}) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(BtcBaseTable, &storeProof, key...)
	if err != nil {
		logger.Error("get btc base proof error:%v %v", key, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreBtcMiddleProof(proof, witness []byte, key ...interface{}) error {
	return fs.Store(BtcMiddleTable, common.BtcMiddleType, proof, witness, key...)
}

func (fs *FileStorage) CheckBtcMiddleProof(key ...interface{}) (bool, error) {
	return fs.Check(BtcMiddleTable, key...)
}

func (fs *FileStorage) GetBtcMiddleProof(key ...interface{}) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(BtcMiddleTable, &storeProof, key...)
	if err != nil {
		logger.Error("get btc middle proof error:%v %v", key, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreBtcUpperProof(proof, witness []byte, key ...interface{}) error {
	return fs.Store(BtcUpperTable, common.BtcUpperType, proof, witness, key...)
}

func (fs *FileStorage) CheckBtcUpperProof(key ...interface{}) (bool, error) {
	return fs.Check(BtcUpperTable, key...)
}

func (fs *FileStorage) GetBtcUpperProof(key ...interface{}) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(BtcUpperTable, &storeProof, key...)
	if err != nil {
		logger.Error("get btc upper proof error:%v %v", key, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreBtcGenesisProof(proof, witness []byte, key ...interface{}) error {
	return fs.Store(BtcGenesisTable, common.BtcGenesisType, proof, witness, key...)
}

func (fs *FileStorage) CheckBtcGenesisProof(key ...interface{}) (bool, error) {
	return fs.Check(BtcGenesisTable, key...)
}

func (fs *FileStorage) GetBtcGenesisProof(key ...interface{}) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(BtcGenesisTable, &storeProof, key...)
	if err != nil {
		logger.Error("get btc genesis proof error:%v %v", key, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) StoreBtcRecursiveProof(proof, witness []byte, key ...interface{}) error {
	return fs.Store(BtcRecursiveTable, common.BtcRecursiveType, proof, witness, key...)
}

func (fs *FileStorage) CheckBtcRecursiveProof(key ...interface{}) (bool, error) {
	return fs.Check(BtcRecursiveTable, key...)
}

func (fs *FileStorage) GetBtcRecursiveProof(key ...interface{}) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(BtcRecursiveTable, &storeProof, key...)
	if err != nil {
		logger.Error("get btc recursive proof error:%v %v", key, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func newStoreProof(proofType common.ZkProofType, id string, proof, witness []byte) *StoreProof {
	return &StoreProof{
		ProofType: proofType,
		Id:        id,
		Proof:     hex.EncodeToString(proof),   // todo
		Witness:   hex.EncodeToString(witness), // todo
	}
}

func (fs *FileStorage) NeedBtcUpEndIndexes(height uint64) ([]uint64, error) {
	fileStore, ok := fs.GetFileStore(BtcUpperTable)
	if !ok {
		return nil, fmt.Errorf("get file store error %v", BtcUpperTable)
	}
	indexes, err := fileStore.Indexes(getEndIndex)
	if err != nil {
		logger.Error("get update indexes error:%v", err)
		return nil, err
	}
	// todo
	var endIndexes []uint64
	for index := fs.btcGenesisHeight + common.BtcUpperDistance; index <= height; index = index + common.BtcUpperDistance {
		if _, ok := indexes[index]; !ok {
			endIndexes = append(endIndexes, index)
		}
	}
	return endIndexes, nil

}

func (fs *FileStorage) NeedBtcRecursiveEndIndex(height uint64) ([]uint64, error) {
	fileStore, ok := fs.GetFileStore(BtcRecursiveTable)
	if !ok {
		return nil, fmt.Errorf("get file store error %v", BtcRecursiveTable)
	}
	indexes, err := fileStore.Indexes(getEndIndex)
	if err != nil {
		logger.Error("get update indexes error:%v", err)
		return nil, err
	}
	/*
				example:
				up1: 0~2, up2 2~4, up3 4~6  up4 6~8
				genesis: 0~4(up1,up2)
				recursive1: 0~6(genesis,up3)
			    recursive2: 0~8(recursive1,up4)
			    ....
		 we are record startIndex,not endIndex
	*/
	var endIndexes []uint64
	for index := fs.btcGenesisHeight + common.BtcUpperDistance*3; index <= height; index = index + common.BtcUpperDistance {
		if _, ok := indexes[index]; !ok {
			endIndexes = append(endIndexes, index)
		}
	}
	return endIndexes, nil
}

func (fs *FileStorage) GetNearTxSlotFinalizedSlot(txSlot uint64) (uint64, bool, error) {
	// todo  more efficient
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
	var tmpIndexes []uint64
	for key, _ := range indexes {
		tmpIndexes = append(tmpIndexes, key)
	}
	sort.SliceStable(tmpIndexes, func(i, j int) bool {
		return tmpIndexes[i] < tmpIndexes[j]
	})
	var finalizedSlot uint64
	for _, index := range tmpIndexes {
		if index >= txSlot {
			finalizedSlot = index
			break
		}
	}
	if finalizedSlot-txSlot > 64 { // todo
		//logger.Warn("get txSlot nearest finalizedSlot %v txSlot %v, more than 64 ", finalizedSlot, txSlot)
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
	latestPeriod, ok, err := fs.GetLatestPeriod()
	if err != nil {
		logger.Error("get latest Index error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest Index error")
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
	latestPeriod, ok, err := fs.GetLatestPeriod()
	if err != nil {
		logger.Error("get latest Index error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest Index error")
	}
	indexes, err := fileStore.Indexes(getStartIndex)
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
	latestPeriod, ok, err := fs.GetLatestPeriod()
	if err != nil {
		logger.Error("get latest Index error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest Index error")
	}
	indexes, err := fileStore.AllIndexes()
	if err != nil {
		logger.Error("get update indexes error:%v", err)
		return nil, err
	}
	/*
	 	unit: u0,u1,u2,u3
	    genesis: 0~1 (u0,u1)
	    recursive1:0~2 (genesis,u2)
	    recursive2:0~3 (recursive1,u3)
	    ...
	*/
	var needUpdateIndex []uint64
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
	sort.SliceStable(needUpdateIndexes, func(i, j int) bool {
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

// Clear be careful when you use it,
func (fs *FileStorage) Clear() error {
	return os.RemoveAll(fs.RootPath)
}

func genKey(prefix Table, args ...interface{}) string {
	name := fmt.Sprintf("%v", prefix)
	for _, arg := range args {
		name = name + fmt.Sprintf("_%v", arg)
	}
	return name
}

func getStartIndex(fileName string) (uint64, error) {
	keys := strings.Split(fileName, "_")
	if len(keys) < 2 {
		return 0, fmt.Errorf("invalid file name %v", fileName)
	}
	/*
		base_200_100 -> 200
	*/
	index, err := strconv.ParseUint(keys[1], 10, 64)
	if err != nil {
		return 0, err
	}
	return index, nil
}

func getEndIndex(fileName string) (uint64, error) {
	keys := strings.Split(fileName, "_")
	if len(keys) != 3 {
		return 0, fmt.Errorf("invalid file name %v", fileName)
	}
	/*
		base_200_100 -> 100
	*/
	index, err := strconv.ParseUint(keys[2], 10, 64)
	if err != nil {
		return 0, err
	}
	return index, nil
}
