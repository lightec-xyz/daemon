package bitcoin

import "github.com/lightec-xyz/daemon/rpc/bitcoin/types"

func (c *Client) GetNetworkInfo() (types.NetworkInfo, error) {
	var result types.NetworkInfo
	err := c.call(GETNETWORKINFO, &result)
	if err != nil {
		return result, err
	}
	return result, err
}
