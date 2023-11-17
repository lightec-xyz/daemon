package types

type LightClientUpdateResp struct {
	Version string            `json:"version"`
	Data    LightClientUpdate `json:"data"`
}

type NextSyncCommittee struct {
	Pubkeys         []string `json:"pubkeys"`
	AggregatePubkey string   `json:"aggregate_pubkey"`
}

type LightClientUpdate struct {
	AttestedHeader          AttestedHeader    `json:"attested_header"`
	NextSyncCommittee       NextSyncCommittee `json:"next_sync_committee"`
	FinalizedHeader         FinalizedHeader   `json:"finalized_header"`
	SyncAggregate           SyncAggregate     `json:"sync_aggregate"`
	NextSyncCommitteeBranch []string          `json:"next_sync_committee_branch"`
	FinalityBranch          []string          `json:"finality_branch"`
	SignatureSlot           string            `json:"signature_slot"`
}
