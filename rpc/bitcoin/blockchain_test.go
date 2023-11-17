package bitcoin

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestClient_scantxoutset(t *testing.T) {

	result, err := client.Scantxoutset("tb1qn9fpljh5ggp407z02jx8x76pemzclgd6rla0qp")
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	for _, utxo := range result.Unspents {
		fmt.Printf("txid:%v,vout:%v,scriptPubKey:%v,amount:%0.8f\n", utxo.Txid, utxo.Vout, utxo.ScriptPubKey, utxo.Amount)
	}

}
