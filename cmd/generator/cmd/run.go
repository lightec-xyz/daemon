package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/proof"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/spf13/cobra"
	"os"
)

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "generate zk proof node",
	Example: "./generator --config ./client_config.json run",
	Run: func(cmd *cobra.Command, args []string) {
		cfgBytes, err := os.ReadFile(cfgFile)
		if err != nil {
			fmt.Printf("read config error: %v %v \n", cfgFile, err)
			return
		}
		fmt.Printf("%v\n", string(cfgBytes))
		var config proof.Config
		err = json.Unmarshal(cfgBytes, &config)
		if err != nil {
			fmt.Printf("unmarshal config error: %v %v \n", cfgFile, err)
			return
		}

		client, err := rpc.NewNodeClient("https://testnet.zkbtc.money/api")
		if err != nil {
			fmt.Printf("new client error: %v %v \n", cfgFile, err)
			return
		}
		info, err := client.ProofInfo([]string{"0x4438c9e843b35e549173658a1409c4577ad78dae5b2cda70008cb31a541c4458"})
		if err != nil {
			fmt.Printf("proof info error: %v %v \n", cfgFile, err)
			return
		}
		fmt.Printf("proof info: %v \n", info)
		zkProofTypes := common.GetEnvZkProofTypes()
		fmt.Printf("zk proof types: %v \n", zkProofTypes)
		return

		node, err := proof.NewNode(config)
		if err != nil {
			fmt.Printf("new node error: %v %v \n", cfgFile, err)
			return
		}
		err = node.Init()
		if err != nil {
			fmt.Printf("init node error: %v %v \n", cfgFile, err)
			return
		}
		err = node.Start()
		if err != nil {
			fmt.Printf("start node error: %v %v \n", cfgFile, err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

}
