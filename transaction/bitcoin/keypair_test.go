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
	fmt.Printf("regtestMultiAddress: %v\n", address)
	address, err = MultiScriptAddress(2, TestNet, pubBytesList)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("testnetMultiAddress: %v\n", address)
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
