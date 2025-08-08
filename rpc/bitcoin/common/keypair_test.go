package bitcoin

import (
	"fmt"
	ethCommon "github.com/ethereum/go-ethereum/common"
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
	pubk1 := ethCommon.FromHex("0x02e510903d107b5594a5ee854eecb2858aeed5e84838c36fb12041676f71a17eba")
	pubk2 := ethCommon.FromHex("0x03ea366ed0cfa0f48ee1e40ae321dab882f017fc8e7cb6a2159ced6fc42c6746da")
	pubk3 := ethCommon.FromHex("0218e65b47da2b63543f5b40c4b98532a97c737fa39c7e18d117bd3351eabbdc6c")
	//0x522102e510903d107b5594a5ee854eecb2858aeed5e84838c36fb12041676f71a17eba2103ea366ed0cfa0f48ee1e40ae321dab882f017fc8e7cb6a2159ced6fc42c6746da210218e65b47da2b63543f5b40c4b98532a97c737fa39c7e18d117bd3351eabbdc6c53ae
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
	t.Logf("regtest MultiSig EthAddress: %v\n", address)

	addrScript, err := GenPayToAddrScript(address, RegTest)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("regtest MultiSig Lock Script: %v\n", addrScript)

	address, err = MultiScriptAddress(2, TestNet, pubBytesList)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("testnet MultiSig EthAddress: %v\n", address)

	addrScript, err = GenPayToAddrScript(address, TestNet)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("testnet MultiSig Lock Script: %v\n", addrScript)
}

func TestGenPayToAddrScript(t *testing.T) {
	lock, err := GenPayToAddrScript("tb1qn9fpljh5ggp407z02jx8x76pemzclgd6rla0qp", TestNet)
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
