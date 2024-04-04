package bitcoin

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lightec-xyz/daemon/rpc/bitcoin/types"
	ethrpc "github.com/lightec-xyz/daemon/rpc/ethereum"
	"github.com/lightec-xyz/daemon/rpc/oasis"
	"github.com/lightec-xyz/daemon/transaction/bitcoin"
	"github.com/lightec-xyz/daemon/transaction/ethereum"
	"github.com/stretchr/testify/require"
)

var client *Client
var err error

func init() {
	url := "https://go.getblock.io/d54c59f635654cc082de1f3fd14e5d02"
	//url := "http://127.0.0.1:8332"
	user := "lightec"
	pwd := "abcd1234"
	network := "regtest"
	client, err = NewClient(url, user, pwd, network)
	if err != nil {
		panic(err)
	}
}

func TestClient_GetBlockHeader(t *testing.T) {
	header, err := client.GetBlockHeader("6dd9258faef568db4816b9cd61bd0b920772bc5f3e440f1f6c12baae64f44144")
	if err != nil {
		panic(err)
	}
	fmt.Println(header)
}

func TestClient_GetBlockCount1(t *testing.T) {
	blockCount, err := client.GetBlockCount()
	if err != nil {
		panic(err)
	}
	fmt.Println(blockCount)
}

func TestClient_GetBlockHash(t *testing.T) {
	hash, err := client.GetBlockHash(2585162)
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
}

func TestClient_GetBlockTx(t *testing.T) {
	hash, err := client.GetBlockHash(2585162)
	if err != nil {
		panic(err)
	}
	blockWithTx, err := client.GetBlock(hash)
	if err != nil {
		panic(err)
	}
	fmt.Println(blockWithTx.Hash)
}

func TestMultiTransactionBuildTxData(t *testing.T) {
	publicKeys := []string{
		"0x0377e958f7a5636e92375dce8fa9d35ed4397b1d25eaa76bdc4c2f0b49ec0e0efe",
		"0x028b4f7f78afe170a8c3896997cd3780a9367c6d653772687bce54bb28f35a28af",
		"0x02765e2e1e204f6b0894b193e2a80768f8e0fd8f2c5a751e38b5955b1df7d00a13",
	}
	var pubKeys [][]byte
	for _, pub := range publicKeys {
		pubKey, err := hexutil.Decode(pub)
		if err != nil {
			t.Fatal(err)
		}
		pubKeys = append(pubKeys, pubKey)
	}

	txInputs := []bitcoin.TxIn{
		{
			Hash:     "6b9efd78750765f65e8733ef7a5869c4ca8c68ef35e1a8177a638e54786faa3f",
			VOut:     0,
			PkScript: "0020d3eaa044dbfccbbf98d192a35d456694991bb037d3c4915e4aedf26d0f12dcd2",
			Amount:   20000,
		},
	}

	txOutputs := []bitcoin.TxOut{
		{
			Address: "tb1qvn2x35f0vy543q4slrmrce943c40qk8y0snkj7",
			Amount:  19000,
		},
	}

	builder := bitcoin.NewMultiTransactionBuilder()
	err = builder.NetParams(bitcoin.RegTest)
	if err != nil {
		t.Fatal(err)
	}
	err = builder.AddMultiPublicKey(pubKeys, 2)
	if err != nil {
		t.Fatal(err)
	}
	err = builder.AddTxIn(txInputs)
	if err != nil {
		t.Fatal(err)
	}
	err = builder.AddTxOut(txOutputs)
	if err != nil {
		t.Fatal(err)
	}
	txData, err := builder.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("txData: %x \n", txData)

}

