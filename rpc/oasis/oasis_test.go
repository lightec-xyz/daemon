package oasis

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var err error
var client *Client

func init() {
	client, err = NewClient("https://testnet.sapphire.oasis.io", "0xe46A4519c9FD97EDcdAA895464a5B8953f4Fa9D3")
	if err != nil {
		panic(err)
	}
}

func TestClient_PublicKey(t *testing.T) {
	publicKey, err := client.PublicKey()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range publicKey {
		t.Logf("%v\n", hexutil.Encode(item))
	}
}

func TestClient_AlphaPublicKey(t *testing.T) {
	publicKey, err := client.PublicKey()
	if err != nil {
		t.Fatal(err)
	}
	for _, pub := range publicKey {
		t.Logf("%x\n", pub)
	}
}

func TestClient_TestnetSignBtcTx(t *testing.T) {
	//6740a187943efa84149a0757ee3e30cb9d22a43edbcaa5bda47de03737c52998
	sigs, err := client.SignBtcTx(
		"50740e4f18a8a1f2f04743c2433b758f1a8eb11e76bc164efc3e41df4894220c",
		"e8f03ee8f69f1ada70ca510d8c27a85c713be52c7402b37a8173a40bc81c7770",
		"2b7e1b475748d9ba20d27ca1987e322414cba36bdde04695e63875cc3eba4cf42a316fc6af8dd3a8e0a0d065db8c679caf9520adc9e1c4880747311797fa6b9701e3c94e4afb1813a6f3114191c4e6ffac1dfaf438eb9196e66fd0203cf85446172f55ee0802663591b7c174589f529f5bfb55c4d01d5d19185227e4a8b1735603d3aac4e91b0a673fee1cc5f906a0bf023a0c4d1b16a7ea652cf1dc34a557a3166a8dcd666b5c6b11820c435dd18e8eb04e6ae4da4aefe8327647987806f17b1ba47b09fa0f4e2b4acbc7debb2eb838a2d37053ca818864f54e5ffc7f02bb6e2c21ee3123e4cb8b95a02e1e3b0c5440943dd793e1bf556fbb0dbd4a7fff68b601d0abcd28a18b3066dc2213ea9ebc02a0af59827a5b1a8515defd6df9029a3d0a751cb963fb711afc9930514f8a5c217212e4b98e18321e134f9d3e9a040cd7196bed359f6b1ffbb969493e2cb7296d28d94b069e3333fcce3100127c30879908ae6822f55ff0950d7b48b368624f499279d6bae6d1e2266be7aca43910de6d2f86517bf4ffb94b98222ccc3ec19d894ef9ac9eb7a8fef3a2dc3ec625b2684c260002a36dd0ca799b6e9a0ac853bfcebf6b7c411d13418b8578166aa8720fb70817c3c2e06b0e465d3258d10e3da8d421da04a662b1eb9af41b5487b626b869190f9ac308b9df5533d63e97f96f506c26c2c80ee90865408d6d532089418c202139ec5c9738c6c17b9c4bfe62827628e779253084953bd4453c94ddf5b3d6ba1cb16a9bee64fdcb26f005625ed12983ea268ec6eaaf8885f77071f7edbf6b9805844c9ad88737a1404e60e73cadf1b560ed8c53adb81184fd75e024f0fe91c026bf20fea6af6e9ed7c135fe9e384befac5e5bb940061d9c36219756ff68419b07af7711342e63761fd5652c4a2bb276d6e22b6dd3b4c5014462f7022133ba1d2cbc451b2a27e0d1fb3a1a4a7229c5547ab83da3f2c9f22a42216e6ee8e6cee015f0f329b129381c532f47e5804db51cd1831ddea271aa1105505f53247943f927d4c320f06ba8843b201ec4d9e210217ddbd8f5865ed65b3839a994ee342fcc1177f603d6fc7c6d8cb91386d261c122b76bac0f00f778e18f6296d554173578182a77bc6b9b2ab85de425a95da9435bade9e817fa2704f18c7cfda3ca94e4651fe335eaaae3812126a478491821835d287f18ba99c6b252b5315db444097f16",
		[]string{"015f4b691b31091a1cee0a6c2191a6ebc2d45a516f8c3fe50280ad610edf0147"},
		big.NewInt(0).SetInt64(7012000000000000000),
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range sigs {
		for _, sig := range item {
			t.Logf("%x\n", sig)
		}
		t.Log("--------------------")
	}
}
