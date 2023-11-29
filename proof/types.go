package proof

type Config struct {
	RpcBind      string `json:"rpcbind"`
	RpcPort      string `json:"rpcport"`
	ParallelNums int    `json:"parallelNums"`
	Network      string `json:"network"`
	Datadir      string `json:"datadir"`
}
