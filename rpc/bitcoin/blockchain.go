package bitcoinClient

import (
	"bitcoinClient/types"
	"fmt"
)

func (client *Client) Scantxoutset(address string) (types.ScanUtxoSet, error) {
	var result types.ScanUtxoSet
	err := client.Call(SCANTXOUTSET, &result, "start", []string{fmt.Sprintf("addr(%v)", address)})
	if err != nil {
		return result, err
	}
	return result, err
}
