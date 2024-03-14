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

type DbTx struct {
	TxHash    string
	Height    int64
	TxType    TxType
	ChainType ChainType

	DestHash string
	BtcTxId  string
	Amount   int64
	EthAddr  string
	Utxo     []Utxo
	Inputs   []Utxo
	Outputs  []TxOut
}

type DbProof struct {
	TxId      string      `json:"txId"`
	ProofType ZkProofType `json:"type"`
	Proof     string      `json:"Proof"`
	Status    ProofStatus `json:"status"`
}

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
