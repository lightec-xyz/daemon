package bitcoinClient

import (
	"bitcoinClient/types"
)

func (client *Client) Createrawtransaction(inputs []types.TxIn, outputs []types.TxOut) (string, error) {
	var result string
	err := client.Call(CREATERAWTRANSACTION, &result, inputs, outputParseParam(outputs))
	if err != nil {
		return "", err
	}
	return result, err
}

func (client *Client) Signrawtransactionwithkey(hexDaa string, privateKeys []string, inputs []types.TxIn) (types.SignTawTransaction, error) {
	var result types.SignTawTransaction
	err := client.Call(SIGNRAWTRANSACTIONWITHKEY, &result, hexDaa, privateKeys, inputs)
	if err != nil {
		return result, err
	}
	return result, err
}

func (client *Client) Sendrawtransaction(hexData string) (string, error) {
	var result string
	err := client.Call(SENDRAWTRANSACTION, &result, hexData)
	if err != nil {
		return result, err
	}
	return result, err
}

func (client *Client) GetRawtransaction(txHash string) (types.RawTransaction, error) {
	var result types.RawTransaction
	err := client.Call(GETRAWTRANSACTION, &result, txHash, true)
	if err != nil {
		return result, err
	}
	return result, err
}

func outputParseParam(outputs []types.TxOut) types.Param {
	param := types.Param{}
	for _, item := range outputs {
		param.Add(item.Address, item.Amount)
	}
	return param
}
