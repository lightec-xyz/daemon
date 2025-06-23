package node

import (
	"bytes"
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	"math/big"
	"strconv"
	"sync"
	"time"
)

type WorkerCount struct {
	worker map[string]time.Time
	sync.Mutex
}

func NewWorkerCount() *WorkerCount {
	return &WorkerCount{
		worker: make(map[string]time.Time),
	}
}

func (w *WorkerCount) AddWorker(workerId string) {
	w.Lock()
	defer w.Unlock()
	w.worker[workerId] = time.Now()
}

func (w *WorkerCount) Len() int {
	w.Lock()
	defer w.Unlock()
	for id, t := range w.worker {
		//todo
		if time.Now().Sub(t) > 2*time.Hour {
			delete(w.worker, id)
		}
	}
	return len(w.worker)
}

type Notify struct {
}

type WrapSyncCommittee struct {
	*proverType.SyncCommittee
	Version string
}

type UpdateCp struct {
	Height    uint64
	BlockTime uint64
	TxId      string
}

type BlockHeader struct {
	Headers []string
}

type ReScnSignal struct {
	Height uint64
}

type MinerPower struct {
	Address    string
	Power      uint64
	CreateTime time.Time
}

func (m *MinerPower) AddConstant(constant uint64) {
	m.Power += constant
}

func (m *MinerPower) AvgConstantPerHour() float64 {
	hours := time.Now().Sub(m.CreateTime).Hours()
	if hours <= 0 {
		return 0
	}
	return float64(m.Power) / hours
}

func NewMinerPower(address string, power uint64, createTime time.Time) *MinerPower {
	return &MinerPower{
		Address:    address,
		Power:      power,
		CreateTime: createTime,
	}
}

type FetchType int

const (
	GenesisUpdateType FetchType = iota + 1
	PeriodUpdateType
	FinalityUpdateType
)

func (ft FetchType) String() string {
	switch ft {
	case GenesisUpdateType:
		return "genesisUpdateType"
	case PeriodUpdateType:
		return "periodUpdateType"
	case FinalityUpdateType:
		return "finalityUpdateType"
	default:
		return "unknown"
	}
}

type DownloadStatus int

type FetchRequest struct {
	UpdateType FetchType
	Status     DownloadStatus
	period     uint64
}

type FetchResponse struct {
	FetchId    string
	Index      uint64
	UpdateType FetchType
	data       interface{}
}

func NewFetchResponse(updateType FetchType, index uint64, data interface{}) *FetchResponse {
	return &FetchResponse{
		FetchId:    NewFetchId(updateType, index),
		Index:      index,
		UpdateType: updateType,
		data:       data,
	}
}

func NewFetchId(updateType FetchType, index uint64) string {
	return fmt.Sprintf("%v_%v", updateType.String(), index)
}

func (f *FetchResponse) Id() string {
	return f.FetchId
}

type Utxo struct {
	TxId  string `json:"txId"`
	Index uint32 `json:"FIndex"`
}

type TxOut struct {
	Value    int64
	PkScript []byte
}

func formatUtxo(utxos []Utxo) string {
	var buf bytes.Buffer
	for _, vin := range utxos {
		buf.WriteString(vin.TxId)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(int(vin.Index)))
		buf.WriteString(",")
	}
	return buf.String()
}
func formatOut(outputs []TxOut) string {
	var buf bytes.Buffer
	for _, out := range outputs {
		buf.WriteString(fmt.Sprintf("%x", out.PkScript))
		buf.WriteString(":")
		buf.WriteString(fmt.Sprintf("%v", out.Value))
		buf.WriteString(",")
	}
	return buf.String()
}

type Transaction struct {
	Height    uint64
	TxIndex   uint
	Hash      string
	BlockHash string
	BlockTime uint64
	TxType    common.TxType
	ChainType common.ChainType
	ProofType common.ProofType

	Proved bool
	Amount int64
	// bitcoin chain
	EthAddr string
	BtcFrom []string
	Utxo    []Utxo
	// ethereum chain
	LogIndex  uint
	UtxoId    string
	UtxoIndex int64
	Sender    string
	Receiver  string
}

type ChainIndex struct {
	Genesis uint64
	Start   uint64
	End     uint64
	Step    uint64
	PreStep uint64
}

type RedeemReward struct {
	TxId   string
	Reward *big.Int
}

type UniqueList struct {
	Ids []string
	Map map[string]bool
}

func (ul *UniqueList) Add(value string) {
	if _, ok := ul.Map[value]; !ok {
		ul.Ids = append(ul.Ids, value)
		ul.Map[value] = true
	}
}
func (ul *UniqueList) List() []string {
	return ul.Ids
}

func NewUniqueList() *UniqueList {
	return &UniqueList{
		Ids: make([]string, 0),
		Map: make(map[string]bool),
	}
}
