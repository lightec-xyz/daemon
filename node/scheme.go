package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
)

const (
	ProofPrefix           = "p_"  // p_ + hash
	TxPrefix              = "t_"  // t_ + hash
	DestChainHashPrefix   = "d_"  // d_ + hash
	UnGenProofPrefix      = "u_"  // u_ + hash
	UnSubmitTxPrefix      = "us_" // s_ + hash
	BeaconSlotPrefix      = "bs_"
	BeaconEthNumberPrefix = "bh_"

	TxSlotPrefix           = "ts_"
	TxFinalizeSlotPrefix   = "tfs_"
	PendingReqPrefix       = "pr_"
	PendingProofRespPrefix = "prp_"
)

const (
	workerIdKey = "workerIdKey"
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
}

type DbUnSubmitTx struct {
	Hash      string
	ProofType common.ZkProofType
	Proof     string
	Timestamp int64
}

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
	pTxID := fmt.Sprintf("%s%s", ProofPrefix, trimOx(txId))
	return pTxID
}

func DbBtcHeightPrefix(height int64, txId string) string {
	pTxID := fmt.Sprintf("%db_%s", height, trimOx(txId))
	return pTxID
}

func DbEthHeightPrefix(height int64, txId string) string {
	pTxID := fmt.Sprintf("%de_%s", height, trimOx(txId))
	return pTxID
}

func DbTxId(txId string) string {
	pTxID := fmt.Sprintf("%s%s", TxPrefix, trimOx(txId))
	return pTxID
}

func DbDestId(txId string) string {
	pTxID := fmt.Sprintf("%s%s", DestChainHashPrefix, trimOx(txId))
	return pTxID
}

func DbUnGenProofId(chain ChainType, txId string) string {
	pTxID := fmt.Sprintf("%s%d_%s", UnGenProofPrefix, chain, trimOx(txId))
	return pTxID
}

func DbUnSubmitTxId(txId string) string {
	pTxID := fmt.Sprintf("%s%s", UnSubmitTxPrefix, trimOx(txId))
	return pTxID
}

func DbAddrPrefixTxId(addr string, txId string) string {
	key := fmt.Sprintf("%s_%s", addr, trimOx(txId))
	return key
}

func DbBeaconSlotId(slot uint64) string {
	key := fmt.Sprintf("%s%d", BeaconSlotPrefix, slot)
	return key
}
func DbBeaconEthNumberId(number uint64) string {
	key := fmt.Sprintf("%s%d", BeaconEthNumberPrefix, number)
	return key
}

func DbTxSlotId(slot uint64, hash string) string {
	key := fmt.Sprintf("%s%d_%s", TxSlotPrefix, slot, hash)
	return key
}

func DbTxFinalizeSlotId(slot uint64, hash string) string {
	key := fmt.Sprintf("%s%d_%s", TxFinalizeSlotPrefix, slot, hash)
	return key
}

func DbPendingRequestId(id string) string {
	key := fmt.Sprintf("%s%s", PendingReqPrefix, id)
	return key
}
