package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/circuits"
	"github.com/lightec-xyz/daemon/logger"
	beacon_header_finality "github.com/lightec-xyz/provers/circuits/beacon-header-finality"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	"github.com/lightec-xyz/provers/common"
	reLightCommon "github.com/lightec-xyz/reLight/circuits/common"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var bhfCmd = &cobra.Command{
	Use:   "bhf",
	Short: "A brief description of your command",
	Long:  `./genproof bhf --paramFile <paramFile> --datadir <datadir>`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(paramFile)
		if err != nil {
			fmt.Printf("read config error: %v %v \n", paramFile, err)
			return
		}
		var param BhfParam
		err = json.Unmarshal(data, &param)
		if err != nil {
			fmt.Printf("unmarshal config error: %v %v \n", paramFile, err)
			return
		}
		proof, err := BhfProve(param, datadir)
		if err != nil {
			fmt.Printf("gen proof error: %v %v \n", paramFile, err)
			return
		}
		fmt.Printf("success generate bhf proof: %v \n", proof)

	},
}

func init() {
	rootCmd.AddCommand(bhfCmd)
}

type BhfParam struct {
	GenesisScRoot    string                         `json:"genesis_sc_root"`
	RecursiveProof   string                         `json:"recursive_proof"`
	RecursiveWitness string                         `json:"recursive_witness"`
	OuterProof       string                         `json:"outer_proof"`
	OuterWitness     string                         `json:"outer_witness"`
	FinalityUpdate   proverType.FinalityUpdate      `json:"finality_update"`
	ScUpdate         proverType.SyncCommitteeUpdate `json:"sc_update"`
}

func BhfProve(param BhfParam, datadir string) (*common.Proof, error) {
	scRecursiveVk, err := utils.ReadVk(filepath.Join(datadir, reLightCommon.RecursiveVkFile))
	if err != nil {
		logger.Error("read vk error:%v", err)
		return nil, err
	}
	recursiveProof, err := hex.DecodeString(param.RecursiveProof)
	if err != nil {
		logger.Error("decode proof error:%v", err)
		return nil, err
	}
	scRecursiveProof, err := circuits.ParseProof(recursiveProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	recursiveWitness, err := hex.DecodeString(param.RecursiveWitness)
	if err != nil {
		logger.Error("decode witness error:%v", err)
		return nil, err
	}
	scRecursiveWitness, err := circuits.ParseWitness(recursiveWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	scOuterVk, err := utils.ReadVk(filepath.Join(datadir, reLightCommon.OuterVkFile))
	if err != nil {
		logger.Error("read vk error:%v", err)
		return nil, err
	}
	outerProof, err := hex.DecodeString(param.OuterProof)
	if err != nil {
		logger.Error("decode proof error:%v", err)
		return nil, err
	}
	scOuterProof, err := circuits.ParseProof(outerProof)
	if err != nil {
		logger.Error("parse proof error:%v", err)
		return nil, err
	}
	outerWitness, err := hex.DecodeString(param.OuterWitness)
	if err != nil {
		logger.Error("decode witness error:%v", err)
		return nil, err
	}
	scOuterWitness, err := circuits.ParseWitness(outerWitness)
	if err != nil {
		logger.Error("parse witness error:%v", err)
		return nil, err
	}
	proof, err := beacon_header_finality.Prove(datadir, param.GenesisScRoot, scRecursiveVk, scRecursiveProof,
		scRecursiveWitness, scOuterVk, scOuterProof, scOuterWitness, &param.FinalityUpdate, &param.ScUpdate)
	if err != nil {
		logger.Error("beacon header finality update prove error:%v", err)
		return nil, err
	}
	return proof, nil
}
