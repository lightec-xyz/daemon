package proof

import "testing"

func TestLocalProofNode(t *testing.T) {
	node, err := NewNode(LocalDevConfig())
	if err != nil {
		t.Fatal(err)
	}
	err = node.Start()
	if err != nil {
		t.Fatal(err)
	}
}
