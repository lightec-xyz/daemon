package common

import (
	"github.com/lightec-xyz/daemon/store"
)

const (
	LatestPeriodKey = "latestPeriod"
	LatestSlotKey   = "latestFinalitySlot"
)
const (
	IndexTable    store.Table = "index"
	UpdateTable   store.Table = "update"
	FinalityTable store.Table = "finalityUpdate"
	RequestTable  store.Table = "request"

	InnerTable             store.Table = "inner"
	OuterTable             store.Table = "outer"
	UnitTable              store.Table = "unit"
	GenesisTable           store.Table = "genesis"
	RecursiveTable         store.Table = "recursive"
	DutyTable              store.Table = "duty"
	TxesTable              store.Table = "txInEth2"
	BeaconHeaderTable      store.Table = "beaconHeader"
	BhfTable               store.Table = "bhf"
	RedeemTable            store.Table = "redeem"
	SgxRedeemTable         store.Table = "sgxRedeem"
	BackendRedeemTable     store.Table = "backendRedeem"
	BtcTimestampTable      store.Table = "btcTimestamp"
	BtcBaseTable           store.Table = "btcBase"
	BtcMiddleTable         store.Table = "btcMiddle"
	BtcUpperTable          store.Table = "btcUpper"
	BtcBulkTable           store.Table = "btcBulk"
	BtcDuperRecursiveTable store.Table = "btcDuperRecursive"
	BtcDepositTable        store.Table = "btcDeposit"
	BtcChangeTable         store.Table = "btcChange"
	BtcDepthRecursiveTable store.Table = "btcDepthRecursive"
	BtcUpdateCpTable       store.Table = "btcUpdateCp"
)

func GenKey(zkType ProofType, prefix, index, end uint64, hash string) store.FileKey {
	table := ProofTypeToTable(zkType)
	switch zkType {
	case SyncComInnerType:
		return store.GenFileKey(table, prefix, index)
	case SyncComOuterType:
		return store.GenFileKey(table, index)
	case SyncComUnitType:
		return store.GenFileKey(table, index)
	case SyncComRecursiveType:
		return store.GenFileKey(table, index)
	case SyncComDutyType:
		return store.GenFileKey(table, index)
	case TxInEth2Type:
		return store.GenFileKey(table, hash)
	case BeaconHeaderType:
		return store.GenFileKey(table, index, end)
	case BeaconHeaderFinalityType:
		return store.GenFileKey(table, index)
	case RedeemTxType:
		return store.GenFileKey(table, hash)
	case SgxRedeemTxType:
		return store.GenFileKey(table, hash)
	case BackendRedeemTxType:
		return store.GenFileKey(table, hash)
	case BtcBulkType:
		return store.GenFileKey(table, index, end)
	case BtcBaseType:
		return store.GenFileKey(table, index, end)
	case BtcMiddleType:
		return store.GenFileKey(table, index, end)
	case BtcUpperType:
		return store.GenFileKey(table, index, end)
	case BtcDuperRecursiveType:
		return store.GenFileKey(table, index, end)
	case BtcDepthRecursiveType:
		return store.GenFileKey(table, prefix, index, end)
	case BtcDepositType:
		return store.GenFileKey(table, hash)
	case BtcChangeType:
		return store.GenFileKey(table, hash)
	case BtcUpdateCpType:
		return store.GenFileKey(table, hash)
	case BtcTimestampType:
		return store.GenFileKey(table, index, end)
	default:
		panic("unknown table type" + table) //todo
	}
}

func ProofTypeToTable(pType ProofType) store.Table {
	switch pType {
	case SyncComInnerType:
		return InnerTable
	case SyncComOuterType:
		return OuterTable
	case SyncComUnitType:
		return UnitTable
	case SyncComRecursiveType:
		return RecursiveTable
	case SyncComDutyType:
		return DutyTable
	case TxInEth2Type:
		return TxesTable
	case BeaconHeaderType:
		return BeaconHeaderTable
	case BeaconHeaderFinalityType:
		return BhfTable
	case RedeemTxType:
		return RedeemTable
	case SgxRedeemTxType:
		return SgxRedeemTable
	case BackendRedeemTxType:
		return BackendRedeemTable
	case BtcBaseType:
		return BtcBaseTable
	case BtcMiddleType:
		return BtcMiddleTable
	case BtcUpperType:
		return BtcUpperTable
	case BtcDuperRecursiveType:
		return BtcDuperRecursiveTable
	case BtcBulkType:
		return BtcBulkTable
	case BtcDepthRecursiveType:
		return BtcDepthRecursiveTable
	case BtcDepositType:
		return BtcDepositTable
	case BtcChangeType:
		return BtcChangeTable
	case BtcUpdateCpType:
		return BtcUpdateCpTable
	case BtcTimestampType:
		return BtcTimestampTable
	default:
		panic("unknown proof type" + pType.Name()) //todo
	}
}
