package beacon

import "testing"

var endpoint = "http://127.0.0.1:8970"

func TestClient(t *testing.T) {
	// todo more test
	client, err := NewClient(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	latestSyncPeriod, err := client.GetLatestSyncPeriod()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(latestSyncPeriod)

	latestSlot, err := client.GetLatestSlot()
	if err != nil {
		t.Fatal(err)
	}

	bootstrap, err := client.Bootstrap(latestSlot)
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
