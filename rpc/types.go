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

type EthereumTx struct {
	Hash string `json:"hash"`
}

type BitcoinTx struct {
	Hash string `json:"hash"`
}

type NodeInfo struct {
	Version string
	Desc    string
}
type Utxo struct {
	TxId  string `json:"txId"`
	Index uint32 `json:"index"`
}

type TxOut struct {
	Value    int64
	PkScript []byte
}

type ProofRequest struct {
	// redeem
	Inputs  []Utxo  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`
	BtcTxId string  `json:"btcTxId"`

	// deposit
	Utxos   []Utxo
	Amount  int64  `json:"amount"`
	EthAddr string `json:"ethAddr"`

	// other
	Height    int64  `json:"height"`
	BlockHash string `json:"blockHash"`
	TxId      string `json:"txId"`
	ProofType int    `json:"type"` // todo
	Proof     string `json:"proof"`
	Msg       string `json:"msg"`
}

type ProofResponse struct { // redeem
	// redeem
	Inputs  []Utxo  `json:"inputs"`
	Outputs []TxOut `json:"outputs"`
	BtcTxId string  `json:"btcTxId"`

	// deposit
	Utxos   []Utxo
	Amount  int64  `json:"amount"`
	EthAddr string `json:"ethAddr"`

	// other
	Height    int64  `json:"height"`
	BlockHash string `json:"blockHash"`
	TxId      string `json:"txId"`
	ProofType int    `json:"type"` // todo
	Proof     string `json:"proof"`
	Msg       string `json:"msg"`
}

type ProofInfo struct {
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
