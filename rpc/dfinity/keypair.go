package dfinity

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/aviate-labs/agent-go/identity"
)

type KeyPair struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	Identity   *identity.Ed25519Identity
}

func (k *KeyPair) PrincipalId() string {
	return k.Identity.Sender().String()
}

func (k *KeyPair) AccountId() string {
	hexPub := hex.EncodeToString(k.PublicKey)
	return hexPub
}

func NewKeyPairFromPrivateKey(secret string) (*KeyPair, error) {
	privBytes, err := hex.DecodeString(secret)
	if err != nil {
		return nil, err
	}
	privateKey := ed25519.PrivateKey(privBytes)
	publicKey, ok := privateKey.Public().(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to get public key from private key")
	}
	ed25519Identity, err := identity.NewEd25519Identity(publicKey, privateKey)
	if err != nil {
		return nil, err
	}
	return &KeyPair{PrivateKey: privateKey, PublicKey: publicKey, Identity: ed25519Identity}, nil
}

func NewKeyPair(seed string) (*KeyPair, error) {
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		return nil, err
	}
	privateKey := ed25519.NewKeyFromSeed(seedBytes)
	publicKey, ok := privateKey.Public().(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to get public key from private key")
	}
	ed25519Identity, err := identity.NewEd25519Identity(publicKey, privateKey)
	if err != nil {
		return nil, err
	}
	return &KeyPair{PrivateKey: privateKey, PublicKey: publicKey, Identity: ed25519Identity}, nil
}
