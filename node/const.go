package node

// testnet

const (
	TestnetBtcOperatorAddress = ""
	TestnetEthZkBridgeAddress
	TestnetEthZkBtcAddress
)

const (
	Deposit = "deposit"
	Redeem  = "redeem"
)

const RpcRegisterName = "zkbtc"

const (
	InitBitcoinHeight  = 2540942
	InitEthereumHeight = 10127532
)

const ProofPrefix = "p"

var (
	btcCurHeightKey = []byte("btcCurHeight")
	ethCurHeightKey = []byte("ethCurHeight")
)
