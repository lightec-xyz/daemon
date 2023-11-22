package codec

import (
	"github.com/ethereum/go-ethereum/rlp"
	"io"
)

// todo maybe remove

func Decode(r io.Reader, val interface{}) error {
	return rlp.Decode(r, val)
}
func DecodeBytes(b []byte, val interface{}) error {
	return rlp.DecodeBytes(b, val)
}

func Encode(w io.Writer, val interface{}) error {
	return rlp.Encode(w, val)
}

func EncodeToBytes(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
