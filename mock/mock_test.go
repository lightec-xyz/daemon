package main

import (
	"github.com/lightec-xyz/daemon/logger"
	"math/big"
	"testing"
)

var err error
var mock *Mock

func init() {
	logger.InitLogger()
	mock, err = NewMock("local")
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

func TestMock_DepositBtcToEth(t *testing.T) {
	err := mock.DepositBtcToEth("0b346341a54aca2d7b86d1b6d6a44c318650d4e311bfb12628ada949a3648dfa",
		"0x771815eFD58e8D6e66773DB0bc002899c00d5b0c", 1, big.NewInt(1199998950))
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
