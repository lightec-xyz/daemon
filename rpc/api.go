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
	GenDepositProof(req DepositRequest) (DepositResponse, error)
	GenRedeemProof(req RedeemRequest) (RedeemResponse, error)
	GenVerifyProof(req VerifyRequest) (VerifyResponse, error)
	GenSyncCommGenesisProof(req SyncCommGenesisRequest) (SyncCommGenesisResponse, error)
	GenSyncCommitUnitProof(req SyncCommUnitsRequest) (SyncCommUnitsResponse, error)
	GenSyncCommRecursiveProof(req SyncCommRecursiveRequest) (SyncCommRecursiveResponse, error)
	AddReqNum()
	DelReqNum()
	MaxNums() int
	CurrentNums() int
}
