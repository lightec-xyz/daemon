package dfinity

import (
	"encoding/hex"
	"fmt"
)

type BlockSignature struct {
	Hash      string
	Height    uint32
	Signature string
}

func (bs *BlockSignature) ToRS() ([]byte, []byte, error) {
	sigBytes, err := hex.DecodeString(bs.Signature)
	if err != nil {
		return nil, nil, err
	}
	if len(sigBytes) < 64 {
		return nil, nil, fmt.Errorf("invalid signature length")
	}
	return sigBytes[:32], sigBytes[32:], nil
}

type TxSignature struct {
	Signed    bool
	Signature []string
}
