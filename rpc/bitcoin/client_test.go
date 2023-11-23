package bitcoinClient

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"strconv"
	"strings"
	"testing"
)

var client *Client
var err error

func init() {
	url := "http://localhost:8332"
	user := "ascendex"
	pwd := "Abcd1234"
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

func TestClient_PrivateKeyToAddress(t *testing.T) {

	seed := "cUgkKZ7JhaeDaNckcAsuL4zvmwTmkAD4cLVrHcWREHSDMzjVHwJm"
	decode, _, err := base58.CheckDecode(seed)
	if err != nil {
		panic(err)
	}
	privateKey, publicKey := btcec.PrivKeyFromBytes(decode)
	netParams := &chaincfg.RegressionNetParams
	from, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey.SerializeCompressed()), netParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x %x \n", privateKey.Serialize(), publicKey.SerializeCompressed())
	fmt.Println(from.EncodeAddress())
	if from.String() == "bcrt1q4mmnjm0nykr6atzs9np9kecdenp2qe7f5wulfa" {
		fmt.Println("success")
	}

}

func TestClient_GetBlockCount(t *testing.T) {
	chainInfo, err := client.GetBlockCount()
	if err != nil {
		panic(err)
	}
	fmt.Println(chainInfo)
}

func getAddressFromSeed(seed string, netParams *chaincfg.Params) (btcutil.Address, error) {
	//seed := "cUgkKZ7JhaeDaNckcAsuL4zvmwTmkAD4cLVrHcWREHSDMzjVHwJm"
	decode, _, err := base58.CheckDecode(seed)
	if err != nil {
		panic(err)
	}
	_, publicKey := btcec.PrivKeyFromBytes(decode)
	address, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey.SerializeCompressed()), netParams)
	if err != nil {
		return address, err
	}
	return address, nil

}

func hex2int(hexStr string) int64 {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)

	// base 16 for hexadecimal
	result, _ := strconv.ParseInt(cleaned, 16, 64)
	return result
}
