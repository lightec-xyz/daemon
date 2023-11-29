package node

// todo

const (
	BitcoinChain  = "bitcoin"
	EthereumChain = "ethereum"
)

const (
	InitBitcoinHeight  = 100
	InitEthereumHeight = 10000
)

const ProofPrefix = "p"

var (
	btcCurHeightKey = []byte("btcCurHeight")
	ethCurHeightKey = []byte("ethCurHeight")
)

type ProofStatus int

const (
	ProofDefault ProofStatus = iota
	ProofPending
	ProofSuccess
	ProofFailed
)
