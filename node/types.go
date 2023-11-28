package node

type DepositTx struct {
	TxId    string
	Addr    string
	EthAddr string
	Height  int64
	Amount  string
	Extra   string
}

type RedeemTx struct {
	TxId   string
	Addr   string
	Height string
	Amount string
	Extra  string
}

type TxProof struct {
	PTxId  string
	Proof  string
	Status int
}

// todo

type ProofRequest struct {
	TxId   string    `json:"txId"`
	PType  ProofType `json:"type"`
	Proof  string    `json:"proof"`
	ToAddr string    `json:"toAddr"`
	Amount string    `json:"amount"`
	Msg    string    `json:"msg"`
}

//todo

type ProofResponse struct {
	TxId   string    `json:"txId"`
	PType  ProofType `json:"type"`
	Proof  string    `json:"proof"`
	ToAddr string    `json:"toAddr"`
	Amount string    `json:"amount"`
	Msg    string    `json:"msg"`
}
