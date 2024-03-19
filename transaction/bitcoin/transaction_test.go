package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/assert"
)

func TestCreateDepositTx_1(t *testing.T) {
	//secret, _, err := base58.CheckDecode("cPbyxXLqYqAjAHDtbvKq7ETd6BsQBbS643RLHH4u3k1YeVAXkAqR")
	//if err != nil {
	//	panic(err)
	//}
	secret, _ := hex.DecodeString("f9962336ca15bdd2acd61edfc6857fe733ef36a3c1380acf5f91c17347df93e5")
	privKey, pubKey := btcec.PrivKeyFromBytes(secret)
	//fmt.Printf("priv: %v\n", privKey)03d4c6fac559b9e8182288fde7d4e42d6050910c6b0fbcc6bf3ba261e4168ca2d1
	fmt.Printf("pub: %v\n", hex.EncodeToString(pubKey.SerializeCompressed()))
	netParams := &chaincfg.TestNet3Params
	from, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), netParams)
	fmt.Printf("from addr:%v, from sriptAddress:%v\n", from.EncodeAddress(), hex.EncodeToString(from.ScriptAddress()))

	// Create the transaction to redeem the fake transaction.
	to, _ := btcutil.DecodeAddress("bcrt1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq2aw6ze", netParams)
	fmt.Printf("to address:%v, to sriptAddress:%v\n", to.String(), hex.EncodeToString(to.ScriptAddress()))

	depositTx := wire.NewMsgTx(wire.TxVersion)

	hash, err := chainhash.NewHashFromStr("1210f7219953e0b3168ccf2f9b1cd9b815af076dac5ee9c1b6ec35373cc5ce4e")
	txIn := wire.NewTxIn(wire.NewOutPoint(hash, 1), nil, nil)
	depositTx.AddTxIn(txIn)

	//TxIn's pkScript and value
	txInPkScript, err := hex.DecodeString("0014d7fae4fbdc8bf6c86a08c7177c5d06683754ea71")
	txInValue := btcutil.Amount(1199999859)

	//
	txOutScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(1199999559, txOutScript)
	depositTx.AddTxOut(txOut)

	var buf bytes.Buffer
	depositTx.Serialize(&buf)
	fmt.Printf("before sign deposit TxHash: %v\n", depositTx.TxHash())
	fmt.Printf("before sign deposit WitnessHash: %v\n", depositTx.WitnessHash())
	fmt.Printf("before sign deposit: %v\n", hex.EncodeToString(buf.Bytes()))

	err = createTx(depositTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue}, []*btcec.PrivateKey{privKey}, netParams)
	assert.NoError(t, err)

	err = validateMsgTx(depositTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue})
	assert.NoError(t, err)

	buf.Reset()
	depositTx.Serialize(&buf)
	fmt.Printf("after sign deposit TxHash: %v\n", depositTx.TxHash())
	fmt.Printf("after sign deposit WitnessHash: %v\n", depositTx.WitnessHash())
	fmt.Printf("after sign deposit: %v\n", hex.EncodeToString(buf.Bytes()))
}

func TestCreateDepositTx_2(t *testing.T) {
	secret, _ := hex.DecodeString("b29c3157e4a68b240ec821515fc77181c7a828259efbb3c1ab1df9b67d03c645")
	privKey, pubKey := btcec.PrivKeyFromBytes(secret)
	from, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), &chaincfg.TestNet3Params)
	fmt.Printf("from addr:%v, from sriptAddress:%v\n", from.EncodeAddress(), hex.EncodeToString(from.ScriptAddress()))

	// Create the transaction to redeem the fake transaction.
	to, _ := btcutil.DecodeAddress("tb1q7yc8ncrxy6wsdlhvhd6gglpfatg07835uses5mpsc2rfv7zulhcqy0m979", &chaincfg.TestNet3Params)
	fmt.Printf("to address:%v, to sriptAddress:%v\n", to.String(), hex.EncodeToString(to.ScriptAddress()))

	depositTx := wire.NewMsgTx(2)

	hash, err := chainhash.NewHashFromStr("77014d81ac23bdb3f4646c29afe3e8803e291ea9f69e585aee331832c0a62581")
	txIn := wire.NewTxIn(wire.NewOutPoint(hash, 0), nil, nil)
	depositTx.AddTxIn(txIn)

	//TxIn's pkScript and value
	txInPkScript, err := hex.DecodeString("0014f97a2ead90717062357c8c1ee15d3ed0a5324efd")
	txInValue := btcutil.Amount(7000)

	//
	txOutScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(6500, txOutScript)
	depositTx.AddTxOut(txOut)

	err = createTx(depositTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue}, []*btcec.PrivateKey{privKey}, &chaincfg.TestNet3Params)
	assert.NoError(t, err)

	err = validateMsgTx(depositTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue})
	assert.NoError(t, err)

	var buf bytes.Buffer
	depositTx.Serialize(&buf)
	fmt.Printf("after sign deposit: %v\n", hex.EncodeToString(buf.Bytes()))
}

