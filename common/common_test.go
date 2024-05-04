package common

import (
	"strings"
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
	//0x622af9392653f10797297e2fa72c6236db55d28234fad5a12c098349a8c5bd3f
	// 1434666
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

func TestDemo(t *testing.T) {
	res := "http://18.116.118.39:18332"
	url := strings.Replace(res, "http://", "", 1)
	t.Log(url)

}
