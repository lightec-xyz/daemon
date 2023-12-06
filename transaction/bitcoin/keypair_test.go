package bitcoin

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/base58"
	"testing"
)

func TestKeyPair(t *testing.T) {
	keyPair, err := NewRandSeed()
	if err != nil {
		t.Fatal(err)
	}
	privateKey := keyPair.PrivateKey()
	t.Log(hex.EncodeToString(privateKey))
	publicKey := keyPair.PublicKey()
	t.Log(hex.EncodeToString(publicKey))
	address, err := keyPair.Address(P2WPKH, RegTest)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(address)
	msg := []byte("hello")
	signatrue := keyPair.Sign(msg)
	t.Log(hex.EncodeToString(signatrue))
	verify, err := keyPair.Verify(msg, signatrue)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(verify)
}

func TestMultiAddress(t *testing.T) {
	//bcrt1q7yc8ncrxy6wsdlhvhd6gglpfatg07835uses5mpsc2rfv7zulhcqfk3rtl
	scretes := []string{
		"23c9cdb2685d0905c0969dbbbfd27fdc1791e16e43b0352d9f11a89053d268ac",
		"47b38c30407286330562e228a73bf84f0c6d5d9593bd16b2dfc66ca1654ab83d",
		"968b40431da7f3aba9dfea20f0c9790ca38117d884ce47ef03d36829cfc48f49",
	}
	var pubBytesList [][]byte
	for _, secret := range scretes {
		s, _ := hex.DecodeString(secret)
		key, _ := btcec.PrivKeyFromBytes(s)
		pubBytesList = append(pubBytesList, key.PubKey().SerializeCompressed())
	}
	address, err := MultiScriptAddress(2, RegTest, pubBytesList)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(address)

}

func TestMultiAddress01(t *testing.T) {
	pubk1, _ := hex.DecodeString("03bd96c4d06aa773e5d282f0b6bccd1fb91268484918648ccda1ae768209edb050")
	pubk2, _ := hex.DecodeString("03aa9c4245340a02864c903f7f9e7bc9ef1cc374093aacbf72b614002f6d8c8c22")
	pubk3, _ := hex.DecodeString("03351a7971bf7ed886fca99aebdc3b195fc79ffe93b499e2309a4e69ab115405e0")
	var pubBytesList [][]byte
	pubBytesList = append(pubBytesList, pubk1)
	pubBytesList = append(pubBytesList, pubk2)
	pubBytesList = append(pubBytesList, pubk3)
	address, err := MultiScriptAddress(2, RegTest, pubBytesList)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(address)
}

func TestGenerateKey(t *testing.T) {
	//03d4c6fac559b9e8182288fde7d4e42d6050910c6b0fbcc6bf3ba261e4168ca2d1
	//03d4c6fac559b9e8182288fde7d4e42d6050910c6b0fbcc6bf3ba261e4168ca2d1
	result, _, err := base58.CheckDecode("cSwe7Np3o7eCec6hgKFrwqGs9bb6x2dubKBucLYqQ6mJu5JH1aCn")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%x\n", result)
	keyPair, err := NewKeyPairFromSecret(fmt.Sprintf("%x", result))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%x\n", keyPair.PrivateKey())
	fmt.Printf("%x\n", keyPair.PublicKey())
}

func TestNewKeyPairFromSecret(t *testing.T) {
	keyPair, err := NewKeyPairFromSecret("3c5579d538347d56ed5ef6d56e7e36ae453dcbd6ff6586783f82c17c8190d716")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%x\n", keyPair.PrivateKey())
	fmt.Printf("%x\n", keyPair.PublicKey())
	address, err := keyPair.Address(P2WPKH, RegTest)
	if err != nil {
		t.Fatal(err)
	}
	//bcrt1q6lawf77u30mvs6sgcuthchgxdqm4f6n3kvx4z5
	fmt.Printf("%v\n", address)

}
