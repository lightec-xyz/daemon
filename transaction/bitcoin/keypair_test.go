package bitcoin

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
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
