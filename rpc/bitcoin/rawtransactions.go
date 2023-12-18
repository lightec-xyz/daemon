package bitcoin

import (
	"fmt"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
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

func (client *Client) CheckTx(txHash string) (bool, error) {
	var result types.RawTransaction
	err := client.Call(GETRAWTRANSACTION, &result, txHash, true)
	if err != nil {
		return false, err
	}
	return true, err
}

func (client *Client) GetRawTransaction(txHash string) (types.RawTransaction, error) {
	var result types.RawTransaction
	err := client.Call(GETRAWTRANSACTION, &result, txHash, true)
	if err != nil {
		return result, err
	}
	return result, err
}

func (client *Client) GetTransaction(txHash string) (types.RawTransaction, error) {
	var result types.RawTransaction
	err := client.Call(GETTRANSACTION, &result, txHash, true)
	if err != nil {
		return result, err
	}
	return result, err
}
func (client *Client) GetUtxoByTxId(txId string, vout int) (types.Unspents, error) {
	var result types.Unspents
	transaction, err := client.GetRawTransaction(txId)
	if err != nil {
		return result, err
	}
	for index, out := range transaction.Vout {
		if index == vout {
			result.ScriptPubKey = out.ScriptPubKey.Hex
			result.Amount = out.Value
			result.Txid = txId
			result.Vout = index
			return result, nil
		}
	}
	return result, fmt.Errorf("no find %v %v", txId, vout)
}

func outputParseParam(outputs []types.TxOut) types.Param {
	param := types.Param{}
	for _, item := range outputs {
		param.Add(item.Address, item.Amount)
	}
	return param
}
