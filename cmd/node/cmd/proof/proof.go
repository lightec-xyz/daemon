package proof

import (
	"fmt"
	"github.com/lightec-xyz/daemon/logger"

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
	err := logger.InitLogger(nil)
	if err != nil {
		fmt.Printf("init logger error: %v \n", err)
		return
	}
	ProofCmd.AddCommand(importCmd)
	ProofCmd.AddCommand(removeBtcProofCmd)
	ProofCmd.AddCommand(readVkCmd)
}
