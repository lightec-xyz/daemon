package node

import (
	"fmt"
	"strings"

	"github.com/lightec-xyz/daemon/common"
)

const protocolSeparator = "_"

const (
	beaconEthNumberPrefix      = "bh"
	btcTxHeightPrefix          = "bth"
	btcBlockHashPrefix         = "bb"
	btcTxPrefix                = "bt"
	btcBlockDepositCountPrefix = "bc"
	btcHeaderPrefix            = "bhp"
	btcBlockPrefix             = "btb"
	btcTxParamPrefix           = "btp"
	checkpointPrefix           = "cp"
	chainForkPrefix            = "cf"
	destChainHashPrefix        = "d" // d_ + hash
	dfinityBlockSigPrefix      = "db"
	ethTxHeightPrefix          = "eth"
	ethBlockHashPrefix         = "eb"
	finalityUpdateSlotPrefix   = "fu"
	minerAddrPrefix            = "ma"
	minerPowerPrefix           = "mp"
	noncePrefix                = "n"
	proofPrefix                = "p" // p_ + hash
	pendingReqPrefix           = "pr"
	pendingProofRespPrefix     = "prp"
	txProvedPrefix             = "tp"
	txSlotPrefix               = "ts"
	txFinalizeSlotPrefix       = "tfs"
	txPrefix                   = "t" // t_ + hash
	txMinerRewardPrefix        = "tm"
	unGenProofPrefix           = "u"  // u_ + hash
	unSubmitTxPrefix           = "us" // s_ + hash
	updateUtxoDestPrefix       = "uu"
)

const (
	workerIdKey               = "workerIdKey"
	zkVerifyKey               = "zkVerifyKey"
	latestCheckPointHeightKey = "latestCheckPointHeight"
	latestUpdateCpKey         = "latestUpdateCpKey"
	latestIcpSignatureKey     = "latestIcpSignatureKey"
	maxGasPriceKey            = "maxGasPriceKey"
	submitMaxValueKey         = "submitMaxValue"
	submitMinValueKey         = "submitMinValue"
)

var (
	btcCurHeightKey = []byte("btcCurHeight")
	ethCurHeightKey = []byte("ethCurHeight")
	beaconLatestKey = []byte("beaconLatest")
)

type DbTx struct {
	Height       uint64
	TxIndex      uint
	Hash         string
	BlockHash    string
	BlockTime    uint64
	TxType       common.TxType
	ChainType    common.ChainType
	ProofType    common.ProofType
	Proved       bool
	Amount       int64
	GenProofNums int

	// bitcoin chain
	EthAddr string

	// ethereum chain
	LogIndex  uint
	UtxoId    string
	UtxoIndex int64
	Sender    string
	Receiver  string

	TxSlot        uint64 `cbor:"omitzero"`
	FinalizedSlot uint64 `cbor:"omitzero"`

	// for btc
	CheckPointHeight uint64
	LatestHeight     uint64
	CpMinDepth       uint64
	SigSigned        bool // flag: if icp block hash signed
}

func (t *DbTx) GenReset() {
	t.FinalizedSlot = 0
	t.TxSlot = 0
	t.LatestHeight = 0
	t.CheckPointHeight = 0
	t.GenProofNums = t.GenProofNums + 1
	t.SigSigned = false
}

type DbProof struct {
	TxHash    string           `json:"txId"`
	ProofType common.ProofType `json:"type"`
	Status    int              `json:"status"`
	Proof     string           `json:"Proof"`
}

type DbUnGenProof struct {
	Hash      string
	ProofType common.ProofType
	ChainType common.ChainType
	Height    uint64
	TxIndex   uint
	Amount    uint64
}

type DbUnSubmitTx struct {
	Hash        string
	ProofType   common.ProofType
	Proof       string
	Timestamp   int64
	ConfirmHash string //
	Status      int
}

type DbMiner struct {
	Miner     string
	Power     uint64
	Timestamp uint64
}
type ChainFork struct {
	ForkHeight uint64
	Chain      common.ChainType
	Timestamp  int64
}

type DbIcpSignature struct {
	Signature string
	Hash      string
	Height    uint64
}

