package proof

import (
	"testing"
)

func TestClientModeProof(t *testing.T) {
	config := NewClientModeConfig()
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
