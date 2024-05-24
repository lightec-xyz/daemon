package common

import "fmt"

// env

const (
	BeaconHeaderSlot = 32
	SlotPerPeriod    = 8192
)

const (
	ZkDebugEnv     = "ZkDebug"
	ZkParameterDir = "ZkParameterDir"
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

var DefaultProofTypes = []ZkProofType{
	DepositTxType,
	RedeemTxType,
	VerifyTxType,
	SyncComGenesisType,
	SyncComUnitType,
	SyncComRecursiveType,
	BeaconHeaderFinalityType,
	BeaconHeaderType,
	TxInEth2, // todo
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
	default:
		return 0, fmt.Errorf("uKnown zk proof type %v", str)
	}
}
