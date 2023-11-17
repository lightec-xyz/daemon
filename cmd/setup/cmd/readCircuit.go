package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var readCircuitCmd = &cobra.Command{
	Use:   "readCircuit",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Printf("get path error: %v \n", err)
			return
		}
		if path == "" {
			fmt.Printf("path is empty \n")
			return
		}
		fmt.Printf("path: %s \n", path)
		fvk, err := os.Open(path)
		if err != nil {
			fmt.Printf("open file error: %v \n", err)
			return
		}
		defer fvk.Close()
		buffer := bytes.NewBuffer(nil)
		_, err = buffer.ReadFrom(fvk)
		if err != nil {
			fmt.Printf("read file error: %v \n", err)
			return
		}
		fmt.Printf("%x\n", buffer.Bytes())
		fmt.Printf("read vk success")

	},
}

func init() {
	readCircuitCmd.Flags().String("file", "", "vk file")
	readCircuitCmd.Flags().String("type", "vk", "")
	rootCmd.AddCommand(readCircuitCmd)
}
