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
		initSlot, err := cmd.Flags().GetUint64("initSlot")
		if err != nil {
			panic(err)
		}
		cfg, err := node.NewLightLocalDaemonConfig(true, datadir, "testnet", rpcbind,
			rpcport, beaconUrl, initSlot)
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
	recursiveCmd.Flags().Uint64("initSlot", 1253376, "init slot height")
	rootCmd.AddCommand(recursiveCmd)
}
