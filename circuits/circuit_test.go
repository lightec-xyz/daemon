package circuits

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

type StoreProof struct {
	Period  uint64 `json:"period"`
	Proof   []byte `json:"proof"`
	Witness []byte `json:"witness"`
}

func TestCircuit(t *testing.T) {
	data, err := ioutil.ReadFile("/Users/red/lworkspace/lightec/daemon/circuits/test/unit/153")
	if err != nil {
		t.Fatal(err)
	}
	proof := StoreProof{}
	err = json.Unmarshal(data, &proof)
	if err != nil {
		t.Fatal(err)
	}
	proofBytes, err := ParseProof(proof.Proof)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(proofBytes)
}