func Test_MergeAndSendTx(t *testing.T) {
	// https://holesky.etherscan.io/tx/0x3db1bb46352898a1ff0349274d0dcc7c8e78020ab2268c2bfa0863ab0e9de001
	txRaw, err := hexutil.Decode("0x0200000001b095e7df80a9e758ae67fa9c7d9b8c5464dc41249fffe9b164d4be56b404dfcc0000000000ffffffff02802e00000000000016001464d468d12f61295882b0f8f63c64b58e2af058e4401d0000000000002200207c907704071d036924b69db3b98b683cc405384828b55b1b3d25ffd8d04381bf00000000")
	require.NoError(t, err)

	transaction := bitcoin.NewTestnetMultiTxBuilder()
	err = transaction.Deserialize(txRaw)
	require.NoError(t, err)

	multiSigScript, err := hexutil.Decode("0x522103930b4b1bbe4b98128262f1ca2ac431f726b8390962eb581a60566010936637022103c9745bde1ed2f977704638c4bbe34ba764347d78340641341e2f9c1602bfecfc21031908d0ff51de51205b102e28430b2e81d5d5a36758b43cdc9cf497230983e8e853ae")
	require.NoError(t, err)
	transaction.AddMultiScript(multiSigScript, 2, 3)

	const nKey = 3
	const nTxin = 1
	sigs := make([][][]byte, nKey)
	for i := 0; i < nKey; i++ {
		sigs[i] = make([][]byte, nTxin)
	}

	sigs[0][0], err = hex.DecodeString("30450221009f293e400073c2eafccfe6869d2b3a5bca8991f9dab9169ea71b4ddef0b92f770220200d4c1c616fe25323082855954b72d2f5f750f829c37e27b57c4a9b265ed24b01")
	require.NoError(t, err)

	sigs[1][0], err = hex.DecodeString("3045022100d24fd371d6e2cfc8d5c0f83877a15accae81402b1cf576a00efb1f9dcfdbb72102201974c8a21f53a9194a591140ece5130bade30e9422f648e3f2451494788e3aef01")
	require.NoError(t, err)

	sigs[2][0], err = hex.DecodeString("3045022100c35f6795a89d4639d89f335fcd70e87ad0d64e3b8110f7fc3fe27f77f4f7f97502205bb7ecbb742e61c5193fb304a417e9b0c195eabfa2ff14c8688f90c83829595e01")
	require.NoError(t, err)

	err = transaction.MergeSignature(sigs[1:])
	require.NoError(t, err)

	tx, err := transaction.Serialize()
	require.NoError(t, err)

	txHex := hex.EncodeToString(tx)
	fmt.Printf("tx: %v\n", txHex)

	hash := transaction.TxHash()
	fmt.Printf("hash: %v\n", hash)

	txHash, err := client.Sendrawtransaction(txHex)
	require.NoError(t, err)
	fmt.Println(txHash)
}

func Test_GetMultiSigScriptRelatedsFromOasis(t *testing.T) {
	btcSignerContract := "0x99e514Dc90f4Dd36850C893bec2AdC9521caF8BB"
	oasisClient, err := oasis.NewClient("https://testnet.sapphire.oasis.io", btcSignerContract)
	require.NoError(t, err)

	publicKeys, err := oasisClient.PublicKey()
	require.NoError(t, err)

	fmt.Println("publicKeys: ")
	for _, publicKey := range publicKeys {
		fmt.Printf("    %v\n", hexutil.Encode(publicKey))
	}

	multiSigScript, walletAddr, lockScript, err :=
		bitcoin.GetMultiSigScriptRelateds(2, &chaincfg.TestNet3Params, publicKeys)
	require.NoError(t, err)

	fmt.Printf("MultiSig Script: %v\n", hexutil.Encode(multiSigScript))
	fmt.Printf("Wallet Address: %v\n", walletAddr.EncodeAddress())
	fmt.Printf("Lock Script: %v\n", hexutil.Encode(lockScript))
}

