package node

import (
	"bytes"
	"fmt"
	"github.com/consensys/gnark/frontend"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/provers/circuits/fabric/receipt-proof"
	"github.com/lightec-xyz/provers/circuits/fabric/tx-proof"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	"strconv"

	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
)

type DepositProofParam struct {
	Version   string
	Body      interface{}
	TxHash    string
	BlockHash string
}

type VerifyProofParam struct {
	Version   string
	TxHash    string
	BlockHash string
}

type TxInEth2Param struct {
	Version string
	TxHash  string
	TxData  *ethblock.TxInEth2ProofData
}

type BeaconHeaderParam struct {
	Index   uint64
	Headers proverType.BeaconHeaderChain
}

type RedeemProofParam struct {
	TxHash  string
	Version string
	TxProof, TxWitness, BhProof, BhWitness, BhfProof, BhfWitness, BeginId, EndId, GenesisScRoot,
	CurrentSCSSZRoot []byte
	TxVar      *[tx.MaxTxUint128Len]frontend.Variable
	ReceiptVar *[receipt.MaxReceiptUint128Len]frontend.Variable
}

type GenesisProofParam struct {
	Version       string
	FirstProof    []byte
	SecondProof   []byte
	FirstWitness  []byte
	SecondWitness []byte
	GenesisId     []byte
	FirstId       []byte
	SecondId      []byte
	RecursiveFp   []byte
}

type UnitProofParam struct {
	Version                 string                     `json:"version"`
	AttestedHeader          *structs.BeaconBlockHeader `json:"attested_header"`
	CurrentSyncCommittee    *structs.SyncCommittee     `json:"current_sync_committee,omitempty"`     //current_sync_committee
	SyncAggregate           *structs.SyncAggregate     `json:"sync_aggregate"`                       //sync_aggregate for attested_header, signed by current_sync_committee
	FinalizedHeader         *structs.BeaconBlockHeader `json:"finalized_header,omitempty"`           //finalized_header in attested_header.state_root
	FinalityBranch          []string                   `json:"finality_branch,omitempty"`            // finality_branch in attested_header.state_root
	NextSyncCommittee       *structs.SyncCommittee     `json:"next_sync_committee,omitempty"`        //next_sync_committee in finalized_header.state_root
	NextSyncCommitteeBranch []string                   `json:"next_sync_committee_branch,omitempty"` //next_sync_committee branch in finalized_header.state_root
	SignatureSlot           string                     `json:"signature_slot"`
}

type RecursiveProofParam struct {
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

type FinalityBeaconHeaderParam struct {
	GenesisSCSSZRoot string
	RecursiveProof, RecursiveWitness, OuterProof,
	OuterWitness []byte
	FinalityUpdate *proverType.FinalityUpdate
	ScUpdate       *proverType.SyncCommitteeUpdate
}

type FetchType int

const (
	GenesisUpdateType FetchType = iota + 1
	PeriodUpdateType
)

func (ft FetchType) String() string {
	switch ft {
	case GenesisUpdateType:
		return "GenesisUpdateType"
	case PeriodUpdateType:
		return "PeriodUpdateType"
	default:
		return "unknown"
	}
}

type DownloadStatus int

type FetchRequest struct {
	UpdateType FetchType
	Status     DownloadStatus
	period     uint64
}

type FetchDataResponse struct {
	period     uint64
	UpdateType FetchType
}

type Utxo struct {
	TxId  string `json:"txId"`
	Index uint32 `json:"index"`
}

type TxOut struct {
	Value    int64
	PkScript []byte
}

func formatUtxo(utxos []Utxo) string {
	var buf bytes.Buffer
	for _, vin := range utxos {
		buf.WriteString(vin.TxId)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(int(vin.Index)))
		buf.WriteString(",")
	}
	return buf.String()
}
func formatOut(outputs []TxOut) string {
	var buf bytes.Buffer
	for _, out := range outputs {
		buf.WriteString(fmt.Sprintf("%x", out.PkScript))
		buf.WriteString(":")
		buf.WriteString(fmt.Sprintf("%v", out.Value))
		buf.WriteString(",")
	}
	return buf.String()
}

type Transaction struct {
	TxHash    string
	Height    int64
	BlockHash string
	TxType    TxType
	ChainType ChainType

	// bitcoin
	EthAddr string
	Amount  int64
	Utxo    []Utxo

	// ethereum
	BtcTxId string
}

type Proof struct {
	TxHash    string             `json:"txHash"`
	ProofType common.ZkProofType `json:"type"`
	Status    int                `json:"status"`
	Proof     string             `json:"Proof"`
}

// todo
type UnGenPreProof struct {
	TxId      string
	ChainType ChainType
}
