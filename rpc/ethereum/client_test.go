package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lightec-xyz/daemon/logger"
	"log"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"
)

var err error
var client *Client

// var endpoint = "https://1rpc.io/54japjRWgXHfp58ud/sepolia"
var endpoint = "https://rpc.holesky.ethpandaops.io"
var zkBridgeAddr = "0xa7becea4ce9040336d7d4aad84e684d1daeabea1"
var zkBtcAddr = "0x5898953ff9c1c11a8a6bc578bd6c93aabcd1f083"

func init() {
	//https://sepolia.publicgoods.network
	client, err = NewClient(endpoint, zkBridgeAddr, zkBtcAddr)
	if err != nil {
		panic(err)
	}
}

func TestClient_GetPendingNonce(t *testing.T) {
	nonce, err := client.GetPendingNonce("0x771815eFD58e8D6e66773DB0bc002899c00d5b0c")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(nonce)
}

func TestClient_GetEstimateGasLimit(t *testing.T) {
	gasLimit, err := client.GetEstimateGaslimit(
		"0x771815eFD58e8D6e66773DB0bc002899c00d5b0c",
		"0xbdfb7b89e9c77fe647ac1628416773c143ca4b51",
		"c937229bbd89dadb76e6f7285220e765c6e195b553cf0df6cd3e7505077df970",
		1,
		big.NewInt(9000000),
		[]byte("wo heni"),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gasLimit)
}

func TestClient_ZkbtcBalance(t *testing.T) {
	balance, err := client.GetZkBtcBalance("0x771815eFD58e8D6e66773DB0bc002899c00d5b0c")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance)
}

func TestClient_TestEth(t *testing.T) {
	result, err := client.EthGetBlockByNumber(451228, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestClient_GetLogs(t *testing.T) {
	//563180
	//563166
	block, err := client.GetBlock(576047)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log(block)
	address := []string{"0x52ebc075616195cc7deb79d5c21bd9b04acc33ee", "0x8b404b735afe5bcdce85a1ce753c79715f86062c"}
	topic := []string{"0x975dbbd59299029fdfc12db336ede29e2e2b2d117effa1a45be55f0b4f9cfbce", "0xb28ad0403b0a341130002b9eef334c5daa3c1002a73dd90d4626f7079d0a804a"}
	logs, err := client.GetLogs(block.Hash().Hex(),
		address, topic)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log(logs)
	for _, log := range logs {
		t.Log(log.Address.Hex(), log.Address.String(), log.Index, log.Topics, fmt.Sprintf("%x", log.Data))

	}
}

func parseEthDeposit(log types.Log) {
	txId := strings.ToLower(log.Topics[1].Hex())
	sprintf := strings.TrimPrefix(log.Topics[2].Hex(), "0x")
	vout, err := strconv.ParseInt(strings.ToLower(sprintf), 16, 32)
	if err != nil {
		logger.Error("parse vout error:%v", err)
		panic(err)
	}
	amount, err := strconv.ParseInt(fmt.Sprintf("%x", log.Data), 16, 64)
	if err != nil {
		logger.Error("parse amount error:%v", err)
		panic(err)
	}
	fmt.Println(txId, vout, amount)
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
	redeemAmount := big.NewInt(2199999800)
	minerFee := big.NewInt(300)

	fromAddr := "0x771815eFD58e8D6e66773DB0bc002899c00d5b0c"
	balance, err := client.GetZkBtcBalance(fromAddr)
	if err != nil {
		t.Fatal(err)
	}
	if balance.Cmp(big.NewInt(10000)) < 0 {
		t.Fatal("not enough balance")
	}

	redeemLockScript, err := hex.DecodeString("0014d7fae4fbdc8bf6c86a08c7177c5d06683754ea71")
	if err != nil {
		t.Fatal(err)
	}
	from := common.HexToAddress(fromAddr)
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
	txHash, err := client.Deposit(privateKey, txId, "", index,
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
	txHash, err := client.UpdateUtxoChange(privateKey, []string{txId},
		nonce, uint64(gasLimit), chainID, gasPrice, proofBytes)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}
