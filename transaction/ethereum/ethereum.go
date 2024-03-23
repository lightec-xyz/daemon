package ethereum

import (
	"bytes"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

func DecodeRedeemLog(logData []byte) (btcRawTx []byte, sigHashs [][32]byte, err error) {
	t1, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return nil, nil, err
	}
	t2, err := abi.NewType("bytes32[]", "", nil)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	arguments := abi.Arguments{
		abi.Argument{Type: t1},
		abi.Argument{Type: t2},
	}
	decoded, err := arguments.UnpackValues(logData)
	if err != nil {
		log.Fatal(err)
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
