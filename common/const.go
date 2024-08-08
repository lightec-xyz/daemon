package common

import (
	"fmt"
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
	BitcoinChain ChainType = iota + 1
	EthereumChain
)

func (ct *ChainType) String() string {
	switch *ct {
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

func (ps *ProofStatus) String() string {
	switch *ps {
	case ProofDefault:
		return "proofDefault"
	case ProofFinalized:
		return "proofFinalized"
	case ProofQueued:
		return "proofQueued"
	case ProofGenerating:
		return "proofGenerating"
	case ProofSuccess:
		return "proofSuccess"
	case ProofFailed:
		return "proofFailed"
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
	DepositTxType ZkProofType = iota + 1
	RedeemTxType
	TxInEth2
	VerifyTxType
	SyncComGenesisType
	SyncComUnitType
	SyncComRecursiveType
	BeaconHeaderFinalityType
	UnitOuter
	BeaconHeaderType
	BtcBulkType
	BtcPackedType
	BtcWrapType
	BtcBaseType
	BtcMiddleType
	BtcUpperType
	BtcGenesisType
	BtcRecursiveType
)

func (z *ZkProofType) String() string {
	switch *z {
	case DepositTxType:
		return "depositTxType"
	case RedeemTxType:
		return "redeemTxType"
	case VerifyTxType:
		return "verifyTxType"
	case SyncComGenesisType:
		return "syncComGenesisType"
	case SyncComUnitType:
		return "syncComUnitType"
	case SyncComRecursiveType:
		return "syncComRecursiveType"
	case TxInEth2:
		return "txInEth2"
	case UnitOuter:
		return "unitOuter"
	case BeaconHeaderFinalityType:
		return "beaconHeaderFinalityType"
	case BeaconHeaderType:
		return "beaconHeaderType"
	case BtcPackedType:
		return "btcPackedType"
	case BtcWrapType:
		return "btcWrapType"
	case BtcBulkType:
		return "btcBulkType"
	case BtcBaseType:
		return "btcBaseType"
	case BtcMiddleType:
		return "btcMiddleType"
	case BtcUpperType:
		return "btcUpperType"
	case BtcGenesisType:
		return "btcGenesisType"
	case BtcRecursiveType:
		return "btcRecursiveType"
	default:
		return "unKnown"
	}
}

func (z *ZkProofType) Weight() ProofWeight {
	switch *z {
	case DepositTxType:
		return Highest
	case RedeemTxType:
		return TenLevel
	case VerifyTxType:
		return NineLevel
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

func ToZkProofType(str string) (ZkProofType, error) {
	switch str {
	case "depositTxType":
		return DepositTxType, nil
	case "redeemTxType":
		return RedeemTxType, nil
	case "verifyTxType":
		return VerifyTxType, nil
	case "syncComGenesisType":
		return SyncComGenesisType, nil
	case "syncComUnitType":
		return SyncComUnitType, nil
	case "syncComRecursiveType":
		return SyncComRecursiveType, nil
	case "txInEth2":
		return TxInEth2, nil
	case "unitOuter":
		return UnitOuter, nil
	case "beaconHeaderFinalityType":
		return BeaconHeaderFinalityType, nil
	case "beaconHeaderType":
		return BeaconHeaderType, nil
	case "btcBulkType":
		return BtcBulkType, nil
	case "btcPackedType":
		return BtcPackedType, nil
	case "btcWrapType":
		return BtcWrapType, nil
	case "btcBaseType":
		return BtcBaseType, nil
	case "btcMiddleType":
		return BtcMiddleType, nil
	case "btcUpperType":
		return BtcUpperType, nil
	case "btcGenesisType":
		return BtcGenesisType, nil
	case "btcRecursiveType":
		return BtcRecursiveType, nil
	default:
		return 0, fmt.Errorf("uKnown zk proof type %v", str)
	}
}
