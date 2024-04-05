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
var beaconUrl string
var ethPrivateKey string

//var enableLocalWorker bool
//var autoSubmit bool

const (
	datadirFlag           = "datadir"
	rpcbindFlag           = "rpcbind"
	rpcportFlag           = "rpcport"
	btcUrlFlag            = "btcUrl"
	btcUserFlag           = "btcUser"
	btcPwdFlag            = "btcPwd"
	ethUrlFlag            = "ethUrl"
	beaconUrlFlag         = "beaconUrl"
	networkFlag           = "network"
	ethPrivateKeyFlag     = "ethPrivateKey"
	enableLocalWorkerFlag = "enableLocalWorker"
	autoSubmitFlag        = "autoSubmit"
)

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "run node",
	Example: "./node run",
	Run: func(cmd *cobra.Command, args []string) {
		_, autoSubmit, btcUrl, btcUser, btcPwd, ethUrl, ethPrivate := getRunConfig()
		//fmt.Printf("datadir:%s, network:%s, rpcbind:%s, rpcport:%s, btcUrl:%s, btcUser:%s, btcPwd:%s, ethUrl:%s, ethPrivateKey:%s \n", datadir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd, ethUrl, ethPrivate)
		config, err := node.NewNodeConfig(true, autoSubmit, datadir, network, rpcbind, rpcport, btcUrl, btcUser, btcPwd,
			beaconUrl, ethUrl, ethPrivate)
		//config := node.TestnetDaemonConfig()
		cobra.CheckErr(err)
		daemon, err := node.NewDaemon(config)
		defer daemon.Close()
		cobra.CheckErr(err)
		err = daemon.Init()
		cobra.CheckErr(err)
		err = daemon.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "node run error:%v", err)
		}
	},
}

func init() {

	runCmd.Flags().StringVar(&btcUrl, btcUrlFlag, "", "bitcoin json rpc endpoint")
	runCmd.Flags().StringVar(&btcUser, btcUserFlag, "", "bitcoin json rpc username")
	runCmd.Flags().StringVar(&btcPwd, btcPwdFlag, "", "bitcoin json rpc password")
	runCmd.Flags().StringVar(&ethUrl, ethUrlFlag, "", "ethereum json rpc endpoint")
	runCmd.Flags().StringVar(&ethUrl, ethUrlFlag, "", "ethereum json rpc endpoint")
	runCmd.Flags().StringVar(&beaconUrl, beaconUrlFlag, "", "eth2 json rpc endpoint")
	//runCmd.Flags().StringVar(&ethPrivateKey, ethPrivateKeyFlag, "", "ethereum private key")
	//runCmd.Flags().BoolVar(&autoSubmit, autoSubmitFlag, false, "autoSubmit eth tx")
	rootCmd.AddCommand(runCmd)
}

func getRunConfig() (tEnableLocalWorker, tAutoSubmit bool, tBtcurl, tBtcUser, tBtcPwd, tEthUrl, tEthPrivateKey string) {
	tBtcurl = viper.GetString(btcUrlFlag)
	tBtcUser = viper.GetString(btcUserFlag)
	tBtcPwd = viper.GetString(btcPwdFlag)
	tEthUrl = viper.GetString(ethUrlFlag)
	tEthPrivateKey = viper.GetString(ethPrivateKeyFlag)
	tEnableLocalWorker = viper.GetBool(enableLocalWorkerFlag)
	tAutoSubmit = viper.GetBool(autoSubmitFlag)
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
	return tEnableLocalWorker, tAutoSubmit, tBtcurl, tBtcUser, tBtcPwd, tEthUrl, tEthPrivateKey
}
