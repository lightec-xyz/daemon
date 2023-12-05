package rpc

type NodeAPI interface {
	Version() (NodeInfo, error)
}

// ProofAPI api between node and proof
type ProofAPI interface {
	GenZkProof(request ProofRequest) (ProofResponse, error)
	ProofInfo(proofId string) (ProofResponse, error)
}
