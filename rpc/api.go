package rpc

type NodeAPI interface {
	Version() (NodeInfo, error)
	AddWorker(endpoint string, max int) (string, error)
	ProofInfo(txId string) (ProofInfo, error)
	Transaction(txHash string) (Transaction, error)
	TransactionsByHeight(height uint64, network string) ([]string, error)
	Transactions(txId []string) ([]Transaction, error)
	Stop() error
}

// ProofAPI api between node and proof
type ProofAPI interface {
	GenZkProof(request ProofRequest) (ProofResponse, error)
	ProofInfo(proofId string) (ProofInfo, error)
}

// TODO(keep), sync committee proof generator interface
type ISyncCommitteeProof interface {
	GenGenesisSyncCommitteeProof(request GenesisSyncCommitteeProofRequest) (SyncCommitteeProofResponse, error)
	GenUnitSyncCommitteeProof(request UnitSyncCommitteeProofRequest) (SyncCommitteeProofResponse, error)
	GenRecursiveSyncCommitteeProof(request RecursiveSyncCommitteeProofRequest) (SyncCommitteeProofResponse, error)
	SyncCommitteeProofInfo(period uint64, proofType SyncCommitteeProofType) (SyncCommitteeProofInfo, error)
}
