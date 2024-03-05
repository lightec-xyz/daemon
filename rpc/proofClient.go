package rpc

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"reflect"
	"time"
)

var _ IProof = (*ProofClient)(nil)

type ProofClient struct {
	*rpc.Client
	timeout time.Duration
}

func (p *ProofClient) GenDepositProof(req DepositRequest) (DepositResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) GenRedeemProof(req RedeemRequest) (RedeemResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) GenVerifyProof(req VerifyRequest) (VerifyResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) GenSyncCommGenesisProof(req SyncCommGenesisRequest) (SyncCommGenesisResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) GenSyncCommitUnitProof(req SyncCommUnitsRequest) (SyncCommUnitsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) GenSyncCommRecursiveProof(req SyncCommRecursiveRequest) (SyncCommRecursiveResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) AddReqNum() {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) DelReqNum() {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) MaxNums() int {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) CurrentNums() int {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) ProofInfo(proofId string) (ProofInfo, error) {
	status := ProofInfo{}
	err := p.call(&status, "zkbtc_proofInfo", proofId)
	if err != nil {
		return status, err
	}
	return status, nil
}

func (p *ProofClient) GenZkProof(request DepositRequest) (DepositResponse, error) {
	response := DepositResponse{}
	err := p.call(&response, "zkbtc_genZkProof", request)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (p *ProofClient) call(result interface{}, method string, args ...interface{}) error {
	vi := reflect.ValueOf(result)
	if vi.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be pointer")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), p.timeout)
	defer cancelFunc()
	return p.CallContext(ctx, result, method, args...)
}

func NewProofClient(url string) (IProof, error) {
	client, err := rpc.DialHTTP(url)
	if err != nil {
		return nil, err
	}
	return &ProofClient{
		Client:  client,
		timeout: 15 * time.Second,
	}, nil
}

func NewWsProofClient(url string) (*ProofClient, error) {
	client, err := rpc.DialWebsocket(context.Background(), url, "")
	if err != nil {
		return nil, err
	}
	return &ProofClient{
		Client:  client,
		timeout: 3 * time.Hour,
	}, nil
}
