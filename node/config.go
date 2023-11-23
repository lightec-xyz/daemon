package node

type Config struct {
	Bitcoin  BtcConfig `yaml:"bitcoin"`
	Ethereum EthConfig `yaml:"ethereum"`
	DFinity  DFinity   `yaml:"dfinity"`
	DbConfig DbConfig  `json:"db_config"`
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
	Url    string `json:"url"`
	Secret Secret `json:"secret"`
}

type DFinity struct {
	Url    string `json:"url"`
	Secret Secret `json:"secret"`
}

type Secret struct {
}
