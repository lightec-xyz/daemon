package bitcoin

import "github.com/lightec-xyz/daemon/rpc/bitcoin/types"

func (client *Client) GetNetworkInfo() (types.NetworkInfo, error) {
	var result types.NetworkInfo
	err := client.Call(GETNETWORKINFO, &result)
	if err != nil {
		return result, err
	}
	return result, err
}
