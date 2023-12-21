package rpc

type Transaction struct {
	EthereumTx EthereumTx `json:"ethereumTx"`
	Bitcoin    BitcoinTx  `json:"bitcoinTx"`
}

type EthereumTx struct {
	Hash string `json:"hash"`
}

type BitcoinTx struct {
	Hash string `json:"hash"`
}

type NodeInfo struct {
	Version string
	Desc    string
}
type Utxo struct {
	TxId  string `json:"txId"`
	Index uint32 `json:"index"`
}

type TxOut struct {
	Value    int64
	PkScript []byte
}

type ProofRequest struct {
	// redeem
	Inputs  []Utxo  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`
	BtcTxId string  `json:"btcTxId"`

	// deposit
	Utxos   []Utxo
	Amount  int64  `json:"amount"`
	EthAddr string `json:"ethAddr"`

	// other
	Height    int64  `json:"height"`
	BlockHash string `json:"blockHash"`
	TxId      string `json:"txId"`
	ProofType int    `json:"type"` // todo
	Proof     string `json:"proof"`
	Msg       string `json:"msg"`
}

type ProofResponse struct { // redeem
	// redeem
	Inputs  []Utxo  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`
	BtcTxId string  `json:"btcTxId"`

	// deposit
	Utxos   []Utxo
	Amount  int64  `json:"amount"`
	EthAddr string `json:"ethAddr"`

	// other
	Height    int64  `json:"height"`
	BlockHash string `json:"blockHash"`
	TxId      string `json:"txId"`
	ProofType int    `json:"type"` // todo
	Proof     string `json:"proof"`
	Msg       string `json:"msg"`
}

type ProofInfo struct {
	Status int    `json:"state"`
	Proof  string `json:"proof"`
	Msg    string `json:"msg"`
}
