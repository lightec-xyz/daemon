package rpc

type Transaction struct {
	TxHash   string
	DestHash string
	Height   int64

	BtcTxId string

	Amount  int64
	EthAddr string
	Utxo    []Utxo

	Inputs  []Utxo
	Outputs []TxOut

	TxType    int
	ChainType int
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
	TxId      string `json:"txId"`
	ProofType int    `json:"type"`
	Proof     string `json:"proof"`
	Status    int    `json:"status"`
}
