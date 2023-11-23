package node

type Config struct {
	Bitcoin  BtcConfig `yaml:"bitcoin"`
	Ethereum BtcConfig `yaml:"ethereum"`
	DFinity  DFinity   `yaml:"dfinity"`
}

type BtcConfig struct {
	Url string `json:"url"`
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
