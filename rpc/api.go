package rpc

type INode interface {
	Version() (NodeInfo, error)
	AddWorker(endpoint string, max int) (string, error)
	ProofInfo(txId string) (ProofInfo, error)
	Transaction(txHash string) (Transaction, error)
	TransactionsByHeight(height uint64, network string) ([]string, error)
	Transactions(txId []string) ([]Transaction, error)
	Stop() error
}

// IProof api between node and proof
type IProof interface {
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
