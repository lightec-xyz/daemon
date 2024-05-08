package common

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lightec-xyz/daemon/logger"
	txineth2 "github.com/lightec-xyz/provers/circuits/tx-in-eth2"
	"strings"
	"testing"
)

func TestHexToTxVar(t *testing.T) {
	txHash := "0x16a4e0568da1d0b53c75a990e74f08996d112ff03c917e9a63b138b7a02d5ec5"
	ethClient, err := ethclient.DialContext(nil, "https://ethereum-holesky-rpc.publicnode.com")
	if err != nil {
		t.Fatal(err)
	}
	txVar, receiptVar, err := txineth2.GenerateTxAndReceiptU128Padded(ethClient, txHash)
	if err != nil {
		logger.Error("get tx and receipt error: %v", err)
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

func TestGetSlot(t *testing.T) {
	//1531905
	//0x622af9392653f10797297e2fa72c6236db55d28234fad5a12c098349a8c5bd3f
	// 1434666
	slot, err := GetSlot(1440262)
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
