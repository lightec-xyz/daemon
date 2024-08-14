package rpc

import (
	"github.com/lightec-xyz/btc_provers/utils/blockdepth"
	"time"

	btcbase "github.com/lightec-xyz/btc_provers/utils/blockchain"
	btcmiddle "github.com/lightec-xyz/btc_provers/utils/blockchain"
	btcupper "github.com/lightec-xyz/btc_provers/utils/blockchain"
	recursiveUtil "github.com/lightec-xyz/btc_provers/utils/blockchain"
	grUtil "github.com/lightec-xyz/btc_provers/utils/txinchain"
	"github.com/lightec-xyz/daemon/common"
	ethblock "github.com/lightec-xyz/provers/circuits/fabric/tx-in-eth2"
	proverType "github.com/lightec-xyz/provers/circuits/types"
	"github.com/lightec-xyz/reLight/circuits/utils"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
)

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
	Desc    string
}

//------

type BtcDuperRecursiveRequest struct {
	Data         *recursiveUtil.RecursiveProofData
	First, Duper Proof
}

type BtcDepthRequest struct {
	Data            *blockdepth.BulksProofData
	Recursive, Unit Proof
}
type BtcChainRequest struct {
	Data                             *recursiveUtil.BlockChainProofData
	Recursive, Base, MidLevel, Upper Proof
}

type BtcDepositRequest struct {
	Data                         *grUtil.TxInChainProofData
	BlockChain, TxDepth, CpDepth Proof
	R, S                         string
	ProverAddr                   string
}

type BtcChangeRequest struct {
	Data                                 *grUtil.TxInChainProofData
	BlockChain, TxDepth, CpDepth, Redeem Proof
	R, S, ProverAddr                     string
}

type BtcGenesisRequest struct {
	Data   *recursiveUtil.RecursiveProofData
	First  Proof
	Second Proof
}

type BtcRecursiveRequest struct {
	Data   *recursiveUtil.RecursiveProofData
	First  Proof
	Second Proof
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
	Data      *blockdepth.BulksProofData
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
	Index                      uint64
	GenesisSCSSZRoot           string
	RecursiveProof, OuterProof Proof
	FinalityUpdate             *proverType.FinalityUpdate
	ScUpdate                   *proverType.SyncCommitteeUpdate
}

type BlockHeaderFinalityResponse struct {
	Proof   []byte
	Witness []byte
}

type RedeemRequest struct {
	TxHash                     string
	Version                    string
	TxProof, BhProof, BhfProof Proof
	BeginId, EndId             []byte
	GenesisScRoot,
	CurrentSCSSZRoot string
	TxVar      []string
	ReceiptVar []string
	//TxVar      *[tx.MaxTxUint128Len]frontend.Variable
	//ReceiptVar *[receipt.MaxReceiptUint128Len]frontend.Variable
}

type RedeemResponse struct {
	Proof   []byte
	Witness []byte
}

type SyncCommGenesisRequest struct {
	Period                  uint64 `json:"period"`
	Version                 string `json:"version"`
	FirstProof, SecondProof Proof
	GenesisID               []byte
	FirstID                 []byte
	SecondID                []byte
	RecursiveFp             []byte
}

type SyncCommGenesisResponse struct {
	Version   string             `json:"version"`
	Period    uint64             `json:"period"`
	ProofType common.ZkProofType `json:"proofType"`
	Proof     []byte
	Witness   []byte
}

type SyncCommUnitsRequest struct {
	Data    *utils.SyncCommitteeUpdate
	Index   uint64
	Version string
}

type SyncCommUnitsResponse struct {
	Version      string             `json:"version"`
	Period       uint64             `json:"period"`
	ProofType    common.ZkProofType `json:"proofType"`
	Proof        []byte             `json:"proof"`
	Witness      []byte             `json:"witness"`
	OuterProof   []byte             `json:"outerProof"`
	OuterWitness []byte             `json:"outerWitness"`
}

type SyncCommRecursiveRequest struct {
	Period                  uint64
	Version                 string
	Choice                  string `json:"choice"`
	FirstProof, SecondProof Proof
	BeginId,
	RelayId,
	EndId,
	RecursiveFp []byte
}

type SyncCommRecursiveResponse struct {
	Version   string             `json:"version"`
	Period    uint64             `json:"period"`
	ProofType common.ZkProofType `json:"proofType"`
	Proof     []byte
	Witness   []byte
}

type ProofInfo struct {
	ProofType int    `json:"-"`
	TxId      string `json:"txId"`
	Proof     string `json:"proof"`
	Status    int    `json:"status"`
}
