package types

import (
	"encoding/json"
	"fmt"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	eth "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"reflect"
)

type SyncAggregate struct {
	SyncCommitteeBits      string `json:"sync_committee_bits"`
	SyncCommitteeSignature string `json:"sync_committee_signature"`
}
type Beacon struct {
	Slot          string `json:"slot"`
	ProposerIndex string `json:"proposer_index"`
	ParentRoot    string `json:"parent_root"`
	StateRoot     string `json:"state_root"`
	BodyRoot      string `json:"body_root"`
}

type AttestedHeader struct {
	Beacon          Beacon    `json:"beacon"`
	Execution       Execution `json:"execution"`
	ExecutionBranch []string  `json:"execution_branch"`
}
type FinalizedHeader struct {
	Beacon          Beacon    `json:"beacon"`
	Execution       Execution `json:"execution"`
	ExecutionBranch []string  `json:"execution_branch"`
}

type SyncCommittee struct {
	Pubkeys         []string `json:"pubkeys"`
	AggregatePubkey string   `json:"aggregate_pubkey"`
}
type Execution struct {
	ParentHash       string `json:"parent_hash"`
	FeeRecipient     string `json:"fee_recipient"`
	StateRoot        string `json:"state_root"`
	ReceiptsRoot     string `json:"receipts_root"`
	LogsBloom        string `json:"logs_bloom"`
	PrevRandao       string `json:"prev_randao"`
	BlockNumber      string `json:"block_number"`
	GasLimit         string `json:"gas_limit"`
	GasUsed          string `json:"gas_used"`
	Timestamp        string `json:"timestamp"`
	ExtraData        string `json:"extra_data"`
	BaseFeePerGas    string `json:"base_fee_per_gas"`
	BlockHash        string `json:"block_hash"`
	TransactionsRoot string `json:"transactions_root"`
	WithdrawalsRoot  string `json:"withdrawals_root"`
	BlobGasUsed      string `json:"blob_gas_used"`
	ExcessBlobGas    string `json:"excess_blob_gas"`
}

func (sc *SyncCommittee) ToConsensus() (*eth.SyncCommittee, error) {
	var value structs.SyncCommittee
	err := ParseObj(sc, &value)
	if err != nil {
		return nil, err
	}
	syncCommittee, err := value.ToConsensus()
	if err != nil {
		return nil, err
	}
	return syncCommittee, nil
}

func ParseObj(src, dst interface{}) error {
	if reflect.ValueOf(dst).Kind() != reflect.Ptr {
		return fmt.Errorf("dst must be a pointer")
	}
	srcBytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(srcBytes, dst)
	if err != nil {
		return err
	}
	return nil
}
