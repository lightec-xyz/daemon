package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var genProofCmd = &cobra.Command{
	Use:     "genProof",
	Short:   "generate zk proof",
	Example: "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("genProof called")
	},
}

func init() {
	rootCmd.AddCommand(genProofCmd)
}
