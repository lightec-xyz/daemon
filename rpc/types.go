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
	Inputs  []TxIn  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`

	EthAddr string `json:"ethAddr"`
	Amount  string `json:"amount"`

	TxId  string `json:"txId"`
	PType string `json:"type"`
	Proof string `json:"proof"`
	Msg   string `json:"msg"`
}

type ProofResponse struct {
	Inputs  []TxIn  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`

	EthAddr string `json:"ethAddr"`
	Amount  string `json:"amount"`

	TxId   string `json:"txId"`
	PType  string `json:"type"`
	Proof  string `json:"proof"`
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type ProofStatus struct {
	State int    `json:"state"`
	Msg   string `json:"msg"`
}
