package main

import (
	"math/big"
	"testing"

	"github.com/lightec-xyz/daemon/logger"
)

var err error
var mock *Mock

func init() {
	logger.InitLogger()
	//mock, err = NewMock("testnet")
	//if err != nil {
	//	panic(err)
	//}
}
func TestMockDeposit(t *testing.T) {
	err = mock.DepositBtc(10000)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMock_DepositBtcToEth(t *testing.T) {
	err := mock.DepositBtcToEth("24335674710be9e120ef40f96c03959960a7eba9b2ddde00b3740046d41d5b4c",
		"0x771815eFD58e8D6e66773DB0bc002899c00d5b0c", 1, big.NewInt(98950))
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
	result := floatToInt(0.123456789)
	t.Log(result)

}
