package node

import (
	"fmt"
	"testing"
)

func TestSafeList(t *testing.T) {
	safeList := NewSafeList()
	request := ProofRequest{
		TxId: "111",
	}
	safeList.PushBack(request)
	proofRequest := ProofRequest{
		TxId: "222",
	}
	safeList.PushBack(proofRequest)
	value := ProofRequest{
		TxId: "333",
	}
	safeList.PushBack(value)
	for e := safeList.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}
