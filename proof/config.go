package proof

import (
	"fmt"
	"os/user"
)

const RpcRegisterName = "zkbtc"

type Config struct {
	RpcBind      string `json:"rpcbind"`
	RpcPort      string `json:"rpcport"`
	ParallelNums int    `json:"parallelNums"`
	Network      string `json:"network"`
	DataDir      string `json:"datadir"`
}

func (c *Config) Info() string {
	return fmt.Sprintf("ip:%v,port:%v,parallelNums:%v,network:%v,datadir:%v",
		c.RpcBind, c.RpcPort, c.ParallelNums, c.Network, c.DataDir)
}

func LocalDevConfig() Config {
	current, err := user.Current()
	if err != nil {
		panic("user.Current() error: " + err.Error())
	}
	return Config{
		RpcBind:      "0.0.0.0",
		RpcPort:      "30001",
		ParallelNums: 1,
		Network:      "testnet",
		DataDir:      fmt.Sprintf("%v/.daemon/node.json", current.HomeDir),
	}
}
