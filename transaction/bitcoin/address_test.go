package bitcoin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	privKey, err := btcec.NewPrivateKey()
	assert.NoError(t, err)
	fmt.Printf("priv: %v\n", privKey)

	pubKey := privKey.PubKey()
	fmt.Printf("pub: %v\n", hex.EncodeToString(pubKey.SerializeCompressed()))

	pkhAddr, _ := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), &chaincfg.MainNetParams)
	fmt.Printf("addr: %v\n", pkhAddr.EncodeAddress())
}

func TestPKHAddressFromBytes(t *testing.T) {
	secret, _ := hex.DecodeString("b29c3157e4a68b240ec821515fc77181c7a828259efbb3c1ab1df9b67d03c645")
	privKey, pubKey := btcec.PrivKeyFromBytes(secret)
	fmt.Printf("priv: %v\n", privKey)
	fmt.Printf("pub: %v\n", hex.EncodeToString(pubKey.SerializeCompressed()))

	pkhAddr, _ := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), &chaincfg.MainNetParams)
	fmt.Printf("addr: %v\n", pkhAddr.EncodeAddress())
}

func TestTestnetPKHAddressFromBytes(t *testing.T) {
	secret, _ := hex.DecodeString("b29c3157e4a68b240ec821515fc77181c7a828259efbb3c1ab1df9b67d03c645")
	privKey, pubKey := btcec.PrivKeyFromBytes(secret)
	fmt.Printf("priv: %v\n", privKey)
	fmt.Printf("pub: %v\n", hex.EncodeToString(pubKey.SerializeCompressed()))

	pkhAddr, _ := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), &chaincfg.TestNet3Params)
	fmt.Printf("addr: %v\n", pkhAddr.EncodeAddress())
}

func TestWPKHAddressFromBytes(t *testing.T) {
	secret, _ := hex.DecodeString("b29c3157e4a68b240ec821515fc77181c7a828259efbb3c1ab1df9b67d03c645")
	privKey, pubKey := btcec.PrivKeyFromBytes(secret)
	fmt.Printf("priv: %v\n", privKey)
	fmt.Printf("pub: %v\n", hex.EncodeToString(pubKey.SerializeCompressed()))

	wpkhAddr, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), &chaincfg.MainNetParams)
	fmt.Printf("addr: %v\n", wpkhAddr.EncodeAddress())
}

// tb1ql9azatvsw9cxydtu3s0wzhf76zjnynhasuy4zy
func TestTestnetWPKHAddressFromBytes(t *testing.T) {
	secret, _ := hex.DecodeString("b29c3157e4a68b240ec821515fc77181c7a828259efbb3c1ab1df9b67d03c645")
	privKey, pubKey := btcec.PrivKeyFromBytes(secret)
	fmt.Printf("priv: %v\n", privKey)
	fmt.Printf("pub: %v\n", hex.EncodeToString(pubKey.SerializeCompressed()))

	wpkhAddr, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), &chaincfg.TestNet3Params)
	fmt.Printf("addr: %v\n", wpkhAddr.EncodeAddress())
}

// 同时对于同一私钥产生的P2WPKH 和 P2PKH，他们的scriptAddress 一样
func TestSameScriptAddressForPKHAddressAndWPKHAddressFromBytes(t *testing.T) {
	secret, _ := hex.DecodeString("b29c3157e4a68b240ec821515fc77181c7a828259efbb3c1ab1df9b67d03c645")
	privKey, pubKey := btcec.PrivKeyFromBytes(secret)
	fmt.Printf("priv: %v\n", privKey)
	fmt.Printf("pub: %v\n", hex.EncodeToString(pubKey.SerializeCompressed()))

	wpkhAddr, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), &chaincfg.MainNetParams)
	fmt.Printf("addr: %v\n", wpkhAddr.EncodeAddress())

	pkhAddr, _ := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), &chaincfg.MainNetParams)
	fmt.Printf("addr: %v\n", pkhAddr.EncodeAddress())
	assert.Equal(t, wpkhAddr.ScriptAddress(), pkhAddr.ScriptAddress())
}

func TestSHAddressFromSerializedScript(t *testing.T) {
	serializedScript, _ := hex.DecodeString("123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0")
	shAddr, _ := btcutil.NewAddressScriptHash(serializedScript, &chaincfg.MainNetParams)

	fmt.Printf("addr: %v\n", shAddr.EncodeAddress())
}

func TestGenerateWSHAddress(t *testing.T) {
	serializedScript, _ := hex.DecodeString("123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0")
	scriptHash := sha256.Sum256(serializedScript)
	wshAddr, _ := btcutil.NewAddressWitnessScriptHash(scriptHash[:], &chaincfg.MainNetParams)

	fmt.Printf("addr: %v\n", wshAddr.EncodeAddress())
}

