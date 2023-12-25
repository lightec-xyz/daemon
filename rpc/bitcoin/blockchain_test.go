package bitcoin

import (
	"fmt"
	"testing"
)

func TestClient_scantxoutset(t *testing.T) {

	result, err := client.Scantxoutset("tb1q6lawf77u30mvs6sgcuthchgxdqm4f6n359lc4a")
	if err != nil {
		panic(err)
	}
	for _, utxo := range result.Unspents {
		fmt.Printf("txid:%v,vout:%v,scriptPubKey:%v,amount:%0.8f\n", utxo.Txid, utxo.Vout, utxo.ScriptPubKey, utxo.Amount)
	}

}
