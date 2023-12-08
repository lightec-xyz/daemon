package ethereum

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"testing"
)

func TestGenerateAddress(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("privateKey:  %x \n", privateKey.D.Bytes())
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("fail")
	}
	compressPubkey := crypto.CompressPubkey(publicKeyECDSA)
	fmt.Printf("Uncompressed:%x \n", crypto.FromECDSAPub(publicKeyECDSA))
	fmt.Printf("Compressed:  %x \n", compressPubkey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("address:     %v \n", address)

}
