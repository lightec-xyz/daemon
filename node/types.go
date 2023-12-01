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
	TxId    string
	From    string
	BtcAddr string
	Height  string
	Amount  string
	Extra   string
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
	TxId   string `json:"txId"`
	PType  string `json:"type"`
	Proof  string `json:"proof"`
	ToAddr string `json:"toAddr"`
	Amount string `json:"amount"`
	Msg    string `json:"msg"`
}

//todo

type ProofResponse struct {
	TxId   string `json:"txId"`
	Index  uint32 `json:"index"`
	PType  string `json:"type"`
	Proof  string `json:"proof"`
	ToAddr string `json:"toAddr"`
	Amount string `json:"amount"`
	Msg    string `json:"msg"`
}
