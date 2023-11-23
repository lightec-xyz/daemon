package node

// todo

const (
	BtcCurHeight = "btcCurHeight"
	EthCurHeight = "ethCurHeight"
)

type ProofStatus int

const (
	Default ProofStatus = iota
	Pending
	Completed
	Failed
)

const ProofPrefix = "p"
