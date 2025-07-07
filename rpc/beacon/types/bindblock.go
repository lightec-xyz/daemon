package types

type BindBlockResp struct {
	Version             string `json:"version"`
	ExecutionOptimistic bool   `json:"execution_optimistic"`
	Finalized           bool   `json:"finalized"`
	Data                struct {
		Message struct {
			Slot          string `json:"slot"`
			ProposerIndex string `json:"proposer_index"`
			ParentRoot    string `json:"parent_root"`
			StateRoot     string `json:"state_root"`
			Body          struct {
				RandaoReveal string `json:"randao_reveal"`
				Eth1Data     struct {
					DepositRoot  string `json:"deposit_root"`
					DepositCount string `json:"deposit_count"`
					BlockHash    string `json:"block_hash"`
				} `json:"eth1_data"`
				Graffiti          string `json:"graffiti"`
				ProposerSlashings []any  `json:"proposer_slashings"`
				AttesterSlashings []any  `json:"attester_slashings"`
				Attestations      []struct {
					AggregationBits string `json:"aggregation_bits"`
					Data            struct {
						Slot            string `json:"slot"`
						Index           string `json:"index"`
						BeaconBlockRoot string `json:"beacon_block_root"`
						Source          struct {
							Epoch string `json:"epoch"`
							Root  string `json:"root"`
						} `json:"source"`
						Target struct {
							Epoch string `json:"epoch"`
							Root  string `json:"root"`
						} `json:"target"`
					} `json:"data"`
					Signature     string `json:"signature"`
					CommitteeBits string `json:"committee_bits"`
				} `json:"attestations"`
				Deposits       []any `json:"deposits"`
				VoluntaryExits []any `json:"voluntary_exits"`
				SyncAggregate  struct {
					SyncCommitteeBits      string `json:"sync_committee_bits"`
					SyncCommitteeSignature string `json:"sync_committee_signature"`
				} `json:"sync_aggregate"`
				ExecutionPayloadHeader struct {
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
				} `json:"execution_payload_header"`
				BlsToExecutionChanges []any    `json:"bls_to_execution_changes"`
				BlobKzgCommitments    []string `json:"blob_kzg_commitments"`
				ExecutionRequests     struct {
					Deposits       []any `json:"deposits"`
					Withdrawals    []any `json:"withdrawals"`
					Consolidations []any `json:"consolidations"`
				} `json:"execution_requests"`
			} `json:"body"`
		} `json:"message"`
		Signature string `json:"signature"`
	} `json:"data"`
}
