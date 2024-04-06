package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestGenerateAddress(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("privateKey:  %x \n", privateKey.D.Bytes())
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("fail")
	}
	compressPubkey := crypto.CompressPubkey(publicKeyECDSA)
	fmt.Printf("Uncompressed:%x \n", crypto.FromECDSAPub(publicKeyECDSA))
	fmt.Printf("Compressed:  %x \n", compressPubkey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("address:     %v \n", address)

}

func Test_decodeRedeemLog(t *testing.T) {
	// https://holesky.etherscan.io/tx/0x3db1bb46352898a1ff0349274d0dcc7c8e78020ab2268c2bfa0863ab0e9de001#eventlog
	hash := common.HexToHash("0x3db1bb46352898a1ff0349274d0dcc7c8e78020ab2268c2bfa0863ab0e9de001")

	ec, err := ethclient.Dial("https://1rpc.io/holesky")
	require.NoError(t, err)

	receipt, err := ec.TransactionReceipt(context.Background(), hash)
	require.NoError(t, err)

	btcRawTx, sigHashs, err := DecodeRedeemLog(receipt.Logs[3].Data)
	require.NoError(t, err)

	t.Logf("btcRawTx: %v", hexutil.Encode(btcRawTx))
	t.Logf("sigHashs: %x", sigHashs)
}

func Test_getRawTxAndReceipt(t *testing.T) {
	// https://holesky.etherscan.io/tx/0x3db1bb46352898a1ff0349274d0dcc7c8e78020ab2268c2bfa0863ab0e9de001
	hash := common.HexToHash("0x3db1bb46352898a1ff0349274d0dcc7c8e78020ab2268c2bfa0863ab0e9de001")

	ec, err := ethclient.Dial("https://1rpc.io/holesky")
	require.NoError(t, err)

	tx, _, err := ec.TransactionByHash(context.Background(), hash)
	require.NoError(t, err)

	receipt, err := ec.TransactionReceipt(context.Background(), hash)
	require.NoError(t, err)

	rawTx, rawReceipt := GetRawTxAndReceipt(tx, receipt)
	t.Logf("rawTx: %v", hexutil.Encode(rawTx))
	t.Logf("rawReceipt: %v", hexutil.Encode(rawReceipt))
}
