package codec

import (
	"github.com/fxamacker/cbor/v2"
	"testing"
)

func TestCbor(t *testing.T) {
	type Animal struct {
		Age    int
		Name   string
		Owners []string
		Male   bool
	}
	animal := Animal{Age: 4, Name: "Candy", Owners: []string{"Mary", "Joe"}}
	b, err := cbor.Marshal(animal)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)
	err = cbor.Unmarshal(b, &animal)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(animal)

}
