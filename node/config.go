package node

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
	BtcWhiteList     []string         `json:"btcWhiteList"`
	MultiAddressInfo MultiAddressInfo `json:"multiAddressInfo"`
	EthWhiteList     []string         `json:"ethWhiteList"`
	EthUrl           string           `json:"ethUrl"`
	ZkBridgeAddr     string           `json:"zkBridgeAddr"`
	EthBlockTime     int64            `json:"ethBlockTime"`
	EthPrivateKey    string           `json:"ethPrivateKey"`
	ProofUrl         string           `json:"proofUrl"`

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

func localDevDaemonConfig() NodeConfig {
	return NodeConfig{
		DataDir:         "/Users/red/.daemon",
		Network:         "testnet",
		Rpcbind:         "127.0.0.1",
		RpcPort:         "8899",
		BtcUrl:          "https://go.getblock.io/d54c59f635654cc082de1f3fd14e5d02",
		BtcUser:         "lightec",
		BtcPwd:          "abcd1234",
		BtcNetwork:      "regtest",
		BTcBtcBlockTime: 15,
		BtcOperatorAddr: "testOperatorAddr",
		BtcWhiteList: []string{
			"",
		},
		EthUrl:        "https://rpc.notadegen.com/eth/sepolia",
		ZkBridgeAddr:  "0x8dda72ee36ab9c91e92298823d3c0d4d73894081",
		EthBlockTime:  15,
		EthPrivateKey: "c0781e4ca498e0ad693751bac014c0ab00c2841f28903e59cdfe1ab212438e49",
		EthWhiteList: []string{
			"",
		},
		MultiAddressInfo: MultiAddressInfo{
			PublicKeyList: [][]byte{
				{1, 2, 3, 4},
				{1, 2, 3, 5},
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
