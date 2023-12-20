package rpc

type NodeAPI interface {
	Version() (NodeInfo, error)
	AddWorker(endpoint string, max int) (string, error)
	ProofInfo(proofId string) (ProofInfo, error)
	Transaction(txHash string) (Transaction, error)
	Stop() error
}

// ProofAPI api between node and proof
type ProofAPI interface {
	GenZkProof(request ProofRequest) (ProofResponse, error)
	ProofInfo(proofId string) (ProofInfo, error)
}
