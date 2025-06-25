package miner

import (
	"github.com/spf13/cobra"
)

var MinerCmd = &cobra.Command{
	Use:   "miner",
	Short: "A brief description of your command",
	Long:  ``,
}

func init() {
	MinerCmd.AddCommand(nonce)
}
