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
}

func NewProofClient(url string) (*ProofClient, error) {
	client, err := rpc.DialHTTP(url)
	if err != nil {
		return nil, err
	}
	return &ProofClient{
		client,
	}, nil
}

func NewWsProofClient(url string) (*ProofClient, error) {
	client, err := rpc.DialWebsocket(context.Background(), url, "")
	if err != nil {
		return nil, err
	}
	return &ProofClient{
		client,
	}, nil
}

func (p *ProofClient) ProofStatus(proofId string) (ProofStatus, error) {
	status := ProofStatus{}
	err := p.call(&status, "proof_status", proofId)
	if err != nil {
		return status, err
	}
	return status, nil
}

func (p *ProofClient) GenZkProof(request ProofRequest) (ProofResponse, error) {
	//todo
	response := ProofResponse{
		TxId:   request.TxId,
		Status: 0,
		Msg:    "ok",
		Proof:  "test proof",
	}
	err := p.call(&response, "zkbtc_genBtcProof", request)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (p *ProofClient) Info() (ProofInfo, error) {
	info := ProofInfo{}
	err := p.call(&info, "proof_info")
	if err != nil {
		return info, err
	}
	return info, nil
}

func (p *ProofClient) call(result interface{}, method string, args ...interface{}) error {
	vi := reflect.ValueOf(result)
	if vi.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be pointer")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelFunc()
	return p.CallContext(ctx, result, method, args...)
}
