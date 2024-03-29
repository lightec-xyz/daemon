package cmd

import (
	"github.com/lightec-xyz/daemon/node"
	"github.com/spf13/cobra"
)

var recursiveCmd = &cobra.Command{
	Use:   "recursive",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		beaconUrl, err := cmd.Flags().GetString("beaconUrl")
		if err != nil {
			panic(err)
		}
		// todo
		//datadir = "/Users/red/lworkspace/lightec/daemon/node/test"
		//beaconUrl = "http://37.120.151.183:8970"
		cfg, err := node.NewLightLocalDaemonConfig(true, datadir, "testnet", rpcbind,
			rpcport, beaconUrl)
		if err != nil {
			panic(err)
		}
		daemon, err := node.NewRecursiveLightDaemon(cfg)
		if err != nil {
			panic(err)
		}
		err = daemon.Init()
		if err != nil {
			panic(err)
		}
		err = daemon.Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	recursiveCmd.Flags().String("beaconUrl", "", "rpc beaconUrl")
	rootCmd.AddCommand(recursiveCmd)
}
