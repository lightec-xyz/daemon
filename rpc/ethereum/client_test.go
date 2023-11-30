package ethereum

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"testing"
)

var err error
var client *Client

// var endpoint = "https://1rpc.io/54japjRWgXHfp58ud/sepolia"
var endpoint = "https://rpc.notadegen.com/eth/sepolia"
var zkBridgeAddr = "0x8dda72ee36ab9c91e92298823d3c0d4d73894081"

func init() {
	//https://sepolia.publicgoods.network
	client, err = NewClient(endpoint, zkBridgeAddr)
	if err != nil {
		panic(err)
	}
}

func TestClient_TestEth(t *testing.T) {
	result, err := client.EthGetBlockByNumber(4794370, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestPrivateKey(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%x \n", privateKey.D.Bytes())
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("fail")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	t.Log(address)

}

func TestRedeemTx(t *testing.T) {
	privateKey := "c0781e4ca498e0ad693751bac014c0ab00c2841f28903e59cdfe1ab212438e49"
	from := ""
	redeemAmount := big.NewInt(100)
	minerFee := big.NewInt(100)
	redeemLockScript := []byte{}
	txhash, err := client.Redeem(from, privateKey, redeemAmount, minerFee, redeemLockScript)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txhash)
}

func TestTransaction(t *testing.T) {
	privateKey := "c0781e4ca498e0ad693751bac014c0ab00c2841f28903e59cdfe1ab212438e49"
	txId := "7d8f46b43caebfc8f5940b3bbab189aa96d6569580e7328f19d5542de2a51467"
	ethAddr := "0x771815eFD58e8D6e66773DB0bc002899c00d5b0c"
	index := uint32(1)
	amount := big.NewInt(123456)
	proofBytes := []byte("test proof")
	txHash, err := client.Deposit(privateKey, txId, ethAddr, index, amount, proofBytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}
