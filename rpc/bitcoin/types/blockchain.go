package types

type BlockHeader struct {
	Hash          string `json:"hash"`
	Confirmations int    `json:"confirmations"`
	Height        int    `json:"height"`
	Version       int    `json:"version"`
	VersionHex    string `json:"versionHex"`
	Merkleroot    string `json:"merkleroot"`
	Time          int    `json:"time"`
	//Mediantime    int    `json:"mediantime"`
	//Nonce         int    `json:"nonce"`
	//Bits          string `json:"bits"`
	//Difficulty        string `json:"difficulty"`
	Chainwork string `json:"chainwork"`
	//NTx               int    `json:"nTx"`
	Previousblockhash string `json:"previousblockhash"`
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

type Block struct {
	Bits              string  `json:"bits"`
	Chainwork         string  `json:"chainwork"`
	Confirmations     int     `json:"confirmations"`
	Difficulty        float64 `json:"difficulty"`
	Hash              string  `json:"hash"`
	Height            int     `json:"height"`
	Mediantime        int     `json:"mediantime"`
	Merkleroot        string  `json:"merkleroot"`
	NTx               int     `json:"nTx"`
	Nextblockhash     string  `json:"nextblockhash"`
	Nonce             int     `json:"nonce"`
	Previousblockhash string  `json:"previousblockhash"`
	Size              int     `json:"size"`
	Strippedsize      int     `json:"strippedsize"`
	Time              int     `json:"time"`
	Tx                []Tx    `json:"tx"`
	Version           int     `json:"version"`
	VersionHex        string  `json:"versionHex"`
	Weight            int     `json:"weight"`
}

type Tx struct {
	Hash     string   `json:"hash"`
	Hex      string   `json:"hex"`
	Locktime int      `json:"locktime"`
	Size     int      `json:"size"`
	Txid     string   `json:"txid"`
	Version  int      `json:"version"`
	Vin      []TxVin  `json:"vin"`
	Vout     []TxVout `json:"vout"`
	Vsize    int      `json:"vsize"`
	Weight   int      `json:"weight"`
	Fee      float64  `json:"fee,omitempty"`
}

type TxVin struct {
	TxId        string    `json:"txid"`
	Vout        int       `json:"vout"`
	ScriptSig   ScriptSig `json:"scriptSig"`
	Prevout     Prevout   `json:"prevout"`
	Coinbase    string    `json:"coinbase"`
	Sequence    int64     `json:"sequence"`
	Txinwitness []string  `json:"txinwitness"`
}

type TxVout struct {
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
	Value        float64      `json:"value"`
}

type ScriptPubKey struct {
	Address string `json:"address"`
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Type    string `json:"type"`
}

type Prevout struct {
	Generated    bool         `json:"generated"`
	Height       int          `json:"height"`
	Value        float64      `json:"value"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}
