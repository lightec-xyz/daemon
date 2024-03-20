package node

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func Test_decodeRedeemLog(t *testing.T) {
	// https://holesky.etherscan.io/tx/0x578b6beae157c2ebc22331f68314e68cb727f90eca0bdb1252bb083e4a1b1a83#eventlog
	dataStr := "02000000013faa6f78548e637a17a8e135ef688ccac469587aef33875ef665077578fd9e6b0000000000ffffffff01384a00000000000016001464d468d12f61295882b0f8f63c64b58e2af058e400000000"

	dataBytes, err := hexutil.Decode(dataStr)
	require.NoError(t, err)

	btcRawTx, sigHashs, err := decodeRedeemLog(dataBytes)
	require.NoError(t, err)

	t.Logf("btcRawTx: %v", hexutil.Encode(btcRawTx))
	t.Logf("sigHashs: %x", sigHashs)
}

func Test_getRawTxAndReceipt(t *testing.T) {
	// https://holesky.etherscan.io/tx/0x578b6beae157c2ebc22331f68314e68cb727f90eca0bdb1252bb083e4a1b1a83
	hash := common.HexToHash("0x578b6beae157c2ebc22331f68314e68cb727f90eca0bdb1252bb083e4a1b1a83")

	ec, err := ethclient.Dial("https://1rpc.io/holesky")
	require.NoError(t, err)

	tx, _, err := ec.TransactionByHash(context.Background(), hash)
	require.NoError(t, err)

	receipt, err := ec.TransactionReceipt(context.Background(), hash)
	require.NoError(t, err)

	rawTx, rawReceipt := getRawTxAndReceipt(tx, receipt)
	t.Logf("rawTx: %v", hexutil.Encode(rawTx))
	t.Logf("rawReceipt: %v", hexutil.Encode(rawReceipt))
}
