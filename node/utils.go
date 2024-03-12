package node

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"math/big"
	"strings"
)

func UUID() string {
	newV7, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return newV7.String()
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

func NewRat() *big.Rat {
	rat := new(big.Rat)
	return rat
}
