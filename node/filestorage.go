package node

import (
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/store"
)

type StoreProof struct {
	Id        string `json:"id"`
	ProofType string `json:"type"`
	Proof     string `json:"proof"`
	Witness   string `json:"witness"`
}

type StoreKey struct {
	PType          common.ProofType
	Hash           string
	FIndex, SIndex uint64
	Prefix         uint64
	isCp           bool
	BlockTime      uint64
	TxIndex        uint32
}

func (sk *StoreKey) FileKey() store.FileKey {
	key := common.GenKey(sk.PType, sk.Prefix, sk.FIndex, sk.SIndex, sk.Hash)
	return key
}
func (sk *StoreKey) ProofId() string {
	return sk.FileKey().String()
}

var initStoreTables = []store.Table{common.IndexTable, common.UpdateTable, common.FinalityTable, common.RequestTable, common.InnerTable, common.OuterTable, common.UnitTable, common.GenesisTable,
	common.RecursiveTable, common.DutyTable, common.TxesTable, common.BeaconHeaderTable, common.BhfTable, common.RedeemTable, common.SgxRedeemTable, common.BackendRedeemTable, common.BtcBaseTable, common.BtcMiddleTable,
	common.BtcUpperTable, common.BtcBulkTable, common.BtcTimestampTable, common.BtcDuperRecursiveTable, common.BtcDepositTable, common.BtcChangeTable, common.BtcDepthRecursiveTable, common.BtcUpdateCpTable}

type FileStorage struct {
	RootPath         string
	FileMaps         *sync.Map // store.Table -> store.IFileStore
	genesisSlot      uint64
	genesisPeriod    uint64
	btcGenesisHeight uint64
	memoryStore      *MemoryStore
}

func NewFileStorage(path string, genesisSlot, btcGenesisHeight uint64, tables ...store.Table) (*FileStorage, error) {
	if btcGenesisHeight%common.BtcUpperDistance != 0 {
		return nil, fmt.Errorf("btcGenesisHeight must be a multiple of %d", common.BtcUpperDistance)
	}
	fileStoreMap := new(sync.Map)
	rootPath := fmt.Sprintf("%s/proofData", path)
	logger.Info("fileStorage path: %v", rootPath)
	if tables == nil {
		tables = initStoreTables
	}
	for _, key := range tables {
		filePath := fmt.Sprintf("%s/%s", rootPath, key)
		fileStore, err := store.NewFileStore(filePath, func(name store.FileKey) (uint64, bool) {
			_, _, end, err := FileKeyToIndex(name)
			if err != nil {
				return 0, false
			}
			return end, true
		})
		if err != nil {
			logger.Error("create table store error")
			return nil, err
		}
		fileStoreMap.Store(key, fileStore)
	}
	return &FileStorage{
		RootPath:         rootPath,
		FileMaps:         fileStoreMap,
		genesisSlot:      genesisSlot,
		genesisPeriod:    genesisSlot / common.SlotPerPeriod,
		btcGenesisHeight: btcGenesisHeight,
	}, nil
}

func (fs *FileStorage) GetGenesisPeriod() uint64 {
	return fs.genesisPeriod
}
func (fs *FileStorage) StoreRequest(req *common.ProofRequest, extraIds ...string) error {
	proofId := req.ProofId()
	for _, id := range extraIds {
		proofId = proofId + "_" + id
	}
	return fs.StoreObj(common.RequestTable, store.FileKey(proofId), req)
}

func (fs *FileStorage) StoreLatestPeriod(period uint64) error {
	return fs.StoreObj(common.IndexTable, common.LatestPeriodKey, period)
}

func (fs *FileStorage) GetLatestPeriod() (uint64, bool, error) {
	var period uint64
	exists, err := fs.GetObj(common.IndexTable, common.LatestPeriodKey, &period)
	if err != nil {
		logger.Error("get FIndex error:%v", err)
		return 0, false, err
	}
	if !exists {
		return fs.genesisPeriod, true, nil
	}
	return period, exists, nil
}

