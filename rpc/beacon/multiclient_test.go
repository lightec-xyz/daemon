package beacon

import "testing"

func TestMultiClient_FinalizedSlot(t *testing.T) {
	client, err := NewMultiClient("", "")
	if err != nil {
		t.Fatal(err)
	}
	for {
		latestSlot, err := client.FinalizedSlot()
		if err != nil {
			t.Fatal(err)
		}
		client.Next()
		t.Log(latestSlot)
	}
}
