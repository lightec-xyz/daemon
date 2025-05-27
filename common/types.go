package common

import (
	"bytes"
	"encoding/json"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/store"
	"time"
)

type TaskRequest struct {
	Id        string
	MaxNums   int
	ProofType []ProofType
	Version   int
}

type TaskResponse struct {
	CanGen bool
	Data   string
}

func (r *TaskResponse) ToRequest() (*ProofRequest, error) {
	zkRequest := &ProofRequest{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(r.Data)))
	decoder.UseNumber()
	err := decoder.Decode(zkRequest)
	if err != nil {
		logger.Error("unmarshal request error:%v", err)
		return nil, err
	}
	return zkRequest, nil
}

type SubmitProof struct {
	Responses []*ProofResponse
	Requests  []*ProofRequest
	Status    bool
	WorkerId  string
	Id        string
	Version   int
}

type ProofRequest struct {
	FileKey    store.FileKey `json:"fileKey"`
	ProofType  ProofType     `json:"proofType"`
	Data       interface{}   `json:"data"`
	FIndex     uint64        `json:"index"`
	SIndex     uint64        `json:"sIndex"`
	Prefix     uint64        `json:"prefix"`
	Hash       string        `json:"hash"`
	Status     ProofStatus   `json:"status"`
	Weight     ProofWeight   `json:"weight"`
	CreateTime time.Time     `json:"createTime"`
	StartTime  time.Time     `json:"startTime"`
	EndTime    time.Time     `json:"endTime"`
	BlockTime  uint64        `json:"-"`
	TxIndex    uint32        `json:"-"`
}

func NewProofRequest(reqType ProofType, data interface{}, prefix, fIndex, sIndex uint64, hash string, blockTime uint64, txIndex uint32) *ProofRequest {
	return &ProofRequest{
		FileKey:    GenKey(reqType, prefix, fIndex, sIndex, hash),
		ProofType:  reqType,
		Data:       data,
		FIndex:     fIndex,
		SIndex:     sIndex,
		Hash:       hash,
		Weight:     reqType.Weight(),
		Prefix:     prefix,
		Status:     ProofDefault,
		BlockTime:  blockTime,
		TxIndex:    txIndex,
		CreateTime: time.Now(),
	}
}

func (r *ProofRequest) SetStartTime(t time.Time) {
	r.StartTime = t
}

func (r *ProofRequest) ProofId() string {
	return r.FileKey.String()
}

func (r *ProofRequest) Key() store.FileKey {
	return r.FileKey
}

type ProofResponse struct {
	FileKey       store.FileKey `json:"fileKey"`
	ProofType     ProofType     `json:"proofType"`
	Status        ProofStatus   `json:"status"`
	Proof         []byte        `json:"proof"`
	Witness       []byte        `json:"witness"`
	FIndex        uint64        `json:"index"`
	SIndex        uint64        `json:"sIndex"`
	Prefix        uint64        `json:"prefix"`
	Hash          string        `json:"hash"`
	CreateTime    time.Time     `json:"createTime"`
	ReqCreateTime time.Time     `json:"reqCreateTime"`
}

func NewProofResponse(reqType ProofType, proof []byte, witness []byte, prefix, fIndex, sIndex uint64, hash string, reqCreateTime time.Time) *ProofResponse {
	return &ProofResponse{
		FileKey:       GenKey(reqType, prefix, fIndex, sIndex, hash),
		ProofType:     reqType,
		Status:        ProofSuccess,
		Proof:         proof,
		Witness:       witness,
		FIndex:        fIndex,
		SIndex:        sIndex,
		Hash:          hash,
		Prefix:        prefix,
		CreateTime:    time.Now(),
		ReqCreateTime: reqCreateTime,
	}
}

func (r *ProofResponse) ProofId() string {
	return r.FileKey.String()
}

func (r *ProofResponse) Key() store.FileKey {
	return r.FileKey
}
