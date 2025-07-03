package proof

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/node"
	"github.com/lightec-xyz/daemon/store"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "export proof to db",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Printf("get path error: %v \n", err)
			return
		}
		if name == "" {
			fmt.Printf("name is empty \n")
			return
		}
		datadir, err := cmd.Flags().GetString("datadir")
		if err != nil {
			fmt.Printf("get path error: %v \n", err)
			return
		}
		if datadir == "" {
			fmt.Printf("datadir is empty \n")
			return
		}

		proofPath, err := cmd.Flags().GetString("proof")
		if err != nil {
			fmt.Printf("get path error: %v \n", err)
			return
		}
		witPath, err := cmd.Flags().GetString("witness")
		if err != nil {
			fmt.Printf("get path error: %v \n", err)
			return
		}
		proofBytes, err := os.ReadFile(proofPath)
		if err != nil {
			fmt.Printf("read proof error: %v \n", err)
			return
		}
		witnessByts, err := os.ReadFile(witPath)
		if err != nil {
			fmt.Printf("read witness error: %v \n", err)
			return
		}
		fileStorage, err := node.NewFileStorage(datadir, 0, 0)
		if err != nil {
			fmt.Printf("new file storage error: %v \n", err)
			return
		}
		proofType, err := getProofType(name)
		if err != nil {
			fmt.Printf("get proof type error: %v \n", err)
			return
		}
		prefix, fIndex, sIndex, err := node.FileKeyToIndex(store.FileKey(strings.ToLower(name)))
		if err != nil {
			fmt.Printf("get file key error: %v \n", err)
			return
		}
		fmt.Printf("prefix: %v fIndex: %v sIndex: %v \n", prefix, fIndex, sIndex)
		err = fileStorage.StoreProof(node.NewStoreKey(proofType, "", prefix, fIndex, sIndex), proofBytes, witnessByts)
		if err != nil {
			fmt.Printf("store proof error: %v \n", err)
			return
		}
		fmt.Printf("import proof success \n")

	},
}

func init() {
	importCmd.Flags().String("proof", "", "gnark proof path")
	importCmd.Flags().String("witness", "", "gnark witness path")
	importCmd.Flags().String("name", "", "circuit name")
	importCmd.Flags().String("datadir", "", "datadir path")
}

func getProofType(name string) (common.ProofType, error) {
	ids := strings.Split(name, "_")
	if len(ids) < 2 {
		return 0, fmt.Errorf("invalid circuit name")
	}
	proofType, err := common.ToZkProofType(ids[0])
	if err != nil {
		return 0, err
	}
	return proofType, nil

}
