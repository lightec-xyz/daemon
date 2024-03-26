package cmd

import (
	"github.com/spf13/cobra"
)

var recursiveCmd = &cobra.Command{
	Use:   "recursive",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(recursiveCmd)
}
