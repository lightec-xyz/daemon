package beacon

import "testing"

var endpoint = "http://127.0.0.1:9870"

// var endpoint = "http://37.120.151.183:8970"
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

func TestClient_Bootstrap(t *testing.T) {
	bootstrap, err := client.Bootstrap(153 * 8192)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bootstrap)
}

func TestClient_GetLightClientUpdates(t *testing.T) {
	updates, err := client.GetLightClientUpdates(153, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(updates)
}

func TestClient_RetrieveBeaconHeaders(t *testing.T) {
	latestSlot, err := client.GetLatestFinalizedSlot()
	if err != nil {
		t.Fatal(err)
	}
	latestSlot = latestSlot - 1
	headers, err := client.RetrieveBeaconHeaders(latestSlot-100, latestSlot)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(headers))
	for _, header := range headers {
		t.Log(header.ParentRoot)
	}
}

func TestClient_GetFinalityUpdate(t *testing.T) {
	slot, err := client.GetLatestFinalizedSlot()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(slot)
	t.Log(slot / 8192)
	update, err := client.GetFinalityUpdate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(update.Data)
}
