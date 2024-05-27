package common

import (
	"os"
	"testing"
)

func TestCheckZkParametersMd5(t *testing.T) {
	data, err := os.ReadFile("/Users/red/lworkspace/lightec/daemon/node/test/redeem.vk")
	if err != nil {
		t.Fatal(err)
	}
	hexMd5 := HexMd5(data)
	t.Log(hexMd5)

}
