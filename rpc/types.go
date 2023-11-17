package rpc

import (
	"encoding/json"
	"github.com/consensys/gnark/std/math/uints"
	"time"

	btcproverType "github.com/lightec-xyz/btc_provers/circuits/types"
	blockchainUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	btcbase "github.com/lightec-xyz/btc_provers/utils/blockchain"
	btcmiddle "github.com/lightec-xyz/btc_provers/utils/blockchain"
	btcupper "github.com/lightec-xyz/btc_provers/utils/blockchain"
	recursiveUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	"github.com/lightec-xyz/btc_provers/utils/blockdepth"
	blockDepthUtil "github.com/lightec-xyz/btc_provers/utils/blockdepth"
	grUtil "github.com/lightec-xyz/btc_provers/utils/txinchain"
	"github.com/lightec-xyz/daemon/common"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	txineth2Utils "github.com/lightec-xyz/provers/utils/tx-in-eth2"
)

type WrapSyncCommitteeUpdate struct {
	*proverType.SyncCommitteeUpdate
	CurrentSyncCommitteeBranch []string
}

type MinerInfo struct {
	Address   string `json:"address"`
	Power     uint64 `json:"power"`
	Timestamp uint64 `json:"timestamp"`
}

type ProofTaskInfo struct {
	Id             string    `json:"id"`
	QueueTime      time.Time `json:"queueTime"`
	GeneratingTime time.Time `json:"generatingTime"`
	EndTime        time.Time `json:"endTime"`
}

type Transaction struct {
	Height    uint64        `json:"height"`
	TxIndex   uint          `json:"txIndex"`
	Hash      string        `json:"hash"`
	ChainType string        `json:"chainType"`
	TxType    string        `json:"txType"`
	Amount    int64         `json:"amount"`
	DestChain DestChainInfo `json:"destChain"`
	Proof     ProofInfo     `json:"proof"`
	Tasks     interface{}   `json:"tasks"`
}

type BtcDuperRecursiveRequest struct {
	BlockChainData  *recursiveUtil.BlockChainProofData
	HybridChainData *blockchainUtil.HybridProofData
	First, Second   Proof
	Start, End      uint64

	FirstType, SecondType string
	FirstStep, SecondStep uint64
}

type DestChainInfo struct {
	Hash string `json:"hash"`
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
	Network string
	Desc    string
}

type BtcDepthRecursiveRequest struct {
	Data                *blockdepth.RecursiveBulksProofData
	First               Proof
	Genesis, Start, End uint64
	PreStep             uint64
	IsRecursive         bool
}

type BtcChainRequest struct {
	Data                *blockchainUtil.HybridProofData
	First               Proof
	Genesis, Start, End uint64
}

type BtcDepositRequest struct {
	Data                                    *grUtil.TxInChainProofData
	BlockChain, TxDepth, CpDepth, SigVerify Proof
	TxRecursive, CpRecursive                bool
	ProverAddr                              string
	ChainType                               string
	ChainStep                               uint64
	TxDepthStep                             uint64
	CpDepthStep                             uint64

	CpFlag            uint8
	SmoothedTimestamp uint32
	SigVerifyData     *blockDepthUtil.SigVerifProofData
}

type BtcTimestampRequest struct {
	CpTime     *blockDepthUtil.CptimestampProofData
	SmoothData *blockDepthUtil.SmoothedTimestampProofData
}

func (bt *BtcTimestampRequest) Check() error {
	if bt.CpTime != nil {
		var cpHeaders []btcproverType.BlockHeader
		for _, item := range bt.CpTime.BlockHeaders {
			cpHeaders = append(cpHeaders, btcproverType.BlockHeader(toU8Array(item[:])))
		}
		copy(bt.CpTime.BlockHeaders[:], cpHeaders)
	}
	if bt.SmoothData != nil {
		var smoothHeaders []btcproverType.BlockHeader
		for _, item := range bt.SmoothData.BlockHeaders {
			smoothHeaders = append(smoothHeaders, btcproverType.BlockHeader(toU8Array(item[:])))
		}
		copy(bt.SmoothData.BlockHeaders[:], smoothHeaders)
	}
	return nil
}

func toU8Array(values []uints.U8) []uints.U8 {
	var newValue []uints.U8
	for _, item := range values {
		newValue = append(newValue, jsonNumberToU8(item))
	}
	return newValue
}

func jsonNumberToU8(value uints.U8) uints.U8 {
	if v, ok := value.Val.(json.Number); ok {
		i, _ := v.Int64()
		return uints.NewU8(uint8(i))
	}
	return uints.NewU8(0)
}

