package common

import (
	"fmt"
	"time"
)

type TaskRequest struct {
	Id        string
	MaxNums   int
	ProofType []ZkProofType
}

type TaskResponse struct {
	CanGen  bool
	Request *ZkProofRequest
}

type SubmitProof struct {
	Data     []*ZkProofResponse
	WorkerId string
	Id       string
}

// todo
const ZkProofLength = 928

type ZkProof []byte

type ZkProofRequest struct {
	ZkId    string // todo
	ReqType ZkProofType
	Data    interface{}
	Index   uint64
	TxHash  string

	Status     ProofStatus
	Weight     ProofWeight // todo
	CreateTime time.Time
	StartTime  time.Time
	EndTime    time.Time
}

func NewZkProofRequest(reqType ZkProofType, data interface{}, period uint64, txHash string) *ZkProofRequest {
	return &ZkProofRequest{
		ZkId:       NewProofId(reqType, period, txHash), // todo
		ReqType:    reqType,
		Data:       data,
		Index:      period,
		TxHash:     txHash,
		Weight:     reqType.Weight(),
		Status:     ProofDefault,
		CreateTime: time.Now(),
	}
}

func (zk *ZkProofRequest) Id() string {
	return zk.ZkId
}

func (r *ZkProofRequest) String() string {
	return fmt.Sprintf("ZkProofRequest{ReqType:%v,Index:%v,Data:%v}", r.ReqType, r.Index, r.Data)
}

type ZkProofResponse struct {
	ZkProofType ZkProofType
	Status      ProofStatus
	Proof       ZkProof
	Witness     []byte
	Period      uint64
	TxHash      string
}

func (zkp *ZkProofResponse) Id() string {
	return NewProofId(zkp.ZkProofType, zkp.Period, zkp.TxHash)
}

func (zkResp *ZkProofResponse) String() string {
	return fmt.Sprintf("ZkProofType:%v Index:%v Proof:%v", zkResp.ZkProofType, zkResp.Period, zkResp.Proof)
}

func NewProofId(reqType ZkProofType, period uint64, txHash string) string {
	return fmt.Sprintf("%v_%v_%v", reqType.String(), period, txHash)
}
