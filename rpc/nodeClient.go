package rpc

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lightec-xyz/daemon/common"
	"reflect"
	"time"
)

var _ INode = (*NodeClient)(nil)

type NodeClient struct {
	*rpc.Client
	timeout time.Duration
}

func (c *NodeClient) ProofTask(id string) (*ProofTaskInfo, error) {
	var result ProofTaskInfo
	err := c.call(&result, "zkbtc_proofTask")
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *NodeClient) PendingTask() ([]*ProofTaskInfo, error) {
	var result []*ProofTaskInfo
	err := c.call(&result, "zkbtc_pendingTask")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *NodeClient) TxesByAddr(addr, txType string) ([]*Transaction, error) {
	var result []*Transaction
	err := c.call(&result, "zkbtc_txesByAddr", addr, txType)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *NodeClient) GetZkProofTask(request common.TaskRequest) (*common.TaskResponse, error) {
	var result common.TaskResponse
	err := c.call(&result, "zkbtc_getZkProofTask", request)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *NodeClient) SubmitProof(req *common.SubmitProof) (string, error) {
	var result string
	err := c.call(&result, "zkbtc_submitProof", req)
	if err != nil {
		return "", err
	}
	return result, nil

}

func (c *NodeClient) TransactionsByHeight(height uint64, network string) ([]string, error) {
	var result []string
	err := c.call(&result, "zkbtc_transactionsByHeight", height, network)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (c *NodeClient) Transactions(txIds []string) ([]*Transaction, error) {
	var result []*Transaction
	err := c.call(&result, "zkbtc_transactions", txIds)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (c *NodeClient) Transaction(txHash string) (*Transaction, error) {
	var result Transaction
	err := c.call(&result, "zkbtc_transaction", txHash)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *NodeClient) ProofInfo(proofId []string) ([]ProofInfo, error) {
	var result []ProofInfo
	err := c.call(&result, "zkbtc_proofInfo", proofId)
	if err != nil {
		return result, err
	}
	return result, nil
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
