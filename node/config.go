package node

type Config struct {
	NodeConfig  NodeConfig   `json:"node"`
	DFinity     DFinity      `json:"dfinity"`
	DbConfig    DbConfig     `json:"db"`
	SeverConfig ServerConfig `json:"sever_config"`
}

type ServerConfig struct {
	IP   string
	Port string
}

type DbConfig struct {
	Path    string
	Cache   int
	Handler int
}

type NodeConfig struct {
	DataDir         string `json:"datadir"`
	Network         string `json:"network"`
	Rpcbind         string `json:"rpcbind"`
	RpcPort         string `json:"rpcport"`
	BtcUrl          string `json:"btcUrl"`
	BtcUser         string `json:"btcUser"`
	BtcPwd          string `json:"btcPwd"`
	BtcNetwork      string `json:"btcNetwork"`
	BTcBtcBlockTime int64  `json:"btcBlockTime"`
	BtcOperatorAddr string `json:"btcOperatorAddr"`
	EthUrl          string `json:"ethUrl"`
	EthBlockTime    int64  `json:"ethBlockTime"`
	ProofUrl        string `json:"proofUrl"`
}

type EthConfig struct {
	Url       string `json:"url"`
	BlockTime int64  `json:"block_time"`
	Secret    Secret `json:"secret"`
}

type DFinity struct {
	Url    string `json:"url"`
	Secret Secret `json:"secret"`
}

type Secret struct {
}

func localDevDaemonConfig() Config {
	return Config{
		NodeConfig: NodeConfig{
			BtcUrl:          "http://127.0.0.1:8332",
			BtcUser:         "lightec",
			BtcPwd:          "abcd1234",
			BtcNetwork:      "regtest",
			BTcBtcBlockTime: 15,
			BtcOperatorAddr: "testOperatorAddr",

			EthUrl:       "",
			EthBlockTime: 15,
		},
		DFinity: DFinity{
			Url:    "http://127.0.0.1:8000",
			Secret: Secret{},
		},
		DbConfig: DbConfig{
			Path:    "/Users/red/.daemon/testnet/data",
			Cache:   1000,
			Handler: 1000,
		},
		SeverConfig: ServerConfig{
			IP:   "127.0.0.1",
			Port: "8089",
		},
	}
}
