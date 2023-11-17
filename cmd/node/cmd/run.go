package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/node"
	"github.com/spf13/cobra"
	"os"
)

var btcUrl string
var ethUrl string
var beaconUrl string

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "run node",
	Example: "./node run",
	Run: func(cmd *cobra.Command, args []string) {
		cfgBytes, err := os.ReadFile(cfgFile)
		if err != nil {
			fmt.Printf("read config error: %v %v \n", cfgFile, err)
			return
		}
		fmt.Printf("confg data: %v \n", string(cfgBytes))
		var runCfg node.RunConfig
		err = json.Unmarshal(cfgBytes, &runCfg)
		if err != nil {
			fmt.Printf("unmarshal config error: %v %v \n", cfgFile, err)
			return
		}
		config, err := node.NewConfig(runCfg)
		if err != nil {
			fmt.Printf("new config error: %v \n", err)
			return
		}
		daemon, err := node.NewDaemon(config)
		if err != nil {
			fmt.Printf("new daemon error: %v \n", err)
			return
		}
		err = daemon.Init()
		if err != nil {
			fmt.Printf("node init error: %v \n", err)
			return
		}
		err = daemon.Run()
		if err != nil {
			fmt.Printf("node run error: %v \n", err)
		}
	},
}

func init() {
	runCmd.Flags().StringVar(&btcUrl, btcUrlFlag, "", "bitcoin json rpc endpoint")
	runCmd.Flags().StringVar(&ethUrl, ethUrlFlag, "", "ethereum json rpc endpoint")
	runCmd.Flags().StringVar(&beaconUrl, beaconUrlFlag, "", "eth2 json rpc endpoint")
	rootCmd.AddCommand(runCmd)
}
