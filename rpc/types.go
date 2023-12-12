package rpc

type NodeInfo struct {
	Version string
	Desc    string
}
type TxIn struct {
	TxId  string
	Index uint32
}

type TxOut struct {
	Value    int64
	PkScript []byte
}

type ProofRequest struct {
	// redeem
	Inputs  []TxIn  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`

	// deposit
	Amount  string `json:"amount"`
	EthAddr string `json:"ethAddr"`
	Vout    int    `json:"index"`

	TxId  string `json:"txId"`
	PType string `json:"type"`
	Proof string `json:"proof"`
	Msg   string `json:"msg"`
}

type ProofResponse struct { // redeem
	// redeem
	Inputs  []TxIn  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`

	// deposit
	Amount  string `json:"amount"`
	EthAddr string `json:"ethAddr"`
	Vout    int    `json:"index"`

	TxId  string `json:"txId"`
	PType string `json:"type"`
	Proof string `json:"proof"`
	Msg   string `json:"msg"`
}

type ProofInfo struct {
	Status int    `json:"state"`
	Msg    string `json:"msg"`
}
