package node

const (
	ProofPrefix         = "p_" // height + p_ + hash
	TxPrefix            = "t_" // height + t_ + hash
	DestChainHashPrefix = "d_" // height + d_ + hash
)

var (
	btcCurHeightKey = []byte("btcCurHeight")
	ethCurHeightKey = []byte("ethCurHeight")
)
