package node

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strings"
)

func TxIdToProofId(txId string) string {
	pTxID := fmt.Sprintf("%s%s", ProofPrefix, txId)
	return pTxID
}

func BtcToSat(value float64) int64 {
	valueRat := NewRat().Mul(NewRat().SetFloat64(value), NewRat().SetUint64(100000000))
	floatStr := valueRat.FloatString(1)
	valuesStr := strings.Split(floatStr, ".")
	amountBig, ok := big.NewInt(0).SetString(valuesStr[0], 10)
	if !ok {
		panic(fmt.Sprintf("never should happen:%v", value))
	}
	return amountBig.Int64()
}

func privateKeyToEthAddr(secret string) (string, error) {
	privateKey, err := crypto.HexToECDSA(secret)
	if err != nil {
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return address, nil
}

func Str2Big(amount string, decimals int) (*big.Int, error) {
	//todo
	amt, _, err := big.ParseFloat(amount, 10, 256, big.ToNearestEven)
	if err != nil || amt.Sign() < 0 {
		return nil, errors.New("Failed to parse amount: " + amount)
	}
	index := strings.Index(amount, ".")
	decimalsLen := 0
	if index >= 0 {
		decimalsLen = len(amount[index+1:])
		amount = amount[:index] + amount[index+1:]
	}
	if decimals > decimalsLen {
		amount = amount + strings.Repeat("0", decimals-decimalsLen)
	} else {
		amount = amount[:len(amount)-(decimalsLen-decimals)]
	}
	res, _ := big.NewInt(0).SetString(string(amount), 10)
	floatRes, _ := big.NewFloat(0).SetString(string(amount))
	base, _ := big.NewFloat(0).SetString("1" + strings.Repeat("0", decimals))
	amt = amt.Mul(amt, base)
	max := big.NewFloat(0).Mul(amt, big.NewFloat(1.1))
	min := big.NewFloat(0).Mul(amt, big.NewFloat(0.9))
	if floatRes.Cmp(max) <= 0 && floatRes.Cmp(min) >= 0 {
		return res, nil
	} else {
		bigAmt, _ := amt.Int(nil)
		return bigAmt, nil
	}
}
func NewRat() *big.Rat {
	rat := new(big.Rat)
	return rat
}
