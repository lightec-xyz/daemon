/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/spf13/cobra"
	"strconv"
)

var addWorkerCmd = &cobra.Command{
	Use:     "addWorker",
	Short:   "add a new worker to daemon",
	Example: `example: ./daemon --rpcbind 127.0.0.1 --rpcport 9780 addWorker ws://127.0.0.1:30001 1`,
	Run: func(cmd *cobra.Command, args []string) {
		if rpcbind == "" || rpcport == "" || len(args) != 2 {
			fmt.Printf("input data error")
			return
		}
		endpoint := args[0]
		client, err := rpc.NewNodeClient(fmt.Sprintf("http://%s:%s", rpcbind, rpcport))
		cobra.CheckErr(err)
		nums, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			fmt.Printf("args error")
			return
		}
		result, err := client.AddWorker(endpoint, int(nums))
		cobra.CheckErr(err)
		fmt.Printf("result: %s\n", result)
	},
}

func init() {
	rootCmd.AddCommand(addWorkerCmd)
}
