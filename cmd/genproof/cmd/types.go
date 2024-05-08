package cmd

import proverType "github.com/lightec-xyz/provers/circuits/types"

type BhfParam struct {
	GenesisScRoot    string                         `json:"genesis_sc_root"`
	RecursiveProof   string                         `json:"recursive_proof"`
	RecursiveWitness string                         `json:"recursive_witness"`
	OuterProof       string                         `json:"outer_proof"`
	OuterWitness     string                         `json:"outer_witness"`
	FinalityUpdate   proverType.FinalityUpdate      `json:"finality_update"`
	ScUpdate         proverType.SyncCommitteeUpdate `json:"sc_update"`
}
