package bitcoinClient

import (
	"bitcoinClient/types"
	"fmt"
)

func (client *Client) Createmultisig(nRequired int, keys ...string) (types.CreateMultiAddress, error) {
	var result types.CreateMultiAddress
	if nRequired > len(keys) {
		return result, fmt.Errorf("nRequired mustl less than keys len")
	}
	err := client.Call(CREATEMULTISIG, &result, nRequired, keys)
	if err != nil {
		return result, err
	}
	return result, err
}
