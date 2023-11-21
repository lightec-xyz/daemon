package bitcoin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

// Bitcoin: https://blockstream.info/api/
const url = "https://blockstream.info/testnet/api"

func GetLastestBlockHeight(url string) int64 {
	resp, err := http.Get(url + "/blocks/tip/height")
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return 0
	}
	// 打印响应内容
	height, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return 0
	}
	return height
}

func GetBlockHashByHeight(url string, height int64) string {
	resp, err := http.Get(url + "/block-height/" + strconv.FormatInt(height, 10))
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}
	return string(body)
}

func GetTxByHash(url, hash string) string {
	resp, err := http.Get(url + "/tx/" + hash)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}
	// 打印响应内容
	return string(body)
}

// TODO(NeedCheck)
func SendRawTx(url string, rawTx []byte) string {
	resp, err := http.Post(url+"/tx", "application/octet-stream", bytes.NewBuffer(rawTx))
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}
	// 打印响应内容
	return string(body)
}

func GetAddress(url, address string) string {
	resp, err := http.Get(url + "/address/" + address)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}
	// 打印响应内容
	return string(body)
}

func GetAddressTxs(url, address string) string {
	resp, err := http.Get(url + "/address/" + address + "/txs")
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}
	// 打印响应内容
	return string(body)
}

func GetAddressUtxos(url, address string) string {
	resp, err := http.Get(url + "/address/" + address + "/utxo")
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}
	// 打印响应内容
	return string(body)
}

func TestLastestBlockHeight(t *testing.T) {
	res := GetLastestBlockHeight(url)
	fmt.Println(res)
}

func TestGetTxByHash(t *testing.T) {
	res := GetTxByHash(url, "65eb5594eda20b3a2437c2e2c28ba7633f0492cbb33f62ee31469b913ce8a5ca")
	fmt.Println(res)
}

func TestGetAddress(t *testing.T) {
	res := GetAddress(url, "tb1ql9azatvsw9cxydtu3s0wzhf76zjnynhasuy4zy")
	fmt.Println(res)
}

func TestGetAddressTxs(t *testing.T) {
	res := GetAddressTxs(url, "tb1ql9azatvsw9cxydtu3s0wzhf76zjnynhasuy4zy")
	fmt.Println(res)
}

func TestGetAddressUtxos(t *testing.T) {
	res := GetAddressUtxos(url, "tb1ql9azatvsw9cxydtu3s0wzhf76zjnynhasuy4zy")
	fmt.Println(res)
}