func dbBtcDepositCountKey(height uint64) []byte {
	return genKey(btcBlockDepositCountPrefix, height)
}
func dbBtcBlockKey(hash string) []byte {
	return genKey(btcBlockPrefix, hash)
}
func dbBtcHeaderKey(hash string) []byte {
	return genKey(btcHeaderPrefix, hash)
}
func dbChainForkKey(chain string, timestamp int64) []byte {
	return genKey(chainForkPrefix, chain, timestamp)
}
func dbProofId(txId string) []byte {
	return genKey(proofPrefix, txId)
}
func dbTxId(txId string, txType common.TxType, logIndex uint) []byte {
	return genKey(txPrefix, txId, txType, logIndex)
}

func dbDestId(txId string) []byte {
	return genKey(destChainHashPrefix, txId)
}

func dbUnGenProofId(chain common.ChainType, txId string) []byte {
	return genKey(unGenProofPrefix, chain, txId)
}

func dbUnSubmitTxId(txId string) []byte {
	return genKey(unSubmitTxPrefix, txId)
}

func dbAddrPrefixTxId(addr string, txType common.TxType, txId string) []byte {
	return genKey(addr, txType, txId)
}

func dbBeaconEthNumberId(number uint64) []byte {
	return genKey(beaconEthNumberPrefix, number)
}

func dbTxSlotId(slot uint64, hash string) []byte {
	return genKey(txSlotPrefix, slot, hash)
}

func dbTxFinalizeSlotId(slot uint64, hash string) []byte {
	return genKey(txFinalizeSlotPrefix, slot, hash)
}

func dbPendingRequestId(id string) []byte {
	return genKey(pendingReqPrefix, id)
}

func dbAddrNonceId(network, addr string) []byte {
	return genKey(noncePrefix, network, addr)
}

func dbTaskTimeId(id string, status common.ProofStatus) []byte {
	return genKey(id, status)
}

func dbFinalityUpdateSlotId(slot uint64) []byte {
	return genKey(finalityUpdateSlotPrefix, slot)
}

func dbDfinityBlockSigId(height uint64) []byte {
	return genKey(dfinityBlockSigPrefix, height)
}

func dbBtcTxId(hash string) []byte {
	return genKey(btcTxPrefix, hash)
}

func dbBtcBlockHashKey(height uint64) []byte {
	return genKey(btcBlockHashPrefix, height)
}

func dbEthBlockHashKey(height uint64) []byte {
	return genKey(ethBlockHashPrefix, height)
}
func ethTxHeightKey(height uint64, txId string) []byte {
	return genKey(ethTxHeightPrefix, height, txId)
}

func btcTxHeightKey(height uint64, txId string) []byte {
	return genKey(btcTxHeightPrefix, height, txId)
}

func dbMinerPowerKey(addr string) []byte {
	return genKey(minerPowerPrefix, addr)
}

func dbMinerAddrKey(addr string) []byte {
	return genKey(minerAddrPrefix, addr)
}
func dbProofResponseId(requestId string) []byte {
	return genKey(pendingProofRespPrefix, requestId)
}

func dbCheckpointKey(height uint64) []byte {
	return genKey(checkpointPrefix, height)
}
func DbTxMinerReward(txId string) []byte {
	return genKey(txMinerRewardPrefix, txId)
}

func dbTxProvedKey(txId string) []byte {
	return genKey(txProvedPrefix, txId)
}

func dbUpdateUtxoDestKey(hash string) []byte {
	return genKey(updateUtxoDestPrefix, hash)
}

func dbBtcTxParamKey(txId string) []byte {
	return genKey(btcTxParamPrefix, txId)
}

// careful modify this struct if you want to change
func genKey(values ...interface{}) []byte {
	var key string
	for _, value := range values {
		key = key + fmt.Sprintf("%v%v", trimOx(fmt.Sprintf("%v", value)), protocolSeparator)
	}
	fixKey := strings.ToLower(strings.TrimSuffix(key, protocolSeparator))
	return []byte(fixKey)
}

func genPrefix(values ...interface{}) []byte {
	var key string
	for _, value := range values {
		key = key + fmt.Sprintf("%v%v", trimOx(fmt.Sprintf("%v", value)), protocolSeparator)
	}
	return []byte(strings.ToLower(key))
}
