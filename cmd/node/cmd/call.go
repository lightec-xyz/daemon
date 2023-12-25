package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// callCmd represents the call command
var callCmd = &cobra.Command{
	Use:   "call",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("call called")
	},
}

func init() {
	rootCmd.AddCommand(callCmd)

}
