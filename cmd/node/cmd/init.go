package cmd

import (
	"fmt"
	"github.com/lightec-xyz/daemon/cmd/node/cmd/miner"
	"github.com/lightec-xyz/daemon/cmd/node/cmd/proof"
	"github.com/lightec-xyz/daemon/logger"
)

func init() {
	err := logger.InitLogger(nil)
	if err != nil {
		fmt.Printf("init logger error: %v \n", err)
		return
	}
	rootCmd.AddCommand(miner.MinerCmd)
	rootCmd.AddCommand(proof.ProofCmd)
}
