package bitcoin

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	"github.com/ybbus/jsonrpc"
	"net/http"
	"time"
)

type Client struct {
	rpcClient jsonrpc.RPCClient
	debug     bool
}

func (client *Client) GetBlockHeader(hash string) (*types.BlockHeader, error) {
	var header = &types.BlockHeader{}
	err := client.Call(GETBLOCKHEADER, &header, hash)
	if err != nil {
		return nil, err
	}
	return header, err
}

func (client *Client) GetBlockHash(blockCount int64) (string, error) {
	var hash string
	err := client.Call(GETBLOCKHASH, &hash, blockCount)
	if err != nil {
		return "", err
	}
	return hash, err
}

func (client *Client) GetBlockCount() (int64, error) {
	var count int64
	err := client.Call(GETBLOCKCOUNT, &count)
	if err != nil {
		return 0, err
	}
	return count, err
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func NewClient(url, user, pwd, network string) (*Client, error) {
	opts := &jsonrpc.RPCClientOpts{
		HTTPClient:    http.DefaultClient,
		CustomHeaders: make(map[string]string),
	}
	opts.CustomHeaders["Authorization"] = "Basic " + basicAuth(user, pwd)
	opts.CustomHeaders["Connection"] = "close"
	opts.HTTPClient.Timeout = 2 * time.Minute
	return &Client{rpcClient: jsonrpc.NewClientWithOpts(url, opts), debug: true}, nil
}

func (client *Client) Call(method string, result interface{}, args ...interface{}) error {
	//todo

	if client.debug {
		var buff bytes.Buffer
		buff.WriteString(fmt.Sprintf("jsonrpc req  : %s [", method))
		//buff.WriteString(fmt.Sprintf("\tresult: %v", reflect.TypeOf(result)))
		for i, arg := range args {
			data, _ := json.Marshal(arg)
			if i == len(args)-1 {
				buff.WriteString(fmt.Sprintf("%v", string(data)))
			} else {
				buff.WriteString(fmt.Sprintf("%v,", string(data)))
			}
		}
		buff.WriteString("]\n")
		fmt.Printf(buff.String())
	}
	response, err := client.rpcClient.Call(method, args...)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return fmt.Errorf("rpc error : %d %s %v", response.Error.Code, response.Error.Message, response.Error.Data)
	}
	if client.debug {
		var buff bytes.Buffer
		responseData, _ := json.Marshal(response)
		buff.WriteString(fmt.Sprintf("jsonrpc response: %s", method))
		buff.WriteString(fmt.Sprintf("\tresult: %s\n", responseData))
		fmt.Println(buff.String())
	}
	return response.GetObject(result)
}
