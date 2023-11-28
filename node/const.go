package node

// todo

const (
	BitcoinChain  = "bitcoin"
	EthereumChain = "ethereum"
)

const (
	BtcCurHeight = "btcCurHeight"
	EthCurHeight = "ethCurHeight"
)

const (
	InitBitcoinHeight  = 100
	InitEthereumHeight = 10000
)

type ProofStatus int

const (
	Default ProofStatus = iota
	Pending
	Completed
	Failed
)

const ProofPrefix = "p"
