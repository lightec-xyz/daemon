package rpc

type NodeInfo struct {
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
	TxId   string `json:"txId"`
	PType  string `json:"type"`
	Proof  string `json:"proof"`
	ToAddr string `json:"toAddr"`
	Amount string `json:"amount"`
	Msg    string `json:"msg"`
}

type ProofResponse struct {
	TxId   string `json:"txId"`
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	PType  string `json:"type"`
	Proof  string `json:"proof"`
}

type EthProofRequest struct {
	TxId    string `json:"txId"`
	EthAddr string `json:"ethAddr"`
	Proof   string `json:"proof"`
	Msg     string `json:"msg"`
}

type EthProofResponse struct {
	TxId   string
	Status int
	Msg    string
	Proof  string
}

type ProofInfo struct {
	Version string
}

type ProofStatus struct {
	State int    `json:"state"`
	Msg   string `json:"msg"`
}
