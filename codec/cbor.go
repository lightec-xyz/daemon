package codec

import "github.com/fxamacker/cbor/v2"

func Marshal(val interface{}) ([]byte, error) {
	return cbor.Marshal(val)
}
func Unmarshal(b []byte, val interface{}) error {
	return cbor.Unmarshal(b, val)
}
