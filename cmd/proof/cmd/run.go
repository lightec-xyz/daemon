package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/proof"
	"github.com/spf13/cobra"
)

var ip *string
var port *string
var maxNums *int

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "generate zk proof node",
	Long:  "example: ./proof  run --ip 0.0.0.0 --port 30001 --nums 1",
	Run: func(cmd *cobra.Command, args []string) {
		//config := proof.LocalDevConfig()
		config := newConfig()
		node, err := proof.NewNode(config)
		cobra.CheckErr(err)
		err = node.Start()
		if err != nil {
			fmt.Printf(" %v\n", err)
		}
	},
}

func newConfig() proof.Config {
	return proof.Config{
		RpcBind: *ip,
		RpcPort: *port,
		MaxNums: *maxNums,
	}
}

func init() {
	ip = runCmd.Flags().String("ip", "127.0.0.1", "rpc server host")
	port = runCmd.Flags().String("port", "30001", "rpc server port")
	maxNums = runCmd.Flags().Int("nums", 1, "The maximum number of proofs that can be generated at the same time")
	rootCmd.AddCommand(runCmd)

}
