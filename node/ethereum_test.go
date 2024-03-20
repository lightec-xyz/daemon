package node

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
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
