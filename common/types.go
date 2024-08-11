package common

import (
	"fmt"
	"strings"
	"time"
)

type TaskRequest struct {
	Id        string
	MaxNums   int
	ProofType []ZkProofType
	Version   string
}

type TaskResponse struct {
	CanGen  bool
	Request *ZkProofRequest
}

type SubmitProof struct {
	Data     []*ZkProofResponse
	WorkerId string
	Id       string
	Version  string
}

type ZkProofRequest struct {
	ZkId       string // todo
	ReqType    ZkProofType
	Data       interface{}
	Index      uint64
	SIndex     uint64
	TxHash     string
	Status     ProofStatus
	Weight     ProofWeight // todo
	ProvedNum  int
	CreateTime time.Time
	StartTime  time.Time
	EndTime    time.Time
}

func NewZkProofRequest(reqType ZkProofType, data interface{}, fIndex, sIndex uint64, txHash string) *ZkProofRequest {
	return &ZkProofRequest{
		ZkId:       NewProofId(reqType, fIndex, sIndex, txHash), // todo
		ReqType:    reqType,
		Data:       data,
		Index:      fIndex,
		SIndex:     sIndex,
		TxHash:     txHash,
		Weight:     reqType.Weight(),
		Status:     ProofDefault,
		CreateTime: time.Now(),
	}
}

func (zk *ZkProofRequest) SetStartTime(t time.Time) {
	zk.StartTime = t
}

func (zk *ZkProofRequest) SetEndTime(t time.Time) {
	zk.EndTime = t
}

func (zk *ZkProofRequest) Id() string {
	return zk.ZkId
}

func (r *ZkProofRequest) String() string {
	return fmt.Sprintf("ZkProofRequest{ReqType:%v,Index:%v,Data:%v}", r.ReqType, r.Index, r.Data)
}

type ZkProofResponse struct {
	RespId      string
	ZkProofType ZkProofType
	Status      ProofStatus
	Proof       []byte
	Witness     []byte
	Index       uint64
	End         uint64
	TxHash      string
}

func (zkp *ZkProofResponse) Id() string {
	return zkp.RespId
}

func (zkResp *ZkProofResponse) String() string {
	return fmt.Sprintf("ZkProofType:%v Index:%v Proof:%v", zkResp.ZkProofType, zkResp.Index, zkResp.Proof)
}

func NewProofId(reqType ZkProofType, fIndex, sIndex uint64, hash string) string {
	/* example
	1. type_hash
	2. type_index
	3. type_index_end
	4. type_index_end_hash
	*/

	id := reqType.String()
	if fIndex != 0 {
		id = fmt.Sprintf("%v_%v", id, fIndex)
	}
	if sIndex != 0 {
		id = fmt.Sprintf("%v_%v", id, sIndex)
	}
	if hash != "" {
		id = fmt.Sprintf("%v_%v", id, hash)
	}
	return id
}

func ParseProofId(id string) (ZkProofType, uint64, string, error) {
	// todo
	if len(id) == 0 {
		return ZkProofType(0), 0, "", fmt.Errorf("proof id is empty")
	}
	params := strings.Split(id, "_")
	if len(params) == 2 {

	} else if len(params) == 3 {

	} else {
		return ZkProofType(0), 0, "", fmt.Errorf("proof id format error: %v", id)
	}

	var reqType ZkProofType
	var period uint64
	var txHash string
	_, err := fmt.Sscanf(id, "%v_%v", &reqType, &period)
	if err != nil {
		return ZkProofType(0), 0, "", err
	}
	return reqType, period, txHash, nil
}
