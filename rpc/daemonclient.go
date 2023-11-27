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
}

func (c *NodeClient) Version() (*DaemonInfo, error) {
	info := DaemonInfo{}
	err := c.call(info, "version")
	if err != nil {
		return nil, err
	}
	return nil, err

}

func NewClient(url string) (*NodeClient, error) {
	client, err := rpc.DialHTTP(url)
	if err != nil {
		return nil, err
	}
	return &NodeClient{
		client,
	}, nil
}

func (c *NodeClient) call(result interface{}, method string, args ...interface{}) error {
	vi := reflect.ValueOf(result)
	if vi.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be pointer")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc()
	return c.CallContext(ctx, result, method, args)
}