func TestMultiTransactionSignFromOasis(t *testing.T) {
	// https://holesky.etherscan.io/tx/0xe88fa618634a210f2b2a5c32393d75c44474057a77fc2062288c078b005a4a12
	ethTxHash := common.HexToHash("0x48262094098497cbbe7246487601c141d0427aae294505b0a949419f649c699b")

	zkBridgeAddr, zkBtcAddr := "0x4f413972ab2d4e53714a479db1519ce0e89ea30c", "0xbf3041e37be70a58920a6fd776662b50323021c9"
	ec, err := ethrpc.NewClient("https://1rpc.io/holesky", zkBridgeAddr, zkBtcAddr)
	require.NoError(t, err)

	ethTx, _, err := ec.TransactionByHash(context.Background(), ethTxHash)
	require.NoError(t, err)

	receipt, err := ec.TransactionReceipt(context.Background(), ethTxHash)
	require.NoError(t, err)

	btcRawTx, _, err := ethereum.DecodeRedeemLog(receipt.Logs[3].Data)
	require.NoError(t, err)
	fmt.Printf("btcRawTx: %v\n", hexutil.Encode(btcRawTx))

	rawTx, rawReceipt := ethereum.GetRawTxAndReceipt(ethTx, receipt)
	fmt.Printf("rawTx: %v\n", hexutil.Encode(rawTx))
	fmt.Printf("rawReceipt: %v\n", hexutil.Encode(rawReceipt))

	proofData, err := hexutil.Decode("0x233fd9fd21e614b801719516313a177c21222c28e1ab9b9bd2ed6db4d6f3b92823be74df686b218c75f8f880cdc5f73555d155d027a4ac74a259c23437652e571c0c35956868e4369bc057e61247cefaac3e2eb94a17f9d450f03d80dd32cd780e4fe637196404a0fd627321428a8ce2e3cb1180b15c8c1d792b99ac9140059a0a216a23f068afff3c128210f761d1e7cb607708ce78089085c57ad7f84a50411309fcd392aa0c8c9f6c7397b092b74308d3fb4548d410a7466427cc21df9ccf22c0c09a663e7467cb39541c80c710df0de9b0ff5dd6de21736d7dd0397b2e572e973a19be35130a1a609e1dfc31f1c6ea28ba726dd030dd64d192d68f737c4813e96f5e44c5e765cb9c1f3cc67a3fa5579f678b5cfee24441049336c54cf915155892072acceb10137b31fb5c814abe95501c026a124d28f5271996a51013620511de4dd95d4bf40198421e9f763057ea12759de557e767447b85ac925d4f9012a220765cb87f07f91bce290761d5ef844a097d4f02a4147b3228e6ecf6ab42204988fb87256e71c9f8568ee35c47f3324e0df0df38c306c0f689278fd9cb1824313cbfe50fa627ce6bce8cc93cc7a55fbb28bc54f3bdf7a7193b9ec921611c107d019f0eb237195a9809db1c40ddd33898fdc0969f63650fec75eca891a5111ceb4dcfd61c4f5728d21e9fea1b5a06e1fe9e572f4920433f1f53dbaed96de0094f921d4db89fd3dd439bfb9e40c505b12772b0969467098e56ac69aab018f8162d0617d2f0c141da8e50579d13fc66d24476f5901c513286704d2a1abe010c27248ba61827bc5632d406db69c7c297a5fadd889aaa7056b49c3743e70622b21bb0e9af0696ccc99697635022b3c2a902b78ff6cfea5c50424a6d89443fabfe21ce37b2e681be73fe20e675a2cd41d03295d3486d5c2143a2f6ea8032793cb82c97a09950a2013e46251bdfe0eeb62ac4ffca710234bd0f02dbf4ba44d211f11c72c19b8e5727d0b1e95fd69ef1fb6441ba6b6eb9771d71e21f6b7ae377c86503674893d4034ebd83e611dceead450c9b646d01bc9f59873b17da46834db4d11f23bb03da1d90ba41bc50367ba0ae47a1d8e42a61ec647577c684e5ac53d977287cc37d5a8d8489188b4ba0659d0b712ba682256338a8b1876460fe5ee743f405db2a4b48334b29a11429b4e56b2c559fc3dd9ae361ff37d3ec52a432ff503617add429e949fe248fd2b1b4b6af3e12aee4a58088d35eb6b77c4eb9ffbb944011c041500303ee2d6846c29b0e08bd5ce8c5c76d8181c4ebf638dbf9d6e417cd")
	require.NoError(t, err)

	btcSignerContract := "0x99e514Dc90f4Dd36850C893bec2AdC9521caF8BB"
	oasisClient, err := oasis.NewClient("https://testnet.sapphire.oasis.io", btcSignerContract)
	require.NoError(t, err)

	sigs, err := oasisClient.SignBtcTx(rawTx, rawReceipt, proofData)
	require.NoError(t, err)
	fmt.Printf("btx sigs: %x\n", sigs)

	transaction := bitcoin.NewMultiTransactionBuilder()
	err = transaction.Deserialize(btcRawTx)
	require.NoError(t, err)

	multiSigScript, err := ec.GetMultiSigScript()
	require.NoError(t, err)
	fmt.Printf("multiSigScript: %v\n", hexutil.Encode(multiSigScript))

	nTotal, nRequred := 3, 2
	transaction.AddMultiScript(multiSigScript, nRequred, nTotal)

	err = transaction.MergeSignature(sigs[:nRequred])
	require.NoError(t, err)

	btxTx, err := transaction.Serialize()
	require.NoError(t, err)

	txHex := hex.EncodeToString(btxTx)
	fmt.Printf("btx Tx: %v\n", txHex)

	txHash, err := client.Sendrawtransaction(txHex)
	require.NoError(t, err)
	fmt.Printf("btc hash: %v\n", txHash)
}

