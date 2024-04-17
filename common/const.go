package common

// env

const BeaconHeaderSlot = 32 // todo

const SlotPerPeriod = 8192

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
	ProofPending
	ProofSuccess
	ProofFailed
)

type Mode string

const (
	Client  Mode = "client"
	Cluster Mode = "cluster"
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
	BlockHeaderFinalityType //BeaconHeaderFinalityUpdate
	UnitOuter
	BlockHeaderType
)

func (zkpr *ZkProofType) String() string {
	switch *zkpr {
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
	case BlockHeaderFinalityType:
		return "BlockHeaderFinalityType"
	case BlockHeaderType:
		return "BlockHeaderType"
	default:
		return ""
	}
}
