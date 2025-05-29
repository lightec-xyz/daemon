package common

import (
	"fmt"
	"github.com/lightec-xyz/daemon/store"
	"time"

	btcproverCommon "github.com/lightec-xyz/btc_provers/circuits/common"
	proverCommon "github.com/lightec-xyz/provers/common"
)

const (
	ZkDebugEnv   = "zkDebug"
	ZkProofTypes = "zkProofTypes"
	DbNameSpace  = "zkbtc"
)

var BlockChainPlan = ReverseU32(btcproverCommon.BlockChainPlan[:])
var BtcBlockDepthPlan = []uint64{btcproverCommon.MinNewCpDepth, btcproverCommon.HalfMinNewCpDepth, btcproverCommon.CapacityAbsorbedBulk}

const (
	BtcBaseDistance   = btcproverCommon.CapacityBaseLevel
	BtcMiddleDistance = btcproverCommon.CapacitySuperBatch
	BtcUpperDistance  = btcproverCommon.CapacityDifficultyBlock
	CapacityMiniLevel = btcproverCommon.CapacityMiniLevel
	BtcTxMinDepth     = btcproverCommon.FirstMinTxDepth
	BtcTxUnitMaxDepth = btcproverCommon.LastMinTxDepth
	BtcCpMinDepth     = btcproverCommon.MinNewCpDepth
	SyncInnerNum      = proverCommon.NbBatches // todo

)

const (
	SlotPerPeriod         = 8192
	MaxDiffTxFinalitySlot = 32
	BtcLatestBlockMaxDiff = 24 // todo
)

type TxMode int

const (
	NormalTx = iota
	OnlyMigrateTx
)

type TxType int

const (
	DepositTx TxType = iota + 1
	RedeemTx
	UpdateUtxoTx
	DepositRewardTx
	RedeemRewardTx
)

func (tt TxType) String() string {
	switch tt {
	case DepositTx:
		return "deposit"
	case RedeemTx:
		return "redeem"
	case UpdateUtxoTx:
		return "updateutxo"
	case DepositRewardTx:
		return "depositReward"
	case RedeemRewardTx:
		return "redeemReward"
	default:
		return "unknown"
	}
}

func ToTxType(value string) (TxType, error) {
	switch value {
	case "deposit":
		return DepositTx, nil
	case "redeem":
		return RedeemTx, nil
	default:
		return 0, fmt.Errorf("unknown tx type")
	}

}

type ChainType int

const (
	BitcoinChain ChainType = iota + 1
	EthereumChain
)

func (ct ChainType) String() string {
	switch ct {
	case BitcoinChain:
		return "bitcoin"
	case EthereumChain:
		return "ethereum"
	default:
		return "unknown"
	}
}

type Network string

const (
	ETH   Network = "eth"
	ICP   Network = "icp"
	Oasis Network = "oasis"
)

func (n Network) String() string {
	switch n {
	case ETH:
		return "eth"
	case Oasis:
		return "oasis"
	case ICP:
		return "icp"
	default:
		return "unknown"
	}
}

type ProofStatus int

const (
	ProofDefault ProofStatus = iota
	ProofFinalized
	ProofQueued
	ProofGenerating
	ProofSuccess
	ProofFailed
)

func (ps ProofStatus) String() string {
	switch ps {
	case ProofDefault:
		return "default"
	case ProofFinalized:
		return "finalized"
	case ProofQueued:
		return "queued"
	case ProofGenerating:
		return "generating"
	case ProofSuccess:
		return "success"
	case ProofFailed:
		return "failed"
	default:
		return "unknown"
	}
}

type Mode string

const (
	Client  Mode = "client"
	Cluster Mode = "cluster"
	Custom  Mode = "custom"
)

type ProofWeight int

const (
	WeightDefault ProofWeight = iota
	OneLevel
	TwoLevel
	ThreeLevel
	FourLevel
	FiveLevel
	SixLevel
	SevenLevel
	EightLevel
	NineLevel
	TenLevel
	Highest
)

type ProofType int

const (
	SyncComOuterType ProofType = iota + 1
	SyncComInnerType
	SyncComUnitType
	SyncComRecursiveType
	SyncComDutyType
	TxInEth2Type
	BeaconHeaderType
	BeaconHeaderFinalityType
	RedeemTxType
	BackendRedeemTxType
	BtcFreeBulkType
	BtcBulkType
	BtcBaseType
	BtcMiddleType
	BtcUpperType
	BtcDuperRecursiveType
	BtcDepthRecursiveType
	BtcDepositType
	BtcChangeType
	BtcUpdateCpType
	BtcTimestampType
	SgxRedeemTxType
)

