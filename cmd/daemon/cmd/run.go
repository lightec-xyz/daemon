package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/node"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var btcUrl string
var btcUser string
var btcPwd string
var ethUrl string
var ethPrivateKey string
var enableLocalWorker bool

const (
	datadirFlag           = "datadir"
	rpcbindFlag           = "rpcbind"
	rpcportFlag           = "rpcport"
	btcUrlFlag            = "btcUrl"
	btcUserFlag           = "btcUser"
	btcPwdFlag            = "btcPwd"
	ethUrlFlag            = "ethUrl"
	networkFlag           = "network"
	ethPrivateKeyFlag     = "ethPrivateKey"
	enableLocalWorkerFlag = "enableLocalWorker"
)

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "run daemon",
	Example: "./daemon run",
	Run: func(cmd *cobra.Command, args []string) {
		enableLocalWorker, btcUrl, btcUser, btcPwd, ethUrl, ethPrivate := getRunConfig()
		//fmt.Printf("datadir:%s, network:%s, rpcbind:%s, rpcport:%s, btcUrl:%s, btcUser:%s, btcPwd:%s, ethUrl:%s, ethPrivateKey:%s \n", datadir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd, ethUrl, ethPrivate)
		config, err := node.NewNodeConfig(enableLocalWorker, datadir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd, ethUrl, ethPrivate)
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

	runCmd.Flags().StringVar(&btcUrl, btcUrlFlag, "", "bitcoin json rpc endpoint")
	runCmd.Flags().StringVar(&btcUser, btcUserFlag, "", "bitcoin json rpc username")
	runCmd.Flags().StringVar(&btcPwd, btcPwdFlag, "", "bitcoin json rpc password")
	runCmd.Flags().StringVar(&ethUrl, ethUrlFlag, "", "ethereum json rpc endpoint")
	runCmd.Flags().StringVar(&ethPrivateKey, ethPrivateKeyFlag, "", "ethereum private key")
	rootCmd.AddCommand(runCmd)
}

func getRunConfig() (tEnableLocalWorker bool, tBtcurl, tBtcUser, tBtcPwd, tEthUrl, tEthPrivateKey string) {
	tBtcurl = viper.GetString(btcUrlFlag)
	tBtcUser = viper.GetString(btcUserFlag)
	tBtcPwd = viper.GetString(btcPwdFlag)
	tEthUrl = viper.GetString(ethUrlFlag)
	tEthPrivateKey = viper.GetString(ethPrivateKeyFlag)
	tEnableLocalWorker = viper.GetBool(enableLocalWorkerFlag)
	if btcUrl != "" {
		tBtcurl = btcUrl
	}
	if btcUser != "" {
		tBtcUser = btcUser
	}
	if btcPwd != "" {
		tBtcPwd = btcPwd
	}
	if ethUrl != "" {
		tEthUrl = ethUrl
	}
	if ethPrivateKey != "" {
		tEthPrivateKey = ethPrivateKey
	}
	return tEnableLocalWorker, tBtcurl, tBtcUser, tBtcPwd, tEthUrl, tEthPrivateKey
}
