package bitcoinClient

import (
	"fmt"
	"testing"
)

func TestClient_CreateMultiAddress(t *testing.T) {
	createmultisig, err := client.Createmultisig(2,
		"033dad71ad8c37910aa5fe0d6ef9bb30a7ddab298291f4cc6cfb01fb022cc27f49",
		"033dad71ad8c37910aa5fe0d6ef9bb30a7ddab298291f4cc6cfb01fb022cc27f49",
		"033dad71ad8c37910aa5fe0d6ef9bb30a7ddab298291f4cc6cfb01fb022cc27f49")
	if err != nil {
		panic(err)
	}
	fmt.Println(createmultisig)
}
