package ethereum

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/onrik/ethrpc"
)

type Client struct {
	*rpc.Client
	*ethrpc.EthRPC
}

func NewClient(url string) (*Client, error) {
	ethRPC := ethrpc.New(url)
	client, err := rpc.DialHTTP(url)
	if err != nil {
		return nil, err
	}
	return &Client{
		client,
		ethRPC,
	}, nil
}