func (fs *FileStorage) StoreLatestFinalizedSlot(slot uint64) error {
	return fs.StoreObj(common.IndexTable, common.LatestSlotKey, slot)
}

func (fs *FileStorage) GetLatestFinalizedSlot() (uint64, bool, error) {
	var slot uint64
	exists, err := fs.GetObj(common.IndexTable, common.LatestSlotKey, &slot)
	if err != nil {
		logger.Error("get slot error:%v", err)
		return 0, false, err
	}
	return slot, exists, nil
}
func (fs *FileStorage) StoreFinalityUpdate(slot uint64, data interface{}) error {
	return fs.StoreObj(common.FinalityTable, store.GenFileKey(common.FinalityTable, slot), data)
}
func (fs *FileStorage) StoreErrorFinalityUpdate(key string, data interface{}) error {
	return fs.StoreObj(common.FinalityTable, store.GenFileKey(common.FinalityTable, key), data)
}

func (fs *FileStorage) CheckFinalityUpdate(slot uint64) (bool, error) {
	return fs.CheckObj(common.FinalityTable, store.GenFileKey(common.FinalityTable, slot))
}

func (fs *FileStorage) GetFinalityUpdate(slot uint64, value interface{}) (bool, error) {
	return fs.GetObj(common.FinalityTable, store.GenFileKey(common.FinalityTable, slot), value)
}

func (fs *FileStorage) StoreUpdate(period uint64, value interface{}) error {
	return fs.StoreObj(common.UpdateTable, store.GenFileKey(common.UpdateTable, period), value)
}

func (fs *FileStorage) CheckUpdate(period uint64) (bool, error) {
	return fs.CheckObj(common.UpdateTable, store.GenFileKey(common.UpdateTable, period))
}
func (fs *FileStorage) GetUpdate(period uint64, value interface{}) (bool, error) {
	return fs.GetObj(common.UpdateTable, store.GenFileKey(common.UpdateTable, period), value)
}

func (fs *FileStorage) StoreBootStrapBySlot(slot uint64, data interface{}) error {
	return fs.StoreObj(common.GenesisTable, store.GenFileKey(common.GenesisTable, slot), data)
}

func (fs *FileStorage) GetBootStrapBySlot(slot uint64, value interface{}) (bool, error) {
	return fs.GetObj(common.GenesisTable, store.GenFileKey(common.GenesisTable, slot), value)
}

func (fs *FileStorage) CheckBootStrapBySlot(slot uint64) (bool, error) {
	return fs.CheckObj(common.GenesisTable, store.GenFileKey(common.GenesisTable, slot))
}

