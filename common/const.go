package common

import (
	"fmt"
)

// env

const (
	ZkDebugEnv     = "ZkDebug"
	ZkParameterDir = "ZkParameterDir"
	ZkProofTypes   = "ZkProofTypes"
	DbNameSpace    = "zkbtc"
)

const (
	//BtcBaseDistance   = btcproverCommon.CapacityBaseLevel
	//BtcMiddleDistance = btcproverCommon.CapacityMidLevel * btcproverCommon.CapacityBaseLevel
	//BtcUpperDistance  = btcproverCommon.CapacityDifficultyBlock
	BtcBaseDistance   = 2
	BtcMiddleDistance = 2 * BtcBaseDistance
	BtcUpperDistance  = 2 * BtcMiddleDistance
)

const (
	BeaconHeaderSlot = 32
	SlotPerPeriod    = 8192
)

type TxType = int

const (
	DepositTx TxType = iota + 1
	RedeemTx
)

type ChainType = int

const (
	Bitcoin ChainType = iota + 1
	Ethereum
)

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

type TaskStatusFlag int

const (
	Start TaskStatusFlag = iota + 1
	Prove
	End
)

func (ts TaskStatusFlag) String() string {
	switch ts {
	case Start:
		return "start"
	case Prove:
		return "prove"
	case End:
		return "end"
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
		return "ProofDefault"
	case ProofFinalized:
		return "ProofFinalized"
	case ProofQueued:
		return "ProofQueued"
	case ProofGenerating:
		return "ProofGenerating"
	case ProofSuccess:
		return "ProofSuccess"
	case ProofFailed:
		return "ProofFailed"
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
	BeaconHeaderFinalityType //BeaconHeaderFinalityUpdate
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

func (zkpt *ZkProofType) String() string {
	switch *zkpt {
	case DepositTxType:
		return "DepositTxType"
	case RedeemTxType:
		return "RedeemTxType"
	case VerifyTxType:
		return "VerifyTxType"
	case SyncComGenesisType:
		return "SyncComGenesisType"
	case SyncComUnitType:
		return "SyncComUnitType"
	case SyncComRecursiveType:
		return "SyncComRecursiveType"
	case TxInEth2:
		return "TxInEth2"
	case UnitOuter:
		return "UnitOuter"
	case BeaconHeaderFinalityType:
		return "BeaconHeaderFinalityType"
	case BeaconHeaderType:
		return "BeaconHeaderType"
	case BtcPackedType:
		return "BtcPackedType"
	case BtcWrapType:
		return "BtcWrapType"
	case BtcBulkType:
		return "BtcBulkType"
	case BtcBaseType:
		return "BtcBaseType"
	case BtcMiddleType:
		return "BtcMiddleType"
	case BtcUpperType:
		return "BtcUpperType"
	case BtcGenesisType:
		return "BtcGenesisType"
	case BtcRecursiveType:
		return "BtcRecursiveType"
	default:
		return "unKnown"
	}
}

func (zkpt *ZkProofType) Weight() ProofWeight {
	// todo
	switch *zkpt {
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
	case "DepositTxType":
		return DepositTxType, nil
	case "RedeemTxType":
		return RedeemTxType, nil
	case "VerifyTxType":
		return VerifyTxType, nil
	case "SyncComGenesisType":
		return SyncComGenesisType, nil
	case "SyncComUnitType":
		return SyncComUnitType, nil
	case "SyncComRecursiveType":
		return SyncComRecursiveType, nil
	case "TxInEth2":
		return TxInEth2, nil
	case "UnitOuter":
		return UnitOuter, nil
	case "BeaconHeaderFinalityType":
		return BeaconHeaderFinalityType, nil
	case "BeaconHeaderType":
		return BeaconHeaderType, nil
	case "BtcBulkType":
		return BtcBulkType, nil
	case "BtcPackedType":
		return BtcPackedType, nil
	case "BtcWrapType":
		return BtcWrapType, nil
	case "BtcBaseType":
		return BtcBaseType, nil
	case "BtcMiddleType":
		return BtcMiddleType, nil
	case "BtcUpperType":
		return BtcUpperType, nil
	case "BtcGenesisType":
		return BtcGenesisType, nil
	case "BtcRecursiveType":
		return BtcRecursiveType, nil
	default:
		return 0, fmt.Errorf("uKnown zk proof type %v", str)
	}
}
