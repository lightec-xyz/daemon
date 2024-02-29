package beacon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
)

// Client defines typed wrappers for the Beacon RPC API.
type Client struct {
	ctx context.Context
	url string
}

// Dial connects a client to the given URL.
func NewClient(rawurl string) (*Client, error) {
	return &Client{
		ctx: context.Background(),
		url: rawurl,
	}, nil
}

func (c *Client) GetBootstrap(slot uint64) (*structs.LightClientBootstrapResponse, error) {
	resp := &structs.GetBlockHeaderResponse{}

	url := c.url + fmt.Sprintf("/eth/v1/beacon/headers/%v", slot)
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to get beacon block header, failed: %s", err)
	}

	r, err := http.DefaultClient.Do(req)
	defer func() {
		err = r.Body.Close()
	}()
	if err != nil {
		return nil, fmt.Errorf("failt to get beacon block header, failed: %s", err)
	}

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failt to get beacon block header failed: %s", err)
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("fail to get beacon block header, failed: %s", err)
	}
	if data == nil {
		return nil, fmt.Errorf("empty beacon block header")
	}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, fmt.Errorf("fail to get beacon block header, failed: %s", err)
	}

	url = c.url + fmt.Sprintf("/eth/v1/beacon/light_client/bootstrap/%s", resp.Data.Root)
	req, err = http.NewRequestWithContext(c.ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	r, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r.Body.Close()
	}()
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("requesting light client bootstrap, bad status code %d", r.StatusCode)
	}
	data, err = io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("requesting light client bootstrap,failed: %s", err)
	}

	bootstrap := &structs.LightClientBootstrapResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal light client bootstrap response, failed: %s", err)
	}
	return bootstrap, nil
}

func (c *Client) GetLatestSyncPeriod() (uint64, error) {
	resp := &structs.GetBlockHeaderResponse{}

	url := c.url + fmt.Sprintf("/eth/v1/beacon/headers/finalized")
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("fail to get latest sync committee period, failed: %s", err)
	}

	r, err := http.DefaultClient.Do(req)
	defer func() {
		err = r.Body.Close()
	}()
	if err != nil {
		return 0, fmt.Errorf("failt to get latest sync committee period, failed: %s", err)
	}

	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failt to get latest sync committee period, failed: %s", err)
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, fmt.Errorf("fail to get latest sync committee period, failed: %s", err)
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return 0, fmt.Errorf("fail to get latest sync committee period, failed: %s", err)
	}

	slot, ok := big.NewInt(0).SetString(resp.Data.Header.Message.Slot, 10)
	if !ok {
		return 0, fmt.Errorf("fail to get latest sync committee period")
	}

	period := slot.Uint64() / 8192

	return period, nil
}

func (c *Client) GetBeaconBlockHeader(slot uint64) (*structs.GetBlockHeaderResponse, error) {
	resp := &structs.GetBlockHeaderResponse{}

	url := c.url + fmt.Sprintf("/eth/v1/beacon/headers/%v", slot)
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to get beacon block header, failed: %s", err)
	}

	r, err := http.DefaultClient.Do(req)
	defer func() {
		err = r.Body.Close()
	}()
	if err != nil {
		return nil, fmt.Errorf("failt to get beacon block header, failed: %s", err)
	}

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failt to get beacon block header failed: %s", err)
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("fail to get beacon block header, failed: %s", err)
	}
	if data == nil {
		return nil, fmt.Errorf("empty beacon block header")
	}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, fmt.Errorf("fail to get beacon block header, failed: %s", err)
	}

	return resp, nil
}

func (c *Client) GetLightClientUpdates(start uint64, count uint64) ([]structs.LightClientUpdateWithVersion, error) {
	uri := c.url + fmt.Sprintf("/eth/v1/beacon/light_client/updates?start_period=%d&count=%d", start, count)
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("requesting light client updates %s\n", err)
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r.Body.Close()
	}()
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("requesting light client updates, failed: %v", r.StatusCode)
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("requesting light client updates,failed: %v", err)
	}

	updates := []structs.LightClientUpdateWithVersion{}
	err = json.Unmarshal(data, &updates)
	if err != nil {
		return nil, fmt.Errorf("unmarshal light client updates response, failed: %v", err)
	}

	return updates, nil
}
