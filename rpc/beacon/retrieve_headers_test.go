package beacon

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"github.com/prysmaticlabs/prysm/v5/container/slice"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
	"io"
	"math/big"
	"net/http"
	"os"
	"testing"
)

type BeaconHeader struct {
	Slot          string `json:"slot"`
	ProposerIndex string `json:"proposer_index"`
	ParentRoot    string `json:"parent_root"`
	StateRoot     string `json:"state_root"`
	BodyRoot      string `json:"body_root"`
}

type BeaconHeaderChain struct {
	BeginSlot           uint64         `json:"begin_slot"`
	BeginRoot           string         `json:"begin_root"`
	MiddleBeaconHeaders []BeaconHeader `json:"middle_beacon_block_headers"`
	EndSlot             uint64         `json:"end_slot"`
	EndRoot             string         `json:"end_root"`
}

func retrieveBeaconHeaderBySlot(ctx context.Context, uri string, slot int) (*structs.BeaconBlockHeader, error) {
	resp := &structs.GetBlockHeaderResponse{}
	uri = uri + fmt.Sprintf("/eth/v1/beacon/headers/%v", slot)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	r, err := http.DefaultClient.Do(req)
	defer func() {
		err = r.Body.Close()
	}()
	if err != nil {
		return nil, err
	}

	if r.StatusCode != http.StatusOK {
		return nil, err
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}
	return resp.Data.Header.Message, nil
}

func retrieveBeaconHeaderByParentRoot(ctx context.Context, uri string, parentRoot string) (*structs.BeaconBlockHeader, error) {
	resp := &structs.GetBlockHeaderResponse{}
	uri = uri + fmt.Sprintf("/eth/v1/beacon/headers/%v", parentRoot)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	r, err := http.DefaultClient.Do(req)
	defer func() {
		err = r.Body.Close()
	}()
	if err != nil {
		return nil, err
	}

	if r.StatusCode != http.StatusOK {
		return nil, err
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}
	return resp.Data.Header.Message, nil
}

func RetrieveBeaconHeaders(ctx context.Context, uri string, start, end int) ([]structs.BeaconBlockHeader, error) {
	headers := make([]structs.BeaconBlockHeader, 0)
	header, err := retrieveBeaconHeaderBySlot(ctx, uri, end)
	if err != nil {
		return nil, err
	}
	headers = append(headers, *header)

	slot, ok := big.NewInt(0).SetString(header.Slot, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse slot")
	}
	if slot.Int64() == int64(start) {
		return headers, nil
	}

	found := false
	for i := end; i > start; {
		header, err = retrieveBeaconHeaderByParentRoot(ctx, uri, header.ParentRoot)
		if err != nil {
			return nil, err
		}

		slot, ok = big.NewInt(0).SetString(header.Slot, 10)
		if !ok {
			return nil, fmt.Errorf("failed to parse slot")
		}
		headers = append(headers, *header)
		i = int(slot.Int64())
		if i == start {
			found = true
		}
	}

	if found {
		return slice.Reverse(headers), nil
	}
	return nil, fmt.Errorf("failed to %v headers", start)
}

var rpcURL = ""

func TestRetrieveBeaconHeaders(t *testing.T) {
	ctx := context.Background()
	headers, err := RetrieveBeaconHeaders(ctx, rpcURL, 1315329, 1315360)

	require.NoError(t, err)

	headersChain := &BeaconHeaderChain{}
	beginHeader, err := headers[0].ToConsensus()
	require.NoError(t, err)
	endHeader, err := headers[len(headers)-1].ToConsensus()
	require.NoError(t, err)

	beginRootBytes, err := beginHeader.HashTreeRoot()
	require.NoError(t, err)
	endRootBytes, err := endHeader.HashTreeRoot()
	require.NoError(t, err)

	headersChain.BeginRoot = hex.EncodeToString(beginRootBytes[:])
	headersChain.BeginSlot = uint64(beginHeader.Slot)
	headersChain.EndRoot = hex.EncodeToString(endRootBytes[:])
	headersChain.EndSlot = uint64(endHeader.Slot)

	middleHeaders := make([]BeaconHeader, 0)
	for i := 1; i < len(headers)-1; i++ {
		middleHeaders = append(middleHeaders, BeaconHeader{
			Slot:          headers[i].Slot,
			ProposerIndex: headers[i].ProposerIndex,
			ParentRoot:    headers[i].ParentRoot,
			StateRoot:     headers[i].StateRoot,
			BodyRoot:      headers[i].BodyRoot,
		})
	}
	headersChain.MiddleBeaconHeaders = middleHeaders
	fn := fmt.Sprintf("beacon_headers_%v_%v.json", headersChain.BeginSlot, headersChain.EndSlot)

	data, err := json.Marshal(headersChain)
	require.NoError(t, err)

	err = os.WriteFile(fn, data, 0644)
	require.NoError(t, err)
}

func TestRetrieveBeaconHeaderBySlot(t *testing.T) {
	ctx := context.Background()
	header, err := retrieveBeaconHeaderBySlot(ctx, rpcURL, 1315329)
	require.NoError(t, err)
	fmt.Printf("header: %v\n", header)
	require.Equal(t, "1315329", header.Slot)
}

func TestRetrieveBeaconHeaderByParentRoot(t *testing.T) {
	ctx := context.Background()
	header, err := retrieveBeaconHeaderByParentRoot(ctx, rpcURL, "0x54c88fe9971560cc02503ba296e0526366ca1c5ee6b89450783596a68454dcce")
	require.NoError(t, err)
	fmt.Printf("header: %v\n", header)
	require.Equal(t, "1315328", header.Slot)
}
