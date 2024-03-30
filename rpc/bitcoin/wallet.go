package bitcoin

import "github.com/lightec-xyz/daemon/rpc/bitcoin/types"

func (c *Client) Getaddressinfo(address string) (types.AddressInfo, error) {
	var result types.AddressInfo
	err := c.call(GETADDRESSINFO, &result, address)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) DumpPrivkey(address string) (string, error) {
	var result string
	err := c.call(DUMPPRIVKEY, &result, address)
	if err != nil {
		return "", err
	}
	return result, err
}

func (c *Client) GetRawChangeAddress(param ...AddrType) (string, error) {
	var result string
	addrType := BECH32
	if len(param) != 0 {
		addrType = param[0]
	}
	err := c.call(GETRAWCHANGEADDRESS, &result, addrType)
	if err != nil {
		return "", err
	}
	return result, err
}
