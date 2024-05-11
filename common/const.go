package common

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
	ProofPending             // already not use now,placeholder compatible
	ProofSuccess
	ProofFailed
	ProofGenerating
	ProofQueueWait
)

func (ps *ProofStatus) String() string {
	switch *ps {
	case ProofDefault:
		return "ProofDefault"
	case ProofPending:
		return "ProofPending"
	case ProofSuccess:
		return "ProofSuccess"
	case ProofFailed:
		return "ProofFailed"
	case ProofGenerating:
		return "ProofGenerating"
	case ProofQueueWait:
		return "ProofQueueWait"
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
	WeightLow
	WeightMedium
	WeightHigh
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
		return ""
	}
}

func (zkpt *ZkProofType) Weight() ProofWeight {
	// todo
	switch *zkpt {
	case SyncComRecursiveType, SyncComGenesisType, RedeemTxType:
		return Highest
	case BeaconHeaderFinalityType:
		return WeightHigh
	case SyncComUnitType:
		return WeightMedium
	case VerifyTxType:
		return WeightLow
	default:
		return WeightDefault
	}
}