func TestCreateRedeemTx_1(t *testing.T) {
	expectedHash, _ := hex.DecodeString("f889a5df00f886ba0c932ab668a28ab8b9f60ef8332a628065a61561de515585")
	scretes := []string{
		"b26dbaab82d9ebd8f37c88bbe56e22bf9cb21150c96dfb35ece4b787d3710d3301",
		"62dd5835dc2ce7f4f40eea1b88c816043d288532c8bb91964adef9bc0f0b4b7201",
		"9ff573d948c80fa1a50da6f66229b4bede9ec3fb482dd126f58d3acfb4b2979801",
	}

	privKeys := []*btcec.PrivateKey{}
	pubKeys := []*btcec.PublicKey{}
	addrPubKeys := []*btcutil.AddressPubKey{}
	params := &chaincfg.TestNet3Params
	for _, secret := range scretes {
		s, _ := hex.DecodeString(secret)
		privKey, pubKey := btcec.PrivKeyFromBytes(s)
		privKeys = append(privKeys, privKey)
		pubKeys = append(pubKeys, pubKey)
		addrPubKey, _ := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), params)
		addrPubKeys = append(addrPubKeys, addrPubKey)
	}

	multiSigScript, _ := txscript.MultiSigScript(addrPubKeys, 2)
	fmt.Printf("multiSigScript: %v\n", hex.EncodeToString(multiSigScript))

	scriptHash := sha256.Sum256(multiSigScript)
	from, _ := btcutil.NewAddressWitnessScriptHash(scriptHash[:], params)
	fmt.Printf("from addr: %v\n", from.EncodeAddress())

	sc, _ := txscript.PayToAddrScript(from)
	fmt.Printf("script: %v\n", hex.EncodeToString(sc))

	to, _ := btcutil.DecodeAddress("bcrt1q6lawf77u30mvs6sgcuthchgxdqm4f6n3kvx4z5", params)

	reedeemTx := wire.NewMsgTx(2)
	hash1, err := chainhash.NewHashFromStr("72c5d616a7d5564d6330d691777db44427c48227111a5ea6d7c30ff3d6d40c74")
	txIn1 := wire.NewTxIn(wire.NewOutPoint(hash1, 1), nil, nil)
	reedeemTx.AddTxIn(txIn1)

	txInPkScript1, err := hex.DecodeString("0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58")
	txInValue1 := btcutil.Amount(1199999780)

	hash2, err := chainhash.NewHashFromStr("769c8b678d0ae7bee83058fdc5e0da38fe0ce7304312e88a8acb43a9b6cbaef4")
	txIn2 := wire.NewTxIn(wire.NewOutPoint(hash2, 0), nil, nil)
	reedeemTx.AddTxIn(txIn2)

	txInPkScript2, err := hex.DecodeString("0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58")
	txInValue2 := btcutil.Amount(1199999800)

	//TxOut
	txOutScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(2399999280, txOutScript)
	reedeemTx.AddTxOut(txOut)

	hashes, err := CalWitnessSigHash(reedeemTx, [][]byte{txInPkScript1, txInPkScript2}, []btcutil.Amount{txInValue1, txInValue2}, from, params, [][]byte{multiSigScript, multiSigScript})
	assert.NoError(t, err)
	fmt.Printf("hash:%v\n", hex.EncodeToString(hashes[0]))
	assert.Equal(t, expectedHash, hashes[0])

	for index, hash := range hashes {
		var sigs [][]byte
		for _, priv := range privKeys {
			sig := ecdsa.Sign(priv, hash)

			sigWithType := append(sig.Serialize(), byte(txscript.SigHashAll))
			sigs = append(sigs, sigWithType)
		}

		witnessScript, err := MergeMultiSignatures(2, multiSigScript, sigs)
		assert.NoError(t, err)

		reedeemTx.TxIn[index].Witness = witnessScript
	}

	var buf bytes.Buffer
	err = reedeemTx.Serialize(&buf)
	fmt.Printf("signed reedeem: %v\n", hex.EncodeToString(buf.Bytes()))
	//correct: 020000000001010270eb309e31d2a0b8ac505f297cd413501742f8c3a79f83c5cbd1d7cdd403a40000000000ffffffff01401f000000000000160014f97a2ead90717062357c8c1ee15d3ed0a5324efd04004830450221009a2ccd91d89bf37c556863f13ed939aed04694e34dc97e0ea9f1c35018e46d23022055e657a3d93ceb693a4983773d6907ffdc8325798ce977546e1c87f43a67bf5b014730440220765f46fcb6bc52d24ee6fe593661d414c26242aac6ec8c17e7b61c9e1d8fbacc02202c0ddd31048508f2b9ed73d81540055c868b2e78df14f93c137c0bf2baaa39e001695221028fa190883221d93c3ecd3d9a7c7afa130393d56826acc811b3d27834b4986f3221033e8d41a47d121a6a4ac4e05db8967b47ff3036507e7d95a6b912483bea9ab7162103d78e3a9b9b1b966b930e13acf2eb90eb9b9c87c044e6f05a49b6bc0c3d5c5a2b53ae00000000
	//err = validateMsgTx(reedeemTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue})
	assert.NoError(t, err)
}

