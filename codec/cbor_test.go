package codec

import (
	"encoding/json"
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

func TestCbor_Demo(t *testing.T) {
	data := "woheni"
	t.Log([]byte(data))
	bytes, err := cbor.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	var result string
	err = cbor.Unmarshal(bytes, &result)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	t.Log(bytes)
	marshal, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(marshal)
	var tmp string
	err = json.Unmarshal(marshal, &tmp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tmp)

}

func TestCodec_Demo(t *testing.T) {
	value := 100
	result, err := cbor.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	var tmp int
	err = cbor.Unmarshal(result, &tmp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tmp)
}
