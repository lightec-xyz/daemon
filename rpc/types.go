package rpc

type JsonRpcReq struct {
	Id     float64     `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type JsonRpcResp struct {
	Id     float64     `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}
