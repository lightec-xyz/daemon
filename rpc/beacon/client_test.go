package beacon

import "testing"

// var endpoint = "http://127.0.0.1:8970"
var endpoint = "http://37.120.151.183:8970"
var err error
var client *Client

func init() {
	client, err = NewClient(endpoint)
	if err != nil {
		panic(err)
	}
}

func TestClient(t *testing.T) {
	latestSyncPeriod, err := client.GetLatestSyncPeriod()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(latestSyncPeriod)

	latestSlot, err := client.GetLatestSlot()
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