func (p ProofType) Name() string {
	switch p {
	case SyncComInnerType:
		return "unitInner"
	case SyncComOuterType:
		return "unitOuter"
	case SyncComUnitType:
		return "syncComUnitType"
	case SyncComRecursiveType:
		return "syncComRecursiveType"
	case SyncComDutyType:
		return "syncComDutyType"
	case TxInEth2Type:
		return "txInEth2"
	case BeaconHeaderType:
		return "beaconHeaderType"
	case BeaconHeaderFinalityType:
		return "beaconHeaderFinalityType"
	case RedeemTxType:
		return "redeemTxType"
	case BackendRedeemTxType:
		return "backendRedeemTxType"
	case BtcBulkType:
		return "btcBulkType"
	case BtcFreeBulkType:
		return "btcFreeBulkType"
	case BtcBaseType:
		return "btcBaseType"
	case BtcMiddleType:
		return "btcMiddleType"
	case BtcUpperType:
		return "btcUpperType"
	case BtcDuperRecursiveType:
		return "btcDuperRecursive"
	case BtcDepthRecursiveType:
		return "btcDepthRecursiveType"
	case BtcDepositType:
		return "btcDepositType"
	case BtcChangeType:
		return "btcChangeType"
	case BtcUpdateCpType:
		return "btcUpdateCpType"
	case BtcTimestampType:
		return "btcTimestampType"
	default:
		return "unKnown"
	}
}

func (p ProofType) Weight() ProofWeight {
	switch p {
	case RedeemTxType, BtcDepositType, BtcChangeType:
		return Highest
	case SyncComDutyType, BackendRedeemTxType, BtcFreeBulkType:
		return NineLevel
	case SyncComUnitType:
		return EightLevel
	case SyncComOuterType:
		return SevenLevel
	case SyncComInnerType, BtcDuperRecursiveType, BtcDepthRecursiveType, BtcTimestampType:
		return SixLevel
	case BtcUpperType:
		return FiveLevel
	case BtcMiddleType:
		return FourLevel
	case BeaconHeaderFinalityType, TxInEth2Type, BtcBulkType, BeaconHeaderType, BtcBaseType:
		return ThreeLevel
	default:
		return WeightDefault
	}
}

func (p ProofType) ProveTime() time.Duration {
	// todo
	switch p {
	case SyncComUnitType:
		return 1*time.Hour + 30*time.Minute
	case BtcBaseType, BtcBulkType, BeaconHeaderType, RedeemTxType, BackendRedeemTxType:
		return 15 * time.Minute
	case BtcMiddleType, BtcUpperType, BeaconHeaderFinalityType, TxInEth2Type:
		return 15 * time.Minute
	case BtcDuperRecursiveType, BtcDepthRecursiveType, SyncComRecursiveType, SyncComDutyType:
		return 15 * time.Minute
	case BtcDepositType:
		return 20 * time.Minute
	case BtcChangeType:
		return 25 * time.Minute
	default:
		return 20 * time.Minute
	}
}

func (p ProofType) ConstraintQuantity() uint64 {
	switch p { // todo
	case SyncComInnerType:
		return 100
	case SyncComOuterType:
		return 100
	case SyncComUnitType:
		return 100
	case SyncComRecursiveType:
		return 100
	case SyncComDutyType:
		return 100
	case TxInEth2Type:
		return 100
	case BeaconHeaderType:
		return 100
	case BeaconHeaderFinalityType:
		return 100
	case RedeemTxType:
		return 100
	case BackendRedeemTxType:
		return 100
	case BtcBulkType:
		return 100
	case BtcBaseType:
		return 100
	case BtcMiddleType:
		return 100
	case BtcUpperType:
		return 100
	case BtcDuperRecursiveType:
		return 100
	case BtcDepthRecursiveType:
		return 100
	case BtcDepositType:
		return 100
	case BtcChangeType:
		return 100
	case BtcUpdateCpType:
		return 100
	case BtcTimestampType:
		return 100
	default:
		return 0

	}
}

func ToZkProofType(str string) (ProofType, error) {
	switch store.Table(str) {
	case InnerTable:
		return SyncComInnerType, nil
	case OuterTable:
		return SyncComOuterType, nil
	case UnitTable:
		return SyncComUnitType, nil
	case RecursiveTable:
		return SyncComRecursiveType, nil
	case DutyTable:
		return SyncComDutyType, nil
	case TxesTable:
		return TxInEth2Type, nil
	case BeaconHeaderTable:
		return BeaconHeaderType, nil
	case BhfTable:
		return BeaconHeaderFinalityType, nil
	case RedeemTable:
		return RedeemTxType, nil
	case BtcBulkTable:
		return BtcBulkType, nil
	case BtcBaseTable:
		return BtcBaseType, nil
	case BtcMiddleTable:
		return BtcMiddleType, nil
	case BtcUpperTable:
		return BtcUpperType, nil
	case BtcDuperRecursiveTable:
		return BtcDuperRecursiveType, nil
	case BtcDepthRecursiveTable:
		return BtcDepthRecursiveType, nil
	case BtcDepositTable:
		return BtcDepositType, nil
	case BtcChangeTable:
		return BtcChangeType, nil
	case BtcTimestampTable:
		return BtcTimestampType, nil
	default:
		return 0, fmt.Errorf("uKnown zk proof type %v", str)
	}
}

func IsBtcProofType(proofType ProofType) bool {
	switch proofType {
	case BtcBulkType, BtcTimestampType, BtcBaseType, BtcMiddleType, BtcUpperType, BtcDuperRecursiveType, BtcDepthRecursiveType,
		BtcDepositType, BtcChangeType, BtcUpdateCpType:
		return true
	default:
		return false
	}
}
