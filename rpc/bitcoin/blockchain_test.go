package bitcoin

import (
	"fmt"
	"testing"
)

func TestClient_scantxoutset(t *testing.T) {

	result, err := client.Scantxoutset("bcrt1q6lawf77u30mvs6sgcuthchgxdqm4f6n3kvx4z5")
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

}