func (fs *FileStorage) GetOuterProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.OuterTable, &storeProof, period)
	if err != nil {
		logger.Error("get outer proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetUnitProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.UnitTable, &storeProof, period)
	if err != nil {
		logger.Error("get unit proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetRecursiveProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.RecursiveTable, &storeProof, period)
	if err != nil {
		logger.Error("get recursive proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetDutyProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.DutyTable, &storeProof, period)
	if err != nil {
		logger.Error("get recursive proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetBtcTimestampProof(txHeight, latestHeight uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.BtcTimestampTable, &storeProof, txHeight, latestHeight)
	if err != nil {
		logger.Error("get recursive proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetTxProof(txHash string) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.TxesTable, &storeProof, txHash)
	if err != nil {
		logger.Error("get tx proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetBeaconHeaderProof(start, end uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.BeaconHeaderTable, &storeProof, start, end)
	if err != nil {
		logger.Error("get recursive proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetBhfProof(period uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.BhfTable, &storeProof, period)
	if err != nil {
		logger.Error("get recursive proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetSgxRedeemProof(hash string) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.SgxRedeemTable, &storeProof, hash)
	if err != nil {
		logger.Error("get recursive proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetRedeemProof(txHash string) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.RedeemTable, &storeProof, txHash)
	if err != nil {
		logger.Error("get Redeem proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetBackendRedeemProof(txHash string) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.BackendRedeemTable, &storeProof, txHash)
	if err != nil {
		logger.Error("get tx proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetBtcBulkProof(index, end uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exists, err := fs.Get(common.BtcBulkTable, &storeProof, index, end)
	if err != nil {
		logger.Error("get btc bulk proof error:%v %v %v", index, end, err)
		return nil, false, err
	}
	return &storeProof, exists, nil
}

func (fs *FileStorage) GetBtcBaseProof(start, end uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.BtcBaseTable, &storeProof, start, end)
	if err != nil {
		logger.Error("get btc base proof error:%v_%v %v", start, end, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetBtcMiddleProof(start, end uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.BtcMiddleTable, &storeProof, start, end)
	if err != nil {
		logger.Error("get btc middle proof error:%v %v %v", start, end, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetBtcUpperProof(start, end uint64) (*StoreProof, bool, error) {
	var storeProof StoreProof
	exist, err := fs.Get(common.BtcUpperTable, &storeProof, start, end)
	if err != nil {
		logger.Error("get btc upper proof error:%v %v %v", start, end, err)
		return nil, false, err
	}
	return &storeProof, exist, nil
}

func (fs *FileStorage) GetSyncInnerProof(period, index uint64) (*StoreProof, bool, error) {
	key := NewStoreKey(common.SyncComInnerType, "", period, index, 0)
	return fs.GetProof(key)

}

func (fs *FileStorage) FindBtcChainProof(height uint64) (*WrapStorageProof, bool, error) {
	fileStore, ok := fs.GetFileStore(common.BtcDuperRecursiveTable)
	if !ok {
		return nil, false, fmt.Errorf("get file store error %v", common.BtcDuperRecursiveTable)
	}
	fileKey, ok := keyBinarySearch(fileStore, height)
	if !ok {
		return nil, false, nil
	}
	_, start, end, err := FileKeyToIndex(fileKey)
	if err != nil {
		return nil, false, nil
	}
	var storeProof StoreProof
	exists, err := fileStore.Get(fileKey, &storeProof)
	if err != nil {
		return nil, false, err
	}
	return &WrapStorageProof{
		ChainIndex: ChainIndex{
			Genesis: fs.btcGenesisHeight,
			Start:   start,
			End:     end,
			Step:    end - start,
		},
		StoreProof: storeProof,
	}, exists, nil

}

func (fs *FileStorage) FindDepthProof(prefix, height uint64) (*WrapStorageProof, bool, error) {
	fileStore, ok := fs.GetSubFileStore(common.BtcDepthRecursiveTable, prefixToTable(prefix))
	if !ok {
		return nil, false, nil
	}
	fileKey, ok := keyBinarySearch(fileStore, height)
	if !ok {
		return nil, false, nil
	}
	_, start, end, err := FileKeyToIndex(fileKey)
	if err != nil {
		return nil, false, nil
	}
	var storeProof StoreProof
	exists, err := fileStore.Get(fileKey, &storeProof)
	if err != nil {
		return nil, false, err
	}
	return &WrapStorageProof{
		ChainIndex: ChainIndex{
			Genesis: fs.btcGenesisHeight,
			Start:   start,
			End:     end,
			Step:    end - start,
		},
		StoreProof: storeProof,
	}, exists, nil
}

func (fs *FileStorage) CurrentBtcCpDepthIndex(cpHeight uint64) (*ChainIndex, bool, error) {
	fileStore, ok := fs.GetSubFileStore(common.BtcDepthRecursiveTable, prefixToTable(cpHeight), true)
	if !ok {
		return nil, false, fmt.Errorf("get file store error %v", common.BtcDuperRecursiveTable)
	}
	key, value := fileStore.MaxIndex()
	_, fileKey, ok := toMaxIndex(key, value)
	index := ChainIndex{
		End: cpHeight + common.BtcCpMinDepth,
	}
	if !ok {
		return &index, true, nil
	}
	_, start, end, err := FileKeyToIndex(fileKey)
	if err != nil {
		return nil, false, err
	}
	if end >= index.End {
		index.Start = start
		index.End = end
		index.Step = end - start
	}
	return &index, true, nil
}

func (fs *FileStorage) BtcDepthIndex(cpHeight, genesis, height uint64) (uint64, bool, error) {
	fileStore, ok := fs.GetSubFileStore(common.BtcDepthRecursiveTable, prefixToTable(cpHeight), true)
	if !ok {
		return 0, false, fmt.Errorf("get file store error %v", common.BtcDuperRecursiveTable)
	}
	if len(fileStore.Keys()) == 0 {
		return cpHeight + genesis, true, nil
	}
	index, ok := lessThanOrEqualIndex(fileStore.Keys(), height)
	if !ok {
		return 0, false, nil
	}
	return uint64(index), true, nil

}

func (fs *FileStorage) BtcChainIndex(height uint64) (uint64, bool, error) {
	fileStore, ok := fs.GetFileStore(common.BtcDuperRecursiveTable)
	if !ok {
		return 0, false, fmt.Errorf("get file store error %v", common.BtcDuperRecursiveTable)
	}
	index, ok := lessThanOrEqualIndex(fileStore.Keys(), height)
	if !ok {
		return 0, false, nil
	}
	return index, true, nil

}

func (fs *FileStorage) CurrentBtcChainIndex() (*ChainIndex, bool, error) {
	fileStore, ok := fs.GetFileStore(common.BtcDuperRecursiveTable)
	if !ok {
		return nil, false, fmt.Errorf("get file store error %v", common.BtcDuperRecursiveTable)
	}
	key, value := fileStore.MaxIndex()
	_, fileKey, ok := toMaxIndex(key, value)
	index := ChainIndex{
		End: fs.btcGenesisHeight,
	}
	if !ok {
		return &index, true, nil
	}
	_, start, end, err := FileKeyToIndex(fileKey)
	if err != nil {
		logger.Error("%v", err)
		return nil, false, err
	}
	if end >= index.End {
		index.Start = start
		index.End = end
		index.Step = end - start
	}
	return &index, true, nil

}

func (fs *FileStorage) BtcBaseIndexes(start, height uint64) ([]uint64, error) {
	fileStore, ok := fs.GetFileStore(common.BtcBaseTable)
	if !ok {
		return nil, fmt.Errorf("get file store error %v", common.BtcBaseTable)
	}
	var tmpIndexes []uint64
	for index := start; index <= height-common.BtcBaseDistance; index = index + common.BtcBaseDistance {
		exists, err := fileStore.CheckExists(store.GenFileKey(common.BtcBaseTable, index, index+common.BtcBaseDistance))
		if err != nil {
			return nil, err
		}
		if !exists {
			tmpIndexes = append(tmpIndexes, index)
		}
	}
	return tmpIndexes, nil
}

func (fs *FileStorage) BtcMiddleIndexes(start, height uint64) ([]Index, error) {
	fileStore, ok := fs.GetFileStore(common.BtcMiddleTable)
	if !ok {
		return nil, fmt.Errorf("get file store error %v", common.BtcMiddleTable)
	}
	var tmpIndexes []Index
	for index := start; index <= height-common.BtcMiddleDistance; index = index + common.BtcMiddleDistance {
		exists, err := fileStore.CheckExists(store.GenFileKey(common.BtcMiddleTable, index, index+common.BtcMiddleDistance))
		if err != nil {
			return nil, err
		}
		if !exists {
			tmpIndexes = append(tmpIndexes, Index{
				Start: index,
				End:   index + common.BtcMiddleDistance,
			})
		}
	}
	return tmpIndexes, nil
}

func (fs *FileStorage) SyncComInnerIndexes() ([]Index, error) {
	latestPeriod, ok, err := fs.GetLatestPeriod()
	if err != nil {
		logger.Error("get latest period error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest period error")
	}
	var tmpIndexes []Index
	for period := fs.genesisPeriod; period <= latestPeriod; period++ {
		fileStore, ok := fs.getSubFileStore(common.SyncComInnerType, prefixToTable(period), false)
		if !ok {
			for index := 0; index < common.SyncInnerNum; index++ {
				tmpIndexes = append(tmpIndexes, Index{Prefix: period, Start: uint64(index)})
			}
		} else {
			for index := 0; index < common.SyncInnerNum; index++ {
				exists, err := fileStore.CheckExists(store.GenFileKey(common.InnerTable, period, index))
				if err != nil {
					return nil, err
				}
				if !exists {
					tmpIndexes = append(tmpIndexes, Index{Prefix: period, Start: uint64(index)})

				}
			}
		}
	}
	return tmpIndexes, nil
}

func (fs *FileStorage) GetTxFinalizedSlot(txSlot uint64) (uint64, bool, error) {
	fileStore, ok := fs.GetFileStore(common.FinalityTable)
	if !ok {
		logger.Error("get file store error %v", common.FinalityTable)
		return 0, false, fmt.Errorf("get file store error %v", common.FinalityTable)
	}
	index, ok := overValueIndex(fileStore.Keys(), txSlot)
	if !ok {
		return 0, false, nil
	}
	return index, true, nil

}

func (fs *FileStorage) NeedUpdateIndexes() ([]uint64, error) {
	fileStore, ok := fs.GetFileStore(common.UpdateTable)
	if !ok {
		logger.Error("get file store error %v", common.UpdateTable)
		return nil, fmt.Errorf("get file store error %v", common.UpdateTable)
	}
	latestPeriod, ok, err := fs.GetLatestPeriod()
	if err != nil {
		logger.Error("get latest FIndex error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest FIndex error")
	}
	var needUpdateIndex []uint64
	for index := fs.genesisPeriod; index <= latestPeriod; index++ {
		exists, err := fileStore.CheckExists(store.GenFileKey(common.UpdateTable, index))
		if err != nil {
			return nil, err
		}
		if !exists {
			needUpdateIndex = append(needUpdateIndex, index)
		}
	}
	return needUpdateIndex, nil
}

func (fs *FileStorage) GenOuterIndexes() ([]uint64, error) {
	fileStore, ok := fs.GetFileStore(common.OuterTable)
	if !ok {
		logger.Error("get file store error %v", common.OuterTable)
		return nil, fmt.Errorf("get file store error %v", common.OuterTable)
	}
	latestPeriod, ok, err := fs.GetLatestPeriod()
	if err != nil {
		logger.Error("get latest FIndex error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest FIndex error")
	}
	var indexes []uint64
	for index := fs.genesisPeriod; index <= latestPeriod; index++ {
		exists, err := fileStore.CheckExists(store.GenFileKey(common.OuterTable, index))
		if err != nil {
			return nil, err
		}
		if !exists {
			indexes = append(indexes, index)
		}
	}
	return indexes, nil
}

func (fs *FileStorage) NeedGenUnitProofIndexes() ([]uint64, error) {
	fileStore, ok := fs.GetFileStore(common.UnitTable)
	if !ok {
		logger.Error("get file store error %v", common.UnitTable)
		return nil, fmt.Errorf("get file store error %v", common.UnitTable)
	}
	latestPeriod, ok, err := fs.GetLatestPeriod()
	if err != nil {
		logger.Error("get latest FIndex error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest FIndex error")
	}
	var needUpdateIndex []uint64
	for index := fs.genesisPeriod; index <= latestPeriod; index++ {
		exists, err := fileStore.CheckExists(store.GenFileKey(common.UnitTable, index))
		if err != nil {
			return nil, err
		}
		if !exists {
			needUpdateIndex = append(needUpdateIndex, index)
		}
	}
	return needUpdateIndex, nil
}

func (fs *FileStorage) NeedDutyIndexes() ([]uint64, error) {
	fileStore, ok := fs.GetFileStore(common.DutyTable)
	if !ok {
		logger.Error("get file store error %v", common.DutyTable)
		return nil, fmt.Errorf("get file store error %v", common.DutyTable)
	}
	latestPeriod, ok, err := fs.GetLatestPeriod()
	if err != nil {
		logger.Error("get latest FIndex error:%v", err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("get latest FIndex error")
	}
	var indexes []uint64
	for index := fs.genesisPeriod + 1; index <= latestPeriod; index++ {
		exists, err := fileStore.CheckExists(store.GenFileKey(common.DutyTable, index))
		if err != nil {
			return nil, err
		}
		if !exists {
			indexes = append(indexes, index)
		}
	}
	return indexes, nil
}

func (fs *FileStorage) StoreObj(table store.Table, key store.FileKey, value interface{}) error {
	fileStore, ok := fs.GetFileStore(table)
	if !ok {
		logger.Error("get file store error %v", table)
		return fmt.Errorf("get file store error %v", table)
	}
	return fileStore.Store(key, value)
}
func (fs *FileStorage) GetObj(table store.Table, key store.FileKey, value interface{}) (bool, error) {
	fileStore, ok := fs.GetFileStore(table)
	if !ok {
		logger.Error("get file store error %v", table)
		return false, fmt.Errorf("get file store error %v", table)
	}
	return fileStore.Get(key, value)
}

func (fs *FileStorage) CheckObj(table store.Table, key store.FileKey) (bool, error) {
	fileStore, ok := fs.GetFileStore(table)
	if !ok {
		logger.Error("get file store error %v", table)
		return false, fmt.Errorf("get file store error %v", table)
	}
	return fileStore.CheckExists(key)
}

func (fs *FileStorage) GetSubFileStore(table, sub store.Table, create ...bool) (store.IFileStore, bool) {
	fileStore, ok := fs.GetFileStore(table)
	if !ok {
		return nil, false
	}
	subFileStore, ok := fileStore.SubFileStore(sub, create...) //todo
	if !ok {
		return nil, false
	}
	return subFileStore, true

}

func (fs *FileStorage) Get(table store.Table, value interface{}, key ...interface{}) (bool, error) {
	return fs.GetObj(table, store.GenFileKey(table, key...), value)
}

func (fs *FileStorage) GetFileStore(table store.Table) (store.IFileStore, bool) {
	if value, ok := fs.FileMaps.Load(table); ok {
		filestore, ok := value.(*store.FileStore)
		if !ok {
			return nil, false
		}
		return filestore, true
	}
	return nil, false
}

func (fs *FileStorage) GetRootPath() string {
	return fs.RootPath
}

func (fs *FileStorage) Clear() error {
	return os.RemoveAll(fs.RootPath)
}

func (fs *FileStorage) GetProof(key StoreKey) (*StoreProof, bool, error) {
	storageKey := common.GenKey(key.PType, key.Prefix, key.FIndex, key.SIndex, key.Hash)
	fileStore, ok := fs.FileStore(key)
	if !ok {
		return nil, false, nil
	}
	storeProof := StoreProof{}
	exists, err := fileStore.Get(storageKey, &storeProof)
	if err != nil {
		logger.Error("get proof error:%v", err)
		return nil, false, err
	}
	return &storeProof, exists, nil
}
func (fs *FileStorage) DelProof(key StoreKey) error {
	fileStoreKey := common.GenKey(key.PType, key.Prefix, key.FIndex, key.SIndex, key.Hash)
	fileStore, ok := fs.FileStore(key)
	if !ok {
		logger.Error("get file store error %v", key.PType)
		return fmt.Errorf("get file store error %v", key.PType)
	}
	return fileStore.Del(fileStoreKey)
}

func (fs *FileStorage) CheckProof(key StoreKey) (bool, error) {
	fileStoreKey := common.GenKey(key.PType, key.Prefix, key.FIndex, key.SIndex, key.Hash)
	fileStore, ok := fs.FileStore(key)
	if !ok {
		return false, nil
	}
	return fileStore.CheckExists(fileStoreKey)

}

func (fs *FileStorage) StoreProof(key StoreKey, proof, witness []byte) error {
	fileStoreKey := common.GenKey(key.PType, key.Prefix, key.FIndex, key.SIndex, key.Hash)
	fileStore, ok := fs.FileStore(key, true)
	if !ok {
		logger.Error("get file store error %v", key.PType)
		return fmt.Errorf("get file store error %v", key.PType)
	}
	storeProof := newStoreProof(key.PType, fileStoreKey.String(), proof, witness)
	return fileStore.Store(fileStoreKey, storeProof)
}

func (fs *FileStorage) getFileStore(pType common.ProofType, create ...bool) (store.IFileStore, bool) {
	table := common.ProofTypeToTable(pType)
	if value, ok := fs.FileMaps.Load(table); ok {
		filestore, ok := value.(*store.FileStore)
		if !ok {
			if len(create) > 0 && create[0] {
				path := fmt.Sprintf("%s/%s", fs.RootPath, table)
				fileStore, err := store.NewFileStore(path, func(name store.FileKey) (uint64, bool) {
					_, _, end, err := FileKeyToIndex(name)
					if err != nil {
						return 0, false
					}
					return end, true
				})
				if err != nil {
					return nil, false
				}
				fs.FileMaps.Store(table, fileStore)
				return fileStore, true
			} else {
				return nil, false
			}

		}
		return filestore, true
	}
	return nil, false
}

func (fs *FileStorage) FileStore(key StoreKey, create ...bool) (store.IFileStore, bool) {
	if key.Prefix != 0 { // todo fix
		return fs.getSubFileStore(key.PType, prefixToTable(key.Prefix), create...)
	} else {
		return fs.getFileStore(key.PType, create...)
	}
}

func (fs *FileStorage) getSubFileStore(pType common.ProofType, sub store.Table, create ...bool) (store.IFileStore, bool) {
	fileStore, ok := fs.getFileStore(pType)
	if !ok {
		return nil, false
	}
	subFileStore, ok := fileStore.SubFileStore(sub, create...) //todo
	if !ok {
		return nil, false
	}
	return subFileStore, true
}

type Index struct {
	Prefix uint64
	Start  uint64
	End    uint64
}

func (fs *FileStorage) RemoveBtcProof(height uint64) error {
	logger.Warn("remove btc proof height <= %v", height)
	tables := []store.Table{common.BtcBaseTable, common.BtcTimestampTable, common.BtcMiddleTable, common.BtcUpperTable, common.BtcDuperRecursiveTable, common.BtcBulkTable}
	for _, table := range tables {
		fileStore, ok := fs.GetFileStore(table)
		if !ok {
			logger.Error("no find table %v", fileStore.RootPath())
			return fmt.Errorf("get file store error %v", table)
		}
		err := fs.removeFiles(fileStore, height)
		if err != nil {
			logger.Error("remove btc proof %v fileStore error %v", fileStore.RootPath(), err)
			continue
		}
	}
	depthStore, ok := fs.GetFileStore(common.BtcDepthRecursiveTable)
	if !ok {
		logger.Error("no find depth table")
		return fmt.Errorf("no find table")
	}
	subFileStores, err := depthStore.SubFileStores()
	if err != nil {
		logger.Error("get all subFileStores error %v", err)
		return err
	}
	for _, fileStore := range subFileStores {
		err := fs.removeFiles(fileStore, height)
		if err != nil {
			logger.Error("remove btc depth %v fileStore error %v", fileStore.RootPath(), err)
			continue
		}
	}
	return nil
}

func (fs *FileStorage) removeFiles(fileStore store.IFileStore, height uint64) error {
	indexes := fileStore.Keys()
	for _, index := range indexes {
		if index.(uint64) >= height {
			logger.Warn("delete fileStore %v : >=%v proof", fileStore.RootPath(), index.(uint64))
			pattern := fmt.Sprintf("*_%v", index.(uint64))
			err := fileStore.DelMatch(pattern, index.(uint64))
			if err != nil {
				logger.Error("del match error %v", err)
				return err
			}
		}
	}
	return nil
}

func newStoreProof(proofType common.ProofType, id string, proof, witness []byte) *StoreProof {
	return &StoreProof{
		ProofType: proofType.Name(),
		Id:        id,
		Proof:     hex.EncodeToString(proof),   // todo
		Witness:   hex.EncodeToString(witness), // todo
	}
}

func prefixToTable(prefix interface{}) store.Table {
	return store.Table(fmt.Sprintf("%v", prefix))
}

func NewStoreKey(proofType common.ProofType, hash string, prefix, fIndex, sIndex uint64) StoreKey {
	return StoreKey{PType: proofType, Hash: hash, Prefix: prefix, FIndex: fIndex, SIndex: sIndex}
}

func NewHashStoreKey(proofType common.ProofType, hash string) StoreKey {
	return StoreKey{PType: proofType, Hash: hash}
}

func NewHeightStoreKey(proofType common.ProofType, height uint64) StoreKey {
	return StoreKey{PType: proofType, FIndex: height}
}

func NewDoubleStoreKey(proofType common.ProofType, fIndex, sIndex uint64) StoreKey {
	return StoreKey{PType: proofType, FIndex: fIndex, SIndex: sIndex}
}
func NewPrefixStoreKey(proofType common.ProofType, prefix, fIndex, sIndex uint64) StoreKey {
	return StoreKey{PType: proofType, Prefix: prefix, FIndex: fIndex, SIndex: sIndex, isCp: true}
}

func FileKeyToIndex(fileKey store.FileKey) (uint64, uint64, uint64, error) {
	ids := strings.Split(fileKey.String(), "_")
	if len(ids) == 2 { //btcchain_3195552
		height, err := strconv.ParseUint(ids[1], 10, 64)
		if err != nil {
			return 0, 0, 0, err
		}
		return 0, 0, height, nil
	} else if len(ids) == 3 { //btcbase_3192280_3192336
		start, err := strconv.ParseUint(ids[1], 10, 64)
		if err != nil {
			return 0, 0, 0, err
		}
		end, err := strconv.ParseUint(ids[2], 10, 64)
		if err != nil {
			return 0, 0, 0, err
		}
		return 0, start, end, nil
	} else if len(ids) == 4 { //btcdepthrecursive_3191403_3195519_3195603
		prefix, err := strconv.ParseUint(ids[1], 10, 64)
		if err != nil {
			return 0, 0, 0, err
		}
		start, err := strconv.ParseUint(ids[2], 10, 64)
		if err != nil {
			return 0, 0, 0, err
		}

		end, err := strconv.ParseUint(ids[3], 10, 64)
		if err != nil {
			return 0, 0, 0, err
		}
		return prefix, start, end, nil
	} else {
		//txineth2_cbad3bd4ac28531a21c7d05c21d78b7684adcaa779e974fb5c59cd50fe8dcbbe
		return 0, 0, 0, fmt.Errorf("unexpected proof id %v", fileKey)
	}
}

func toMaxIndex(key, value interface{}) (uint64, store.FileKey, bool) {
	if key == nil || value == nil {
		return 0, "", false
	}
	return key.(uint64), value.(store.FileKey), true
}

func keyBinarySearch(fileStore store.IFileStore, height uint64) (store.FileKey, bool) {
	keys := fileStore.Keys()
	i := sort.Search(len(keys), func(i int) bool {
		return keys[i].(uint64) >= height
	})
	if i >= len(keys) {
		return "", false
	}
	index := keys[i].(uint64)
	if index != height {
		return "", false
	}
	value, ok := fileStore.GetValue(index)
	if !ok {
		return "", false
	}
	return value.(store.FileKey), true
}

type WrapStorageProof struct {
	StoreProof
	ChainIndex
}

func overValueIndex(indexes []interface{}, value uint64) (uint64, bool) {
	for _, num := range indexes {
		v := num.(uint64)
		if v > value {
			return v, true
		}
	}
	return 0, false
}

func lessThanOrEqualIndex(indexes []interface{}, value uint64) (uint64, bool) {
	result := uint64(0)
	for _, num := range indexes {
		v := num.(uint64)
		if v > value {
			break
		}
		result = v
	}
	return result, result != 0
}
