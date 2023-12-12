package cmd

import (
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "generate zk proof node",
	Long:  `example: ./proof run`,
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
