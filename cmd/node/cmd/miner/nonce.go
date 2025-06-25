package miner

import (
	"fmt"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/store"
	"github.com/spf13/cobra"
)

var nonce = &cobra.Command{
	Use:   "nonce",
	Short: "update miner nonce",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := logger.InitLogger(nil)
		if err != nil {
			fmt.Printf("init logger error: %v \n", err)
			return
		}
		cfgFile, err := cmd.Root().PersistentFlags().GetString("config")
		if err != nil {
			fmt.Printf("get config error: %v \n", err)
			return
		}
		fmt.Printf("config file: %v \n", cfgFile)
		var cfg node.RunConfig
		err = common.ReadObj(cfgFile, &cfg)
		if err != nil {
			fmt.Printf("read config error: %v %v \n", cfgFile, err)
			return
		}
		miner, err := cmd.Flags().GetString("addr")
		if err != nil {
			fmt.Printf("get miner error: %v \n", err)
			return
		}
		if !ethCommon.IsHexAddress(miner) {
			fmt.Printf("miner addr error:%v \n", miner)
			return
		}

		nonce, err := cmd.Flags().GetInt64("nonce")
		if err != nil {
			fmt.Printf("get nonce error: %v \n", err)
			return
		}
		if nonce < 0 {
			fmt.Printf("nonce is error: %v \n", nonce)
			return
		}

		path := fmt.Sprintf("%s/%s", cfg.Datadir, cfg.Network)
		storeDb, err := store.NewStore(path, 0, 0, common.DbNameSpace, false)
		if err != nil {
			fmt.Printf("new store error: %v \n", err)
			return
		}
		chainStore := node.NewChainStore(storeDb)
		err = chainStore.WriteNonce(common.EthereumChain.String(), miner, uint64(nonce))
		if err != nil {
			fmt.Printf("write nonce error: %v \n", err)
			return
		}
		minerNonce, ok, err := chainStore.ReadNonce(common.EthereumChain.String(), miner)
		if err != nil {
			fmt.Printf("read nonce error: %v \n", err)
			return
		}
		if !ok {
			fmt.Printf("read nonce error: %v \n", err)
			return
		}
		if uint64(nonce) != minerNonce {
			fmt.Printf("read nonce error: %v \n", err)
			return
		}
		fmt.Printf("write nonce success %v %v \n", miner, nonce)
	},
}

func init() {
	nonce.Flags().String("addr", "", "miner address")
	nonce.Flags().Int64("nonce", -1, "miner nonce")
}
