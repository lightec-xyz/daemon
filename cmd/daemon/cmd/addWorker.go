/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// addWorkerCmd represents the addWorker command
var addWorkerCmd = &cobra.Command{
	Use:   "addWorker",
	Short: "add a new worker to daemon",
	Long:  `example: daemon addWorker http://127.0.0.1:8485 1`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addWorker called", args)

	},
}

func init() {
	rootCmd.AddCommand(addWorkerCmd)
}
