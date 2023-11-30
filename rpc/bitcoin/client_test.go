package bitcoin

import (
	"fmt"
	"testing"
)

var client *Client
var err error

func init() {
	url := "https://go.getblock.io/d54c59f635654cc082de1f3fd14e5d02"
	user := "lightec"
	pwd := "abcd1234"
	network := "regtest"
	client, err = NewClient(url, user, pwd, network)
	if err != nil {
		panic(err)
	}
}

func TestClient_GetBlockHeader(t *testing.T) {
	header, err := client.GetBlockHeader("")
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
	hash, err := client.GetBlockHash(200)
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
}

func TestClient_GetBlockTx(t *testing.T) {
	hash, err := client.GetBlockHash(2540940)
	if err != nil {
		panic(err)
	}
	blockWithTx, err := client.GetBlockWithTx(hash)
	if err != nil {
		panic(err)
	}
	fmt.Println(blockWithTx)
}
