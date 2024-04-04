package rpc

import (
	"github.com/lightec-xyz/daemon/common"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
)

type Transaction struct {
	TxHash   string
	DestHash string
	Height   int64

	BtcTxId string

	Amount  int64
	EthAddr string
	Utxo    []Utxo

	Inputs  []Utxo
	Outputs []TxOut

	TxType    int
	ChainType int
}
type Utxo struct {
	TxId  string `json:"txId"`
	Index uint32 `json:"index"`
}

type TxOut struct {
	Value    int64
	PkScript []byte
}

type NodeInfo struct {
	Version string
	Desc    string
}

//------

type TxInEth2ProveReq struct {
	Version string
	TxHash  string
	TxData  *ethblock.TxInEth2ProofData
}

type TxInEth2ProveResp struct {
	ProofStr string
	Proof    []byte
	Witness  []byte
}

type TxBlockIsParentOfCheckPointProveReq struct {
}

type TxBlockIsParentOfCheckPointResp struct {
}
type CheckPointFinalityProveReq struct {
}

type CheckPointFinalityProveResp struct {
}

type DepositRequest struct {
	Version   string
	TxHash    string
	BlockHash string
}

type DepositResponse struct {
	TxHash   string
	Proof    common.ZkProof
	ProofStr string
	Witness  []byte
}

type RedeemRequest struct {
	Version string
	TxHash  string
	TxData  *ethblock.TxInEth2ProofData
}

type RedeemResponse struct {
	Proof   common.ZkProof
	Witness []byte
}

type VerifyRequest struct {
	Version string
}

type VerifyResponse struct {
	Proof common.ZkProof
}

type SyncCommGenesisRequest struct {
	Period        uint64 `json:"period"`
	Version       string `json:"version"`
	FirstProof    []byte
	FirstWitness  []byte
	SecondProof   []byte
	SecondWitness []byte
	GenesisID     []byte
	FirstID       []byte
	SecondID      []byte
	RecursiveFp   []byte
}

type SyncCommGenesisResponse struct {
	Version   string             `json:"version"`
	Period    uint64             `json:"period"`
	ProofType common.ZkProofType `json:"proofType"`
	Proof     common.ZkProof
	Witness   []byte
}

type SyncCommUnitsRequest struct {
	Version                 string                     `json:"version"`
	Period                  uint64                     `json:"period"`
	AttestedHeader          *structs.BeaconBlockHeader `json:"attested_header"`
	CurrentSyncCommittee    *structs.SyncCommittee     `json:"current_sync_committee"`
	SyncAggregate           *structs.SyncAggregate     `json:"sync_aggregate"`
	NextSyncCommittee       *structs.SyncCommittee     `json:"next_sync_committee"`
	NextSyncCommitteeBranch []string                   `json:"next_sync_committee_branch"`
	FinalizedHeader         *structs.BeaconBlockHeader `json:"finalized_header,omitempty"`
	FinalityBranch          []string                   `json:"finality_branch,omitempty"`
	SignatureSlot           string                     `json:"signature_slot"`
}

type SyncCommUnitsResponse struct {
	Version   string             `json:"version"`
	Period    uint64             `json:"period"`
	ProofType common.ZkProofType `json:"proofType"`
	Proof     common.ZkProof     `json:"proof"`
	Witness   []byte             `json:"witness"`
}

type SyncCommRecursiveRequest struct {
	Period        uint64
	Version       string
	Choice        string `json:"choice"`
	FirstProof    []byte
	FirstWitness  []byte
	SecondProof   []byte
	SecondWitness []byte
	BeginId       []byte
	RelayId       []byte
	EndId         []byte
	RecursiveFp   []byte
}

type SyncCommRecursiveResponse struct {
	Version   string             `json:"version"`
	Period    uint64             `json:"period"`
	ProofType common.ZkProofType `json:"proofType"`
	Proof     common.ZkProof
	Witness   []byte
}

type ProofInfo struct {
	reqType   int
	TxId      string `json:"txId"`
	ProofType int    `json:"type"`
	Proof     string `json:"proof"`
	Status    int    `json:"status"`
}
