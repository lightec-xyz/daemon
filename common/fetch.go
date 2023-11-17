package common

import (
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	eth "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
)

/*
careful modify this struct if you want to change
the struct is store in datadir
*/

type LightClientUpdateResponse struct {
	Version string             `json:"version"`
	Data    *LightClientUpdate `json:"data"`
}

type LightClientFinalityUpdateEvent struct {
	Version string                     `json:"version"`
	Data    *LightClientFinalityUpdate `json:"data"`
}

type BootstrapResponse struct {
	Version string     `json:"version"`
	Data    *Bootstrap `json:"data"`
}

type LightClientUpdate struct {
	AttestedHeader          *BeaconBlockHeader `json:"attested_header"`
	NextSyncCommittee       *SyncCommittee     `json:"next_sync_committee,omitempty"`
	FinalizedHeader         *BeaconBlockHeader `json:"finalized_header,omitempty"`
	SyncAggregate           *SyncAggregate     `json:"sync_aggregate"`
	NextSyncCommitteeBranch []string           `json:"next_sync_committee_branch,omitempty"`
	FinalityBranch          []string           `json:"finality_branch,omitempty"`
	SignatureSlot           string             `json:"signature_slot"`
}

type Bootstrap struct {
	Header                     *BeaconBlockHeader `json:"header"`
	CurrentSyncCommittee       *SyncCommittee     `json:"current_sync_committee"`
	CurrentSyncCommitteeBranch []string           `json:"current_sync_committee_branch"`
}

type LightClientFinalityUpdate struct {
	AttestedHeader  *BeaconBlockHeader `json:"attested_header"`
	FinalizedHeader *BeaconBlockHeader `json:"finalized_header"`
	FinalityBranch  []string           `json:"finality_branch"`
	SyncAggregate   *SyncAggregate     `json:"sync_aggregate"`
	SignatureSlot   string             `json:"signature_slot"`
}

type BeaconBlockHeader struct {
	Slot          string `json:"slot"`
	ProposerIndex string `json:"proposer_index"`
	ParentRoot    string `json:"parent_root"`
	StateRoot     string `json:"state_root"`
	BodyRoot      string `json:"body_root"`
}

func (bh *BeaconBlockHeader) ToConsensus() (*eth.BeaconBlockHeader, error) {
	var value structs.BeaconBlockHeader
	err := ParseObj(bh, &value)
	if err != nil {
		return nil, err
	}
	blockHeader, err := value.ToConsensus()
	if err != nil {
		return nil, err
	}
	return blockHeader, nil
}

type SyncAggregate struct {
	SyncCommitteeBits      string `json:"sync_committee_bits"`
	SyncCommitteeSignature string `json:"sync_committee_signature"`
}

type SyncCommittee struct {
	Pubkeys         []string `json:"pubkeys"`
	AggregatePubkey string   `json:"aggregate_pubkey"`
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
