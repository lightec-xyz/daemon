package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lightec-xyz/daemon/transaction/ethereum/zkbridge"
	"testing"
)

var err error
var client *Client

// var endpoint = "https://1rpc.io/54japjRWgXHfp58ud/sepolia"
var endpoint = "https://rpc.notadegen.com/eth/sepolia"

func init() {
	//https://sepolia.publicgoods.network
	client, err = NewClient(endpoint)
	if err != nil {
		panic(err)
	}
}

func TestClient_TestEth(t *testing.T) {
	result, err := client.EthRPC.EthGetBlockByNumber(4794370, false)
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
	t.Fatalf("%x", privateKey.D.Bytes())
}

func TestTransaction(t *testing.T) {
	rpcDial, err := rpc.Dial(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	client := ethclient.NewClient(rpcDial)

	zkBridgeCall, err := zkbridge.NewZkbridge(common.HexToAddress("0x8dda72ee36ab9c91e92298823d3c0d4d73894081"), client)
	if err != nil {
		t.Fatal(err)
	}
	feeRate, err := zkBridgeCall.FeeRate(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(feeRate)
	//transaction, err := zkBridgeCall.Deposit()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//eip155Signer := types.NewEIP155Signer(big.NewInt(11155111))
	//privateKey, err := crypto.HexToECDSA("")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//signTx, err := types.SignTx(transaction, eip155Signer, privateKey)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//err = client.SendTransaction(context.TODO(), signTx)
	//if err != nil {
	//	t.Fatal(err)
	//}

}
