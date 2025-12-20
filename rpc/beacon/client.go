package beacon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lightec-xyz/daemon/rpc/beacon/types"
	"github.com/prysmaticlabs/prysm/v5/container/slice"

	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
)

type IBeacon interface {
	GetBlindedBlock(slot uint64) (types.BindBlockResp, error)
	Eth1MapToEth2(slot uint64) (*Eth1MapToEth2, error)
	Bootstrap(slot uint64) (*types.BootstrapResp, error)
	BootstrapByRoot(root string) (*types.BootstrapResp, error)
	BeaconHeaders(slot uint64) (*structs.GetBlockHeaderResponse, error)
	FinalizedSlot() (uint64, error)
	FinalizedPeriod() (uint64, error)
	LightClientUpdates(period, count uint64) ([]types.LightClientUpdateResp, error)
	BeaconHeaderBySlot(slot uint64) (*structs.GetBlockHeaderResponse, error)
	BeaconHeaderByRoot(root string) (*structs.GetBlockHeaderResponse, error)
	GetFinalityUpdate() (types.LightClientFinalityUpdateResp, error)
	RetrieveBeaconHeaders(start, end uint64) ([]*structs.BeaconBlockHeader, error)
}

type Client struct {
	ctx        context.Context
	endpoints  []string
	timeout    time.Duration
	debug      bool
	httpClient *http.Client
}

func NewClient(urls ...string) (*Client, error) {
	if len(urls) == 0 {
		return nil, errors.New("no url")
	}
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     6 * time.Hour, // todo
			DisableKeepAlives:   false,
			MaxIdleConnsPerHost: 10,
		},
	}
	return &Client{
		ctx:        context.Background(),
		endpoints:  urls,
		timeout:    6 * time.Hour, // todo
		debug:      false,
		httpClient: client,
	}, nil
}

func (c *Client) GetBlindedBlock(slot uint64) (types.BindBlockResp, error) {
	var result types.BindBlockResp
	path := fmt.Sprintf("/eth/v1/beacon/blinded_blocks/%d", slot)
	err := c.get(path, nil, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (c *Client) Eth1MapToEth2(slot uint64) (*Eth1MapToEth2, error) {
	blindedBlock, err := c.GetBlindedBlock(slot)
	if err != nil {
		return nil, err
	}
	blockNumber, err := strconv.ParseUint(blindedBlock.Data.Message.Body.ExecutionPayloadHeader.BlockNumber, 10, 64)
	if err != nil {
		return nil, err
	}
	slot, err = strconv.ParseUint(blindedBlock.Data.Message.Slot, 10, 64)
	if err != nil {
		return nil, err
	}

	slotMapInfo := Eth1MapToEth2{
		BlockNumber: blockNumber,
		BlockHash:   blindedBlock.Data.Message.Body.ExecutionPayloadHeader.BlockHash,
		BlockSlot:   slot,
		BlockRoot:   blindedBlock.Data.Message.StateRoot,
	}
	return &slotMapInfo, nil
}

func (c *Client) Bootstrap(slot uint64) (*types.BootstrapResp, error) {
	beaconHeaders, err := c.BeaconHeaders(slot)
	if err != nil {
		return nil, err
	}
	bootstrap, err := c.BootstrapByRoot(beaconHeaders.Data.Root)
	if err != nil {
		return nil, err
	}
	return bootstrap, nil
}

func (c *Client) BootstrapByRoot(root string) (*types.BootstrapResp, error) {
	bootstrap := &types.BootstrapResp{}
	path := fmt.Sprintf("/eth/v1/beacon/light_client/bootstrap/%v", root)
	err := c.get(path, nil, &bootstrap)
	if err != nil {
		return nil, err
	}
	return bootstrap, nil
}

func (c *Client) BeaconHeaders(slot uint64) (*structs.GetBlockHeaderResponse, error) {
	resp := &structs.GetBlockHeaderResponse{}
	path := fmt.Sprintf("/eth/v1/beacon/headers/%v", slot)
	err := c.get(path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) FinalizedSlot() (uint64, error) {
	resp := &structs.GetBlockHeaderResponse{}
	err := c.get("/eth/v1/beacon/headers/finalized", nil, &resp)
	if err != nil {
		return 0, err
	}
	slot, ok := big.NewInt(0).SetString(resp.Data.Header.Message.Slot, 10)
	if !ok {
		return 0, fmt.Errorf("fail to get latest finalized slot")
	}
	return slot.Uint64(), nil
}

func (c *Client) FinalizedPeriod() (uint64, error) {
	finalizedSlot, err := c.FinalizedSlot()
	if err != nil {
		return 0, err
	}
	period := finalizedSlot / 8192
	return period, nil
}

func (c *Client) LightClientUpdates(start uint64, count uint64) ([]types.LightClientUpdateResp, error) {
	var updates []types.LightClientUpdateResp
	param := Param{}
	param.Add("count", count)
	param.Add("start_period", start)
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

func (c *Client) GetFinalityUpdate() (types.LightClientFinalityUpdateResp, error) {
	var result types.LightClientFinalityUpdateResp
	err := c.get("/eth/v1/beacon/light_client/finality_update", nil, &result)
	if err != nil {
		return result, err
	}
	return result, nil
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
	msg := "request error "
	for _, url := range c.endpoints {
		err := c.httpReq(http.MethodGet, url, path, param, value, headers...)
		if err != nil {
			msg = msg + err.Error() + "\n"
			continue
		}
		return nil
	}
	return errors.New(msg)
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

func (c *Client) httpReq(httpMethod, url, method string, param Param, value interface{}, headers ...Header) (err error) {
	vi := reflect.ValueOf(value)
	if vi.Kind() != reflect.Ptr {
		return fmt.Errorf("value must be pointer")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.timeout)
	defer cancelFunc()
	req, err := c.newRequest(ctx, httpMethod, url, method, param, headers...)
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
		data, _ := io.ReadAll(resp.Body)
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
		log.Printf("httpReq response: %v\n %v \n", method, string(data))
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
