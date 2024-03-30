package bitcoin

import (
	"fmt"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"strings"
)

func (c *Client) Createrawtransaction(inputs []types.TxIn, outputs []types.TxOut) (string, error) {
	var result string
	err := c.call(CREATERAWTRANSACTION, &result, inputs, outputParseParam(outputs))
	if err != nil {
		return "", err
	}
	return result, err
}

func (c *Client) Signrawtransactionwithkey(hexDaa string, privateKeys []string, inputs []types.TxIn) (types.SignTawTransaction, error) {
	var result types.SignTawTransaction
	err := c.call(SIGNRAWTRANSACTIONWITHKEY, &result, hexDaa, privateKeys, inputs)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) Sendrawtransaction(hexData string) (string, error) {
	var result string
	err := c.call(SENDRAWTRANSACTION, &result, hexData)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) CheckTx(txHash string) (bool, error) {
	//todo
	txId := strings.TrimPrefix(txHash, "0x")
	var result types.RawTransaction
	err := c.call(GETTRANSACTION, &result, txId, true)
	if err != nil {
		return false, nil
	}
	return true, err
}

func (c *Client) GetRawTransaction(txHash string) (types.RawTransaction, error) {
	var result types.RawTransaction
	err := c.call(GETRAWTRANSACTION, &result, txHash, true)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) GetTransaction(txHash string) (types.RawTransaction, error) {
	var result types.RawTransaction
	err := c.call(GETTRANSACTION, &result, txHash, true)
	if err != nil {
		return result, err
	}
	return result, err
}
func (c *Client) GetUtxoByTxId(txId string, vout int) (types.Unspents, error) {
	var result types.Unspents
	transaction, err := c.GetRawTransaction(txId)
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