func TestCreateRedeemTx_2(t *testing.T) {
	expectedHash, _ := hex.DecodeString("00a165032df678da63f0e024bff6f6593a3d3af2bb318a992fc5c88fb2bcf613")
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
	}

	multiSigScript, _ := txscript.MultiSigScript(addrPubKeys, 2)

	scriptHash := sha256.Sum256(multiSigScript)
	from, _ := btcutil.NewAddressWitnessScriptHash(scriptHash[:], &chaincfg.RegressionNetParams)

	to, _ := btcutil.DecodeAddress("tb1ql9azatvsw9cxydtu3s0wzhf76zjnynhasuy4zy", &chaincfg.TestNet3Params)

	reedeemTx := wire.NewMsgTx(2)
	hash, err := chainhash.NewHashFromStr("a0b391b03d17c3a07a65652b5807931bcbb31d63894b8fd46538fc50602948c3")
	txIn := wire.NewTxIn(wire.NewOutPoint(hash, 0), nil, nil)
	reedeemTx.AddTxIn(txIn)

	//TxIn's pkScript and value
	txInPkScript, err := hex.DecodeString("0020f13079e066269d06feecbb74847c29ead0ff1e34e4330a6c30c28696785cfdf0")
	txInValue := btcutil.Amount(6500)

	//TxOut
	txOutScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(6000, txOutScript)
	reedeemTx.AddTxOut(txOut)

	hashes, err := CalWitnessSigHash(reedeemTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue}, from, &chaincfg.TestNet3Params, [][]byte{multiSigScript})
	assert.NoError(t, err)
	fmt.Printf("hash:%v\n", hex.EncodeToString(hashes[0]))
	assert.Equal(t, expectedHash, hashes[0])

	var sigs [][]byte
	for _, priv := range privKeys {
		sig := ecdsa.Sign(priv, hashes[0])

		sigWithType := append(sig.Serialize(), byte(txscript.SigHashAll))
		sigs = append(sigs, sigWithType)
	}

	witnessScript, err := MergeMultiSignatures(2, multiSigScript, sigs)
	assert.NoError(t, err)

	reedeemTx.TxIn[0].Witness = witnessScript
	var buf bytes.Buffer
	err = reedeemTx.Serialize(&buf)
	//fmt.Printf("signed reedeem TxHash: %v\n", reedeemTx.TxHash())
	//fmt.Printf("signed reedeem WitnessHash: %v\n", reedeemTx.WitnessHash())
	fmt.Printf("signed reedeem: %v\n", hex.EncodeToString(buf.Bytes()))
	//correct: 020000000001010270eb309e31d2a0b8ac505f297cd413501742f8c3a79f83c5cbd1d7cdd403a40000000000ffffffff01401f000000000000160014f97a2ead90717062357c8c1ee15d3ed0a5324efd04004830450221009a2ccd91d89bf37c556863f13ed939aed04694e34dc97e0ea9f1c35018e46d23022055e657a3d93ceb693a4983773d6907ffdc8325798ce977546e1c87f43a67bf5b014730440220765f46fcb6bc52d24ee6fe593661d414c26242aac6ec8c17e7b61c9e1d8fbacc02202c0ddd31048508f2b9ed73d81540055c868b2e78df14f93c137c0bf2baaa39e001695221028fa190883221d93c3ecd3d9a7c7afa130393d56826acc811b3d27834b4986f3221033e8d41a47d121a6a4ac4e05db8967b47ff3036507e7d95a6b912483bea9ab7162103d78e3a9b9b1b966b930e13acf2eb90eb9b9c87c044e6f05a49b6bc0c3d5c5a2b53ae00000000
	err = validateMsgTx(reedeemTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue})
	assert.NoError(t, err)
}

