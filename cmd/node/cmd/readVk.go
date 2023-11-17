package cmd

import (
	"bytes"
	"fmt"
	"github.com/lightec-xyz/common/operations"
	"github.com/spf13/cobra"
)

// readVkCmd represents the readVk command
var readVkCmd = &cobra.Command{
	Use:   "readVk",
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
		vkey, err := operations.ReadVk(path)
		if err != nil {
			fmt.Printf("read vk error: %v \n", err)
			return
		}
		var data []byte
		buffer := bytes.NewBuffer(data)
		_, err = vkey.WriteTo(buffer)
		if err != nil {
			fmt.Printf("write to buffer error: %v \n", err)
			return
		}
		fmt.Printf("vk info: %v \n", path)
		fmt.Printf("vk data: %x \n", buffer.Bytes())

	},
}

func init() {
	readVkCmd.Flags().String("file", "", "vk file")
	rootCmd.AddCommand(readVkCmd)

}
