package node

import "testing"

func TestTxIdIsEmpty(t *testing.T) {
	txId := [32]byte{}
	res := TxIdIsEmpty(txId)
	t.Log(res)
}
