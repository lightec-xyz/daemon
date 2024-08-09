package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"strings"
	"time"
)

const ProtocolSeparator = "_"

const (
	ProofPrefix           = "p"  // p_ + hash
	TxPrefix              = "t"  // t_ + hash
	DestChainHashPrefix   = "d"  // d_ + hash
	UnGenProofPrefix      = "u"  // u_ + hash
	UnSubmitTxPrefix      = "us" // s_ + hash
	BeaconSlotPrefix      = "bs"
	BeaconEthNumberPrefix = "bh"

	TxSlotPrefix             = "ts"
	TxFinalizeSlotPrefix     = "tfs"
	PendingReqPrefix         = "pr"
	PendingProofRespPrefix   = "prp"
	NoncePrefix              = "n"
	UnConfirmTxPrefix        = "un"
	TaskTimePrefix           = "tt"
	FinalityUpdateSlotPrefix = "fu"
	BtcHeightPrefix          = "b"
	EthHeightPrefix          = "e"
)

const (
	workerIdKey = "workerIdKey"
	zkVerifyKey = "zkVerifyKey"
)

var (
	btcCurHeightKey = []byte("btcCurHeight")
	ethCurHeightKey = []byte("ethCurHeight")
	beaconLatestKey = []byte("beaconLatest")
)

type DbTx struct {
	TxHash    string
	Height    uint64
	TxIndex   uint
	TxType    common.TxType
	ChainType common.ChainType
	Amount    int64
}

type DbProof struct {
	TxHash    string             `json:"txId"`
	ProofType common.ZkProofType `json:"type"`
	Status    int                `json:"status"`
	Proof     string             `json:"Proof"`
}

type DbUnGenProof struct {
	TxHash    string
	ProofType common.ZkProofType
	ChainType common.ChainType
	Height    uint64
	TxIndex   uint
	Amount    uint64
}

type DbUnSubmitTx struct {
	Hash      string
	ProofType common.ZkProofType
	Proof     string
	Timestamp int64
}

type DbUnConfirmTx struct {
	Hash    string
	ProofId string
	Network string
}

type DbTask struct {
	Id        string
	ProofType common.ZkProofType
	StartTime time.Time
	EndTime   time.Time
}

func DbProofId(txId string) string {
	key := fmt.Sprintf("%s%s%s", ProofPrefix, ProtocolSeparator, trimOx(txId))
	return strings.ToLower(key)
}

func DbBtcHeightPrefix(height int64) string {
	key := fmt.Sprintf("%s%s%d", BtcHeightPrefix, ProtocolSeparator, height)
	return strings.ToLower(key)
}

func DbEthHeightPrefix(height int64) string {
	key := fmt.Sprintf("%s%s%d", EthHeightPrefix, ProtocolSeparator, height)
	return strings.ToLower(key)
}

func DbTxId(txId string) string {
	key := fmt.Sprintf("%s%s%s", TxPrefix, ProtocolSeparator, trimOx(txId))
	return strings.ToLower(key)
}

func DbDestId(txId string) string {
	key := fmt.Sprintf("%s%s%s", DestChainHashPrefix, ProtocolSeparator, trimOx(txId))
	return strings.ToLower(key)
}

func DbUnGenProofId(chain common.ChainType, txId string) string {
	key := fmt.Sprintf("%s%s%d%s%s", UnGenProofPrefix, ProtocolSeparator, chain, ProtocolSeparator, trimOx(txId))
	return strings.ToLower(key)
}

func DbUnSubmitTxId(txId string) string {
	key := fmt.Sprintf("%s%s%s", UnSubmitTxPrefix, ProtocolSeparator, trimOx(txId))
	return strings.ToLower(key)
}

func DbAddrPrefixTxId(addr string, txType common.TxType, txId string) string {
	key := fmt.Sprintf("%s%s%d%s%s", addr, ProtocolSeparator, txType, ProtocolSeparator, trimOx(txId))
	return strings.ToLower(key)
}

func DbBeaconSlotId(slot uint64) string {
	key := fmt.Sprintf("%s%s%d", BeaconSlotPrefix, ProtocolSeparator, slot)
	return strings.ToLower(key)
}
func DbBeaconEthNumberId(number uint64) string {
	key := fmt.Sprintf("%s%s%d", BeaconEthNumberPrefix, ProtocolSeparator, number)
	return strings.ToLower(key)
}

func DbTxSlotId(slot uint64, hash string) string {
	key := fmt.Sprintf("%s%s%d%s%s", TxSlotPrefix, ProtocolSeparator, slot, ProtocolSeparator, trimOx(hash))
	return strings.ToLower(key)
}

func DbTxFinalizeSlotId(slot uint64, hash string) string {
	key := fmt.Sprintf("%s%s%d%s%s", TxFinalizeSlotPrefix, ProtocolSeparator, slot, ProtocolSeparator, trimOx(hash))
	return strings.ToLower(key)
}

func DbPendingRequestId(id string) string {
	key := fmt.Sprintf("%s%s%s", PendingReqPrefix, ProtocolSeparator, id)
	return strings.ToLower(key)
}

func DbAddrNonceId(network, addr string) string {
	key := fmt.Sprintf("%s%s%s%s%s", NoncePrefix, ProtocolSeparator, network, ProtocolSeparator, addr)
	return strings.ToLower(key)
}

func DbUnConfirmTxId(txId string) string {
	key := fmt.Sprintf("%s%s%s", UnConfirmTxPrefix, ProtocolSeparator, trimOx(txId))
	return strings.ToLower(key)
}

// proof task time

func DbTaskTimeId(flag common.ProofStatus, id string) string {
	key := fmt.Sprintf("%s%s%s%s%s", TaskTimePrefix, ProtocolSeparator, id, ProtocolSeparator, flag.String())
	return strings.ToLower(key)
}

func DbFinalityUpdateSlotId(slot uint64) string {
	key := fmt.Sprintf("%s%s%d", FinalityUpdateSlotPrefix, ProtocolSeparator, slot)
	return strings.ToLower(key)
}
