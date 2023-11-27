package rpc

type NodeAPI interface {
	Version() (*DaemonInfo, error)
}

type ProofAPI interface {
	Info() (ProofInfo, error)
	GenBtcProof(request ProofRequest) (BtcProofResponse, error)
	GenEthProof(request EthProofRequest) (EthProofResponse, error)
	ProofStatus(proofId string) (ProofStatus, error)
}
