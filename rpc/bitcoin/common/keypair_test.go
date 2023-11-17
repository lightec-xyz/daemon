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
	pubk1, _ := hexutil.Decode("0x0363f549d250342df02ee8b51ad6c9148dabc587c6569761ab58aa68488bd2e2c5")
	pubk2, _ := hexutil.Decode("0x031cbb294f9955d80f65d9499feaeb5cb29d44c070adddd75cd48a40791d39b971")
	pubk3, _ := hexutil.Decode("0x035c54e8287a7f7ba31886249fc89f295a4cb74cebf0d925f1eafe87f22fba57f9")

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
