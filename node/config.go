package node

import (
	"encoding/hex"
	"fmt"
	"os/user"
)

type NodeConfig struct {
	DataDir          string           `json:"datadir"`
	Network          string           `json:"network"`
	Rpcbind          string           `json:"rpcbind"`
	RpcPort          string           `json:"rpcport"`
	BtcUrl           string           `json:"btcUrl"`
	BtcUser          string           `json:"btcUser"`
	BtcPwd           string           `json:"btcPwd"`
	BtcNetwork       string           `json:"btcNetwork"`
	BTcBtcBlockTime  int64            `json:"btcBlockTime"`
	BtcOperatorAddr  string           `json:"btcOperatorAddr"`
	BtcPrivateKeys   []string         `json:"btcPrivateKeys"`
	BtcWhiteList     []string         `json:"btcWhiteList"`
	BtcInitHeight    int64            `json:"btcInitHeight"`
	MultiAddressInfo MultiAddressInfo `json:"multiAddressInfo"`
	EthInitHeight    int64            `json:"ethInitHeight"`
	EthWhiteList     []string         `json:"ethWhiteList"`
	EthUrl           string           `json:"ethUrl"`
	ZkBridgeAddr     string           `json:"zkBridgeAddr"`
	EthBlockTime     int64            `json:"ethBlockTime"`
	EthPrivateKey    string           `json:"ethPrivateKey"`
	LogAddr          []string         `json:"logAddr"`
	LogTopic         []string         `json:"logTopic"`

	Workers []WorkerConfig `json:"workers"`
}

type MultiAddressInfo struct {
	PublicKeyList [][]byte
	NRequired     int
}

type WorkerConfig struct {
	ParallelNums int    `json:"parallelNums"`
	ProofUrl     string `json:"proofUrl"`
}

func LocalDevDaemonConfig() NodeConfig {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	pub1, err := hex.DecodeString("03bd96c4d06aa773e5d282f0b6bccd1fb91268484918648ccda1ae768209edb050")
	if err != nil {
		panic(err)
	}
	pub2, err := hex.DecodeString("03aa9c4245340a02864c903f7f9e7bc9ef1cc374093aacbf72b614002f6d8c8c22")
	if err != nil {
		panic(err)
	}
	pub3, err := hex.DecodeString("03351a7971bf7ed886fca99aebdc3b195fc79ffe93b499e2309a4e69ab115405e0")
	if err != nil {
		panic(err)
	}
	return NodeConfig{
		DataDir:         fmt.Sprintf("%v/.daemon", user.HomeDir),
		Network:         "devnet",
		Rpcbind:         "127.0.0.1",
		RpcPort:         "8899",
		BtcUrl:          "http://127.0.0.1:8332",
		BtcUser:         "lightec",
		BtcPwd:          "abcd1234",
		BtcNetwork:      "RegTest",
		BTcBtcBlockTime: 15,
		BtcOperatorAddr: "bcrt1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq2aw6ze",
		BtcPrivateKeys: []string{
			"b26dbaab82d9ebd8f37c88bbe56e22bf9cb21150c96dfb35ece4b787d3710d3301",
			"62dd5835dc2ce7f4f40eea1b88c816043d288532c8bb91964adef9bc0f0b4b7201",
			"9ff573d948c80fa1a50da6f66229b4bede9ec3fb482dd126f58d3acfb4b2979801",
		},
		BtcInitHeight: 2542024,
		EthInitHeight: 481270,
		EthUrl:        "https://ethereum-holesky.publicnode.com",
		ZkBridgeAddr:  "0xc5529b337c2ec0cb075299f68134d9e4d12227c6",
		EthBlockTime:  10,
		EthPrivateKey: "c0781e4ca498e0ad693751bac014c0ab00c2841f28903e59cdfe1ab212438e49",
		LogAddr:       []string{"0xc5529b337c2ec0cb075299f68134d9e4d12227c6"},
		LogTopic:      []string{"0xb28ad0403b0a341130002b9eef334c5daa3c1002a73dd90d4626f7079d0a804a"},
		MultiAddressInfo: MultiAddressInfo{
			PublicKeyList: [][]byte{
				pub1, pub2, pub3,
			},
			NRequired: 2,
		},
		Workers: []WorkerConfig{
			{
				ParallelNums: 3,
				ProofUrl:     "http://127.0.0.1:8485",
			},
		},
	}
}

func TestnetDaemonConfig() NodeConfig {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	pub1, err := hex.DecodeString("03bd96c4d06aa773e5d282f0b6bccd1fb91268484918648ccda1ae768209edb050")
	if err != nil {
		panic(err)
	}
	pub2, err := hex.DecodeString("03aa9c4245340a02864c903f7f9e7bc9ef1cc374093aacbf72b614002f6d8c8c22")
	if err != nil {
		panic(err)
	}
	pub3, err := hex.DecodeString("03351a7971bf7ed886fca99aebdc3b195fc79ffe93b499e2309a4e69ab115405e0")
	if err != nil {
		panic(err)
	}
	return NodeConfig{
		DataDir:         fmt.Sprintf("%v/.daemon", user.HomeDir),
		Network:         "testnet",
		Rpcbind:         "127.0.0.1",
		RpcPort:         "8899",
		BtcUrl:          "https://go.getblock.io/d54c59f635654cc082de1f3fd14e5d02",
		BtcUser:         "lightec",
		BtcPwd:          "abcd1234",
		BtcNetwork:      "TestNet",
		BTcBtcBlockTime: 10 * 60,
		BtcOperatorAddr: "tb1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq8yyuhr",
		BtcPrivateKeys: []string{
			"b26dbaab82d9ebd8f37c88bbe56e22bf9cb21150c96dfb35ece4b787d3710d3301",
			"62dd5835dc2ce7f4f40eea1b88c816043d288532c8bb91964adef9bc0f0b4b7201",
			"9ff573d948c80fa1a50da6f66229b4bede9ec3fb482dd126f58d3acfb4b2979801",
		},
		BtcInitHeight: 2542051,
		EthInitHeight: 481270,
		EthUrl:        "https://ethereum-holesky.publicnode.com",
		ZkBridgeAddr:  "0xc5529b337c2ec0cb075299f68134d9e4d12227c6",
		EthBlockTime:  10,
		EthPrivateKey: "c0781e4ca498e0ad693751bac014c0ab00c2841f28903e59cdfe1ab212438e49",
		LogAddr:       []string{"0xc5529b337c2ec0cb075299f68134d9e4d12227c6"},
		LogTopic:      []string{"0xb28ad0403b0a341130002b9eef334c5daa3c1002a73dd90d4626f7079d0a804a"},
		MultiAddressInfo: MultiAddressInfo{
			PublicKeyList: [][]byte{
				pub1, pub2, pub3,
			},
			NRequired: 2,
		},
		Workers: []WorkerConfig{
			{
				ParallelNums: 3,
				ProofUrl:     "http://127.0.0.1:8485",
			},
		},
	}
}
