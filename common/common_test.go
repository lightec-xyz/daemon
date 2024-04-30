package common

import (
	"testing"
)

func TestUuid(t *testing.T) {
	uuid, err := Uuid()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(uuid)
}

func TestGetSlot(t *testing.T) {
	//1531905
	slot, err := GetSlot(1434666)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(slot)
}

func TestProofId(t *testing.T) {
	proofId := NewProofId(DepositTxType, 100, "")
	t.Log(proofId)
	proofId1 := NewProofId(DepositTxType, 0, "test")
	t.Log(proofId1)
	proofId2 := NewProofId(DepositTxType, 100, "test")
	t.Log(proofId2)
	zkProofType, index, txHash, err := ParseProofId(proofId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(zkProofType.String(), index, txHash)
}
