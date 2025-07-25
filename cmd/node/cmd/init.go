package cmd

import (
	"github.com/lightec-xyz/daemon/cmd/node/cmd/miner"
	"github.com/lightec-xyz/daemon/cmd/node/cmd/proof"
)

func init() {
	rootCmd.AddCommand(miner.MinerCmd)
	rootCmd.AddCommand(proof.ProofCmd)
}
