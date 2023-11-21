package rpc

import (
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	*rpc.Client
}

func NewClient(url string) (*Client, error) {
	client, err := rpc.DialHTTP(url)
	if err != nil {
		return nil, err
	}
	return &Client{
		client,
	}, nil
}
