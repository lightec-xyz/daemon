package sgx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

type Client struct {
	endpoint string
	timeout  time.Duration
	imp      *http.Client
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		imp:      http.DefaultClient,
		timeout:  60 * time.Second,
	}
}

func (c *Client) SgxKeyInfo() (*KeyInfo, error) {
	var result KeyInfo
	err := c.post("/zkbtcKeyInfo", nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) BtcTxSignature(currentScRoot, minerReward, txId, proof string, sigHashes []string) (*TxSignature, error) {
	var result TxSignature
	param := Param{}
	param.Add("proof", proof)
	param.Add("txId", txId)
	param.Add("currentScRoot", currentScRoot)
	param.Add("minerReward", minerReward)
	param.Add("sigHashes", sigHashes)
	err := c.post("/signBtc", param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
func (c *Client) BtcTxSignatureV1(currentScRoot, txId, proof, minerReward string, sigHashes []string) (*TxSignature, error) {
	var result TxSignature
	param := Param{}
	param.Add("proof", proof)
	param.Add("txId", txId)
	param.Add("sigHashes", sigHashes)
	param.Add("currentScRoot", currentScRoot)
	param.Add("minerReward", minerReward)
	err := c.post("/signBtc", param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) post(method string, param Param, value interface{}) error {
	return c.httpReq(http.MethodPost, method, param, value)
}

func (c *Client) newRequest(ctx context.Context, httpMethod, url, method string, param interface{}) (*http.Request, error) {
	reqData, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, httpMethod, fmt.Sprintf("%s%s", url, method), bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) httpReq(httpMethod, method string, param Param, value interface{}) (err error) {
	vi := reflect.ValueOf(value)
	if vi.Kind() != reflect.Ptr {
		return fmt.Errorf("value must be pointer")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	req, err := c.newRequest(ctx, httpMethod, c.endpoint, method, param)
	if err != nil {
		return err
	}
	resp, err := c.imp.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		return err
	}
	if resp == nil || resp.StatusCode < http.StatusOK || resp.StatusCode > 300 {
		data, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("response err: %v %v %v", resp.StatusCode, resp.Status, string(data))
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("data is empty")
	}
	jsonResp := Response{}
	err = json.Unmarshal(data, &jsonResp)
	if err != nil {
		return err
	}
	if jsonResp.Code != 200 {
		return fmt.Errorf("code is not 1: %v", jsonResp.Message)
	}
	err = json.Unmarshal(jsonResp.Data, value)
	if err != nil {
		return err
	}
	return nil

}
