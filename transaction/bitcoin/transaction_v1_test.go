package bitcoin

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestDemo(t *testing.T) {
	//src:="e8c84a631D71E1Bb7083D3a82a3a74870a286B97"
	data, err := hex.DecodeString("6a14e8c84a631d71e1bb7083d3a82a3a74870a286b97")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v : %x\n", len(data[2:]), data[2:])
}
