package dfinity

import (
	"encoding/hex"
	"fmt"
	"github.com/aviate-labs/agent-go/identity"
	"github.com/aviate-labs/agent-go/principal"
)

type Options struct {
	WalletCanisterId string
	TxCanisterId     string
	BlockCanisterId  string
	identity         identity.Identity
}

func NewOption(walletCanisterId, txCanisterId, blockCanisterId string, identity identity.Identity) *Options {
	return &Options{
		WalletCanisterId: walletCanisterId,
		TxCanisterId:     txCanisterId,
		BlockCanisterId:  blockCanisterId,
		identity:         identity,
	}
}

type WalletCallArg struct {
	Canister   principal.Principal `ic:"canister" json:"canister"`
	MethodName string              `ic:"method_name" json:"method_name"`
	Args       []byte              `ic:"args" json:"args"`
	Cycles     uint64              `ic:"cycles" json:"cycles"`
}

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

func (ts *TxSignature) SignatureBytes() ([][]byte, error) {
	var signatureBytes [][]byte
	for _, item := range ts.Signature {
		sigBytes, err := hex.DecodeString(item)
		if err != nil {
			return nil, err
		}
		signatureBytes = append(signatureBytes, sigBytes)
	}
	return signatureBytes, nil
}
