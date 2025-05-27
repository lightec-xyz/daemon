package ethereum

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"os"
)

func GetTestSecret() string {
	return os.Getenv("EthTestSecret")
}

func Fix32Bytes(data []byte) [32]byte {
	var fixCp [32]byte
	copy(fixCp[:], data)
	return fixCp
}

func DecodeRedeemLog(logData []byte) (btcRawTx []byte, sigHashs [][32]byte, err error) {
	t1, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return nil, nil, err
	}
	t2, err := abi.NewType("bytes32[]", "", nil)
	if err != nil {
		return nil, nil, err
	}

	arguments := abi.Arguments{
		abi.Argument{Type: t1},
		abi.Argument{Type: t2},
	}
	decoded, err := arguments.UnpackValues(logData)
	if err != nil {
		return nil, nil, err
	}

	btcRawTx, ok := decoded[0].([]byte)
	if !ok {
		return nil, nil, err
	}
	sigHashs, ok = decoded[1].([][32]byte)
	if !ok {
		return nil, nil, err
	}

	return btcRawTx, sigHashs, nil
}

func GetRawTxAndReceipt(tx *types.Transaction, receipt *types.Receipt) (rawTx, rawReceipt []byte) {
	buf1 := new(bytes.Buffer)
	types.Transactions{tx}.EncodeIndex(0, buf1)
	rawTx = buf1.Bytes()

	buf2 := new(bytes.Buffer)
	types.Receipts{receipt}.EncodeIndex(0, buf2)
	rawReceipt = buf2.Bytes()

	return
}
func privateKeyToAddr(secret string) (string, error) {
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
