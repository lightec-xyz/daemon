package common

import (
	"fmt"
	"time"

	//btcproverCommon"github.com/lightec-xyz/btc_provers/circuits/common"
	btcproverCommon "github.com/lightec-xyz/btc_provers/circuits/common"
)

// env

const (
	ZkDebugEnv     = "zkDebug"
	ZkParameterDir = "zkParameterDir"
	ZkProofTypes   = "zkProofTypes"
	DbNameSpace    = "zkbtc"
)

const (
	BtcBaseDistance   = btcproverCommon.CapacityBaseLevel
	BtcMiddleDistance = btcproverCommon.CapacityMidLevel * btcproverCommon.CapacityBaseLevel
	BtcUpperDistance  = btcproverCommon.CapacityDifficultyBlock
	//BtcBaseDistance   = 2
	//BtcMiddleDistance = 2 * BtcBaseDistance
	//BtcUpperDistance  = 2 * BtcMiddleDistance
)

const (
	BeaconHeaderSlot      = 32
	SlotPerPeriod         = 8192
	MaxDiffTxFinalitySlot = 64
)

type TxType int

const (
	DepositTx TxType = iota + 1
	RedeemTx
)

func (tt TxType) String() string {
	switch tt {
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
	BitcoinChain ChainType = iota + 1
	EthereumChain
)

func (ct ChainType) String() string {
	switch ct {
	case BitcoinChain:
		return "bitcoin"
	case EthereumChain:
		return "ethereum"
	default:
		return "unknown"
	}
}

type Network string

const (
	ETH   Network = "eth"
	ICP   Network = "icp"
	Oasis Network = "oasis"
)

func (n Network) String() string {
	switch n {
	case ETH:
		return "eth"
	case Oasis:
		return "oasis"
	case ICP:
		return "icp"
	default:
		return "unknown"
	}
}

type ProofStatus int

const (
	ProofDefault ProofStatus = iota
	ProofFinalized
	ProofQueued
	ProofGenerating
	ProofSuccess
	ProofFailed
)

func (ps ProofStatus) String() string {
	switch ps {
	case ProofDefault:
		return "default"
	case ProofFinalized:
		return "finalized"
	case ProofQueued:
		return "queued"
	case ProofGenerating:
		return "generating"
	case ProofSuccess:
		return "success"
	case ProofFailed:
		return "failed"
	default:
		return "unknown"
	}
}

type Mode string

const (
	Client  Mode = "client"
	Cluster Mode = "cluster"
	Custom  Mode = "custom"
)

type ProofWeight int

const (
	WeightDefault ProofWeight = iota
	OneLevel
	TwoLevel
	ThreeLevel
	FourLevel
	FiveLevel
	SixLevel
	SevenLevel
	EightLevel
	NineLevel
	TenLevel
	Highest
)

type ZkProofType int

const (
	SyncComOuterType ZkProofType = iota + 1
	SyncComUnitType
	SyncComGenesisType
	SyncComRecursiveType
	TxInEth2
	BeaconHeaderType
	BeaconHeaderFinalityType
	RedeemTxType

	BtcBulkType
	BtcPackedType
	BtcBaseType
	BtcMiddleType
	BtcUpperType
	BtcGenesisType
	BtcDuperRecursive
	BtcDepthRecursiveType
	BtcChainType
	BtcDepositType
	BtcChangeType
)

func (z ZkProofType) String() string {
	switch z {

	case SyncComOuterType:
		return "unitOuter"
	case SyncComUnitType:
		return "syncComUnitType"
	case SyncComGenesisType:
		return "syncComGenesisType"
	case SyncComRecursiveType:
		return "syncComRecursiveType"
	case TxInEth2:
		return "txInEth2"
	case BeaconHeaderType:
		return "beaconHeaderType"
	case BeaconHeaderFinalityType:
		return "beaconHeaderFinalityType"
	case RedeemTxType:
		return "redeemTxType"
	case BtcBulkType:
		return "btcBulkType"
	case BtcPackedType:
		return "btcPackedType"
	case BtcBaseType:
		return "btcBaseType"
	case BtcMiddleType:
		return "btcMiddleType"
	case BtcUpperType:
		return "btcUpperType"
	case BtcGenesisType:
		return "btcGenesisType"
	case BtcDuperRecursive:
		return "btcDuperRecursive"
	case BtcDepthRecursiveType:
		return "btcDepthRecursiveType"
	case BtcChainType:
		return "btcChainType"
	case BtcDepositType:
		return "btcDepositType"
	case BtcChangeType:
		return "btcChangeType"
	default:
		return "unKnown"
	}
}

func (z ZkProofType) Weight() ProofWeight {
	switch z {
	case RedeemTxType:
		return TenLevel
	case SyncComRecursiveType, SyncComGenesisType:
		return EightLevel
	case SyncComUnitType:
		return SevenLevel
	case BeaconHeaderFinalityType, TxInEth2, BeaconHeaderType:
		return SixLevel
	default:
		return WeightDefault
	}
}

func (zk ZkProofType) Timeout() time.Duration {
	// todo
	switch zk {
	case SyncComUnitType:
		return 1*time.Hour + 30*time.Minute
	default:
		return 30 * time.Minute
	}
}

func ToZkProofType(str string) (ZkProofType, error) {
	switch str {
	//case "syncUnitOuter":
	//	return SyncComOuterType, nil
	case "syncComUnitType":
		return SyncComUnitType, nil
	case "syncComGenesisType":
		return SyncComGenesisType, nil
	case "syncComRecursiveType":
		return SyncComRecursiveType, nil
	case "txInEth2":
		return TxInEth2, nil
	case "beaconHeaderType":
		return BeaconHeaderType, nil
	case "beaconHeaderFinalityType":
		return BeaconHeaderFinalityType, nil
	case "redeemTxType":
		return RedeemTxType, nil
	case "btcBulkType":
		return BtcBulkType, nil
	case "btcPackedType":
		return BtcPackedType, nil
	case "btcBaseType":
		return BtcBaseType, nil
	case "btcMiddleType":
		return BtcMiddleType, nil
	case "btcUpperType":
		return BtcUpperType, nil
	case "btcDuperRecursive":
		return BtcDuperRecursive, nil
	case "btcDepthRecursiveType":
		return BtcDepthRecursiveType, nil
	case "btcChainType":
		return BtcChainType, nil
	case "btcDepositType":
		return BtcDepositType, nil
	case "btcChangeType":
		return BtcChangeType, nil
	default:
		return 0, fmt.Errorf("uKnown zk proof type %v", str)
	}
}
