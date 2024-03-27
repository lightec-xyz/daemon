package common

// todo

const ZkProofLength = 928

type ZkProof []byte

type ZkProofType int

const (
	DepositTxType ZkProofType = iota + 1
	RedeemTxType
	VerifyTxType
	SyncComGenesisType
	SyncComUnitType
	SyncComRecursiveType
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
	default:
		return ""
	}
}

type CircuitsFP struct {
	RecursiveFp []byte
}
