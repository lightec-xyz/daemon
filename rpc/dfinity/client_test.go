package dfinity

import (
	"crypto/ed25519"
	"encoding/hex"
	"github.com/aviate-labs/agent-go/identity"
	"testing"
)

var client *Client
var err error

var txSignerCanId = "wlkxr-hqaaa-aaaad-aaxaa-cai"
var blockSignerCanId = "xdqo6-dqaaa-aaaal-qsqva-cai"
var walletCanId = ""

func init() {
	//secp256k1Identity, err := identity.NewSecp256k1IdentityFromPEMWithoutParameters([]byte(``))
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("accountId: %x \n", secp256k1Identity.PublicKey())
	//fmt.Printf("sender: %v \n", secp256k1Identity.Sender())
	//fmt.Printf("PrincipalId: %v \n", secp256k1Identity.Sender().String())
	option := NewOption(walletCanId, txSignerCanId, blockSignerCanId, nil)
	client, err = NewClient(option)
	if err != nil {
		panic(err)
	}
}

func TestClient_BtcTxSignatureWithCycle(t *testing.T) {
	sigsHash := []string{
		"91bed04a24ed90fc0a52ae619b18100695268453316f54f99b48dc8206594dc7",
		"42bbf0ceb3b8d9a53356e11db32bfaf6986b8a28730ddf0cf7d21870fc20996c",
		"054ea52a67d4806a786d7a499262c7aa79fe51f9dd94a91644937924d7ffa57e",
	}
	signature, err := client.BtcTxSignWithCycle(
		"210ae266074912483727b69d3afe2603684f5a2341053de484fa74d1ca47165c",
		"9a8a838374b93de88c9b20b7f23c64994dd488ee9bef8ded53818733cfd154bf",
		"7ee2d9372a1e3606cd75774780fc7f51090c394b05e8534ee1fe53c648384aa6",
		"04c4541c2566f110268149fc6f4e2f1de6aa61e743cf979fd767ed2dab2891501dae8c6c615f7830d33780c4017073a12696dcd4713e8cccfe2f4ece3fe554a121bf1edfca18c2d98dee7216da5e8848e4a634ff7c53072a3d3c5ae0e83f15ac0dffb7cfe0c4432efac38573bd65689f5b863c65c8c4befbfa0fb44427294c582c035fea39346b41c00cd579456d3f149fa5bfa85dc732f8bb6290856929316c21bf0c00b75034e9415149b7abaade66808805c9ee6fc3cd4749f606c9c0e10c17871dd6f6d21083b1465c5c2c627bd785fa2fc9fb6600ee3261c8699c9f0794277882760352da2e1db4fcddd347bee338f0e85d58c0405c472d19cb08bc1b4710013a9636413701321349c8c2bf71de3f0593df680e152eade036e0ab6aaf5d1fa2b0943ff2b81ea2b63b6551b2718f88d6f5b6c2f1b4a7d9f92365ac80ee74153960ddc3a3fe1f90f1236f2d9a0e15470bb215e8f38a5a996b40170f72317e0e9b27f9a86a302e4b5b551b2c9e5bc7755c3cb8fe3283d4a63880792f4b294e284d9d08a6a3070c4c7f5c00de08cbbc75a4dec898c0bbd2fa491d5c610e8b30145876dfe2d90a10149e9a7cd6d308a9f222c922fd6070b8404d5c4abd453d7a2e694a0bf6ba28f710eddfaa44eafa39699b3112bee5a70b4cba66bc99d6b9a406deb16ade27346e50a9ceec8adfeb220b04ecc16dfb7613eb1d2568f156574328bcc01423717b864ba79365cca5ba8e8f18d275557f8401bea72790b3aa2c090f66ad0e5377b39ffb948288917817dcb64dab3994ec24aff6f67f7004a9b5381eb64355ed121837cd3293641c74dd129240ca430d6c41930011750c9c800d9b0da3ab72afe3022e1132f25bf5ec02201b9e8d4847978cc10f842d8189536a861a84d0ca38c1f4ca4ac1d2359397a9b63b22f1b6ae1fd54aecf55ad5fe2eed332f52556638ea1ccdab3466ee3a657c6c4b7b083b532ca31900d13456624aa5bc11c5b70677f345312a8c9ec20d52317710e0c0f6df1d68b280b502f843dbf4cd281fba5b0b3c635a55955a13156fa3c5f97df010457f5052c026f06655a83e542e2aacdf834e6a8d6a2bd289c1dace52eafcc08157e01ab0de46f0706e7ec95d299574f378683d8d6ba8e5495916eaf9cf390e505c75fafafd8a8ab627ddb97409368a36157ded4ed3a68353bc41c141d8601345715a8c30dd8bb294527b462b",
		"800000000000000",
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

func TestClient_WalletBalance(t *testing.T) {
	balance, err := client.IcpBalance()
	if err != nil {
		t.Fatal(err)
	}
	//4756213240683
	t.Log(balance)
}
