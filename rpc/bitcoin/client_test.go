package bitcoin

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	bitcoin "github.com/lightec-xyz/daemon/rpc/bitcoin/common"
	"github.com/stretchr/testify/require"
)

var client *Client
var err error

func init() {
	url := "https://rpc.ankr.com/btc/8c933202fbe8dbe6d63377a319b6020f4a4c35bb4424f6368f630b676b4fcc2e"
	client, err = NewClient(url, "", "")
	if err != nil {
		panic(err)
	}
}

func TestClient_CheckTxOnChain(t *testing.T) {
	exists, err := client.CheckTxOnChain("0edef4c17568ee3f6dfdd275c684572a05ffc22283acddc10f84d1b74bc39f82")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exists)
}

func TestClient_GetRawTransaction(t *testing.T) {
	tx, err := client.GetRawTransaction("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tx)
}

func TestClient_Estimatesmartfee(t *testing.T) {
	fee, err := client.Estimatesmartfee(50)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fee.Feerate)
}

func TestClient_GetBlockHeaderByHeight(t *testing.T) {
	header, err := client.GetBlockHeaderByHeight(80091)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(header)
}

func TestClient_Getmempoolentry(t *testing.T) {
	tx, err := client.Getmempoolentry("2242ae94aced9dde965de4648321c17b21c67177ce44b326fd47e2beb2e71aa1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tx)
}

func TestClient_GetBlock(t *testing.T) {
	block, err := client.GetBlock("000000007f1e11a0deb802c7a9df1908e70f349faa38ffd09656fd2b2bde1528")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(block.Time)
	for _, tx := range block.Tx {
		if tx.Txid == "8067f979ca7c23baab0e00311550b8e096a5ec097eb96ff29fdb6e23bfc777e3" {
			t.Log(tx)
		}
	}
}
func TestClient_GetBlockStr(t *testing.T) {
	block, err := client.GetBlockStr("00000000cb5a6a6f3f2dda8ac1c597d307dedaa80a6f131d70cf235d49c78a36")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(block)
}

func TestClient_CheckTx(t *testing.T) {
	tx, err := client.CheckTxOnChain("abd134879e9acd79cdae361ad986b2c1e5832aa28b33bdd4e488a5a01f6e5f05")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tx)
}

func TestClient_Sendrawtransaction(t *testing.T) {
	hash, err := client.Sendrawtransaction("0100000000010108cc56149dac05a4c41d4935a3a5d677135586b853482364c1dc859256d01e7f0000000000ffffffff0120bf020000000000160014a8a337fbfd692a5628e96d44684b9ee35d9e913e0247304402201fae89f0a24b5bb0c31262f76533d73b0e656338ffc5e4c5ebd84b3e18964f7e02205b241e9b989eae138740425849505c240bd33b8037bdf0d3d29397c8775a43d9012103d550bfd354d2c4aec5959c95a4c656a191d93350838195277fe97604d0c0c5ce00000000")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash)
}

func TestClient_GetBlockHeader(t *testing.T) {
	header, err := client.GetBlockHeader("000000000000003179df591c74bc12571d2ec8c68226e7fa010c302435479030")
	if err != nil {
		panic(err)
	}
	fmt.Println(header)
}

func TestClient_GetHexBlockHeader(t *testing.T) {
	header, err := client.GetHexBlockHeader("0000000018221eb554712872295e7c4590696683d726a6bc2d811b07cfae5bb0")
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
	hash, err := client.GetBlockHash(66674)
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
	//000000000002cb956c7eb0f40063b1e1dea575daa5cfc20d5bff55a1c6f46e46
	header, err := client.GetBlockHeader(hash)
	if err != nil {
		panic(err)
	}
	fmt.Println(header.Hash)
}

func TestClient_GetBlockTx(t *testing.T) {
	hash, err := client.GetBlockHash(84073)
	if err != nil {
		panic(err)
	}
	blockWithTx, err := client.GetBlock(hash)
	if err != nil {
		panic(err)
	}
	for _, tx := range blockWithTx.Tx {
		if tx.Txid == "7d8fa15a1368d0fa36952472843d6bbf78dd3376baf1108b9e9555da38d739f0" {
			t.Logf("find tx: %v", tx)
		}
	}
	fmt.Println(blockWithTx.Hash)
}

func Test_GetMultiSigScriptRelateds(t *testing.T) {
	publicKeys := [][]byte{
		common.FromHex("0x02b11c577f0eb7ec10e3af25e2135e9ece2e449ff45189af245bdecc6b7757def3"),
		common.FromHex("0x027f06cc1def813ef9b69cc7f07b79152961467fb3e47bdeb1700094231d38b68e"),
		common.FromHex("0x023f203422be55a3576f46dc6770bdc7865a126381c1963a2d82b49f4158409a2e"),
	}
	multiSigScript, walletAddr, lockScript, err :=
		bitcoin.GetMultiSigScriptRelateds(2, &chaincfg.TestNet3Params, publicKeys)
	require.NoError(t, err)

	fmt.Printf("MultiSig Script: %v\n", hexutil.Encode(multiSigScript))
	fmt.Printf("Wallet Address: %v\n", walletAddr.EncodeAddress())
	fmt.Printf("Lock Script: %v\n", hexutil.Encode(lockScript))
}
