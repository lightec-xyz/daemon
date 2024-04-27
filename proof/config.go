package proof

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"os"
)

type Config struct {
	RpcBind string `json:"rpcbind"`
	RpcPort string `json:"rpcport"`

	Url     string      `json:"url"`
	MaxNums int         `json:"maxNums"`
	Network string      `json:"network"`
	DataDir string      `json:"datadir"`
	Mode    common.Mode `json:"model"` // rpcServer | client
}

func (c *Config) Check() error {
	if c.MaxNums == 0 {
		c.MaxNums = 1
	}
	if c.DataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		c.DataDir = fmt.Sprintf("%s/.generateor", homeDir)
	}
	if c.Network == "" {
		c.Network = "local" // todo
	}
	if c.Mode == "" {
		c.Mode = common.Client
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

const RpcRegisterName = "zkbtc"
