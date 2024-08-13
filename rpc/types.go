package rpc

import (
	"time"

	btcprovertypes "github.com/lightec-xyz/btc_provers/circuits/types"
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
	Data   *btcmiddle.MidLevelProofData
	Proofs []Proof
}

type BtcUpperRequest struct {
	Data   *btcupper.UpperLevelProofData
	Proofs []Proof
}

type BtcBulkRequest struct {
	Data *btcprovertypes.BlockHeaderChain
}

type BtcBulkResponse struct {
	Proof   []byte
	Witness []byte
}

type BtcPackedRequest struct {
	Data *btcprovertypes.BlockHeaderChain
}

type BtcPackResponse struct {
	Proof   []byte
	Witness []byte
}

type BtcWrapRequest struct {
	Flag, Proof, Witness, BeginHash, EndHash string
	NbBlocks                                 uint64
}

type BtcWrapResponse struct {
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
	Index            uint64
	GenesisSCSSZRoot string
	RecursiveProof, RecursiveWitness, OuterProof,
	OuterWitness string
	FinalityUpdate *proverType.FinalityUpdate
	ScUpdate       *proverType.SyncCommitteeUpdate
}

type BlockHeaderFinalityResponse struct {
	Proof   []byte
	Witness []byte
}

type DepositRequest struct {
	TxHash    string
	BlockHash string
	Data      *grUtil.GrandRollupProofData
}

type DepositResponse struct {
	TxHash  string
	Proof   []byte
	Witness []byte
}

type RedeemRequest struct {
	TxHash                                                       string
	Version                                                      string
	TxProof, TxWitness, BhProof, BhWitness, BhfProof, BhfWitness string
	BeginId, EndId, GenesisScRoot,
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

type VerifyRequest struct {
	TxHash    string
	BlockHash string
	Data      *grUtil.GrandRollupProofData
}

type VerifyResponse struct {
	TxHash string
	Proof  []byte
	Wit    []byte
}

type SyncCommGenesisRequest struct {
	Period  uint64 `json:"period"`
	Version string `json:"version"`
	FirstProof,
	FirstWitness,
	SecondProof,
	SecondWitness string
	GenesisID   []byte
	FirstID     []byte
	SecondID    []byte
	RecursiveFp []byte
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
	Period  uint64
	Version string
	Choice  string `json:"choice"`
	FirstProof,
	FirstWitness,
	SecondProof,
	SecondWitness string
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
