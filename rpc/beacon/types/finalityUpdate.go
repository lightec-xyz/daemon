package types

type LightClientFinalityUpdateResp struct {
	Version string                    `json:"version"`
	Data    LightClientFinalityUpdate `json:"data"`
}

type LightClientFinalityUpdate struct {
	AttestedHeader  AttestedHeader  `json:"attested_header"`
	FinalizedHeader FinalizedHeader `json:"finalized_header"`
	FinalityBranch  []string        `json:"finality_branch"`
	SyncAggregate   SyncAggregate   `json:"sync_aggregate"`
	SignatureSlot   string          `json:"signature_slot"`
}
