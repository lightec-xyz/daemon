package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
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
	TxType    TxType
	ChainType ChainType
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
	ChainType ChainType
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

type TxStatus int

const (
	UnSend TxStatus = iota
	Send
	Confirmed
)

type TxType int

const (
	DepositTx TxType = iota + 1
	RedeemTx
)

func (tt *TxType) String() string {
	switch *tt {
	case DepositTx:
		return "deposit"
	case RedeemTx:
		return "redeem"
	default:
		return "unknown"
	}
}

type ChainType int

const (
	Bitcoin ChainType = iota + 1
	Ethereum
)

func (ct *ChainType) String() string {
	switch *ct {
	case Bitcoin:
		return "bitcoin"
	case Ethereum:
		return "ethereum"
	default:
		return "unknown"
	}
}

func DbProofId(txId string) string {
	pTxID := fmt.Sprintf("%s%s%s", ProofPrefix, ProtocolSeparator, trimOx(txId))
	return pTxID
}

func DbBtcHeightPrefix(height int64) string {
	pTxID := fmt.Sprintf("%s%s%d", BtcHeightPrefix, ProtocolSeparator, height)
	return pTxID
}

func DbEthHeightPrefix(height int64) string {
	pTxID := fmt.Sprintf("%s%s%d", EthHeightPrefix, ProtocolSeparator, height)
	return pTxID
}

func DbTxId(txId string) string {
	pTxID := fmt.Sprintf("%s%s%s", TxPrefix, ProtocolSeparator, trimOx(txId))
	return pTxID
}

func DbDestId(txId string) string {
	pTxID := fmt.Sprintf("%s%s%s", DestChainHashPrefix, ProtocolSeparator, trimOx(txId))
	return pTxID
}

func DbUnGenProofId(chain ChainType, txId string) string {
	pTxID := fmt.Sprintf("%s%s%d%s%s", UnGenProofPrefix, ProtocolSeparator, chain, ProtocolSeparator, trimOx(txId))
	return pTxID
}

func DbUnSubmitTxId(txId string) string {
	pTxID := fmt.Sprintf("%s%s%s", UnSubmitTxPrefix, ProtocolSeparator, trimOx(txId))
	return pTxID
}

func DbAddrPrefixTxId(addr string, txId string) string {
	key := fmt.Sprintf("%s%s%s", addr, ProtocolSeparator, trimOx(txId))
	return key
}

func DbBeaconSlotId(slot uint64) string {
	key := fmt.Sprintf("%s%s%d", BeaconSlotPrefix, ProtocolSeparator, slot)
	return key
}
func DbBeaconEthNumberId(number uint64) string {
	key := fmt.Sprintf("%s%s%d", BeaconEthNumberPrefix, ProtocolSeparator, number)
	return key
}

func DbTxSlotId(slot uint64, hash string) string {
	key := fmt.Sprintf("%s%s%d%s%s", TxSlotPrefix, ProtocolSeparator, slot, ProtocolSeparator, trimOx(hash))
	return key
}

func DbTxFinalizeSlotId(slot uint64, hash string) string {
	key := fmt.Sprintf("%s%s%d%s%s", TxFinalizeSlotPrefix, ProtocolSeparator, slot, ProtocolSeparator, trimOx(hash))
	return key
}

func DbPendingRequestId(id string) string {
	key := fmt.Sprintf("%s%s%s", PendingReqPrefix, ProtocolSeparator, id)
	return key
}

func DbAddrNonceId(network, addr string) string {
	key := fmt.Sprintf("%s%s%s%s%s", NoncePrefix, ProtocolSeparator, network, ProtocolSeparator, addr)
	return key
}

func DbUnConfirmTxId(txId string) string {
	key := fmt.Sprintf("%s%s%s", UnConfirmTxPrefix, ProtocolSeparator, trimOx(txId))
	return key
}

// proof task time

func DbTaskTimeId(flag common.TaskStatusFlag, id string) string {
	key := fmt.Sprintf("%s%s%s%s%s", TaskTimePrefix, ProtocolSeparator, id, ProtocolSeparator, flag.String())
	return key
}

func DbFinalityUpdateSlotId(slot uint64) string {
	key := fmt.Sprintf("%s%s%d", FinalityUpdateSlotPrefix, ProtocolSeparator, slot)
	return key
}
