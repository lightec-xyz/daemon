package bitcoin

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
	debug  bool
	urls   []string
	token  string // todo
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func NewClient(user, pwd string, urls ...string) (*Client, error) {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	return &Client{
		client: client,
		urls:   urls,
		token:  BasicAuth(user, pwd),
		debug:  false,
	}, nil
}

func (c *Client) Estimatesmartfee(confirms int) (EstimateSmartFee, error) {
	var fee EstimateSmartFee
	err := c.Call("estimatesmartfee", NewParams(confirms, "CONSERVATIVE"), &fee)
	if err != nil {
		return fee, err
	}
	return fee, err
}

func (c *Client) GetBlockHeader(hash string) (*BlockHeader, error) {
	var header = BlockHeader{}
	err := c.Call("getblockheader", NewParams(hash, true), &header)
	if err != nil {
		return nil, err
	}
	return &header, err
}

func (c *Client) GetBlockHeaderByHeight(height uint64) (*BlockHeader, error) {
	hash, err := c.GetBlockHash(int64(height))
	if err != nil {
		return nil, err
	}
	return c.GetBlockHeader(hash)
}

func (c *Client) GetHexBlockHeader(hash string) (string, error) {
	var header string
	err := c.Call("getblockheader", NewParams(hash, false), &header)
	if err != nil {
		return "", err
	}
	return header, err
}

func (c *Client) Testmempoolaccept(txRaws ...string) ([]TestMempoolAccept, error) {
	var res []TestMempoolAccept
	err := c.Call("testmempoolaccept", NewParams(txRaws), &res)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (c *Client) GetBlockHash(blockCount int64) (string, error) {
	var hash string
	err := c.Call("getblockhash", NewParams(blockCount), &hash)
	if err != nil {
		return "", err
	}
	return hash, err
}

func (c *Client) GetBlockCount() (int64, error) {
	var count int64
	err := c.Call("getblockcount", nil, &count)
	if err != nil {
		return 0, err
	}
	return count, err
}

func (c *Client) GetBlockByNumber(height uint64) (Block, error) {
	blockHash, err := c.GetBlockHash(int64(height))
	if err != nil {
		return Block{}, err
	}
	block, err := c.GetBlock(blockHash)
	if err != nil {
		return Block{}, err
	}
	return block, nil
}

func (c *Client) GetBlock(hash string) (Block, error) {
	res := Block{}
	err := c.Call("getblock", NewParams(hash, 3), &res)
	if err != nil {
		return Block{}, err
	}
	return res, err
}
func (c *Client) GetBlockStr(hash string) (string, error) {
	var res json.RawMessage
	err := c.Call("getblock", NewParams(hash, 3), &res)
	if err != nil {
		return "", err
	}
	return string(res), err
}

func (c *Client) Scantxoutset(address string) (ScanUtxoSet, error) {
	var result ScanUtxoSet
	err := c.Call("scantxoutset", NewParams("start", []string{fmt.Sprintf("addr(%v)", address)}), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) GetNetworkInfo() (NetworkInfo, error) {
	var result NetworkInfo
	err := c.Call("getnetworkinfo", nil, &result)
	if err != nil {
		return result, err
	}
	return result, err
}
func (c *Client) Createrawtransaction(inputs []TxIn, outputs []TxOut) (string, error) {
	var result string
	err := c.Call("createrawtransaction", NewParams(inputs, outputParseParam(outputs)), &result)
	if err != nil {
		return "", err
	}
	return result, err
}

func (c *Client) Signrawtransactionwithkey(hexDaa string, privateKeys []string, inputs []TxIn) (SignTawTransaction, error) {
	var result SignTawTransaction
	err := c.Call("signrawtransactionwithkey", NewParams(hexDaa, privateKeys, inputs), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) Sendrawtransaction(hexData string) (string, error) {
	var result string
	err := c.Call("sendrawtransaction", NewParams(hexData, 0), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) CheckTxOnChain(txHash string) (bool, error) {
	_, err := c.Getmempoolentry(txHash)
	if err == nil {
		return true, nil
	}
	_, err = c.GetRawTransaction(txHash)
	if err == nil {
		return true, nil
	}
	return false, nil
}

func (c *Client) GetRawTransaction(txHash string) (RawTransaction, error) {
	var result RawTransaction
	err := c.Call("getrawtransaction", NewParams(txHash, true), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) GetHexRawTransaction(txHash string) (string, error) {
	var result string
	err := c.Call("getrawtransaction", NewParams(txHash, false), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

//getmempoolentry

func (c *Client) Getmempoolentry(txId string) (RawTransaction, error) {
	var result RawTransaction
	err := c.Call("getmempoolentry", NewParams(txId), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) GetTransaction(txHash string) (RawTransaction, error) {
	var result RawTransaction
	err := c.Call("gettransaction", NewParams(txHash, true), &result)
	if err != nil {
		return result, err
	}
	return result, err
}
func (c *Client) GetUtxoByTxId(txId string, vout int) (Unspents, error) {
	var result Unspents
	transaction, err := c.GetRawTransaction(txId)
	if err != nil {
		return result, err
	}
	for index, out := range transaction.Vout {
		if index == vout {
			result.ScriptPubKey = out.ScriptPubKey.Hex
			result.Amount = out.Value
			result.Txid = txId
			result.Vout = index
			return result, nil
		}
	}
	return result, fmt.Errorf("no find %v %v", txId, vout)
}

func outputParseParam(outputs []TxOut) Param {
	param := Param{}
	for _, item := range outputs {
		param.Add(item.Address, item.Amount)
	}
	return param
}
func (c *Client) Createmultisig(nRequired int, keys ...string) (CreateMultiAddress, error) {
	var result CreateMultiAddress
	if nRequired > len(keys) {
		return result, fmt.Errorf("nRequired mustl less than keys len")
	}
	err := c.Call("createmultisig", NewParams(nRequired, keys), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) Getaddressinfo(address string) (AddressInfo, error) {
	var result AddressInfo
	err := c.Call("getaddressinfo", NewParams(address), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) DumpPrivkey(address string) (string, error) {
	var result string
	err := c.Call("dumpprivkey", NewParams(address), &result)
	if err != nil {
		return "", err
	}
	return result, err
}

func (c *Client) GetRawChangeAddress(param ...AddrType) (string, error) {
	var result string
	addrType := BECH32
	if len(param) != 0 {
		addrType = param[0]
	}
	err := c.Call("getrawchangeaddress", NewParams(addrType), &result)
	if err != nil {
		return "", err
	}
	return result, err
}

func (c *Client) newRequest(url, method string, param Params) (*http.Request, error) {
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
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.token))
	return request, nil
}

func (c *Client) call(url, method string, param Params, result interface{}) error {
	request, err := c.newRequest(url, method, param)
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
	data, err := io.ReadAll(response.Body)
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
		if string(resp.Error) != "null" {
			return fmt.Errorf("%s", string(resp.Error))
		}
	}
	err = json.Unmarshal(resp.Result, result)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Call(method string, param Params, result interface{}) error {
	msg := "call error " + method
	for _, url := range c.urls {
		err := c.call(url, method, param, result)
		if err != nil {
			msg = msg + err.Error()
			continue
		}
		return nil
	}
	return errors.New(msg)
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
