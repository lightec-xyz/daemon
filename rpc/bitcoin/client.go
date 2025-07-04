package bitcoin

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
	debug  bool
	url    string
	token  string // todo
	local  bool
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func NewClient(url, user, pwd, token string) (*Client, error) {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	local := true
	if token != "" {
		local = false
	} else {
		token = BasicAuth(user, pwd)
		local = true
	}
	return &Client{
		client: client,
		url:    url,
		token:  token,
		debug:  false,
		local:  local,
	}, nil
}

func (c *Client) Estimatesmartfee(confirms int) (EstimateSmartFee, error) {
	var fee EstimateSmartFee
	err := c.call("estimatesmartfee", NewParams(confirms, "CONSERVATIVE"), &fee)
	if err != nil {
		return fee, err
	}
	return fee, err
}

func (c *Client) GetBlockHeader(hash string) (*BlockHeader, error) {
	var header = BlockHeader{}
	err := c.call("getblockheader", NewParams(hash, true), &header)
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
	err := c.call("getblockheader", NewParams(hash, false), &header)
	if err != nil {
		return "", err
	}
	return header, err
}

func (c *Client) GetBlockHash(blockCount int64) (string, error) {
	var hash string
	err := c.call("getblockhash", NewParams(blockCount), &hash)
	if err != nil {
		return "", err
	}
	return hash, err
}

func (c *Client) GetBlockCount() (int64, error) {
	var count int64
	err := c.call("getblockcount", nil, &count)
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
	err := c.call("getblock", NewParams(hash, 3), &res)
	if err != nil {
		return Block{}, err
	}
	return res, err
}
func (c *Client) GetBlockStr(hash string) (string, error) {
	var res json.RawMessage
	err := c.call("getblock", NewParams(hash, 3), &res)
	if err != nil {
		return "", err
	}
	return string(res), err
}

func (c *Client) Scantxoutset(address string) (ScanUtxoSet, error) {
	var result ScanUtxoSet
	err := c.call("scantxoutset", NewParams("start", []string{fmt.Sprintf("addr(%v)", address)}), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) GetNetworkInfo() (NetworkInfo, error) {
	var result NetworkInfo
	err := c.call("getnetworkinfo", nil, &result)
	if err != nil {
		return result, err
	}
	return result, err
}
func (c *Client) Createrawtransaction(inputs []TxIn, outputs []TxOut) (string, error) {
	var result string
	err := c.call("createrawtransaction", NewParams(inputs, outputParseParam(outputs)), &result)
	if err != nil {
		return "", err
	}
	return result, err
}

func (c *Client) Signrawtransactionwithkey(hexDaa string, privateKeys []string, inputs []TxIn) (SignTawTransaction, error) {
	var result SignTawTransaction
	err := c.call("signrawtransactionwithkey", NewParams(hexDaa, privateKeys, inputs), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) Sendrawtransaction(hexData string) (string, error) {
	var result string
	err := c.call("sendrawtransaction", NewParams(hexData, 0), &result)
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
	err := c.call("getrawtransaction", NewParams(txHash, true), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) GetHexRawTransaction(txHash string) (string, error) {
	var result string
	err := c.call("getrawtransaction", NewParams(txHash, false), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

//getmempoolentry

func (c *Client) Getmempoolentry(txId string) (RawTransaction, error) {
	var result RawTransaction
	err := c.call("getmempoolentry", NewParams(txId), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) GetTransaction(txHash string) (RawTransaction, error) {
	var result RawTransaction
	err := c.call("gettransaction", NewParams(txHash, true), &result)
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
	err := c.call("createmultisig", NewParams(nRequired, keys), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) Getaddressinfo(address string) (AddressInfo, error) {
	var result AddressInfo
	err := c.call("getaddressinfo", NewParams(address), &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) DumpPrivkey(address string) (string, error) {
	var result string
	err := c.call("dumpprivkey", NewParams(address), &result)
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
	err := c.call("getrawchangeaddress", NewParams(addrType), &result)
	if err != nil {
		return "", err
	}
	return result, err
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
	if c.local {
		request.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.token))
	} else {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
		request.Header.Set("Content-Type", "application/json")
	}
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
