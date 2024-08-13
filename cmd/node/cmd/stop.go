package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:     "stop",
	Short:   "stop node ",
	Example: "./node stop",
	Run: func(cmd *cobra.Command, args []string) {
		if rpcbind == "" || rpcport == "" {
			fmt.Println("rpcbind or rpcport is empty")
			return
		}
		client, err := rpc.NewNodeClient(fmt.Sprintf("http://%s:%s", rpcbind, rpcport), "")
		cobra.CheckErr(err)
		err = client.Stop()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
