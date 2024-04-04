package node

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/lightec-xyz/daemon/common"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
)

type ZkProofType int

const (
	DepositTxType ZkProofType = iota + 1
	RedeemTxType
	TxInEth2
	VerifyTxType
	SyncComGenesisType
	SyncComUnitType
	SyncComRecursiveType
)

func (zkpr *ZkProofType) String() string {
	switch *zkpr {
	case DepositTxType:
		return "DepositTxType"
	case RedeemTxType:
		return "RedeemTxType"
	case VerifyTxType:
		return "VerifyTxType"
	case SyncComGenesisType:
		return "SyncComGenesisType"
	case SyncComUnitType:
		return "SyncComUnitType"
	case SyncComRecursiveType:
		return "SyncComRecursiveType"
	default:
		return ""
	}
}

type ZkProofRequest struct {
	reqType ZkProofType
	data    interface{}
	period  uint64
	TxHash  string
}

func (r *ZkProofRequest) String() string {
	return fmt.Sprintf("ZkProofRequest{reqType:%v,Period:%v,data:%v}", r.reqType, r.period, r.data)
}

type ZkProofResponse struct {
	ZkProofType ZkProofType
	Status      ProofStatus
	Proof       common.ZkProof
	Witness     []byte
	Period      uint64
	TxHash      string
}

func (zkResp *ZkProofResponse) String() string {
	return fmt.Sprintf("ZkProofType:%v Period:%v Proof:%v", zkResp.ZkProofType, zkResp.Period, zkResp.Proof)
}

type DepositProofParam struct {
	Version   string
	Body      interface{}
	TxHash    string
	BlockHash string
}

type TxInEth2Param struct {
	Version string
	TxHash  string
	TxData  *ethblock.TxInEth2ProofData
}

type RedeemProofParam struct {
	Version string
	TxHash  string
	TxData  *ethblock.TxInEth2ProofData
}

type VerifyProofParam struct {
	Version string
	Body    interface{}
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

type ProofStatus int

const (
	ProofDefault ProofStatus = iota
	ProofPending
	ProofSuccess
	ProofFailed
)

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
	TxHash    string      `json:"txId"`
	ProofType ZkProofType `json:"type"`
	Status    int         `json:"status"`
	Proof     string      `json:"Proof"`
}
