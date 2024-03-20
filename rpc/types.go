package rpc

import (
	"github.com/lightec-xyz/daemon/common"
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

type DepositRequest struct {
	Version string
}

type DepositResponse struct {
	Proof common.ZkProof
}

type RedeemRequest struct {
	Version string
}

type RedeemResponse struct {
	Proof common.ZkProof
}

type VerifyRequest struct {
	Version string
}

type VerifyResponse struct {
	Proof common.ZkProof
}

type SyncCommGenesisRequest struct {
	Version                    string `json:"version"`
	AttestedHeader             structs.BeaconBlockHeader
	CurrentSyncCommittee       structs.SyncCommittee
	CurrentSyncCommitteeBranch []string
}

type SyncCommGenesisResponse struct {
	Version   string                           `json:"version"`
	Period    uint64                           `json:"period"`
	ProofType SyncCommitteeProofType           `json:"proofType"`
	Status    SyncCommitteeProofGenerateStatus `json:"status"`
	Proof     common.ZkProof
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
	Version   string                           `json:"version"`
	Period    uint64                           `json:"period"`
	ProofType SyncCommitteeProofType           `json:"proofType"`
	Status    SyncCommitteeProofGenerateStatus `json:"status"`
	Proof     common.ZkProof
}

type SyncCommRecursiveRequest struct {
	Version          string `json:"version"`
	PreProofGOrPoofR string `json:"preProofGOrPoofR"`
	ProofU           string `json:"proofU"`
}

type SyncCommRecursiveResponse struct {
	Version   string                           `json:"version"`
	Period    uint64                           `json:"period"`
	ProofType SyncCommitteeProofType           `json:"proofType"`
	Status    SyncCommitteeProofGenerateStatus `json:"status"`
	Proof     common.ZkProof
}

type ProofInfo struct {
	reqType   int
	TxId      string         `json:"txId"`
	ProofType int            `json:"type"`
	Proof     common.ZkProof `json:"proof"`
	Status    int            `json:"status"`
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
