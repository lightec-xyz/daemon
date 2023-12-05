package proof

type Config struct {
	RpcBind      string `json:"rpcbind"`
	RpcPort      string `json:"rpcport"`
	ParallelNums int    `json:"parallelNums"`
	Network      string `json:"network"`
	Datadir      string `json:"datadir"`
}

func localDevConfig() Config {
	return Config{
		RpcBind:      "127.0.0.1",
		RpcPort:      "88888",
		ParallelNums: 3,
		Network:      "testnet",
		Datadir:      "/Users/red/.daemon",
	}
}
