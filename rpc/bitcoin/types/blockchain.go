package types

type BlockHeader struct {
	Hash   string `json:"hash"`
	Height int64  `json:"height"`
}

type Unspents struct {
	Amount       float64 `json:"amount"`
	Desc         string  `json:"desc"`
	Height       int     `json:"height"`
	ScriptPubKey string  `json:"scriptPubKey"`
	Txid         string  `json:"txid"`
	Vout         int     `json:"vout"`
}

type ScanUtxoSet struct {
	Bestblock   string     `json:"bestblock"`
	Height      int        `json:"height"`
	Success     bool       `json:"success"`
	TotalAmount float64    `json:"total_amount"`
	Txouts      int        `json:"txouts"`
	Unspents    []Unspents `json:"unspents"`
}

type BlockWithTx struct {
	Height int `json:"height"`
	Hash   string
	Size   int64
	Tx     []Tx
}
type Tx struct {
	Hash     string  `json:"hash"`
	Hex      string  `json:"hex"`
	Fee      float64 `json:"fee"`
	Locktime int     `json:"locktime"`
	Size     int     `json:"size"`
	Txid     string  `json:"txid"`
	Version  int     `json:"version"`
	Vin      []TxVin `json:"vin"`
	Vout     []TxOut `json:"vout"`
	Vsize    int     `json:"vsize"`
	Weight   int     `json:"weight"`
}

type TxVin struct {
	TxId        string    `json:"txid"`
	Vout        int       `json:"vout"`
	Coinbase    string    `json:"coinbase"`
	Sequence    int64     `json:"sequence"`
	Txinwitness []string  `json:"txinwitness"`
	ScriptSig   ScriptSig `json:"scriptSig"`
	Prevout     *Prevout  `json:"prevout"`
}

type Prevout struct {
	Generated    bool         `json:"generated"`
	Height       int          `json:"height"`
	Value        float64      `json:"value"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}
