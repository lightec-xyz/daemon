package proof

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
)

type Config struct {
	RpcBind             string      `json:"rpcbind"`
	RpcPort             string      `json:"rpcport"`
	DiscordHookUrl      string      `json:"discordHookUrl"`
	Url                 string      `json:"url"`
	MaxNums             int         `json:"maxNums"`
	CacheCap            int         `json:"cacheCap"`
	Network             string      `json:"network"`
	DataDir             string      `json:"datadir"`
	DisableVerifyZkFile bool        `json:"disableVerifyZkFile"`
	Mode                common.Mode `json:"mode"` // rpcServer | client
	ProofType           []string    `json:"proofType"`
	BtcSetupDir         string      `json:"btcSetupDir"`
	EthSetupDir         string      `json:"ethSetupDir"`
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
		c.Network = "testnet" // todo
	}
	if c.Url == "" {
		c.Url = "https://testnet.zkbtc.money/api"
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

	if c.BtcSetupDir == "" || c.EthSetupDir == "" {
		return fmt.Errorf("btcSetupDir or ethSetupDir is empty")
	}

	return nil
}

func (c *Config) GetZkProofTypes() ([]common.ProofType, error) {
	if len(c.ProofType) != 0 {
		logger.Debug("proof types:%v", c.ProofType)
		return toZkProofType(c.ProofType)
	}
	zkProofTypes := common.GetEnvZkProofTypes()
	if zkProofTypes != "" {
		var zkEnvProofTypes []string
		err := json.Unmarshal([]byte(zkProofTypes), &zkEnvProofTypes)
		if err != nil {
			return nil, err
		}
		logger.Debug("proof types:%v", zkEnvProofTypes)
		return toZkProofType(zkEnvProofTypes)
	}
	return nil, nil

}

func (c *Config) Info() string {
	return fmt.Sprintf("ip:%v,port:%v,Maxnums:%v,network:%v,datadir:%v",
		c.RpcBind, c.RpcPort, c.MaxNums, c.Network, c.DataDir)
}

func NewTestClientModeConfig() Config {
	return Config{
		Url:         "http://127.0.0.1:10977",
		MaxNums:     1,
		Mode:        common.Client,
		Network:     "testnet",
		DataDir:     "./test/generator",
		BtcSetupDir: "/opt/testnet/setup",
		EthSetupDir: "/opt/testnet/setup",
	}
}

func NewTestClusterModeConfig() Config {
	return Config{
		RpcBind: "0.0.0.0",
		RpcPort: "30001",
		MaxNums: 1,
		Mode:    common.Cluster,
		Network: "local",
		DataDir: "./test/generator",
	}
}

func NewTestCustomModeConfig() Config {
	return Config{
		Url:     "ws://127.0.0.1:8970/ws",
		MaxNums: 1,
		Mode:    common.Custom,
		Network: "local",
		DataDir: "./test/generator",
	}
}

const RpcRegisterName = "zkbtc"
