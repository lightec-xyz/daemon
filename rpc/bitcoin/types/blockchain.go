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
