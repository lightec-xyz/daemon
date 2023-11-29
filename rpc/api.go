package rpc

type NodeAPI interface {
	Version() (NodeInfo, error)
}

// ProofAPI api between node and proof
type ProofAPI interface {
	Info() (ProofInfo, error)
	GenZkProof(request ProofRequest) (ProofResponse, error)
	ProofStatus(proofId string) (ProofStatus, error)
}
