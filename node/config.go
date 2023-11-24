package node

type Config struct {
	Bitcoin     BtcConfig    `yaml:"bitcoin"`
	Ethereum    EthConfig    `yaml:"ethereum"`
	DFinity     DFinity      `yaml:"dfinity"`
	DbConfig    DbConfig     `json:"db_config"`
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

type BtcConfig struct {
	Url          string `json:"url"`
	User         string `json:"user"`
	Pwd          string `json:"pwd"`
	Network      string `json:"network"`
	BlockTime    int64  `json:"block_time"`
	OperatorAddr string `json:"operator_addr"`
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

func devDaemonConfig() Config {
	return Config{
		Bitcoin: BtcConfig{
			Url:          "https://bitcoin-mainnet-archive.allthatnode.com",
			User:         "user",
			Pwd:          "pwd",
			Network:      "regtest",
			BlockTime:    10,
			OperatorAddr: "user",
		},
		Ethereum: EthConfig{
			Url:       "http://127.0.0.1:8545",
			BlockTime: 10,
			Secret:    Secret{},
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
