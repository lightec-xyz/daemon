package bitcoin

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestGenerateKeyPair(t *testing.T) {

}

func TestMultiAddress(t *testing.T) {
	pubk1, _ := hexutil.Decode("0x03930b4b1bbe4b98128262f1ca2ac431f726b8390962eb581a6056601093663702")
	pubk2, _ := hexutil.Decode("0x03c9745bde1ed2f977704638c4bbe34ba764347d78340641341e2f9c1602bfecfc")
	pubk3, _ := hexutil.Decode("0x031908d0ff51de51205b102e28430b2e81d5d5a36758b43cdc9cf497230983e8e8")

	var pubBytesList [][]byte
	pubBytesList = append(pubBytesList, pubk1)
	pubBytesList = append(pubBytesList, pubk2)
	pubBytesList = append(pubBytesList, pubk3)

	var addrPubKeys []*btcutil.AddressPubKey
	for _, pubKey := range pubBytesList {
		addressPubKey, err := btcutil.NewAddressPubKey(pubKey, &chaincfg.TestNet3Params)
		if err != nil {
			t.Fatal(err)
		}
		addrPubKeys = append(addrPubKeys, addressPubKey)
	}

	multiSigScript, err := txscript.MultiSigScript(addrPubKeys, 2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("testnet MultiSig Script: %v\n", hexutil.Encode(multiSigScript))

	address, err := MultiScriptAddress(2, RegTest, pubBytesList)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("regtest MultiSig Address: %v\n", address)

	addrScript, err := GenPayToAddrScript(address, RegTest)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("regtest MultiSig Lock Script: %v\n", addrScript)

	address, err = MultiScriptAddress(2, TestNet, pubBytesList)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("testnet MultiSig Address: %v\n", address)

	addrScript, err = GenPayToAddrScript(address, TestNet)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("testnet MultiSig Lock Script: %v\n", addrScript)
}

func TestGenPayToAddrScript(t *testing.T) {
	lock, err := GenPayToAddrScript("tb1qvn2x35f0vy543q4slrmrce943c40qk8y0snkj7", TestNet)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(lock)
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
