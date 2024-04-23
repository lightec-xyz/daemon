package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var genesisCmd = &cobra.Command{
	Use:   "genesis",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("genesis called")
	},
}

func init() {
	rootCmd.AddCommand(genesisCmd)
}
