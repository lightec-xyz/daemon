package rpc

type NodeAPI interface {
	Version() (NodeInfo, error)
}

type ProofAPI interface {
	Info() (ProofInfo, error)
	GenZkProof(request ProofRequest) (ProofResponse, error)
	GenEthProof(request EthProofRequest) (EthProofResponse, error)
	ProofStatus(proofId string) (ProofStatus, error)
}
