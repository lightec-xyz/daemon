package node

type DepositTx struct {
	TxId    string
	BtcTo   string
	EthAddr string
	Height  int64
	Amount  float64
	Extra   string
}

type TxProof struct {
	PTxId  string
	Proof  string
	Status int
}
