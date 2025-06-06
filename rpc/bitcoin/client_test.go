package bitcoin

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	bitcoin "github.com/lightec-xyz/daemon/rpc/bitcoin/common"
	"github.com/stretchr/testify/require"
)

var client *Client
var err error

func init() {
	url := "http://127.0.0.1:9935"
	user := ""
	pwd := ""
	client, err = NewClient(url, user, pwd)
	if err != nil {
		panic(err)
	}
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
	block, err := client.GetBlock("00000000cb5a6a6f3f2dda8ac1c597d307dedaa80a6f131d70cf235d49c78a36")
	if err != nil {
		t.Fatal(err)
	}
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
	tx, err := client.CheckTx("abd134879e9acd79cdae361ad986b2c1e5832aa28b33bdd4e488a5a01f6e5f05")
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
	hash, err := client.GetBlockHash(3397399)
	if err != nil {
		panic(err)
	}
	blockWithTx, err := client.GetBlock(hash)
	if err != nil {
		panic(err)
	}
	for _, tx := range blockWithTx.Tx {
		if tx.Txid == "e946bcee3b6ac0e39fbddc285a0c5b4790dfb154fcee0edb9766753ff1874808" {
			t.Log(tx)
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

func TestDepositTransaction(t *testing.T) {
	utxoSet, err := client.Scantxoutset("tb1qn9fpljh5ggp407z02jx8x76pemzclgd6rla0qp")
	if err != nil {
		t.Fatal(err)
	}
	if len(utxoSet.Unspents) == 0 {
		t.Fatal("no utxo found")
	}
	t.Logf("utxoSet: %v\n", len(utxoSet.Unspents))
	amount := big.NewInt(1000000)
	minerFee := big.NewInt(100000)
	total := big.NewInt(0)
	var inputs []bitcoin.TxIn
	for _, tUtxo := range utxoSet.Unspents {
		fmt.Printf("utxoId:%v index:%v scriptPubKey:%v\n", tUtxo.Txid, tUtxo.Vout, tUtxo.ScriptPubKey)
		inputs = append(inputs, bitcoin.TxIn{
			Hash:     tUtxo.Txid,
			VOut:     uint32(tUtxo.Vout),
			PkScript: tUtxo.ScriptPubKey,
			//Amount:   floatBig.Sign(),
		})
		total = total.Add(total, big.NewInt(int64(tUtxo.Amount*100000000)))
	}
	findChange := total.Sub(total, amount).Sub(total, minerFee)
	outputs := []bitcoin.TxOut{
		{
			Address: "tb1q4sxzxjxuz8lgx0s4g0hspn8v6g8pvx6juj0lgraglq4q6lnn649suxs3ws",
			Amount:  amount.Int64(),
		},
		{
			Address: "tb1qn9fpljh5ggp407z02jx8x76pemzclgd6rla0qp",
			Amount:  findChange.Int64(),
		},
	}
	secret := ethCommon.FromHex("0x084243403ea5c01337388b2068f98d90a845a9f8926fa16631b07dae4e64a5cd")
	ethAddr := ethCommon.FromHex("0x2A6443B5838f9524970c471289AB22f399395Ff6")
	result, err := bitcoin.CreateDepositTransaction(secret, ethAddr, inputs, outputs, bitcoin.TestNet)
	if err != nil {
		t.Fatal(err)
	}
	rawByres := ethCommon.Bytes2Hex(result)
	txHash, err := client.Sendrawtransaction(rawByres)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}
