package proof

type API interface {
	Info() (ProofInfo, error)
	GenBtcProof(request BtcProofRequest) (BtcProofResponse, error)
	GenEthProof(request EthProofRequest) (EthProofResponse, error)
	ProofStatus(proofId string) (ProofStatus, error)
}
