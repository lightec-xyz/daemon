package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var genTxProofCmd = &cobra.Command{
	Use:     "genTxProof",
	Short:   "generate tx zk proof",
	Example: "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("genProof called")
	},
}

var genSyncCommitteeProofCmd = &cobra.Command{
	Use:     "genSyncCommitteeProof",
	Short:   "generate sync committee zk proof",
	Example: "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("genProof called")
	},
}

func init() {
	rootCmd.AddCommand(genTxProofCmd)
	rootCmd.AddCommand(genSyncCommitteeProofCmd)

}
