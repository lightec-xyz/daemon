package types

type CreateMultiAddress struct {
	Address      string `json:"address"`
	Descriptor   string `json:"descriptor"`
	RedeemScript string `json:"redeemScript"`
}
