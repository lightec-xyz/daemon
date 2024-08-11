package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/rpc"
	"github.com/spf13/cobra"
	"os"
)

var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfgBytes, err := os.ReadFile(cfgFile)
		if err != nil {
			fmt.Printf("read config error: %v %v \n", cfgFile, err)
			return
		}
		fmt.Printf("confg data: %v \n", string(cfgBytes))
		var runCfg node.RunConfig
		err = json.Unmarshal(cfgBytes, &runCfg)
		if err != nil {
			fmt.Printf("unmarshal config error: %v %v \n", cfgFile, err)
			return
		}
		secret, err := hex.DecodeString(runCfg.EthPrivateKey)
		if err != nil {
			fmt.Printf("decode private key error: %v %v \n", runCfg.EthPrivateKey, err)
			return
		}
		jwt, err := rpc.CreateJWT(secret, rpc.JwtPermission)
		if err != nil {
			fmt.Printf("create jwt error: %v \n", err)
			return
		}
		fmt.Printf("jwt: %v \n", jwt)
	},
}

func init() {
	rootCmd.AddCommand(jwtCmd)
}
