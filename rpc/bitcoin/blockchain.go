package bitcoin

import (
	"fmt"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
)

func (c *Client) Scantxoutset(address string) (types.ScanUtxoSet, error) {
	var result types.ScanUtxoSet
	err := c.call(SCANTXOUTSET, NewParams("start", []string{fmt.Sprintf("addr(%v)", address)}), &result)
	if err != nil {
		return result, err
	}
	return result, err
}
