package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/rpc/ethereum/zkbridge"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	"math/big"
	"time"
)

type ZkParams struct {
	*zkbridge.IBtcTxVerifierPublicWitnessParams
	amount *big.Int
}

type Notify struct {
}

type WrapSyncCommittee struct {
	*proverType.SyncCommittee
	Version string
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

type ChainIndex struct {
	Genesis uint64
	Start   uint64
	End     uint64
	Step    uint64
	PreStep uint64
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
