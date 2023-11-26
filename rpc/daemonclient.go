package rpc

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lightec-xyz/daemon/node"
)

var _ node.API = (*Client)(nil)

type Client struct {
	*rpc.Client
}

func (c *Client) Version() (node.DaemonInfo, error) {
	var info node.DaemonInfo
	err := c.Client.Call(&info, "version")
	if err != nil {
		return info, err
	}
	return info, nil
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
