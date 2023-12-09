package main

import (
	"github.com/lightec-xyz/daemon/logger"
	"testing"
)

var err error
var mock *Mock

func init() {
	logger.InitLogger()
	mock, err = NewMock("testnet")
	if err != nil {
		panic(err)
	}
}
func TestMockDeposit(t *testing.T) {
	err = mock.DepositBtc(10000)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMockRedeem(t *testing.T) {
	//201982
	err := mock.RedeemTx(10000)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMockMergeTx(t *testing.T) {
	err := mock.MergeBtcTx()
	if err != nil {
		t.Fatal(err)
	}
}

func TestOther(t *testing.T) {
	result := floatToInt(0.00204582)
	t.Log(result)

}