// test data from btcwallet/waddrmgr/manager_test.go/line 943
func TestGenerateWSHAddress1(t *testing.T) {
	serializedScript, _ := hex.DecodeString("52210305a662958b547fe25a71cd28fc7ef1c2" +
		"ad4a79b12f34fc40137824b88e61199d21038552c09d9" +
		"a709c8cbba6e472307d3f8383f46181895a76e01e258f" +
		"09033b4a7821029dd72aba87324af59508380f9564d34" +
		"b0f7b20d864d186e7d0428c9ea241c61653ae")
	scriptHash := sha256.Sum256(serializedScript)
	fmt.Printf("hash:%v\n", hex.EncodeToString(scriptHash[:]))
	wshAddr, _ := btcutil.NewAddressWitnessScriptHash(scriptHash[:], &chaincfg.RegressionNetParams)

	fmt.Printf("addr: %v\n", wshAddr.EncodeAddress())

	//52  //2
	//21  //33
	//0305a662958b547fe25a71cd28fc7ef1c2ad4a79b12f34fc40137824b88e61199d
	//21
	//038552c09d9a709c8cbba6e472307d3f8383f46181895a76e01e258f09033b4a78
	//21
	//029dd72aba87324af59508380f9564d34b0f7b20d864d186e7d0428c9ea241c616
	//53
	//ae

	//52
	//21
	//0305a662958b547fe25a71cd28fc7ef1c2ad4a79b12f34fc40137824b88e61199d
	//21
	//038552c09d9a709c8cbba6e472307d3f8383f46181895a76e01e258f09033b4a78
	//21
	//029dd72aba87324af59508380f9564d34b0f7b20d864d186e7d0428c9ea241c616
	//53
	//ae

}

func TestGenerateWSHAddress2(t *testing.T) {
	serializedScript, _ := hex.DecodeString("02c5389a31ce6149c28ba20d14db8540b2319e5a65000a2919fbf7a6296e7840b5")
	scriptHash := sha256.Sum256(serializedScript)
	fmt.Printf("hash:%v\n", hex.EncodeToString(scriptHash[:]))
	wshAddr, _ := btcutil.NewAddressWitnessScriptHash(scriptHash[:], &chaincfg.MainNetParams)

	fmt.Printf("addr: %v\n", wshAddr.EncodeAddress())
}

func TestGenerateWSHAddress3(t *testing.T) {
	serializedScript, _ := hex.DecodeString("02f7a01e30388dea9673db8cdb48b985441db785382efbcecc05abac079a630481")
	scriptHash := sha256.Sum256(serializedScript)
	fmt.Printf("hash:%v\n", hex.EncodeToString(scriptHash[:]))
	wshAddr, _ := btcutil.NewAddressWitnessScriptHash(scriptHash[:], &chaincfg.TestNet3Params)

	fmt.Printf("addr: %v\n", wshAddr.EncodeAddress())
}

func TestGenerateWSHAddress4(t *testing.T) {
	serializedScript, _ := hex.DecodeString("00201863143c14c5166804bd19203356da136c985678cd4d27a1b8c6329604903262")
	scriptHash := sha256.Sum256(serializedScript)
	fmt.Printf("hash:%v\n", hex.EncodeToString(scriptHash[:]))
	wshAddr, _ := btcutil.NewAddressWitnessScriptHash(scriptHash[:], &chaincfg.TestNet3Params)

	fmt.Printf("addr: %v\n", wshAddr.EncodeAddress())
}

// 于P2WSH 和 P2SH,因为前者没有使用ripemd160, 所以前者的scriptadress 是32bytes， 后者是20bytes
func TestDifferentScriptAddressForSHAddressAndWSHAddressFromBytes(t *testing.T) {
	serializedScript, _ := hex.DecodeString("123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0")
	shAddr, _ := btcutil.NewAddressScriptHash(serializedScript, &chaincfg.MainNetParams)
	fmt.Printf("addr: %v, scriptAddress:%v \n", shAddr.EncodeAddress(), hex.EncodeToString(shAddr.ScriptAddress()))

	scriptHash := sha256.Sum256(serializedScript)
	wshAddr, _ := btcutil.NewAddressWitnessScriptHash(scriptHash[:], &chaincfg.MainNetParams)
	fmt.Printf("addr: %v, scriptAddress:%v \n", wshAddr.EncodeAddress(), hex.EncodeToString(wshAddr.ScriptAddress()))
}