func TestMultiTransactionBuilder(t *testing.T) {
	secrerts := []string{
		"b26dbaab82d9ebd8f37c88bbe56e22bf9cb21150c96dfb35ece4b787d3710d3301",
		"62dd5835dc2ce7f4f40eea1b88c816043d288532c8bb91964adef9bc0f0b4b7201",
		"9ff573d948c80fa1a50da6f66229b4bede9ec3fb482dd126f58d3acfb4b2979801",
	}

	var privateKeys []*btcec.PrivateKey
	var pubkeylist [][]byte
	for _, secret := range secrerts {
		hexPriv, err := hex.DecodeString(secret)
		if err != nil {
			t.Fatal(err)
		}
		privateKey, publikey := btcec.PrivKeyFromBytes(hexPriv)
		if err != nil {
			t.Fatal(err)
		}

		privateKeys = append(privateKeys, privateKey)
		pubkeylist = append(pubkeylist, publikey.SerializeCompressed())

	}
	fmt.Println(privateKeys)

	txInputs := []bitcoin.TxIn{
		{
			Hash:     "9d433edd947173ddfab529ee405d9f06babf1c84fdedfa14d4fea35ba233340b",
			VOut:     1,
			PkScript: "0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58",
			Amount:   99750,
		},
		//{
		//	Hash:     "1f65f8ad3004f73b4c1745328b29d8860ab61cac18506848924ed2a43f7b500d",
		//	VOut:     1,
		//	PkScript: "0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58",
		//	Amount:   1199598600,
		//},
	}
	pk1, _ := hex.DecodeString("0014d7fae4fbdc8bf6c86a08c7177c5d06683754ea71")
	pk2, _ := hex.DecodeString("0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58")
	txOutputs := []bitcoin.TxOut{
		{
			PayScript: pk1,
			Amount:    1000,
		},
		{
			PayScript: pk2,
			Amount:    98450,
		},
	}

	builder := bitcoin.NewMultiTransactionBuilder()
	err = builder.NetParams(bitcoin.RegTest)
	if err != nil {
		t.Fatal(err)
	}
	err = builder.AddMultiPublicKey(pubkeylist, 2)
	if err != nil {
		t.Fatal(err)
	}
	err = builder.AddTxIn(txInputs)
	if err != nil {
		t.Fatal(err)
	}
	err = builder.AddTxOutScript(txOutputs)
	if err != nil {
		t.Fatal(err)
	}
	hash := builder.TxHash()
	t.Logf("before: %v \n", hash)
	//err = builder.Sign(func(hash []byte) ([][]byte, error) {
	//	var sigs [][]byte
	//	for _, privateKey := range privateKeys {
	//		sig := ecdsa.Sign(privateKey, hash)
	//		sigWithType := append(sig.Serialize(), byte(txscript.SigHashAll))
	//		sigs = append(sigs, sigWithType)
	//	}
	//	t.Logf("signature: %x\n", sigs)
	//	return sigs, nil
	//})
	if err != nil {
		t.Fatal(err)
	}
	txBytes, err := builder.Build()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("after: %v \n", hash)
	hexTxBytes := fmt.Sprintf("%x", txBytes)
	txHash, err := client.Sendrawtransaction(hexTxBytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("txHash: %v\n", txHash)

}

func TestMockDemo(t *testing.T) {
	for i := 0; i < 3; i++ {
		address, err := client.GetRawChangeAddress()
		if err != nil {
			t.Fatal(err)
		}
		dumpPrivkey, err := client.DumpPrivkey(address)
		if err != nil {
			t.Fatal(err)
		}
		addressInfo, err := client.Getaddressinfo(address)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("address:%v,privatekey:%v,publickey:%v\n", address, dumpPrivkey, addressInfo.Pubkey)
	}
}

func TestSimpleTx(t *testing.T) {
	privateKey := "cPbyxXLqYqAjAHDtbvKq7ETd6BsQBbS643RLHH4u3k1YeVAXkAqR"
	secret, _, err := base58.CheckDecode(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	inputs := []bitcoin.TxIn{
		{
			Hash:     "820556b04b6955cc0646bd424792c5940753a82d55140d9af3cc2f8060576f21",
			VOut:     1,
			PkScript: "0014d7fae4fbdc8bf6c86a08c7177c5d06683754ea71",
			Amount:   1200000000,
		},
	}

	outputs := []bitcoin.TxOut{
		{
			Address: "bcrt1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq2aw6ze",
			Amount:  1199999800,
		},
	}

	result, err := bitcoin.CreateTransaction(secret, inputs, outputs, bitcoin.RegTest)
	if err != nil {
		t.Fatal(err)
	}
	hexTxData := hex.EncodeToString(result)
	t.Log(hexTxData)
	txHash, err := client.Sendrawtransaction(hexTxData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}

func TestDepositTransaction(t *testing.T) {
	utxoSet, err := client.Scantxoutset("bcrt1q6lawf77u30mvs6sgcuthchgxdqm4f6n3kvx4z5")
	if err != nil {
		t.Fatal(err)
	}
	if len(utxoSet.Unspents) == 0 {
		t.Fatal("no utxo found")
	}
	var utxo types.Unspents
	for _, tUtxo := range utxoSet.Unspents {
		if tUtxo.Amount > 0.0001 {
			utxo = tUtxo
			break
		}
	}
	if utxo.Amount == 0 {
		t.Fatal("no utxo found")
	}
	floatBig := big.NewFloat(0).Mul(big.NewFloat(utxo.Amount), big.NewFloat(100000000))
	amount, _ := floatBig.Int64()
	inputs := []bitcoin.TxIn{
		{
			Hash:     utxo.Txid,
			VOut:     uint32(utxo.Vout),
			PkScript: utxo.ScriptPubKey,
			Amount:   amount,
		},
	}
	outputs := []bitcoin.TxOut{
		{
			Address: "bcrt1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq2aw6ze",
			Amount:  amount - 200,
		},
	}

	//bcrt1q6lawf77u30mvs6sgcuthchgxdqm4f6n3kvx4z5
	privateKey := "cPbyxXLqYqAjAHDtbvKq7ETd6BsQBbS643RLHH4u3k1YeVAXkAqR"
	secret, _, err := base58.CheckDecode(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	ethAddr := "771815eFD58e8D6e66773DB0bc002899c00d5b0c"
	ethAddrBytes, err := hex.DecodeString(ethAddr)
	if err != nil {
		t.Fatal(err)
	}

	result, err := bitcoin.CreateDepositTransaction(secret, ethAddrBytes, inputs, outputs, bitcoin.RegTest)
	if err != nil {
		t.Fatal(err)
	}
	hexTxData := hex.EncodeToString(result)
	t.Log(hexTxData)
	txHash, err := client.Sendrawtransaction(hexTxData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}

func TestMultiTransaction(t *testing.T) {
	privateKey1 := "cTZYUrhkQmpKGf6QjNhiLrw7g52VL449Tgwzo9UmKbnzphfZUUZp"
	secret1, _, err := base58.CheckDecode(privateKey1)
	if err != nil {
		t.Fatal(err)
	}
	privateKey2 := "cQtt5owpfyv79fGpnv2yncKtQTU42rzohXGCSD3Mcn4CGgEkPVVE"
	secret2, _, err := base58.CheckDecode(privateKey2)
	if err != nil {
		t.Fatal(err)
	}
	privateKey3 := "cSwe7Np3o7eCec6hgKFrwqGs9bb6x2dubKBucLYqQ6mJu5JH1aCn"
	secret3, _, err := base58.CheckDecode(privateKey3)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%x\n%x\n%x\n", secret1, secret2, secret3)
	//t.Logf("%x\n%x\n%x", secret1, secret2, secret3)
	secrets := [][]byte{secret1, secret2, secret3}
	inputs := []bitcoin.TxIn{
		{
			Hash:     "b5a43b150d8f9a305b9f19b11411f6a166f57080c8347b258ba48ebce77a2bc9",
			VOut:     1,
			PkScript: "0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58",
			Amount:   203382,
		},
	}

	outputs := []bitcoin.TxOut{
		{
			Address: "tb1q6lawf77u30mvs6sgcuthchgxdqm4f6n359lc4a",
			Amount:  203182,
		},
	}

	txBytes, err := bitcoin.CreateMultiSigTransaction(2, secrets, inputs, outputs, bitcoin.TestNet)
	if err != nil {
		t.Fatal(err)
	}
	hexTxData := hex.EncodeToString(txBytes)
	t.Log(hexTxData)
	txHash, err := client.Sendrawtransaction(hexTxData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}

func TestTxHash(t *testing.T) {
	msgTx := wire.NewMsgTx(2)

	hash, _ := chainhash.NewHashFromStr("d84498c08fadb2276b1010e8a572e351a4969b8bbb58ac9ff71d8fd3748e2faa")
	txIn := wire.NewTxIn(wire.NewOutPoint(hash, 1), nil, nil)
	msgTx.AddTxIn(txIn)

	pk1, _ := hex.DecodeString("0014d7fae4fbdc8bf6c86a08c7177c5d06683754ea71")
	pk2, _ := hex.DecodeString("0020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb58")
	out := &wire.TxOut{
		Value:    99950,
		PkScript: pk1,
	}
	msgTx.AddTxOut(out)
	out1 := &wire.TxOut{
		Value:    1199899500,
		PkScript: pk2,
	}
	msgTx.AddTxOut(out1)
	txId := msgTx.TxHash().String()
	t.Log(txId)

}

func TestDeserializeTransaction(t *testing.T) {
	data := "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000007d0200000001aa2f8e74d38f1df79fac58bb8b9b96a451e372a5e810106b27b2ad8fc09844d80000000000ffffffff026e86010000000000160014d7fae4fbdc8bf6c86a08c7177c5d06683754ea716c03854700000000220020efd9eeb78068445762a9e2e335d6ab826aaecff9cb7da24a39fc4686d468fb5800000000000000"
	hexData, err := hex.DecodeString(data)
	if err != nil {
		t.Fatal(err)
	}
	//version := hexData[0:32]
	l := hexData[32:64]
	length, err := strconv.ParseInt(fmt.Sprintf("%x", l), 16, 64)
	if err != nil {
		t.Fatal(err)
	}
	tx := hexData[64 : 64+length]
	t.Logf("%x\n", tx)
	msgTx := wire.MsgTx{}
	err = msgTx.Deserialize(bytes.NewReader(tx))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msgTx.TxHash().String())
	t.Log(msgTx)

}
