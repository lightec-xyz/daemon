package proof

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
)

const RpcRegisterName = "zkbtc"

type Config struct {
	RpcBind string `json:"rpcbind"`
	RpcPort string `json:"rpcport"`

	Url     string      `json:"url"`
	MaxNums int         `json:"maxNums"`
	Network string      `json:"network"`
	DataDir string      `json:"datadir"`
	Mode    common.Mode `json:"model"` // server | client
}

func (c *Config) Check() error {
	if c.MaxNums == 0 {
		return fmt.Errorf("maxNums is empty")
	}
	if c.DataDir == "" {
		return fmt.Errorf("datadir is empty")
	}
	if c.Network == "" {
		return fmt.Errorf("network is empty")
	}
	if c.Mode == "" {
		return fmt.Errorf("model is empty,please select client or server")
	}
	if c.Mode == common.Client {
		if c.Url == "" {
			return fmt.Errorf("url is empty")
		}
	} else if c.Mode == common.Cluster {
		if c.RpcBind == "" {
			return fmt.Errorf("rpcbind is empty")
		}
		if c.RpcPort == "" {
			return fmt.Errorf("rpcport is empty")
		}
	} else {
		return fmt.Errorf("unknown model:%v", c.Mode)
	}
	return nil
}

func (c *Config) Info() string {
	return fmt.Sprintf("ip:%v,port:%v,Maxnums:%v,network:%v,datadir:%v",
		c.RpcBind, c.RpcPort, c.MaxNums, c.Network, c.DataDir)
}

func NewClientModeConfig() Config {
	return Config{
		Url:     "http://127.0.0.1:9780",
		MaxNums: 1,
		Mode:    common.Client,
		Network: "local",
		DataDir: "/Users/red/lworkspace/lightec/daemon/proof/test",
	}
}

func NewClusterModeConfig() Config {
	return Config{
		RpcBind: "0.0.0.0",
		RpcPort: "30001",
		MaxNums: 1,
		Mode:    common.Cluster,
		Network: "local",
		DataDir: "/Users/red/lworkspace/lightec/daemon/proof/test",
	}
}
