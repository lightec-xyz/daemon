package bitcoin

import "github.com/lightec-xyz/daemon/rpc/bitcoin/types"

func (client *Client) Getaddressinfo(address string) (types.AddressInfo, error) {
	var result types.AddressInfo
	err := client.Call(GETADDRESSINFO, &result, address)
	if err != nil {
		return result, err
	}
	return result, err
}

func (client *Client) DumpPrivkey(address string) (string, error) {
	var result string
	err := client.Call(DUMPPRIVKEY, &result, address)
	if err != nil {
		return "", err
	}
	return result, err
}

func (client *Client) GetRawChangeAddress(param ...AddrType) (string, error) {
	var result string
	addrType := BECH32
	if len(param) != 0 {
		addrType = param[0]
	}
	err := client.Call(GETRAWCHANGEADDRESS, &result, addrType)
	if err != nil {
		return "", err
	}
	return result, err
}
