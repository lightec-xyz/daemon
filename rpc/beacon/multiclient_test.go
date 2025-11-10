package beacon

import "testing"

func TestMultiClient_FinalizedSlot(t *testing.T) {
	client, err := NewMultiClient("a", "b", "c")
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

func TestShift(t *testing.T) {
	res := shift([]string{"a", "b", "c"})
	t.Log(res)
}
