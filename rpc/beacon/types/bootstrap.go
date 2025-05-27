package types

type BootstrapResp struct {
	Version string    `json:"version"`
	Data    Bootstrap `json:"data"`
}
type Header struct {
	Slot          string `json:"slot"`
	ProposerIndex string `json:"proposer_index"`
	ParentRoot    string `json:"parent_root"`
	StateRoot     string `json:"state_root"`
	BodyRoot      string `json:"body_root"`
}
type CurrentSyncCommittee struct {
	Pubkeys         []string `json:"pubkeys"`
	AggregatePubkey string   `json:"aggregate_pubkey"`
}
type Bootstrap struct {
	Header                     Header               `json:"header"`
	CurrentSyncCommittee       CurrentSyncCommittee `json:"current_sync_committee"`
	CurrentSyncCommitteeBranch []string             `json:"current_sync_committee_branch"`
}

type BeaconNodeVersion struct {
	Data Version `json:"data"`
}

type Version struct {
	Version string `json:"version"`
}
