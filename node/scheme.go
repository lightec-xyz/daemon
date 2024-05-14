package node

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
)

const (
	ProofPrefix         = "p_" // p_ + hash
	TxPrefix            = "t_" // t_ + hash
	DestChainHashPrefix = "d_" // d_ + hash
	UnGenProofPrefix    = "u_" // u_ + hash
)

var (
	btcCurHeightKey = []byte("btcCurHeight")
	ethCurHeightKey = []byte("ethCurHeight")
)

type DbTx struct {
	TxHash    string
	Height    int64
	TxType    TxType
	ChainType ChainType
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
		return "Unknown"
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
		return "Unknown"
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

func DbAddrPrefixTxId(addr string, txId string) string {
	key := fmt.Sprintf("%s_%s", addr, trimOx(txId))
	return key
}