func TestCreateRedeemTx_3(t *testing.T) {
	expectedHash, _ := hex.DecodeString("00a165032df678da63f0e024bff6f6593a3d3af2bb318a992fc5c88fb2bcf613")
	scretes := []string{
		"23c9cdb2685d0905c0969dbbbfd27fdc1791e16e43b0352d9f11a89053d268ac",
		"47b38c30407286330562e228a73bf84f0c6d5d9593bd16b2dfc66ca1654ab83d",
		"968b40431da7f3aba9dfea20f0c9790ca38117d884ce47ef03d36829cfc48f49",
	}
	privKeys := []*btcec.PrivateKey{}
	pubKeys := []*btcec.PublicKey{}
	addrPubKeys := []*btcutil.AddressPubKey{}
	netParams := &chaincfg.RegressionNetParams
	for _, secret := range scretes {
		s, _ := hex.DecodeString(secret)
		privKey, pubKey := btcec.PrivKeyFromBytes(s)
		privKeys = append(privKeys, privKey)
		pubKeys = append(pubKeys, pubKey)
		addrPubKey, _ := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), netParams)
		addrPubKeys = append(addrPubKeys, addrPubKey)
	}

	multiSigScript, _ := txscript.MultiSigScript(addrPubKeys, 2)

	scriptHash := sha256.Sum256(multiSigScript)
	from, _ := btcutil.NewAddressWitnessScriptHash(scriptHash[:], netParams)

	to, _ := btcutil.DecodeAddress("tb1ql9azatvsw9cxydtu3s0wzhf76zjnynhasuy4zy", netParams)

	reedeemTx := wire.NewMsgTx(2)
	hash, err := chainhash.NewHashFromStr("a0b391b03d17c3a07a65652b5807931bcbb31d63894b8fd46538fc50602948c3")
	txIn := wire.NewTxIn(wire.NewOutPoint(hash, 0), nil, nil)
	reedeemTx.AddTxIn(txIn)

	//TxIn's pkScript and value
	txInPkScript, err := hex.DecodeString("0020f13079e066269d06feecbb74847c29ead0ff1e34e4330a6c30c28696785cfdf0")
	txInValue := btcutil.Amount(6500)

	//TxOut
	txOutScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(6000, txOutScript)
	reedeemTx.AddTxOut(txOut)

	hashes, err := CalWitnessSigHash(reedeemTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue}, from, netParams, [][]byte{multiSigScript})
	assert.NoError(t, err)
	fmt.Printf("hash:%v\n", hex.EncodeToString(hashes[0]))
	assert.Equal(t, expectedHash, hashes[0])

	var sigs [][]byte
	for _, priv := range privKeys {
		sig := ecdsa.Sign(priv, hashes[0])

		sigWithType := append(sig.Serialize(), byte(txscript.SigHashAll))
		sigs = append(sigs, sigWithType)
	}

	witnessScript, err := MergeMultiSignatures(2, multiSigScript, sigs)
	assert.NoError(t, err)

	reedeemTx.TxIn[0].Witness = witnessScript
	var buf bytes.Buffer
	err = reedeemTx.Serialize(&buf)
	//fmt.Printf("signed reedeem TxHash: %v\n", reedeemTx.TxHash())
	//fmt.Printf("signed reedeem WitnessHash: %v\n", reedeemTx.WitnessHash())
	fmt.Printf("signed reedeem: %v\n", hex.EncodeToString(buf.Bytes()))
	//correct: 020000000001010270eb309e31d2a0b8ac505f297cd413501742f8c3a79f83c5cbd1d7cdd403a40000000000ffffffff01401f000000000000160014f97a2ead90717062357c8c1ee15d3ed0a5324efd04004830450221009a2ccd91d89bf37c556863f13ed939aed04694e34dc97e0ea9f1c35018e46d23022055e657a3d93ceb693a4983773d6907ffdc8325798ce977546e1c87f43a67bf5b014730440220765f46fcb6bc52d24ee6fe593661d414c26242aac6ec8c17e7b61c9e1d8fbacc02202c0ddd31048508f2b9ed73d81540055c868b2e78df14f93c137c0bf2baaa39e001695221028fa190883221d93c3ecd3d9a7c7afa130393d56826acc811b3d27834b4986f3221033e8d41a47d121a6a4ac4e05db8967b47ff3036507e7d95a6b912483bea9ab7162103d78e3a9b9b1b966b930e13acf2eb90eb9b9c87c044e6f05a49b6bc0c3d5c5a2b53ae00000000
	err = validateMsgTx(reedeemTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue})
	assert.NoError(t, err)

}

