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
	Request ZkProofRequest
}

type SubmitProof struct {
	Data ZkProofResponse
}

// todo
const ZkProofLength = 928

type ZkProof []byte

type ZkProofRequest struct {
	Id      string // todo
	ReqType ZkProofType
	Data    interface{}
	Period  uint64
	TxHash  string

	Status     ProofStatus
	Weight     int // todo
	CreateTime time.Time
	StartTime  time.Time
	EndTime    time.Time
}

func NewZkProofRequest(reqType ZkProofType, data interface{}, period uint64, txHash string) *ZkProofRequest {
	return &ZkProofRequest{
		Id:         fmt.Sprintf("%v_%v_%v", reqType.String(), period, txHash), // todo
		ReqType:    reqType,
		Data:       data,
		Period:     period,
		TxHash:     txHash,
		Status:     ProofDefault,
		CreateTime: time.Now(),
	}
}

func (r *ZkProofRequest) String() string {
	return fmt.Sprintf("ZkProofRequest{ReqType:%v,Period:%v,Data:%v}", r.ReqType, r.Period, r.Data)
}

type ZkProofResponse struct {
	ZkProofType ZkProofType
	Status      ProofStatus
	Proof       ZkProof
	Witness     []byte
	Period      uint64
	TxHash      string
}

func (zkResp *ZkProofResponse) String() string {
	return fmt.Sprintf("ZkProofType:%v Period:%v Proof:%v", zkResp.ZkProofType, zkResp.Period, zkResp.Proof)
}

type ProofTask struct {
	Id        string
	Weight    int // todo
	StartTime time.Time
	EndTime   time.Time
	Task      ZkProofRequest
}
