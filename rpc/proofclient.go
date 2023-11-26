package rpc

import (
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lightec-xyz/daemon/proof"
)

var _ proof.API = (*ProofClient)(nil)

type ProofClient struct {
	*rpc.Client
}

func (p *ProofClient) ProofStatus(proofId string) (proof.ProofStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) GenBtcProof(request proof.BtcProofRequest) (proof.BtcProofResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) GenEthProof(request proof.EthProofRequest) (proof.EthProofResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) Info() (proof.ProofInfo, error) {
	//TODO implement me
	panic("implement me")
}
