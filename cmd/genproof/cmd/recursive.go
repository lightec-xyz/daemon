package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var recursiveCmd = &cobra.Command{
	Use:   "recursive",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("recursive called")
	},
}

func init() {
	rootCmd.AddCommand(recursiveCmd)
}
