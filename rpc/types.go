package rpc

import "github.com/prysmaticlabs/prysm/v5/api/server/structs"

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

type DepositRequest struct {
}

type DepositResponse struct {
	Body []byte
}

type RedeemRequest struct {
}

type RedeemResponse struct {
	Body []byte
}

type VerifyRequest struct {
}

type VerifyResponse struct {
	Body []byte
}

type SyncCommGenesisRequest struct {
	Version                    string `json:"version"`
	AttestedHeader             structs.BeaconBlockHeader
	CurrentSyncCommittee       structs.SyncCommittee
	CurrentSyncCommitteeBranch []string
}

type SyncCommGenesisResponse struct {
	Body      []byte
	Version   string                           `json:"version"`
	Period    uint64                           `json:"period"`
	ProofType SyncCommitteeProofType           `json:"proofType"`
	Status    SyncCommitteeProofGenerateStatus `json:"status"`
	Proof     string
}

type SyncCommUnitsRequest struct {
	Version                 string `json:"version"`
	AttestedHeader          structs.BeaconBlockHeader
	CurrentSyncCommittee    structs.SyncCommittee
	SyncAggregate           structs.SyncAggregate
	NextSyncCommittee       structs.SyncCommittee
	NextSyncCommitteeBranch []string
}

type SyncCommUnitsResponse struct {
	Body      []byte
	Version   string                           `json:"version"`
	Period    uint64                           `json:"period"`
	ProofType SyncCommitteeProofType           `json:"proofType"`
	Status    SyncCommitteeProofGenerateStatus `json:"status"`
	Proof     string
}

type SyncCommRecursiveRequest struct {
	Version          string `json:"version"`
	PreProofGOrPoofR string `json:"preProofGOrPoofR"`
	ProofU           string `json:"proofU"`
}

type SyncCommRecursiveResponse struct {
	Body      []byte
	Version   string                           `json:"version"`
	Period    uint64                           `json:"period"`
	ProofType SyncCommitteeProofType           `json:"proofType"`
	Status    SyncCommitteeProofGenerateStatus `json:"status"`
	Proof     string
}

type ProofInfo struct {
	reqType   int
	TxId      string `json:"txId"`
	ProofType int    `json:"type"`
	Proof     string `json:"proof"`
	Status    int    `json:"status"`
}

type SyncCommitteeProofGenerateStatus int

const (
	SyncCommitteeProofGenerateStatus_None       SyncCommitteeProofGenerateStatus = 0
	SyncCommitteeProofGenerateStatus_Generating SyncCommitteeProofGenerateStatus = 1
	SyncCommitteeProofGenerateStatus_Done       SyncCommitteeProofGenerateStatus = 2
)

type SyncCommitteeProofType int

const (
	SyncCommitteeProofType_None      SyncCommitteeProofType = 0
	SyncCommitteeProofType_Genesis   SyncCommitteeProofType = 1
	SyncCommitteeProofType_Unit      SyncCommitteeProofType = 2
	SyncCommitteeProofType_Recursive SyncCommitteeProofType = 3
)

type GenesisSyncCommitteeProofRequest struct {
	Version                    string `json:"version"`
	AttestedHeader             structs.BeaconBlockHeader
	CurrentSyncCommittee       structs.SyncCommittee
	CurrentSyncCommitteeBranch []string
}

type UnitSyncCommitteeProofRequest struct {
	Version                 string `json:"version"`
	AttestedHeader          structs.BeaconBlockHeader
	CurrentSyncCommittee    structs.SyncCommittee
	SyncAggregate           structs.SyncAggregate
	NextSyncCommittee       structs.SyncCommittee
	NextSyncCommitteeBranch []string
}

type RecursiveSyncCommitteeProofRequest struct {
	Version          string `json:"version"`
	PreProofGOrPoofR string `json:"preProofGOrPoofR"`
	ProofU           string `json:"proofU"`
}

type SyncCommitteeProofResponse struct {
	Version   string                           `json:"version"`
	Period    uint64                           `json:"period"`
	ProofType SyncCommitteeProofType           `json:"proofType"`
	Status    SyncCommitteeProofGenerateStatus `json:"status"`
	Proof     string                           `json:"proof"`
}

type SyncCommitteeProofInfo struct {
	Version   string                           `json:"version"`
	Period    uint64                           `json:"period"`
	ProofType SyncCommitteeProofType           `json:"proofType"`
	Status    SyncCommitteeProofGenerateStatus `json:"status"`
	Proof     string                           `json:"proof"`
}
