package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var datadir string
var rpcbind string
var rpcport string
var network string

var rootCmd = &cobra.Command{
	Use:   "daemon",
	Short: "an cross-chain bridge node between Ethereum and Bitcoin",
	Long:  "a node for a cross-chain bridge between Ethereum and Bitcoin implemented in the lightning  protocol",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.daemon/daemon.json)")
	rootCmd.PersistentFlags().StringVar(&datadir, datadirFlag, "", "daemon storage directory")
	rootCmd.PersistentFlags().StringVar(&rpcbind, rpcbindFlag, "", "rpc server host")
	rootCmd.PersistentFlags().StringVar(&rpcport, rpcportFlag, "", "rpc server port")
	rootCmd.PersistentFlags().StringVar(&network, networkFlag, "testnet", "lightec network")
	getRootConfig()
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(fmt.Sprintf("%s/.daemon", home))
		viper.SetConfigType("json")
		viper.SetConfigName("daemon")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
func getRootConfig() {
	tDatadir := viper.GetString(datadirFlag)
	tRpcbind := viper.GetString(rpcbindFlag)
	tRpcport := viper.GetString(rpcportFlag)
	tNetwork := viper.GetString(networkFlag)
	if tDatadir != "" && datadir == "" {
		datadir = tDatadir
	}
	if tRpcbind != "" && rpcbind == "" {
		rpcbind = tRpcbind
	}
	if tRpcport != "" && rpcport == "" {
		rpcport = tRpcport
	}
	if tNetwork != "" && network == "" {
		network = tNetwork
	}
}
