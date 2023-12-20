package node

const (
	ProofPrefix      = "p_"
	DestTxHashPrefix = "d_"
	TxPrefix         = "t_" // height + t_ + hash
	NoncePrefix      = "n_"
)

var (
	btcCurHeightKey = []byte("btcCurHeight")
	ethCurHeightKey = []byte("ethCurHeight")
)
