package dfinity

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/aviate-labs/agent-go/identity"
	"testing"
)

var client *Client
var err error

var txSignerCanId = "wlkxr-hqaaa-aaaad-aaxaa-cai"
var blockSignerCanId = "xdqo6-dqaaa-aaaal-qsqva-cai"
var walletCanId = "mno3q-2iaaa-aaaak-qdewq-cai"
var secp256k1Identity *identity.Secp256k1Identity

func init() {
	secret := ``
	if secret != "" {
		secp256k1Identity, err = identity.NewSecp256k1IdentityFromPEMWithoutParameters([]byte(secret))
		if err != nil {
			panic(err)
		}
		fmt.Printf("accountId: %x \n", secp256k1Identity.PublicKey())
		fmt.Printf("sender: %v \n", secp256k1Identity.Sender())
		fmt.Printf("PrincipalId: %v \n", secp256k1Identity.Sender().String())
	}

	option := NewOption(walletCanId, txSignerCanId, blockSignerCanId, secp256k1Identity)
	client, err = NewClient(option)
	if err != nil {
		panic(err)
	}
}

func TestClient_BtcTxSignatureWithCycle(t *testing.T) {
	sigsHash := []string{"cb0799771fea76af9496fa65df8834529f7cd9ff0cf766ab65219cbef532e5f5", "b5fee61e386aca8671e0b925e597926d9d3c7a8cc65ff73b53a5c09afb13363e"}
	signature, err := client.BtcTxSignWithCycle(
		"6813630637d6a988524df575932262242bd5a4486946ff011bdd58c4397b3d93",
		"ff1ddc991b66997739b70e846d562cc11dd5012487ce9b12b40c555a71bd6f2d",
		"1916ef627d82b4c52d2ec3eb350af63fd9754674892667d6e0b5aa8dd7a850c5",
		"d20903ba2129a1b9812ac3dd37d9ef529996fa5f95305483704f0ed72150bcb0905ab670c87ebd7a064fe64d39ea0b88d62f0975d2ca3a55f2aa8c3b289c56e3f0239f31d190871bba9956dabcf933213abdcf12fd551efd18106807f205fd3cf05abefe1772a98940046338691f2c3e4c6b3050539330b444dd7ff021206a0ea4b38f0a63db9e3b8cd6b904344a95610e9c2298a0379d2727237707d924a35bda1f571539cbafbd9043f6f35eb4f48bf245b2544744e665785a2e2958d75642c28388b5ecc298af8f4e504b55af748863d191fd101ede15109456dcf310274ae1aa8fc6eb24ace3bef362f11ba467686220828b150d3aa9d29afe8d5c77bcff000000071fd5701ed1f973772dc6f286c5f05b4e5cfcc12bb9e8c55bb4b4b6e020747c26026c23d59823a1cf6df3fac1dbbc53df989d32398eca54637060a2064de038941abd7f429913877083fea195480896be241459785790d99a42166123a059f9870f302539d57a188e2780fd795eb12552deddf7abbb12fc13885570e918eb904416820c960fa90c8f5ac8ec7eb1f303cabce812bfae157c5cc3433635809e69e41543b800105f679610a434200fc191173d452530194ed93da965dde1bd52a26a26fa5faa346bb035533ed52b7aee242cb5956ebd1e2bfe1f9e6d61ef8bebc469a75d8b36942085aada9a2b9a7bb18e74f2d69c5d755b84efd26f5219d28643551561e532bbdb7486459397de3deb5663d1a6afbf11db8c51367526ca60e3d72300000001eeebd9c998d847c3afe485f368750d9cfe7e56d72a321d0bba5f61e1c6f3ddde",
		"168000000000000",
		sigsHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(signature)
}

func TestClient_Identity(t *testing.T) {
	seed, err := hex.DecodeString("")
	if err != nil {
		t.Fatal(err)
	}
	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	ed25519Identity, err := identity.NewEd25519Identity(publicKey, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ed25519Identity.Sender())
}

func TestClient_TxPublicKey(t *testing.T) {
	publicKey, err := client.TxPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(publicKey)
}

func TestClient_BlockPublicKey(t *testing.T) {
	publicKey, err := client.BlockPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(publicKey)
}

func TestClient_DummyAddress(t *testing.T) {
	address, err := client.DummyAddress()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(address)
}

func TestClient_Sign(t *testing.T) {
	sigsHash := []string{"cb0799771fea76af9496fa65df8834529f7cd9ff0cf766ab65219cbef532e5f5", "b5fee61e386aca8671e0b925e597926d9d3c7a8cc65ff73b53a5c09afb13363e"}
	signature, err := client.BtcTxSign(
		"6813630637d6a988524df575932262242bd5a4486946ff011bdd58c4397b3d93",
		"ff1ddc991b66997739b70e846d562cc11dd5012487ce9b12b40c555a71bd6f2d",
		"1916ef627d82b4c52d2ec3eb350af63fd9754674892667d6e0b5aa8dd7a850c5",
		"d20903ba2129a1b9812ac3dd37d9ef529996fa5f95305483704f0ed72150bcb0905ab670c87ebd7a064fe64d39ea0b88d62f0975d2ca3a55f2aa8c3b289c56e3f0239f31d190871bba9956dabcf933213abdcf12fd551efd18106807f205fd3cf05abefe1772a98940046338691f2c3e4c6b3050539330b444dd7ff021206a0ea4b38f0a63db9e3b8cd6b904344a95610e9c2298a0379d2727237707d924a35bda1f571539cbafbd9043f6f35eb4f48bf245b2544744e665785a2e2958d75642c28388b5ecc298af8f4e504b55af748863d191fd101ede15109456dcf310274ae1aa8fc6eb24ace3bef362f11ba467686220828b150d3aa9d29afe8d5c77bcff000000071fd5701ed1f973772dc6f286c5f05b4e5cfcc12bb9e8c55bb4b4b6e020747c26026c23d59823a1cf6df3fac1dbbc53df989d32398eca54637060a2064de038941abd7f429913877083fea195480896be241459785790d99a42166123a059f9870f302539d57a188e2780fd795eb12552deddf7abbb12fc13885570e918eb904416820c960fa90c8f5ac8ec7eb1f303cabce812bfae157c5cc3433635809e69e41543b800105f679610a434200fc191173d452530194ed93da965dde1bd52a26a26fa5faa346bb035533ed52b7aee242cb5956ebd1e2bfe1f9e6d61ef8bebc469a75d8b36942085aada9a2b9a7bb18e74f2d69c5d755b84efd26f5219d28643551561e532bbdb7486459397de3deb5663d1a6afbf11db8c51367526ca60e3d72300000001eeebd9c998d847c3afe485f368750d9cfe7e56d72a321d0bba5f61e1c6f3ddde",
		"168000000000000",
		sigsHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(signature)
}

func TestClient_BlockSignature(t *testing.T) {
	result, err := client.BlockSignature()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)

}

func TestClient_BlockSignatureWithCycle(t *testing.T) {
	result, err := client.BlockSignatureWithCycle()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestClient_WalletBalance(t *testing.T) {
	balance, err := client.IcpBalance()
	if err != nil {
		t.Fatal(err)
	}
	//4756213240683
	t.Log(balance)
}
