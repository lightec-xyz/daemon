package proof

type Task struct {
	PoofId uint64
	PTxId  string
	Status int
	Proof  string
}

type BtcProofRequest struct {
	TxId string
}

type BtcProofResponse struct {
}

type EthProofRequest struct {
	TxId string
}

type EthProofResponse struct {
}

type ProofInfo struct {
}
