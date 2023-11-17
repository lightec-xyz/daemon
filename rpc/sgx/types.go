package sgx

import "encoding/json"

type Response struct {
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

type KeyInfo struct {
	PublicKey string `json:"publicKey"`
}

type TxSignature struct {
	Signatures []string `json:"signatures"`
}

type Param map[string]interface{}

func (p *Param) Add(key string, value interface{}) {
	(*p)[key] = value
}
