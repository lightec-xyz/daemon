package rpc

type DaemonInfo struct {
	Version string
	Desc    string
}
type Task struct {
	PoofId uint64
	PTxId  string
	Status int
	Proof  string
}

type ProofRequest struct {
	TxId    string `json:"txId"`
	EthAddr string `json:"ethAddr"`
}

type BtcProofResponse struct {
	TxId   string
	Status int
	Msg    string
	Proof  string
}

type EthProofRequest struct {
	TxId   string
	Status int
}

type EthProofResponse struct {
}

type ProofInfo struct {
}

type ProofStatus struct {
}