func TestCalculateHash_1(t *testing.T) {
	reedeemTx := wire.NewMsgTx(2)
	expectedHash, _ := hex.DecodeString("f889a5df00f886ba0c932ab668a28ab8b9f60ef8332a628065a61561de515585")

	from, _ := btcutil.DecodeAddress("tb1q7yc8ncrxy6wsdlhvhd6gglpfatg07835uses5mpsc2rfv7zulhcqy0m979", &chaincfg.TestNet3Params)
	hash, err := chainhash.NewHashFromStr("a403d4cdd7d1cbc5839fa7c3f842175013d47c295f50acb8a0d2319e30eb7002")
	assert.NoError(t, err)
	txIn := wire.NewTxIn(wire.NewOutPoint(hash, 0), nil, nil)
	reedeemTx.AddTxIn(txIn)

	//TxIn's pkScript and value. get from
	txInOriginalScript, err := hex.DecodeString("5221028fa190883221d93c3ecd3d9a7c7afa130393d56826acc811b3d27834b4986f3221033e8d41a47d121a6a4ac4e05db8967b47ff3036507e7d95a6b912483bea9ab7162103d78e3a9b9b1b966b930e13acf2eb90eb9b9c87c044e6f05a49b6bc0c3d5c5a2b53ae")
	txInPkScript, err := hex.DecodeString("0020f13079e066269d06feecbb74847c29ead0ff1e34e4330a6c30c28696785cfdf0")
	txInValue := btcutil.Amount(9000)

	to, _ := btcutil.DecodeAddress("tb1ql9azatvsw9cxydtu3s0wzhf76zjnynhasuy4zy", &chaincfg.TestNet3Params)
	txOutScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	txOut := wire.NewTxOut(8000, txOutScript)
	reedeemTx.AddTxOut(txOut)

	hashes, err := CalWitnessSigHash(reedeemTx, [][]byte{txInPkScript}, []btcutil.Amount{txInValue}, from, &chaincfg.TestNet3Params, [][]byte{txInOriginalScript})
	assert.Equal(t, expectedHash, hashes[0])
}

func TestCalculateHash_2(t *testing.T) {
	expectedHash, err := hex.DecodeString("f889a5df00f886ba0c932ab668a28ab8b9f60ef8332a628065a61561de515585")
	assert.NoError(t, err)
	raw, _ := hex.DecodeString("02000000081aba67bc0e257c733f0c00dab79853b37e19f617e9a63b271e9e617a16377a3bb13029ce7b1f559ef5e747fcac439f1455a2ec7c5f09b72290795e706650440270eb309e31d2a0b8ac505f297cd413501742f8c3a79f83c5cbd1d7cdd403a400000000695221028fa190883221d93c3ecd3d9a7c7afa130393d56826acc811b3d27834b4986f3221033e8d41a47d121a6a4ac4e05db8967b47ff3036507e7d95a6b912483bea9ab7162103d78e3a9b9b1b966b930e13acf2eb90eb9b9c87c044e6f05a49b6bc0c3d5c5a2b53ae2823000000000000ffffffffedfe2b7591d80ff513a432a6fdb4ae8a410469dd3a700e3749bb05d25251207a0000000001000000")
	hash := chainhash.DoubleHashB(raw)
	assert.Equal(t, expectedHash, hash)
}
