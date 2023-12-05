package node

type DepositTx struct {
	TxId    string
	TxIndex int
	EthAddr string
	Amount  string
}

type RedeemTx struct {
	Inputs  []TxIn
	Outputs []TxOut
	TxIndex uint32
	TxId    string
}

type TxIn struct {
	TxId  string
	Index uint32
}

type TxOut struct {
	Value    int64
	PkScript []byte
}

type TxProof struct {
	PTxId  string      `json:"pTxId"`
	Proof  string      `json:"proof"`
	Status ProofStatus `json:"status"`
	TxId   string      `json:"txId"`
	PType  string      `json:"type"`
	ToAddr string      `json:"toAddr"`
	Amount string      `json:"amount"`
	Msg    string      `json:"msg"`
}

// todo

type ProofRequest struct {
	// redeem
	Inputs  []TxIn  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`

	// deposit
	Amount  string `json:"value"`
	EthAddr string `json:"ethAddr"`

	TxId    string `json:"txId"`
	TxIndex int    `json:"index"`
	PType   string `json:"type"`
	Proof   string `json:"proof"`
	Msg     string `json:"msg"`
}

// todo
type ProofResponse struct {
	// redeem
	Inputs  []TxIn  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`

	// deposit
	Amount  string `json:"value"`
	EthAddr string `json:"ethAddr"`

	TxId    string `json:"txId"`
	TxIndex int    `json:"index"`
	PType   string `json:"type"`
	Proof   string `json:"proof"`
	Msg     string `json:"msg"`
}
