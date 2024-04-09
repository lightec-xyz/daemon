package common

import (
	"fmt"
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

type ZkProofType int

const (
	DepositTxType ZkProofType = iota + 1
	RedeemTxType
	TxInEth2
	VerifyTxType
	SyncComGenesisType
	SyncComUnitType
	SyncComRecursiveType
)

func (zkpr *ZkProofType) String() string {
	switch *zkpr {
	case DepositTxType:
		return "DepositTxType"
	case RedeemTxType:
		return "RedeemTxType"
	case VerifyTxType:
		return "VerifyTxType"
	case SyncComGenesisType:
		return "SyncComGenesisType"
	case SyncComUnitType:
		return "SyncComUnitType"
	case SyncComRecursiveType:
		return "SyncComRecursiveType"
	default:
		return ""
	}
}

type ZkProofRequest struct {
	ReqType ZkProofType
	Data    interface{}
	Period  uint64
	TxHash  string
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
