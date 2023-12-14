package rpc

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"reflect"
	"time"
)

var _ NodeAPI = (*NodeClient)(nil)

type NodeClient struct {
	*rpc.Client
	timeout time.Duration
}

func (c *NodeClient) Stop() error {
	var result string
	err := c.call(&result, "zkbtc_stop")
	if err != nil {
		return err
	}
	return err
}

func (c *NodeClient) AddWorker(endpoint string, max int) (string, error) {
	var result string
	err := c.call(&result, "zkbtc_addWorker", endpoint, max)
	if err != nil {
		return "", err
	}
	return result, err
}

func (c *NodeClient) Version() (NodeInfo, error) {
	info := NodeInfo{}
	err := c.call(&info, "zkbtc_version")
	if err != nil {
		return info, err
	}
	return info, err

}

func NewNodeClient(url string) (*NodeClient, error) {
	client, err := rpc.DialHTTP(url)
	if err != nil {
		return nil, err
	}
	return &NodeClient{
		Client:  client,
		timeout: 15 * time.Second,
	}, nil
}

func (c *NodeClient) call(result interface{}, method string, args ...interface{}) error {
	vi := reflect.ValueOf(result)
	if vi.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be pointer")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	return c.CallContext(ctx, result, method, args...)
}
