package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"math/big"
	"testing"
	"time"
)

var err error
var client *Client

// var endpoint = "https://1rpc.io/54japjRWgXHfp58ud/sepolia"
var endpoint = "https://ethereum-holesky.publicnode.com"
var zkBridgeAddr = "0xada86dce6d7e0d69ce4e25256b58ac1dcbbe2129"

func init() {
	//https://sepolia.publicgoods.network
	client, err = NewClient(endpoint, zkBridgeAddr)
	if err != nil {
		panic(err)
	}
}

func TestClient_TestEth(t *testing.T) {
	result, err := client.EthGetBlockByNumber(451228, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestClient_GetLogs(t *testing.T) {
	block, err := client.GetBlock(472244)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(block)
	logs, err := client.GetLogs(block.Hash().Hex(),
		[]string{"0xada86dce6d7e0d69ce4e25256b58ac1dcbbe2129"},
		[]string{"0xb28ad0403b0a341130002b9eef334c5daa3c1002a73dd90d4626f7079d0a804a"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(logs)

	for _, log := range logs {
		//0000000000000000000000000000000000000000000000000000000000000020
		//0000000000000000000000000000000000000000000000000000000000000071
		//02000000011aae5c5a37f9003aaa12c63dcebdfcd0e5cb6d753c4265ec055d06
		//97e5e0d6100100000000ffffffff026e86010000000000160014d7fae4fbdc8b
		//f6c86a08c7177c5d06683754ea71ecdc7ee202000000160014fb5defb676e7f0
		//a6711e3bc385849572a57fbe7e00000000000000000000000000000000000000
		version := log.Data[0:32]
		length := log.Data[32:64]
		fmt.Printf("%x %x \n", version, length)
		fmt.Printf("%x\n", log.Data)
	}
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
	redeemAmount := big.NewInt(100000)
	minerFee := big.NewInt(300)
	redeemLockScript, err := hex.DecodeString("0014d7fae4fbdc8bf6c86a08c7177c5d06683754ea71")
	if err != nil {
		t.Fatal(err)
	}
	from := common.HexToAddress("0x771815eFD58e8D6e66773DB0bc002899c00d5b0c")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	gasLimit := 500000
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		t.Fatal(err)
	}
	gasPrice = big.NewInt(0).Mul(big.NewInt(2), gasPrice)
	chainID, err := client.ChainID(ctx)
	if err != nil {
		t.Fatal(err)
	}
	nonce, err := client.NonceAt(ctx, from, nil)
	if err != nil {
		t.Fatal(err)
	}
	txhash, err := client.Redeem(privateKey, uint64(gasLimit), chainID, big.NewInt(int64(nonce)), gasPrice, redeemAmount, minerFee, redeemLockScript)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txhash)
}

func TestDepositeTransaction(t *testing.T) {
	privateKey := "c0781e4ca498e0ad693751bac014c0ab00c2841f28903e59cdfe1ab212438e49"
	txId := "31917fbc5da25a5db50a084dcfa4b72c04413e570d60bca338eca1cac70bbb28"
	//ethAddr := "0x771815eFD58e8D6e66773DB0bc002899c00d5b0c"
	index := uint32(1)
	amount := big.NewInt(12390000000)
	proofBytes := []byte("test proof")
	from := common.HexToAddress("0x771815eFD58e8D6e66773DB0bc002899c00d5b0c")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	gasLimit := 500000
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		t.Fatal(err)
	}
	gasPrice = big.NewInt(0).Mul(big.NewInt(2), gasPrice)
	chainID, err := client.ChainID(ctx)
	if err != nil {
		t.Fatal(err)
	}
	nonce, err := client.NonceAt(ctx, from, nil)
	if err != nil {
		t.Fatal(err)
	}
	txHash, err := client.Deposit(privateKey, txId, index,
		nonce, uint64(gasLimit), chainID, gasPrice, amount, proofBytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}

func TestUpdateUtxoChange(t *testing.T) {
	privateKey := "c0781e4ca498e0ad693751bac014c0ab00c2841f28903e59cdfe1ab212438e49"
	txId := "0xd32b0bc770512f49884b1e0e409c2010989c6fc7d76e4e495544a5cdb6cd9f49"
	//ethAddr := "0x771815eFD58e8D6e66773DB0bc002899c00d5b0c"
	proofBytes := []byte("test proof")
	from := common.HexToAddress("0x771815eFD58e8D6e66773DB0bc002899c00d5b0c")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	gasLimit := 500000
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		t.Fatal(err)
	}
	gasPrice = big.NewInt(0).Mul(big.NewInt(2), gasPrice)
	chainID, err := client.ChainID(ctx)
	if err != nil {
		t.Fatal(err)
	}
	nonce, err := client.NonceAt(ctx, from, nil)
	if err != nil {
		t.Fatal(err)
	}
	txHash, err := client.UpdateUtxoChange(privateKey, txId,
		nonce, uint64(gasLimit), chainID, gasPrice, proofBytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}
