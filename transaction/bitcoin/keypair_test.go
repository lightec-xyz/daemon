package bitcoin

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {

}

func TestMultiAddress(t *testing.T) {
	pubk1, _ := hex.DecodeString("02600eb007626c38bfa11aae6e8271a25ec2f272d21960be0cabbfe8b18768fa95")
	pubk2, _ := hex.DecodeString("032bdb4a665e7b8f03d4b4dc33424fe4e3272ec6c15b2af84bf3b8887fc24a0026")
	pubk3, _ := hex.DecodeString("0296808ee7ab89c9982c9ade0e7a41afd350c9ade7d653af3880336f60d36a3f2e")
	var pubBytesList [][]byte
	pubBytesList = append(pubBytesList, pubk1)
	pubBytesList = append(pubBytesList, pubk2)
	pubBytesList = append(pubBytesList, pubk3)
	address, err := MultiScriptAddress(2, RegTest, pubBytesList)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("regtestMultiAddress: %v\n", address)
	addrScript, err := GenPayToAddrScript(address, RegTest)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("regtestMultiAddressScript: %v\n", addrScript)
	address, err = MultiScriptAddress(2, TestNet, pubBytesList)
	if err != nil {
		t.Fatal(err)
	}
	addrScript, err = GenPayToAddrScript(address, TestNet)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("testnetMultiAddress: %v\n", address)
	t.Logf("testnetMultiAddressScript: %v\n", addrScript)
}

func TestGenerateKey(t *testing.T) {
	secrets := []string{
		"cPbyxXLqYqAjAHDtbvKq7ETd6BsQBbS643RLHH4u3k1YeVAXkAqR",
		"cTZYUrhkQmpKGf6QjNhiLrw7g52VL449Tgwzo9UmKbnzphfZUUZp",
		"cQtt5owpfyv79fGpnv2yncKtQTU42rzohXGCSD3Mcn4CGgEkPVVE",
		"cSwe7Np3o7eCec6hgKFrwqGs9bb6x2dubKBucLYqQ6mJu5JH1aCn",
	}
	for _, secret := range secrets {
		privBytes, _, err := base58.CheckDecode(secret)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("secret:         %x\n", secret)
		privatekey := fmt.Sprintf("%x", privBytes)
		fmt.Printf("privateKey:     %v\n", privatekey)
		keyPair, err := NewKeyPairFromSecret(privatekey)
		if err != nil {
			t.Fatal(err)
		}
		publicKey := keyPair.PublicKey()
		fmt.Printf("publicKey:      %x\n", publicKey)
		regtesAddress, lockScript, err := keyPair.Address(P2WPKH, RegTest)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("lockScript:     %v\n", lockScript)
		fmt.Printf("regtestAddress: %v\n", regtesAddress)
		testnetAddress, _, err := keyPair.Address(P2WPKH, TestNet)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("testnetAddress: %v\n", testnetAddress)
		fmt.Printf("*********************************\n")
	}

}
