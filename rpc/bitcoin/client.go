package bitcoin

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
)

type Client struct {
	client *http.Client
	debug  bool
	url    string
	token  string // todo
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func NewClient(url, user, pwd string) (*Client, error) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	return &Client{client: client, url: url, token: BasicAuth(user, pwd), debug: false}, nil
}

func (c *Client) GetBlockHeader(hash string) (*types.BlockHeader, error) {
	var header = &types.BlockHeader{}
	err := c.call(GETBLOCKHEADER, NewParams(hash, true), &header)
	if err != nil {
		return nil, err
	}
	return header, err
}

func (c *Client) GetBlockHeaderByHeight(height uint64) (*types.BlockHeader, error) {
	hash, err := c.GetBlockHash(int64(height))
	if err != nil {
		return nil, err
	}
	return c.GetBlockHeader(hash)
}

func (c *Client) GetHexBlockHeader(hash string) (string, error) {
	var header string
	err := c.call(GETBLOCKHEADER, NewParams(hash, false), &header)
	if err != nil {
		return "", err
	}
	return header, err
}

func (c *Client) GetBlockHash(blockCount int64) (string, error) {
	var hash string
	err := c.call(GETBLOCKHASH, NewParams(blockCount), &hash)
	if err != nil {
		return "", err
	}
	return hash, err
}

func (c *Client) GetBlockCount() (int64, error) {
	var count int64
	err := c.call(GETBLOCKCOUNT, nil, &count)
	if err != nil {
		return 0, err
	}
	return count, err
}

func (c *Client) GetBlockByNumber(height uint64) (*types.Block, error) {
	blockHash, err := c.GetBlockHash(int64(height))
	if err != nil {
		return nil, err
	}
	block, err := c.GetBlock(blockHash)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c *Client) GetBlock(hash string) (*types.Block, error) {
	res := &types.Block{}
	err := c.call(GETBLOCK, NewParams(hash, 3), res)
	if err != nil {
		return nil, err
	}
	return res, err
}
func (c *Client) GetBlockStr(hash string) (string, error) {
	var res json.RawMessage
	err := c.call(GETBLOCK, NewParams(hash, 3), &res)
	if err != nil {
		return "", err
	}
	return string(res), err
}

func (c *Client) newRequest(method string, param Params) (*http.Request, error) {
	jsonRpc := JsonReq{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  param,
		ID:      time.Now().UnixNano(),
	}
	reqData, err := json.Marshal(jsonRpc)
	if err != nil {
		return nil, err
	}
	if c.debug {
		fmt.Printf("%v requst: data: %v \n", method, string(reqData))
	}
	request, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.token))
	return request, nil
}

func (c *Client) call(method string, param Params, result interface{}) error {
	request, err := c.newRequest(method, param)
	if err != nil {
		return err
	}
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("status code error: %d %s", response.StatusCode, response.Status)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if c.debug {
		fmt.Printf("%v rsponse: %v \n", method, string(data))
	}
	var resp JsonResp
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return fmt.Errorf("%s", string(resp.Error))
	}
	err = json.Unmarshal(resp.Result, result)
	if err != nil {
		return err
	}
	return nil
}

type JsonResp struct {
	Result json.RawMessage `json:"result"`
	Error  json.RawMessage `json:"error"`
	ID     int64           `json:"id"`
}

type JsonReq struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int64       `json:"id"`
}

type Params []interface{}

func NewParams(value ...interface{}) Params {
	var param Params
	for _, v := range value {
		param = append(param, v)
	}
	return param
}

func (p *Params) AddValue(value interface{}) {
	*p = append(*p, value)
}
