package node

// todo

const (
	Deposit = "deposit"
	Redeem  = "redeem"
)

const (
	InitBitcoinHeight  = 2540942
	InitEthereumHeight = 10127532
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
