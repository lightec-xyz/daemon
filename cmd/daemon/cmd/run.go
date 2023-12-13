package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/node"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rpcbind *string
var rpcport *string
var btcUrl *string
var btcUser *string
var btcPwd *string
var ethUrl *string
var ethPrivateKey *string
var workers *[]string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run daemon",
	Long:  `Start daemon program`,
	Run: func(cmd *cobra.Command, args []string) {
		//todo
		config, err := toConfig(viper.AllSettings())
		//config := node.TestnetDaemonConfig()
		cobra.CheckErr(err)
		daemon, err := node.NewDaemon(config)
		cobra.CheckErr(err)
		err = daemon.Init()
		cobra.CheckErr(err)
		err = daemon.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "daemon run error:%v", err)
		}
	},
}

func init() {
	rpcbind = runCmd.Flags().String("rpcbind", "", "rpc server host")
	rpcport = runCmd.Flags().String("rpcport", "", "rpc server port")
	btcUrl = runCmd.Flags().String("btcUrl", "", "bitcoin json rpc endpoint")
	btcUser = runCmd.Flags().String("btcUser", "", "bitcoin json rpc username")
	btcPwd = runCmd.Flags().String("btcPwd", "", "bitcoin json rpc password")
	ethUrl = runCmd.Flags().String("ethUrl", "", "ethereum json rpc endpoint")
	ethPrivateKey = runCmd.Flags().String("ethPrivateKey", "", "ethereum private key")
	workers = runCmd.Flags().StringArray("workers", nil, "remote generate proof workers")
	rootCmd.AddCommand(runCmd)
}

func toConfig(data interface{}) (node.NodeConfig, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return node.NodeConfig{}, err
	}
	config := node.NodeConfig{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return node.NodeConfig{}, err
	}
	if *rpcbind != "" {
		config.Rpcbind = *rpcbind
	}
	if *rpcport != "" {
		config.RpcPort = *rpcport
	}
	if *btcUrl != "" {
		config.BtcUrl = *btcUrl
	}
	if *btcUser != "" {
		config.BtcUser = *btcUser
	}
	if *btcPwd != "" {
		config.BtcPwd = *btcPwd
	}
	if *ethUrl != "" {
		config.EthUrl = *ethUrl
	}
	if *ethPrivateKey != "" {
		config.EthPrivateKey = *ethPrivateKey
	}
	return config, nil
}
