package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// unitCmd represents the unit command
var unitCmd = &cobra.Command{
	Use:   "unit",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("unit called")
	},
}

func init() {
	rootCmd.AddCommand(unitCmd)
}