func TestGenerateMultiSignScriptAddress(t *testing.T) {
	scretes := []string{
		"23c9cdb2685d0905c0969dbbbfd27fdc1791e16e43b0352d9f11a89053d268ac",
		"47b38c30407286330562e228a73bf84f0c6d5d9593bd16b2dfc66ca1654ab83d",
		"968b40431da7f3aba9dfea20f0c9790ca38117d884ce47ef03d36829cfc48f49",
	}

	privKeys := []*btcec.PrivateKey{}
	pubKeys := []*btcec.PublicKey{}
	addrPubKeys := []*btcutil.AddressPubKey{}
	for _, secret := range scretes {
		s, _ := hex.DecodeString(secret)
		privKey, pubKey := btcec.PrivKeyFromBytes(s)
		privKeys = append(privKeys, privKey)
		pubKeys = append(pubKeys, pubKey)
		addrPubKey, _ := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), &chaincfg.MainNetParams)
		addrPubKeys = append(addrPubKeys, addrPubKey)
		fmt.Printf("scriptAddr:%v\n", hex.EncodeToString(addrPubKey.ScriptAddress()))
	}

	multiSigScript, _ := txscript.MultiSigScript(addrPubKeys, 2)
	fmt.Printf("script:%v\n", hex.EncodeToString(multiSigScript))

	scriptHash := sha256.Sum256(multiSigScript)
	wshAddr, _ := btcutil.NewAddressWitnessScriptHash(scriptHash[:], &chaincfg.MainNetParams)
	fmt.Printf("addr: %v, scriptAddress:%v \n", wshAddr.EncodeAddress(), hex.EncodeToString(wshAddr.ScriptAddress()))

	/*
		scriptAddr:028fa190883221d93c3ecd3d9a7c7afa130393d56826acc811b3d27834b4986f32
		scriptAddr:033e8d41a47d121a6a4ac4e05db8967b47ff3036507e7d95a6b912483bea9ab716
		scriptAddr:03d78e3a9b9b1b966b930e13acf2eb90eb9b9c87c044e6f05a49b6bc0c3d5c5a2b
		script:
			52  //OP_2, 将2压入
			21028fa190883221d93c3ecd3d9a7c7afa130393d56826acc811b3d27834b4986f32 //将33bytes长度的CompressedPubKey 压入
			21033e8d41a47d121a6a4ac4e05db8967b47ff3036507e7d95a6b912483bea9ab716
			2103d78e3a9b9b1b966b930e13acf2eb90eb9b9c87c044e6f05a49b6bc0c3d5c5a2b
			53 //OP_3
			ae  //OP_CHECKMULTISIG
	*/

}

func TestGenerateTestnetSegWitMultiSignScriptAddress(t *testing.T) {
	scretes := []string{
		"23c9cdb2685d0905c0969dbbbfd27fdc1791e16e43b0352d9f11a89053d268ac",
		"47b38c30407286330562e228a73bf84f0c6d5d9593bd16b2dfc66ca1654ab83d",
		"968b40431da7f3aba9dfea20f0c9790ca38117d884ce47ef03d36829cfc48f49",
	}

	privKeys := []*btcec.PrivateKey{}
	pubKeys := []*btcec.PublicKey{}
	addrPubKeys := []*btcutil.AddressPubKey{}
	for _, secret := range scretes {
		s, _ := hex.DecodeString(secret)
		privKey, pubKey := btcec.PrivKeyFromBytes(s)
		privKeys = append(privKeys, privKey)
		pubKeys = append(pubKeys, pubKey)
		addrPubKey, _ := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), &chaincfg.TestNet3Params)
		addrPubKeys = append(addrPubKeys, addrPubKey)
		fmt.Printf("scriptAddr:%v\n", hex.EncodeToString(addrPubKey.ScriptAddress()))
	}

	multiSigScript, _ := txscript.MultiSigScript(addrPubKeys, 2)
	fmt.Printf("script:%v\n", hex.EncodeToString(multiSigScript))

	scriptHash := sha256.Sum256(multiSigScript)
	wshAddr, _ := btcutil.NewAddressWitnessScriptHash(scriptHash[:], &chaincfg.TestNet3Params)
	fmt.Printf("addr: %v, scriptAddress:%v \n", wshAddr.EncodeAddress(), hex.EncodeToString(wshAddr.ScriptAddress()))

	/*
		scriptAddr:028fa190883221d93c3ecd3d9a7c7afa130393d56826acc811b3d27834b4986f32
		scriptAddr:033e8d41a47d121a6a4ac4e05db8967b47ff3036507e7d95a6b912483bea9ab716
		scriptAddr:03d78e3a9b9b1b966b930e13acf2eb90eb9b9c87c044e6f05a49b6bc0c3d5c5a2b
		script:
			52  //OP_2, 将2压入
			21028fa190883221d93c3ecd3d9a7c7afa130393d56826acc811b3d27834b4986f32 //将33bytes长度的CompressedPubKey 压入
			21033e8d41a47d121a6a4ac4e05db8967b47ff3036507e7d95a6b912483bea9ab716
			2103d78e3a9b9b1b966b930e13acf2eb90eb9b9c87c044e6f05a49b6bc0c3d5c5a2b
			53 //OP_3
			ae  //OP_CHECKMULTISIG
	*/

}

func TestGetScriptAddressFromWSHAddress(t *testing.T) {
	s := "bc1q9h35gj4xmzev04467kv57m4s2fhe2nr05lv6xpln7mmh2r0alp6qfstq84"
	expected, _ := hex.DecodeString("2de3444aa6d8b2c7d6baf5994f6eb0526f954c6fa7d9a307f3f6f7750dfdf874")
	addr, _ := btcutil.DecodeAddress(s, &chaincfg.MainNetParams)
	assert.Equal(t, expected, addr.ScriptAddress())
	assert.Equal(t, s, addr.EncodeAddress())
}
