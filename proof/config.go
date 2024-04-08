package proof

import (
	"fmt"
	"os/user"
)

const RpcRegisterName = "zkbtc"

type Config struct {
	RpcBind string `json:"rpcbind"`
	RpcPort string `json:"rpcport"`
	MaxNums int    `json:"maxNums"`
	Network string `json:"network"`
	DataDir string `json:"datadir"`
	Model   string `json:"model"`
}

func (c *Config) Info() string {
	return fmt.Sprintf("ip:%v,port:%v,Maxnums:%v,network:%v,datadir:%v",
		c.RpcBind, c.RpcPort, c.MaxNums, c.Network, c.DataDir)
}

func LocalDevConfig() Config {
	current, err := user.Current()
	if err != nil {
		panic("user.Current() error: " + err.Error())
	}
	return Config{
		RpcBind: "0.0.0.0",
		RpcPort: "30001",
		MaxNums: 1,
		Network: "testnet",
		DataDir: fmt.Sprintf("%v/.daemon/node.json", current.HomeDir),
	}
}
