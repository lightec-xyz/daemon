package node

const (
	ProofPrefix         = "p_" // p_ + hash
	TxPrefix            = "t_" // t_ + hash
	DestChainHashPrefix = "d_" // d_ + hash
	UnGenProofPrefix    = "u_" // u_ + hash

)

var (
	btcCurHeightKey = []byte("btcCurHeight")
	ethCurHeightKey = []byte("ethCurHeight")
)
