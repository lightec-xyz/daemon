package beacon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/prysmaticlabs/prysm/v5/container/slice"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
)

type Client struct {
	ctx        context.Context
	endpoint   string
	timeout    time.Duration
	debug      bool
	httpClient *http.Client
}

func NewClient(rawurl string) (*Client, error) {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     2 * time.Hour,
			DisableKeepAlives:   false,
			MaxIdleConnsPerHost: 10,
		},
	}
	return &Client{
		ctx:        context.Background(),
		endpoint:   rawurl,
		timeout:    2 * time.Hour,
		debug:      false,
		httpClient: client,
	}, nil
}

func (c *Client) Bootstrap(slot uint64) (*structs.LightClientBootstrapResponse, error) {
	beaconHeaders, err := c.GetBeaconHeaders(slot)
	if err != nil {
		return nil, err
	}
	bootstrap, err := c.GetLightClientBootstrap(beaconHeaders.Data.Root)
	if err != nil {
		return nil, err
	}
	return bootstrap, nil
}

func (c *Client) GetLightClientBootstrap(root string) (*structs.LightClientBootstrapResponse, error) {
	bootstrap := &structs.LightClientBootstrapResponse{}
	path := fmt.Sprintf("/eth/v1/beacon/light_client/bootstrap/%v", root)
	err := c.get(path, nil, &bootstrap)
	if err != nil {
		return nil, err
	}
	return bootstrap, nil
}

func (c *Client) GetBeaconHeaders(slot uint64) (*structs.GetBlockHeaderResponse, error) {
	resp := &structs.GetBlockHeaderResponse{}
	path := fmt.Sprintf("/eth/v1/beacon/headers/%v", slot)
	err := c.get(path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetLatestFinalizedSlot() (uint64, error) {
	resp := &structs.GetBlockHeaderResponse{}
	err := c.get("/eth/v1/beacon/headers/finalized", nil, &resp)
	if err != nil {
		return 0, err
	}
	slot, ok := big.NewInt(0).SetString(resp.Data.Header.Message.Slot, 10)
	if !ok {
		return 0, fmt.Errorf("fail to get latest sync committee period")
	}
	return slot.Uint64(), nil

}

func (c *Client) GetFinalizedSyncPeriod() (uint64, error) {
	resp := &structs.GetBlockHeaderResponse{}
	err := c.get("/eth/v1/beacon/headers/finalized", nil, &resp)
	if err != nil {
		return 0, err
	}
	slot, ok := big.NewInt(0).SetString(resp.Data.Header.Message.Slot, 10)
	if !ok {
		return 0, fmt.Errorf("fail to get latest sync committee period")
	}
	period := slot.Uint64() / 8192
	return period, nil
}

func (c *Client) GetLightClientUpdates(start uint64, count uint64) ([]structs.LightClientUpdateWithVersion, error) {
	var updates []structs.LightClientUpdateWithVersion
	param := Param{}
	param.Add("start_period", start)
	param.Add("count", count)
	err := c.get("/eth/v1/beacon/light_client/updates", param, &updates)
	if err != nil {
		return nil, err
	}
	return updates, nil
}

func (c *Client) BeaconHeaderBySlot(slot uint64) (*structs.GetBlockHeaderResponse, error) {
	result := &structs.GetBlockHeaderResponse{}
	path := fmt.Sprintf("/eth/v1/beacon/headers/%v", slot)
	err := c.get(path, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) BeaconHeaderByRoot(root string) (*structs.GetBlockHeaderResponse, error) {
	result := &structs.GetBlockHeaderResponse{}
	path := fmt.Sprintf("/eth/v1/beacon/headers/%v", root)
	err := c.get(path, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetFinalityUpdate() (*structs.LightClientUpdateWithVersion, error) {
	var result structs.LightClientUpdateWithVersion
	err := c.get("/eth/v1/beacon/light_client/finality_update", nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) RetrieveBeaconHeaders(start, end uint64) ([]*structs.BeaconBlockHeader, error) {
	// todo
	headers := make([]*structs.BeaconBlockHeader, 0)
	response, err := c.BeaconHeaderBySlot(end)
	if err != nil {
		return nil, err
	}
	header := response.Data.Header.Message
	headers = append(headers, header)
	slot, ok := big.NewInt(0).SetString(header.Slot, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse slot")
	}
	if slot.Int64() == int64(start) {
		return headers, nil
	}
	found := false
	for i := end; i > start; {
		response, err = c.BeaconHeaderByRoot(header.ParentRoot)
		if err != nil {
			return nil, err
		}
		header = response.Data.Header.Message
		slot, ok = big.NewInt(0).SetString(header.Slot, 10)
		if !ok {
			return nil, fmt.Errorf("failed to parse slot")
		}
		headers = append(headers, header)
		i = slot.Uint64()
		if i == start {
			found = true
		}
	}
	if found {
		return slice.Reverse(headers), nil
	}
	return nil, fmt.Errorf("failed to %v headers", start)

}

func (c *Client) get(path string, param Param, value interface{}, headers ...Header) error {
	var reqStr string
	for k, v := range param {
		reqStr += fmt.Sprintf("%s=%v&", k, v)
	}
	if len(reqStr) != 0 {
		path = fmt.Sprintf("%s?%s", path, reqStr)
	}
	path = strings.TrimSuffix(path, "&")
	return c.httpReq(http.MethodGet, path, param, value, headers...)
}

func (c *Client) newRequest(ctx context.Context, httpMethod, url, method string, param interface{}, headers ...Header) (*http.Request, error) {
	reqData, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, httpMethod, fmt.Sprintf("%s%s", url, method), bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}
	for _, item := range headers {
		for key, value := range item {
			req.Header.Set(key, value)
		}
	}
	req.Header.Set("Connection", "keep-alive")
	return req, nil
}

func (c *Client) httpReq(httpMethod, method string, param Param, value interface{}, headers ...Header) (err error) {
	vi := reflect.ValueOf(value)
	if vi.Kind() != reflect.Ptr {
		return fmt.Errorf("value must be pointer")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	req, err := c.newRequest(ctx, httpMethod, c.endpoint, method, param, headers...)
	if err != nil {
		return err
	}
	if c.debug {
		if param != nil {
			requestData, err := json.Marshal(param)
			if err != nil {
				return fmt.Errorf("%v", err)
			}
			log.Printf("httpReq request: %v  %v \n", method, string(requestData))
		}
	}
	resp, err := c.httpClient.Do(req)
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
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if c.debug {
		log.Printf("httpReq response: %v %v \n", method, string(data))
	}
	if len(data) == 0 {
		return fmt.Errorf("data is empty")
	}
	err = json.Unmarshal(data, value)
	if err != nil {
		return err
	}
	return nil

}

type Param map[string]interface{}

func (p *Param) Add(key string, value interface{}) {
	(*p)[key] = value
}

type Header map[string]string
