package custom

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	btcProveUtilsClient "github.com/lightec-xyz/btc_provers/utils/client"
)

type Client struct {
	parse *Parse
}

func (c Client) GetBlock(blkHashStr string) (*btcjson.GetBlockVerboseResult, error) {
	//TODO implement me
	panic("implement me")
}

func (c Client) GetHeaderByHashStr(blkHashStr string) (string, error) {
	return c.parse.GetHeaderByHash(blkHashStr)
}

func (c Client) GetHeaderByHash(blkHash *chainhash.Hash) (string, error) {
	return c.parse.GetHeaderByHash(blkHash.String())
}

func (c Client) GetBlockHash(blockHeight int64) (*chainhash.Hash, error) {
	hash, err := c.parse.GetBlockHashByHeight(blockHeight)
	if err != nil {
		return nil, err
	}
	return chainhash.NewHashFromStr(hash)
}

func (c Client) GetHeaderByHeight(blockHeight int64) (string, error) {
	return c.parse.GetHeaderByHeight(blockHeight)
}

func NewCustomClient(path string) (btcProveUtilsClient.IClient, error) {
	parse, err := NewParse(path)
	if err != nil {
		return nil, err
	}
	return &Client{
		parse: parse,
	}, nil
}
