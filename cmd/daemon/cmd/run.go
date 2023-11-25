package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/node"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run daemon",
	Long:  `Start daemon program`,
	Run: func(cmd *cobra.Command, args []string) {
		//todo
		config, err := toConfig(viper.AllSettings())
		if err != nil {
			fmt.Fprintln(os.Stderr, "config file error:%v", err)
			return
		}
		daemon, err := node.NewDaemon(config)
		err = daemon.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "daemon run error:%v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func toConfig(data interface{}) (node.Config, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return node.Config{}, err
	}
	config := node.Config{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return node.Config{}, err
	}
	return config, nil
}
