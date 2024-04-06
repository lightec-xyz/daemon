package bitcoin

import "testing"

func TestGetNetInfo(t *testing.T) {
	result, err := client.GetNetworkInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result.Relayfee)
}
