package proof

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestClientModeProof(t *testing.T) {
	config := NewClientModeConfig()
	marshal, _ := json.Marshal(config)
	fmt.Println(string(marshal))
	node, err := NewNode(config)
	if err != nil {
		t.Fatal(err)
	}
	err = node.Start()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClusterModeProof(t *testing.T) {
	config := NewClusterModeConfig()

	node, err := NewNode(config)
	if err != nil {
		t.Fatal(err)
	}
	err = node.Start()
	if err != nil {
		t.Fatal(err)
	}
}
