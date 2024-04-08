package common

// env

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
