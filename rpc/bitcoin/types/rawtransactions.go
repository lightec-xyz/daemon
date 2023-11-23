package types

type RawTransaction struct {
	Blockhash     string `json:"blockhash"`
	Blocktime     int    `json:"blocktime"`
	Confirmations int    `json:"confirmations"`
	Hash          string `json:"hash"`
	Hex           string `json:"hex"`
	Locktime      int    `json:"locktime"`
	Size          int    `json:"size"`
	Time          int    `json:"time"`
	Txid          string `json:"txid"`
	Version       int    `json:"version"`
	Vin           []struct {
		ScriptSig   ScriptSig `json:"scriptSig"`
		Sequence    int64     `json:"sequence"`
		Txid        string    `json:"txid"`
		Txinwitness []string  `json:"txinwitness"`
		Vout        int       `json:"vout"`
	} `json:"vin"`
	Vout []struct {
		N            int          `json:"n"`
		ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
		Value        float64      `json:"value"`
	} `json:"vout"`
	Vsize  int `json:"vsize"`
	Weight int `json:"weight"`
}
type ScriptPubKey struct {
	Address string `json:"address"`
	Asm     string `json:"asm"`
	Desc    string `json:"desc"`
	Hex     string `json:"hex"`
	Type    string `json:"type"`
}
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type SignTawTransaction struct {
	Complete bool   `json:"complete"`
	Hex      string `json:"hex"`
}

type TxIn struct {
	TxId     string `json:"txid"`
	Vout     int    `json:"vout"`
	Sequence int    `json:"sequence"`

	// extra
	ScriptPubKey  string  `json:"scriptPubKey"`
	RedeemScript  string  `json:"redeemScript"`
	WitnessScript string  `json:"witnessScript"`
	Amount        float64 `json:"amount"`
}

type TxOut struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}
