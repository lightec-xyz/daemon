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

var endpoint = "https://icp0.io"

var canisterId = "xbsjj-qaaaa-aaaai-aqamq-cai"

var walletId = "mno3q-2iaaa-aaaak-qdewq-cai"

func init() {
	secp256k1Identity, err := identity.NewSecp256k1IdentityFromPEMWithoutParameters([]byte(``))
	if err != nil {
		panic(err)
	}
	fmt.Printf("accountId: %x \n", secp256k1Identity.PublicKey())
	fmt.Printf("sender: %v \n", secp256k1Identity.Sender())
	fmt.Printf("PrincipalId: %v \n", secp256k1Identity.Sender().String())
	client, err = NewClientWithIdentity(canisterId, walletId, endpoint, secp256k1Identity)
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

func TestClient_VerifyPublicKey(t *testing.T) {
	publicKey, err := client.PublicKey()
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

func TestClient_EthDecoderCanister(t *testing.T) {
	result, err := client.EthDecoderCanister()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestClient_PlonkVerifierCanister(t *testing.T) {
	result, err := client.PlonkVerifierCanister()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestClient_Sign(t *testing.T) {
	sigsHash := []string{
		"a7f740000b5390f8c05c88d079332a88993ecaef1bea85126cbf9cf1cb9172d4",
	}
	signature, err := client.BtcTxSign(
		"33490dfb08e3bc5cdd28655a6e454e6abf8f79fb0ec5251019f0226e7e6401e0",
		"ff1ddc991b66997739b70e846d562cc11dd5012487ce9b12b40c555a71bd6f2d",
		"d3d4225de9e7d9586d415cdd9797de8a98cccd9db8620750c192aecd041f14cb",
		"00d6e094f862736efa47fa2abe94a0b038ed1a2f0915f7ecfd7f3399b793c8cc25e23fe6226b4a48944c31f90b896ee516b6edcadbba387404f37ba46277971d2156d39eb224717197ca4d6add8b1ace9ba9ec5b595aed5d72196e880067282703ff3314a64055887a35054a97f615ef7ff3076333753a473c5ee85e5b9911c506e341360592ac1b6d501ee8c9e1f96e74af0c5e59ff9b37863d6fa62ce1fb331f8d564db45a3e0f4defb66ec0318f267c6802506c4e69474cb071aeed7bc776142169102a666f00e890160517c4d05e85c5f24045ce631769fce401f12206f70d4d1ffa6ce1d518c0aa7bc910e8a8bf5b2679847133c192b9937a54e4c76aa62a5d2fd06093d0f00d286170b035c3f67b529437e034c569297e8f145401a907036cc63109efd2ece474e0a7da0bc9faa7c1aa66c5714bc2410d3c37f4563f8b115fb988925e04aa612367418539fb1afdc967ba818ba1075e83c65a3315f1e91495bd57b13152120d802babdc8e5e5d11537a45462b716c582f24f4540662ec07f29caebb1c5c0170ffe2df8bf1cb7e1c3f7c06e1a59380cf0c6004ab9f002103aa11963be314078316176e1cd0614d5bbc01ac89a7f978c812e50dda58c9581625c4927b8f6bc9ccc742a9a35646d8cbbcba7da398c995bac3449c2ec3a00511b4c3257e31ba61873bc288a12b2b8c0eab100d4278bdc356fe36a4f53091c81954a19c233511f765e099acbe35d8d1148e8235e4fe7b2a0268feebaa1a4ad50cfad89a3540311c97a51f0b04e8b6630b7b46b92f5a4dcf7d0fec6fac1633ff10418dd9d74f64ca7b5be17d1812c03a59a8c305ccb8c3766f909177eb6cb23a161dd25ae02cdefcba337dc3e000f4b6498962ce9a5898c4fccd09deffe72b602c248c2839f9abeefb772b0c502904760b4307b2b70f7bd84d85d3f855c3425210a06e28755fe1ea2ef8a6a9619f5d1d673227f4d085e48368f54a4737bd316a0c5564d200649bcea6bc1762b7a55e839e4f230446649d574be78cf893bc0e871e8fc1cf22ed7002c02592649338f6ecf1ada071eade1657ab093ad237bbdebf2bf47c6a30e3107fb193bf9ca9d6ad0c4965870f6355a0da6c4f9efe0cc788651924140bd2536c579b2787ce8eeda7be6b799794dea489ee8141ccbba44362362386b6c7788fe302cb3e974818ceff0cb1e72568d3bf9a79ad812be3724f6e15",
		"987600000000000",
		sigsHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(signature)
}

func TestClient_BlockHeight(t *testing.T) {
	height, err := client.BlockSignature()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(height.Height)

}

func TestClient_WalletBalance(t *testing.T) {
	balance, err := client.WalletBalance()
	if err != nil {
		t.Fatal(err)
	}
	//4756213240683
	t.Log(balance)
}
