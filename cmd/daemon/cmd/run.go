package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/node"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var datadir *string
var rpcbind *string
var rpcport *string
var btcUrl *string
var btcUser *string
var btcPwd *string
var ethUrl *string
var ethPrivateKey *string
var network *string

const (
	datadirFlag       = "datadir"
	rpcbindFlag       = "rpcbind"
	rpcportFlag       = "rpcport"
	btcUrlFlag        = "btcUrl"
	btcUserFlag       = "btcUser"
	btcPwdFlag        = "btcPwd"
	ethUrlFlag        = "ethUrl"
	networkFlag       = "network"
	ethPrivateKeyFlag = "ethPrivateKey"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run daemon",
	Long:  `Start daemon program`,
	Run: func(cmd *cobra.Command, args []string) {
		datadir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd, ethUrl, ethPrivate := getConfig()
		config, err := node.NewNodeConfig(datadir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd, "", ethUrl, ethPrivate)
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
	rpcbind = runCmd.Flags().String(rpcbindFlag, "", "rpc server host")
	datadir = runCmd.Flags().String(datadirFlag, "", "rpc server host")
	rpcport = runCmd.Flags().String(rpcportFlag, "", "rpc server port")
	btcUrl = runCmd.Flags().String(btcUrlFlag, "", "bitcoin json rpc endpoint")
	btcUser = runCmd.Flags().String(btcUserFlag, "", "bitcoin json rpc username")
	btcPwd = runCmd.Flags().String(btcPwdFlag, "", "bitcoin json rpc password")
	ethUrl = runCmd.Flags().String(ethUrlFlag, "", "ethereum json rpc endpoint")
	ethUrl = runCmd.Flags().String(ethUrlFlag, "", "ethereum json rpc endpoint")
	network = runCmd.Flags().String(networkFlag, "", "lightec network")
	ethPrivateKey = runCmd.Flags().String(ethPrivateKeyFlag, "", "ethereum private key")
	rootCmd.AddCommand(runCmd)
}

func getConfig() (string, string, string, string, string, string, string, string, string) {
	tDatadir := viper.GetString(datadirFlag)
	tRpcbind := viper.GetString(rpcbindFlag)
	tRpcport := viper.GetString(rpcportFlag)
	tBtcurl := viper.GetString(btcUrlFlag)
	tBtcUser := viper.GetString(btcUserFlag)
	tBtcPwd := viper.GetString(btcPwdFlag)
	tEthUrl := viper.GetString(ethUrlFlag)
	tNetwork := viper.GetString(networkFlag)
	tEthPrivateKey := viper.GetString(ethPrivateKeyFlag)
	if *rpcbind != "" {
		tRpcbind = *rpcbind
	}
	if *rpcport != "" {
		tRpcport = *rpcport
	}
	if *btcUrl != "" {
		tBtcurl = *btcUrl
	}
	if *btcUser != "" {
		tBtcUser = *btcUser
	}
	if *btcPwd != "" {
		tBtcPwd = *btcPwd
	}
	if *ethUrl != "" {
		tEthUrl = *ethUrl
	}
	if *ethPrivateKey != "" {
		tEthPrivateKey = *ethPrivateKey
	}
	return tDatadir, tNetwork, tRpcport, tRpcbind, tBtcurl, tBtcUser, tBtcPwd, tEthUrl, tEthPrivateKey
}
