package common

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	txineth2 "github.com/lightec-xyz/provers/circuits/tx-in-eth2"
)

func TestHexToTxVar(t *testing.T) {
	txHash := "0x16a4e0568da1d0b53c75a990e74f08996d112ff03c917e9a63b138b7a02d5ec5"
	ethClient, err := ethclient.DialContext(nil, "https://ethereum-holesky-rpc.publicnode.com")
	if err != nil {
		t.Fatal(err)
	}
	txVar, receiptVar, err := txineth2.GenerateTxAndReceiptU128Padded(ethClient, txHash)
	if err != nil {
		fmt.Printf("get tx and receipt error: %v", err)
		t.Fatal(err)
	}
	t.Log(txVar[0])
	varToHex, err := TxVarToHex(txVar)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(varToHex)

	bigToTxVar, err := HexToTxVar(varToHex)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bigToTxVar)

	receiptVarToHex, err := ReceiptVarToHex(receiptVar)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(receiptVarToHex)

	hexToTxVar, err := HexToTxVar(varToHex)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hexToTxVar)

	toReceiptVar, err := HexToReceiptVar(receiptVarToHex)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(toReceiptVar)

}

func TestUuid(t *testing.T) {
	uuid, err := Uuid()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(uuid)
}

func TestProofId(t *testing.T) {
	proofId := NewProofId(DepositTxType, 0, 0, "")
	t.Log(proofId)
	proofId1 := NewProofId(DepositTxType, 100, 0, "")
	t.Log(proofId1)
	proofId2 := NewProofId(DepositTxType, 100, 101, "")
	t.Log(proofId2)
	proofId3 := NewProofId(DepositTxType, 100, 101, "sdfsdfsdf")
	t.Log(proofId3)
	proofId4 := NewProofId(DepositTxType, 0, 0, "sdfsdfsdf")
	t.Log(proofId4)
	proofId5 := NewProofId(DepositTxType, 100, 0, "sdfsdfsdf")
	t.Log(proofId5)

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
