package rpc

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"reflect"
	"time"
)

var _ ProofAPI = (*ProofClient)(nil)

type ProofClient struct {
	*rpc.Client
	timeout time.Duration
}

func (p *ProofClient) ProofInfo(proofId string) (ProofInfo, error) {
	status := ProofInfo{}
	err := p.call(&status, "zkbtc_proofInfo", proofId)
	if err != nil {
		return status, err
	}
	return status, nil
}

func (p *ProofClient) GenZkProof(request ProofRequest) (ProofResponse, error) {
	response := ProofResponse{}
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

func NewProofClient(url string) (ProofAPI, error) {
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
