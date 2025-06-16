package ethereum

import (
	"bytes"
	"context"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	btctx "github.com/lightec-xyz/daemon/rpc/bitcoin/common"
	"github.com/lightec-xyz/daemon/rpc/ethereum/zkbridge"
	"math/big"
	"testing"
)

var err error
var client *Client

var endpoint = "http://127.0.0.1:9002"
var zkBridgeAddr = "0x21098979Fc10BBC754C6359E657eA28c52ea1acf"
var utxoManager = "0xD2f892d4Ece281C91Fd5D9f28658F8d445878239"
var btcTxVerifyAddr = "0xB4c6946069Ec022cE06F4C8D5b0d2fb232f8DDa5"
var zkbtcAddr = "0xB4c6946069Ec022cE06F4C8D5b0d2fb232f8DDa5"

func init() {
	//https://sepolia.publicgoods.network
	client, err = NewClient(endpoint, zkBridgeAddr, utxoManager, btcTxVerifyAddr, zkbtcAddr)
	if err != nil {
		panic(err)
	}
}

func TestClient_IsCandidateExist(t *testing.T) {
	exist, err := client.IsCandidateExist("0000000000000001e7a798ae790a9df0befa97d78816d8da1ae46f17b27547ed")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(exist)
}

func TestClient_SuggestBtcMinerFee(t *testing.T) {
	fee, err := client.SuggestBtcMinerFee()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fee)
}

func TestClient_GetRaised(t *testing.T) {
	raised, err := client.GetRaised("0000000000000001e7a798ae790a9df0befa97d78816d8da1ae46f17b27547ed", 100)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(raised)
}

func TestClient_GetCpLatestAddedTime(t *testing.T) {
	time, err := client.GetCpLatestAddedTime()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(time)
}

func TestClient_EthBalance(t *testing.T) {
	balance, err := client.EthBalance("0x771815eFD58e8D6e66773DB0bc002899c00d5b0c")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", balance)
}

func TestClient_SuggestedCP(t *testing.T) {
	cp, err := client.SuggestedCP()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%x", cp)
}

func TestClient_GetMinTxDepth(t *testing.T) {
	depth, err := client.GetDepthByAmount(87392, false, false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(depth)
}

func TestClient_btcRaw(t *testing.T) {
	receipt, err := client.TransactionReceipt(context.Background(), ethcommon.HexToHash("0x72bafb9f6024516acf030eab95a159e8ef0c069e3b92aa9a55fced458b423baf"))
	if err != nil {
		t.Fatal(err)
	}
	//todo
	btcRaw, sigs, err := DecodeRedeemLog(receipt.Logs[3].Data)
	if err != nil {
		t.Fatal(err)
	}
	for _, hash := range sigs {
		t.Logf("%x", hash)
	}
	t.Logf("%x", btcRaw)
}

func TestClient_CheckEndpointHash(t *testing.T) {
	hash, err := client.SuggestedCP()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash)
}