type BtcChangeRequest struct {
	Data                                            *grUtil.TxInChainProofData
	BlockChain, TxDepth, CpDepth, Redeem, SigVerify Proof
	TxRecursive, CpRecursive                        bool
	ProverAddr                                      string
	MinerReward                                     string
	ChainType                                       string
	MinTxDepth, MinCpDepth                          uint32
	ChainStep                                       uint64
	TxDepthStep                                     uint64
	CpDepthStep                                     uint64
	CpFlag                                          uint8
	SmoothedTimestamp                               uint32
	SigVerifyData                                   *blockDepthUtil.SigVerifProofData
}
type ProofResponse struct {
	Proof   []byte
	Witness []byte
}

type Proof struct {
	Proof   string // hex
	Witness string
}

type BtcBaseRequest struct {
	Data *btcbase.BaseLevelProofData
}

type BtcMiddleRequest struct {
	Data   *btcmiddle.BatchedProofData
	Proofs []Proof
}

type BtcUpperRequest struct {
	Data   *btcupper.BatchedProofData
	Proofs []Proof
}

type BtcBulkRequest struct {
	Data *blockdepth.BlockBulkProofData
}

type BtcBulkResponse struct {
	Proof   []byte
	Witness []byte
}

type BtcPackedRequest struct {
	Data      *blockdepth.BlockBulkProofData
	Recursive Proof
	Bulk      Proof
}

type BtcPackResponse struct {
	Proof   []byte
	Witness []byte
}

type TxInEth2ProveRequest struct {
	Version string
	TxHash  string
	TxData  *txineth2Utils.TxInEth2ProofData
}

type TxInEth2ProveResponse struct {
	Proof   []byte
	Witness []byte
}

type BlockHeaderRequest struct {
	Index uint64
	Data  *proverType.BeaconHeaderChain
}

type BlockHeaderResponse struct {
	Proof   []byte
	Witness []byte
}
type BlockHeaderFinalityRequest struct {
	Index          uint64
	FinalityUpdate *proverType.FinalityUpdate
	SyncCommittee  *proverType.SyncCommittee
}

type BlockHeaderFinalityResponse struct {
	Proof   []byte
	Witness []byte
}

type RedeemRequest struct {
	TxHash                           string
	Version                          string
	TxProof, BhProof, BhfProof, Duty Proof
	GenesisScRoot,
	CurrentSCSSZRoot string
	TxId            string
	MinerReward     string
	SigHashes       []string
	IsFront         bool
	NbBeaconHeaders int
}

type RedeemResponse struct {
	Proof   []byte
	Witness []byte

	ProofSgxBytes []byte //for sgx
}

type SyncCommGenesisRequest struct {
	Period     uint64 `json:"period"`
	Version    string `json:"version"`
	FirstProof Proof
	GenesisID  []byte
	//FirstID                 []byte
	SecondID []byte
	//RecursiveFp []byte
}

type SyncCommGenesisResponse struct {
	Version   string           `json:"version"`
	Period    uint64           `json:"period"`
	ProofType common.ProofType `json:"proofType"`
	Proof     []byte
	Witness   []byte
}

type SyncCommInnerRequest struct {
	Data    *proverType.SyncCommittee
	Period  uint64
	Index   uint64
	Version string
}
type SyncCommInnerResponse struct {
	Version   string           `json:"version"`
	Period    uint64           `json:"period"`
	Index     uint64           `json:"index"`
	ProofType common.ProofType `json:"proofType"`
	Proof     []byte           `json:"proof"`
	Witness   []byte           `json:"witness"`
}
type SyncCommOuterRequest struct {
	Data        *proverType.SyncCommittee
	Period      uint64
	Version     string
	InnerProofs []Proof
}

type SyncCommUnitsRequest struct {
	Data    *WrapSyncCommitteeUpdate
	Outer   Proof
	Index   uint64
	Version string
}

type SyncCommUnitsResponse struct {
	Version      string           `json:"version"`
	Period       uint64           `json:"period"`
	ProofType    common.ProofType `json:"proofType"`
	Proof        []byte           `json:"proof"`
	Witness      []byte           `json:"witness"`
	OuterProof   []byte           `json:"outerProof"`
	OuterWitness []byte           `json:"outerWitness"`
}

type SyncCommDutyRequest struct {
	Period                         uint64
	Version                        string
	Choice                         string `json:"choice"`
	FirstProof, SecondProof, Outer Proof
	BeginId, RelayId, EndId        string
	ScIndex                        int
	Update                         *proverType.SyncCommitteeUpdate
}

type SyncCommDutyResponse struct {
	Version          string           `json:"version"`
	Period           uint64           `json:"period"`
	ProofType        common.ProofType `json:"proofType"`
	Proof            []byte
	Witness          []byte
	RecursiveProof   []byte
	RecursiveWitness []byte
}

type ProofInfo struct {
	ProofType int    `json:"-"`
	TxId      string `json:"txId"`
	Proof     string `json:"proof"`
	Status    int    `json:"status"`
}
