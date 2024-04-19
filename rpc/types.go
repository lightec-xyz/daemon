package rpc

import (
	"github.com/consensys/gnark/frontend"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/provers/circuits/fabric/receipt-proof"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	"github.com/lightec-xyz/provers/circuits/fabric/tx-proof"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
)

type ProofTaskRequest struct {
	ProofType []common.ZkProofType
}

type ProfTaskResponse struct {
	ProofType common.ZkProofType
}

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

type TxInEth2ProveRequest struct {
	Version string
	TxHash  string
	TxData  *ethblock.TxInEth2ProofData
}

type TxInEth2ProveResponse struct {
	Proof   []byte
	Witness []byte
}

type BlockHeaderRequest struct {
	Index     uint64
	BeginSlot uint64
	BeginRoot string
	EndSlot   uint64
	EndRoot   string
	Headers   []*structs.BeaconBlockHeader
}

type BlockHeaderResponse struct {
	Proof   []byte
	Witness []byte
}
type BlockHeaderFinalityRequest struct {
	GenesisSCSSZRoot string
	RecursiveProof, RecursiveWitness, OuterProof,
	OuterWitness []byte
	FinalityUpdate *proverType.FinalityUpdate
	ScUpdate       *proverType.SyncCommitteeUpdate
}

type BlockHeaderFinalityResponse struct {
	Proof   []byte
	Witness []byte
}

type DepositRequest struct {
	Version   string
	TxHash    string
	BlockHash string
}

type DepositResponse struct {
	TxHash  string
	Proof   common.ZkProof
	Witness []byte
}

type RedeemRequest struct {
	Version string
	TxProof, TxWitness, BhProof, BhWitness, BhfProof, BhfWitness, BeginId, EndId, GenesisScRoot,
	CurrentSCSSZRoot []byte
	TxVar      *[tx.MaxTxUint128Len]frontend.Variable
	ReceiptVar *[receipt.MaxReceiptUint128Len]frontend.Variable
}

type RedeemResponse struct {
	Proof   []byte
	Witness []byte
}

type VerifyRequest struct {
	Version   string
	TxHash    string
	BlockHash string
}

type VerifyResponse struct {
	TxHash string
	Proof  []byte
	Wit    []byte
}

type SyncCommGenesisRequest struct {
	Period        uint64 `json:"period"`
	Version       string `json:"version"`
	FirstProof    []byte `json:"firstProof"`
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
	Version      string             `json:"version"`
	Period       uint64             `json:"period"`
	ProofType    common.ZkProofType `json:"proofType"`
	Proof        common.ZkProof     `json:"proof"`
	Witness      []byte             `json:"witness"`
	OuterProof   []byte             `json:"outerProof"`
	OuterWitness []byte             `json:"outerWitness"`
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

type CheckReqStatus struct {
	Status int
}

type IDepositRequest struct {
	DepositRequest
	CheckReqStatus
}