func TestClient_CheckUtxo(t *testing.T) {
	result, err := client.GetUtxo("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestClient_Demo001(t *testing.T) {
	receipt, err := client.TransactionReceipt(context.Background(), ethcommon.HexToHash("0xb19639d5c7c5804632f8ed92ca7e16d78cc1c6590a314b0aafee78793be223c6"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(receipt.TransactionIndex)
}

func TestClient_GetTxSender(t *testing.T) {
	sender, err := client.GetTxSender("0xb19639d5c7c5804632f8ed92ca7e16d78cc1c6590a314b0aafee78793be223c6",
		"0xf99ab49c39e77bd6274035cbc1d6db068e014d3dc8e8a6a4c988f327a9b417f1", 39)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sender)
}

func TestClient_Number(t *testing.T) {
	number, err := client.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(number)
}

func TestClient_GetLogs2(t *testing.T) {
	//0x0c7fee5ab535f4842d895aae6de266e0b51b3540327fe03eedb77ad798637e00
	logs, err := client.GetLogs("0xd936a94eabfe6a9cb84382515a99684170271e06c676c1b89c2eed4baf953d08", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, log := range logs {
		if log.TxHash.String() == "0xff1ddc991b66997739b70e846d562cc11dd5012487ce9b12b40c555a71bd6f2d" {
			t.Log(log)
		}
	}
}

func TestClient_GetLogs(t *testing.T) {
	block, err := client.GetBlock(1545882)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log(block)
	address := []string{"0x9d2aaea60dee441981edf44300c26f1946411548", "0x8e4f5a8f3e24a279d8ed39e868f698130777fded"}
	topic := []string{"0xbfb6a0aa850eff6109c854ffb48321dcf37f02d6c7a44c46987a5ddf3419fc07", "0x1e5e2baa6d11cc5bcae8c0d1187d7b9ebf13d6d9b932f7dbbf4e396438845fb8"}
	logs, err := client.GetLogs(block.Hash().Hex(),
		address, topic)
	if err != nil {
		t.Fatal(err)
	}
	//t.Log(logs)
	for _, log := range logs {
		if log.TxHash.String() == "0xea7a29093b228e8d45ba54161689e1ae7c4caa1ce33fd618112eace20e2acf1a" {
			txData, _, err := DecodeRedeemLog(log.Data)
			if err != nil {
				t.Fatal(err)
			}
			transaction := btctx.NewTransaction()
			err = transaction.Deserialize(bytes.NewReader(txData))
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(transaction.TxHash().String())
			//0x0020957ab85b710cb5b577171e23bb3492536c8029cc99511f3920d3cc13871a2327
			for _, out := range transaction.TxOut {
				t.Logf("%x %v \n", out.PkScript, out.Value)
			}

		}
	}
}

func TestClient_UpdateUtxoChange(t *testing.T) {
	secret := GetTestSecret()
	address, err := privateKeyToAddr(secret)
	if err != nil {
		t.Fatal(err)
	}
	minerReward := big.NewInt(71744000000000)
	txId := ethcommon.FromHex("6b60d9aeacc88ba21ccaf4bea35446e801cd26c9bbd334d901e43cb573aa1cb9")
	proof := ethcommon.FromHex("0e32d5593816fb6c9496c53d2f7653e06da18ce4c5fec7ab5e3359d3fe938bf02d7922dc3ab48a1422651cbbb6eb5aa5e5006f25afcfd076cdce2215ef7c60f40958cd7b2fc16833e3fe303d72e873eb8360365452c81325f635b2f885ffd20e108bb3b6028963aaa87eedd7385d1b005785491b69e7cd16322ee2b09bb8cd2f006d246cb1564a6d978da27a9b62c71c0055af2193af7ed117f1c6742987ad591f631b7c34fceb5ef65a641512845755ba25a5af57caf15353f172f711ace4e51a748d17dea4106d26016636d44411a9916a329ed181cfff78e9a37189c45a08287789293a0adc9cc370f72a5bf9976db249b62689970f6efa1ac23d6c415e202f8ec38f0a66ef89f782225b763b8dab6bb5f871d08d3681ced5089216985cbe16d7897d0d2511d1d87949a38c012b728c8cd91a5e3769c95112e46b2481e52a18383b52ca98bd490fc55c4e1b0e7396f37ee63f8c00a0a0083980d095b4070516b871e77b40ef8ac1826441319edaee47bb22531d30b28b5f833d90e1f21f352a26a99f032c5a0788fef7f22be540569ed82f29009aade8592c2bcc7074c8de1f9c1a27cb8cf93f7ad3db71fa7e877fc9f81c7b9b5b8f77aa6ce814dc33101d2a2cafc6258c5b72c0229a6cc5b89792264e3feb7cdc234c3ce2c98b39ee8d4d2bd422cfccda1bc5a11a5ef937d80e705313d6537d9e81abf7f757185dd8eb201987c4aa683c759e63e586b17c8605c888666a3d6af322f2da86191be56695502f0fbe392c4be170748aaa1643c2094500126a250510cf1e2a947d98289d96b11d54224c98a154c53f95dcdd86a48f9e9915b2fcb96934063b6b04d36630bc312d3c5b46f092f698ad2b612e830891c4c25044ed97a59ce610d4053c5ee84fbf10de32ee68d7e364f3e7e4e9fd175ddc3c503b51f33e007e44cdf6e734b4782a267d4d9a835b164b9c99dc770622fcebb9a53ce4c73068500e45f17a31b917c21335917d87d2a1d924221decd7ec87b9f4f6c2475c8a77b7da3ceffce6f12a1b305baa81a230034e0ddeefdf9c1072e5bfadcab301b4f124a2b04db1eb42bc8c208ccb13c5e83664cf04e3943704b714ec31122c6578dcdb869634920e9c064625dff79c6fcf1e086d99a4ead282cfe5009aba3b9b2da2e791fefe1881a4d0fc17867951ab36142debbb92fc9350e8f47ee5c6ae2d1f05648fdd0f585974b6ab")
	nonce, err := client.GetNonce(address)
	if err != nil {
		t.Fatal(err)
	}
	gasPrice, err := client.GetGasPrice()
	if err != nil {
		t.Fatal(err)
	}
	chainId, err := client.GetChainId()
	if err != nil {
		t.Fatal(err)
	}
	gasLimit := uint64(500000)
	params := zkbridge.IBtcTxVerifierPublicWitnessParams{
		Checkpoint:        [32]byte(ethcommon.FromHex("69e71ff8334517efa8d94aaed13087a99b5a484d0aead746da44f90000000000")),
		CpDepth:           1542,
		TxDepth:           24,
		TxBlockHash:       [32]byte(ethcommon.FromHex("50318f29f557d5dcb1347b5af955ce1bbe098b9bceaf8aff0efa506200000000")),
		TxTimestamp:       1749794078,
		ZkpMiner:          ethcommon.HexToAddress(address),
		Flag:              big.NewInt(2),
		SmoothedTimestamp: 1749794078,
	}
	hash, err := client.UpdateUtxoChange(ethcommon.FromHex(secret), &params, nonce, gasLimit, chainId,
		gasPrice, minerReward, txId, proof)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash)
}

func TestClient_Redeem(t *testing.T) {
	privateKey := GetTestSecret()
	redeemAmount := uint64(1000)
	minerFee := uint64(3000)
	gasLimit := uint64(500000)
	from, err := privateKeyToAddr(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	redeemLockScript := ethcommon.FromHex("")
	gasPrice, err := client.GetGasPrice()
	if err != nil {
		t.Fatal(err)
	}
	chainID, err := client.GetChainId()
	if err != nil {
		t.Fatal(err)
	}
	nonce, err := client.GetNonce(from)
	if err != nil {
		t.Fatal(err)
	}
	txhash, err := client.Redeem(privateKey, gasLimit, chainID, big.NewInt(int64(nonce)), gasPrice, redeemAmount, minerFee, redeemLockScript)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txhash)
}

func TestClient_EstimateDepositGasLimitEstDepositeTx(t *testing.T) {
	proof := ethcommon.FromHex("1642857fe0b6d1cb289b97d381b9e78bdd7fd50b9cafd6574baa17219c0d514c103d2cf27ccec6cec2eca3f69e9b31895928d3ee117e45c00a7645acae77526215c8eff0d6e354068226b84fd50c4357e0f65ae0f3fc81a77bd3906c197f780429a465f41dda167e39a9ede8f260ab4bc5b0400cadcaec285b7dcf67b4dd892803b823bbdd9f5d0bdd30779fdde7c31183951d55eaafa9ff30b06317dc2a54220517324dbb2f51df1796816df299063e80e6e09cb82e39fc230ea544f127104d23b37d217e7794259e0607077bb0922fa762037a98e63931900433d88a96c68514731fc6762de17442a3b853fd304cbc9e848e86e3cbe211481b9b4e31b3158b053a852b2d4b9c2fb05733f5f5caf27af42b510f9d142636d7e08a9a2a547d70079d9ee741fcb49d22ba1727d60052ffa206b271b83027597f58973e7eaabe911a4c0f68e9100ed45a558f55f46c83d18f98c5c8cc6a854f7d6c16ca979f4c832a7d7f628735213e76be65a8acc20fa18bffb9bd7fb5e8dcb1fec7eddb2e811b1a049cf212e43f50607e2599303923968f9363bd25d10550cde366ffbdbd06192ca9c088b7cadb673366049b21c589d31ae95cd875d69c524b314964e51e1f3800eca940e234b21de73dedbaf7bf033b3746f1e8b6e6859fa612d96c4f722a0417dc5cde13870ff4fc15ba03204a96c4f1672ad91b7b02ea0cadf1a3ff42b2cb10ef043e9e54ff54e3ee81bace9b6cd7f094dd14ce4d96923ab05855d7dd35870d2ee424844a51e8e9dd66c226e4fc071b9f6b29e7fe8e8936147144f583e2930007db3a1279aa8909e83b7fe011801c80844c7d30f1428359f0341503d0c0731ecdad2d8312d0d7ee54cc9f1cd61fe612f1b4231310f8be18d9d5e8195bdf322537a16c6ae4a45364867c79a4732271edba476514e5127ea99bb12212ef89ef2d19159d9235d1259ad4da38c2dd8fbf7f818cf1e88264e3fc775d261447f705096e32f4c9db2f0134b1e3ff62812ac56ca0555c04b2ca0af7a14e05cef1162a163f41785fcecd057132abdf89c71e652f48fa83d86ebef0933db23358b55c35283a802f9fc7b712c05d35d7dcaaf0aa6241c239a251394f9df389a0d2c9e30e09a8d112876fb17f1e59340b25003dccdc0572635f2dcf9ce43edd006f84d72a10906ddd1789e55dd9af0f67758307ae3b88e11cddd1f3c2ca2151332b755752")
	btcTxRaw := ethcommon.FromHex("020000000001034dfae386f7a6529c47cf9face47462a336ee3708ee7567010b6378af19ab7d1f0000000000ffffffffb8f66f569c76e79e9606d5b91b9038e7286a62f8d3873a6773c227dcb5878c960200000000ffffffff12aee67c0689dcd5f6654b8a4ec754c45cab13c29dc131263c0962a93c1ae7eb0000000000ffffffff03a0860100000000002200200210550035cea0e86c7eaed74de90fd9c91e36c17e6b03f7856ca24c9cb20ebb0000000000000000166a142a6443b5838f9524970c471289ab22f399395ff6324202000000000016001499521fcaf4420357f84f548c737b41cec58fa1ba02483045022100cf41b9504d66ac37c7fa4b635c14f2baf8f173e4b2cac195b2d408dccc3b96e202200cc426ea0958c7d29aa56f8dd358b3a3b4100fca1689f972c5506a8a263fdcf7012103d550bfd354d2c4aec5959c95a4c656a191d93350838195277fe97604d0c0c5ce0247304402204f47737c6b488cd44e9b838d22b583b62e648a89e6df1b35a4be2fc220ad96ea0220071cb26a02b5d068c1ad1f7c182b57c77c0711d21568a9873d0ae03a86d5ec53012103d550bfd354d2c4aec5959c95a4c656a191d93350838195277fe97604d0c0c5ce02483045022100928e2e86c88686a84b0bc7b160e0a70977029c34387caffdcfc01722a7bd5ce00220114a72ee9a95d86c698b058d474a75fa20078a1251f86c6db36fe6aacc73d95d012103d550bfd354d2c4aec5959c95a4c656a191d93350838195277fe97604d0c0c5ce00000000")
	price, err := client.GetGasPrice()
	if err != nil {
		t.Fatal(err)
	}
	checkPoint := ethcommon.FromHex("35325811e9c78be9cfb2db80e52fc06415e38ed0165aa493850026b900000000")
	blockHash := ethcommon.FromHex("0000000000000005bb9ccd58a6cd295772eec9012e840e5de23611db072bc9f1")
	parms := &zkbridge.IBtcTxVerifierPublicWitnessParams{
		Checkpoint:  [32]byte(checkPoint),
		CpDepth:     uint32(1),
		TxDepth:     uint32(1),
		TxBlockHash: [32]byte(blockHash),
		TxTimestamp: uint32(0),
		ZkpMiner:    ethcommon.HexToAddress(""),
	}
	gasLimit, err := client.EstimateDepositGasLimit("", parms, price, btcTxRaw, proof)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gasLimit)
}

func TestClient_EstimateUpdateUtxoGasLimit(t *testing.T) {
	proof := ethcommon.FromHex("0cbdea4d915a8029ec3f10b74d551d36776865751d67946bff9d7d5c21d23800268ef44a9d203b2bd9d8bad202a01608cb9042bd7fe83e52c7221473c96c902f0798d812ff6d8c86b940af8dd819af4b58a69f001852fc140a54f2c6fc066ac918780194fdbe69f050d9d337024323f821df73137056c00418e3eb243064632513a9f7b6dd2ea2177ce447369dd0b30278693fd850467982be87ac29cf9aff7a24c69dba59542d9c129ae2a0c2b49baeea9c5f5dc996502d4822bd57e0a51d4729ba43f3b01d4166eb4f82d0ccc29da82d8e8a855c56c2dd54d5abe8dd3b949502b9a8cbf0113efde166fb4b47bb24b03707010bd0e49474f777a7b7c16e20ef288e5a006b29a08a3de8efd1dd15aae648ee8f7cd88847d003d802389e25df252452d6a78241488cbdff29da7423c43a5016c9453657bbe8757c634da55612d312bd7fe09a30d526e50b262f04f7453149e07779b4ef0dadb3fc88349b272622085e3614729e498ff5d9f361785be1d27bb275189834879dbc1e6fbb7b0aeb142ddac43d8809120f17a4fc0c80d2fca718812fe411cb092708fa5e4f9700af45099078501565188003325368b6d94eb8508236c31504e4c39ad21079d37e8cda19cd027c9b1092dac2e8799bfc77b5d5496c00c896671d947bc2f1d2ea3fd3a528c675088d350ea37e8a8f04c212e1394e75e044e410c4fa61f3d0f6738070440b65264844d68e988677c39571153fea736c4d2582dcf3fa5acc7b8e6a264597248344b66ebd60c2c8ca755b13a395d3d82017e7c51e267d7e33fe3b629743f208e801172655027093572ea547cb9893043f699783cef34f5a9794914ce6062c23c5669895eafd2ac86f63e3f3953f64e869d95f1f1082bdd11782107a1f50911327e8293866b39faafe7d73ba5b11e8dd11db165f34d78b75a6f859901c84b111c3958e8becdebfa3faacdc39a2696c08c2c8bd68b7d7cbd71d9eee7af3b2812dff5138b78612c2f7244d4b3c5344aff656fd932583752b88e0a25e593f1c7f1ab1cb3e6a57a4858c58ccbff7fe586b3f6152092afe02cda5ef7fa751409098042e7aa0adef15a4e21665417dfae3f04c688c438603deaa14859ca6df62041425b532f95f5565c0753b05f934f3cf823446566b57888ba76ae2bc3a4a3421d910ffd6c792ac7a77e7c7ce1b97a09bbcd023a1d3447fd71679b4affb80f56995")
	price, err := client.GetGasPrice()
	if err != nil {
		t.Fatal(err)
	}
	minerReward := big.NewInt(256960000000000)
	checkPoint := ethcommon.FromHex("c044412393615746c969e6ef3ce3addfbf4e084e43add07fcd2c6c9a00000000")
	blockHash := ethcommon.FromHex("8d5e6902fdbc578fefa0e4743d5f280c107662a6c84fea030500000000000000")
	txId := ethcommon.FromHex("53353d34ede405715d258108ebd9f8c3a8b6a2abad923484afcc601e98f4bbc4")

	parms := &zkbridge.IBtcTxVerifierPublicWitnessParams{
		Checkpoint:        [32]byte(checkPoint),
		CpDepth:           uint32(335),
		TxDepth:           uint32(9),
		TxBlockHash:       [32]byte(blockHash),
		TxTimestamp:       uint32(1745902584),
		ZkpMiner:          ethcommon.HexToAddress("0x79e0F79F395CEcF90812ee7beFB883D7210e20E8"),
		Flag:              big.NewInt(0),
		SmoothedTimestamp: 0,
	}
	gasLimit, err := client.EstimateUpdateUtxoGasLimit("0x79e0F79F395CEcF90812ee7beFB883D7210e20E8", parms, price, minerReward, txId, proof)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gasLimit)
}

func TestEthTransfer(t *testing.T) {
	txHash, err := client.EthTransfer("", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}
