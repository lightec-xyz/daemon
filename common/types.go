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
	Id         string      `json:"id"`
	ProofType  ZkProofType `json:"proofType"`
	Data       interface{} `json:"data"`
	Index      uint64      `json:"index"`
	SIndex     uint64      `json:"sIndex"`
	Hash       string      `json:"hash"`
	Status     ProofStatus `json:"status"`
	Weight     ProofWeight `json:"weight"`
	ProvedNum  int         `json:"-"`
	CreateTime time.Time   `json:"createTime"`
	StartTime  time.Time   `json:"startTime"`
	EndTime    time.Time   `json:"endTime"`
}

func NewZkProofRequest(reqType ZkProofType, data interface{}, fIndex, sIndex uint64, txHash string) *ZkProofRequest {
	return &ZkProofRequest{
		Id:         NewProofId(reqType, fIndex, sIndex, txHash), // todo
		ProofType:  reqType,
		Data:       data,
		Index:      fIndex,
		SIndex:     sIndex,
		Hash:       txHash,
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

func (zk *ZkProofRequest) RequestId() string {
	return zk.Id
}

func (r *ZkProofRequest) String() string {
	return fmt.Sprintf("ZkProofRequest{ProofType:%v,Index:%v,Data:%v}", r.ProofType, r.Index, r.Data)
}

type ZkProofResponse struct {
	Id        string      `json:"id"`
	ProofType ZkProofType `json:"proofType"`
	Status    ProofStatus `json:"status"`
	Proof     []byte      `json:"proof"`
	Witness   []byte      `json:"witness"`
	Index     uint64      `json:"index"`
	SIndex    uint64      `json:"sIndex"`
	Hash      string      `json:"hash"`
}

func (zkp *ZkProofResponse) RespId() string {
	return zkp.Id
}

func (zkResp *ZkProofResponse) String() string {
	return fmt.Sprintf("ProofType:%v Index:%v Proof:%v", zkResp.ProofType, zkResp.Index, zkResp.Proof)
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
