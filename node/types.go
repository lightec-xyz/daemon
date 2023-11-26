package node

type DepositTx struct {
	TxId    string
	BtcTo   string
	EthAddr string
	Height  int64
	Amount  float64
	Extra   string
}

type RedeemTx struct {
	TxId    string
	BtcAddr string
	Height  string
	Amount  string
	Extra   string
}

type TxProof struct {
	PTxId  string
	Proof  string
	Status int
}

type DaemonInfo struct {
	Version string
	Desc    string
}
