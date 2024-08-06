package ethereum

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	btctx "github.com/lightec-xyz/daemon/rpc/bitcoin/common"
	"log"
	"math/big"
	"testing"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var err error
var client *Client

// var endpoint = "https://1rpc.io/54japjRWgXHfp58ud/sepolia"
var endpoint = "https://ethereum-holesky-rpc.publicnode.com"

// var endpoint = "http://127.0.0.1:8970"
var zkBridgeAddr = "0xa7becea4ce9040336d7d4aad84e684d1daeabea1"
var zkBtcAddr = "0x5898953ff9c1c11a8a6bc578bd6c93aabcd1f083"
var utxoManager = "0x9d2aaea60dee441981edf44300c26f1946411548"

func init() {
	//https://sepolia.publicgoods.network
	client, err = NewClient(endpoint, zkBridgeAddr, zkBtcAddr, utxoManager)
	if err != nil {
		panic(err)
	}
}

func TestClient_CheckUtxo(t *testing.T) {
	result, err := client.GetUtxo("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestClient_Demo001(t *testing.T) {
	receipt, err := client.TransactionReceipt(context.Background(), ethcommon.HexToHash("0xb19639d5c7c5804632f8ed92ca7e16d78cc1c6590a314b0aafee78793be223c6"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(receipt.TransactionIndex)
}

func TestClient_ChainFork(t *testing.T) {
	fork, err := client.ChainFork(596815)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fork)
}

func TestClient_GetTxSender(t *testing.T) {
	sender, err := client.GetTxSender("0xb19639d5c7c5804632f8ed92ca7e16d78cc1c6590a314b0aafee78793be223c6",
		"0xf99ab49c39e77bd6274035cbc1d6db068e014d3dc8e8a6a4c988f327a9b417f1", 39)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sender)
}

func TestClient_Number(t *testing.T) {
	number, err := client.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(number)
}

func TestClient_GetLogs2(t *testing.T) {
	logs, err := client.GetLogs("", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(logs)
}

func TestClient_GetPendingNonce(t *testing.T) {
	nonce, err := client.GetPendingNonce("0x771815eFD58e8D6e66773DB0bc002899c00d5b0c")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(nonce)
}

func TestClient_GetEstimateGasLimit(t *testing.T) {
	gasLimit, err := client.GetEstimateGasLimit(
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
	numb := big.NewInt(1245780)
	// result, err := client.HeaderByNumber(context.Background(), numb)
	result, err := client.BlockByNumber(context.Background(), numb)
	if err != nil {
		t.Fatal(err)
	}
	marshal, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(marshal))
}

func TestClient_BlockNumber(t *testing.T) {
	result, err := client.GetBlock(607368)
	if err != nil {
		t.Fatal(err)
	}
	marshal, err := json.Marshal(result.Transactions())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(marshal))

}

func TestClient_GetLogs(t *testing.T) {
	//563180
	//563166
	block, err := client.GetBlock(1545882)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log(block)
	address := []string{"0x9d2aaea60dee441981edf44300c26f1946411548", "0x8e4f5a8f3e24a279d8ed39e868f698130777fded"}
	topic := []string{"0xbfb6a0aa850eff6109c854ffb48321dcf37f02d6c7a44c46987a5ddf3419fc07", "0x1e5e2baa6d11cc5bcae8c0d1187d7b9ebf13d6d9b932f7dbbf4e396438845fb8"}
	logs, err := client.GetLogs(block.Hash().Hex(),
		address, topic)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log(logs)
	for _, log := range logs {
		if log.TxHash.String() == "0xea7a29093b228e8d45ba54161689e1ae7c4caa1ce33fd618112eace20e2acf1a" {
			txData, _, err := DecodeRedeemLog(log.Data)
			if err != nil {
				t.Fatal(err)
			}
			transaction := btctx.NewTransaction()
			err = transaction.Deserialize(bytes.NewReader(txData))
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(transaction.TxHash().String())
			//0x0020957ab85b710cb5b577171e23bb3492536c8029cc99511f3920d3cc13871a2327
			for _, out := range transaction.TxOut {
				t.Logf("%x %v \n", out.PkScript, out.Value)
			}

		}
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
	redeemAmount := uint64(2199999800)
	minerFee := uint64(300)

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
	from := ethcommon.HexToAddress(fromAddr)
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
	proof := []byte("test proof")
	from := ethcommon.HexToAddress("0x771815eFD58e8D6e66773DB0bc002899c00d5b0c")
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
	txHash, err := client.Deposit(privateKey,
		nonce, uint64(gasLimit), chainID, gasPrice, nil, proof)
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
	from := ethcommon.HexToAddress("0x771815eFD58e8D6e66773DB0bc002899c00d5b0c")
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

func TestClient_Demo(t *testing.T) {
	ids := TxIdsToFixedIds([]string{"adddd", "dsdsfsd"})
	t.Log(ids)

}

func TestGetTrancaction(t *testing.T) {
	hash := ethcommon.HexToHash("0x9bd7ff0aa08611a2077189fcefb5095eda2e5d28d175cde410540ecc4ec2283b")
	tx, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tx.BlockNumber)
}

func TestClient_Verify(t *testing.T) {
	receipt, err := client.Client.TransactionReceipt(context.Background(), ethcommon.HexToHash("0x291ee31eb6b8cef1ebc571fd090a1e7c96ddac5a1552dae47501581ed7d66641"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(receipt)
}
