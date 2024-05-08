package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/proof"
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
		node, err := proof.NewNode(config)
		if err != nil {
			fmt.Printf("new node error: %v %v \n", cfgFile, err)
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
