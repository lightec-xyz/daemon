package beacon

import (
	"encoding/json"
	"testing"
)

var endpoint = "http://127.0.0.1:9003"

var err error
var client *Client

func init() {
	client, err = NewClient(endpoint)
	if err != nil {
		panic(err)
	}
}

func TestClient(t *testing.T) {
	latestSyncPeriod, err := client.GetFinalizedSyncPeriod()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(latestSyncPeriod)

	latestSlot, err := client.GetLatestFinalizedSlot()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(latestSlot)
	bootstrap, err := client.Bootstrap(latestSlot - 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bootstrap)
	updates, err := client.GetLightClientUpdates(latestSyncPeriod-1, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(updates)
}

func TestCliet_GetLatestPeriod(t *testing.T) {
	period, err := client.GetFinalizedSyncPeriod()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(period)
}

func TestClient_Bootstrap(t *testing.T) {
	bootstrap, err := client.Bootstrap(466 * 8192)
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(bootstrap)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

func TestClient_GetLightClientUpdates(t *testing.T) {
	updates, err := client.GetLightClientUpdates(938, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(updates)
}

func TestClient_RetrieveBeaconHeaders(t *testing.T) {
	start := 2513203
	end := 2513215
	headers, err := client.RetrieveBeaconHeaders(uint64(start), uint64(end))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(headers))
	for index, header := range headers {
		t.Log(index, header.Slot)
	}
}

func TestClient_GetFinalityUpdate(t *testing.T) {
	update, err := client.GetFinalityUpdate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(update.Data)
}

func TestClient_GetBeaconHeaders(t *testing.T) {
	result, err := client.BeaconHeaderBySlot(9192)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	headers, err := client.GetBeaconHeaders(9192)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(headers)

}
