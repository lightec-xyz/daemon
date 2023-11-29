package node

type NodeConfig struct {
	DataDir         string   `json:"datadir"`
	Network         string   `json:"network"`
	Rpcbind         string   `json:"rpcbind"`
	RpcPort         string   `json:"rpcport"`
	BtcUrl          string   `json:"btcUrl"`
	BtcUser         string   `json:"btcUser"`
	BtcPwd          string   `json:"btcPwd"`
	BtcNetwork      string   `json:"btcNetwork"`
	BTcBtcBlockTime int64    `json:"btcBlockTime"`
	BtcOperatorAddr string   `json:"btcOperatorAddr"`
	BtcWhiteList    []string `json:"btcWhiteList"`
	EthWhiteList    []string `json:"ethWhiteList"`
	EthUrl          string   `json:"ethUrl"`
	EthBlockTime    int64    `json:"ethBlockTime"`
	EthPrivateKey   string   `json:"ethPrivateKey"`
	ProofUrl        string   `json:"proofUrl"`

	Workers []WorkerConfig `json:"workers"`
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
		BtcUrl:          "http://127.0.0.1:8332",
		BtcUser:         "lightec",
		BtcPwd:          "abcd1234",
		BtcNetwork:      "regtest",
		BTcBtcBlockTime: 15,
		BtcOperatorAddr: "testOperatorAddr",
		BtcWhiteList: []string{
			"",
		},
		EthUrl:        "http://127.0.0.1:8332",
		EthBlockTime:  15,
		EthPrivateKey: "testprivateKey",
		EthWhiteList: []string{
			"",
		},
		Workers: []WorkerConfig{
			{
				ParallelNums: 3,
				ProofUrl:     "http://127.0.0.1:8485",
			},
		},
	}
}
