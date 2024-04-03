package circuits

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

type StoreProof struct {
	Period  uint64 `json:"period"`
	Proof   []byte `json:"proof"`
	Witness []byte `json:"witness"`
}

func TestCircuit(t *testing.T) {

	for index := 153; index < 162; index++ {
		path := fmt.Sprintf("/Users/red/lworkspace/lightec/daemon/circuits/test/unit/%v", index)
		data, err := ioutil.ReadFile(path)
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
		witness, err := ParseWitness(proof.Witness)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(witness)

	}

}
