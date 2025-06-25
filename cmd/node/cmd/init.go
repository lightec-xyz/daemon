package cmd

import "github.com/lightec-xyz/daemon/cmd/node/cmd/miner"

func init() {
	rootCmd.AddCommand(miner.MinerCmd)
}
