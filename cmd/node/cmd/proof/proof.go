package proof

import (
	"fmt"

	"github.com/spf13/cobra"
)

// proofCmd represents the proof command
var ProofCmd = &cobra.Command{
	Use:   "proof",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proof called")
	},
}

func init() {
	ProofCmd.AddCommand(importCmd)
	ProofCmd.AddCommand(removeBtcProofCmd)
}
