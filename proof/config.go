package proof

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"os"
)

type Config struct {
	RpcBind        string      `json:"rpcbind"`
	RpcPort        string      `json:"rpcport"`
	DiscordHookUrl string      `json:"discordHookUrl"`
	Url            string      `json:"url"`
	MaxNums        int         `json:"maxNums"`
	Network        string      `json:"network"`
	DataDir        string      `json:"datadir"`
	Mode           common.Mode `json:"mode"` // rpcServer | client
	ProofType      []string    `json:"proofType"`
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
	if c.Mode == common.Client || c.Mode == common.Custom {
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

func (c *Config) GetZkProofTypes() ([]common.ZkProofType, error) {
	if len(c.ProofType) != 0 {
		logger.Debug("proof types:%v", c.ProofType)
		return toZkProofType(c.ProofType)
	}
	zkProofTypes := common.GetEnvZkProofTypes()
	var zkEnvProofTypes []string
	err := json.Unmarshal([]byte(zkProofTypes), &zkEnvProofTypes)
	if err != nil {
		return nil, err
	}
	logger.Debug("proof types:%v", zkEnvProofTypes)
	return toZkProofType(zkEnvProofTypes)
}

func (c *Config) Info() string {
	return fmt.Sprintf("ip:%v,port:%v,Maxnums:%v,network:%v,datadir:%v",
		c.RpcBind, c.RpcPort, c.MaxNums, c.Network, c.DataDir)
}

func NewTestClientModeConfig() Config {
	return Config{
		Url:     "http://127.0.0.1:9970",
		MaxNums: 1,
		Mode:    common.Client,
		Network: "local",
		DataDir: "/Users/red/lworkspace/lightec/daemon/proof/test",
	}
}

func NewTestClusterModeConfig() Config {
	return Config{
		RpcBind: "0.0.0.0",
		RpcPort: "30001",
		MaxNums: 1,
		Mode:    common.Cluster,
		Network: "local",
		DataDir: "/Users/red/lworkspace/lightec/daemon/proof/test",
	}
}

func NewTestCustomModeConfig() Config {
	return Config{
		Url:     "ws://127.0.0.1:8970/ws",
		MaxNums: 1,
		Mode:    common.Custom,
		Network: "local",
		DataDir: "/Users/red/lworkspace/lightec/daemon/proof/test",
	}
}

const RpcRegisterName = "zkbtc"
