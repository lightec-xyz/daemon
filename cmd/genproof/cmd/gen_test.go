package cmd

import (
	"fmt"
	"testing"
)

func TestGenProof(t *testing.T) {
	localProof, err := NewLocalProof(genesisSlot, "", "")
	if err != nil {
		fmt.Printf("new local proof error: %v \n", err)
		return
	}
	err = localProof.GenProof(proofType, index)
	if err != nil {
		fmt.Printf("gen proof error: %v \n", err)
		return
	}
}
