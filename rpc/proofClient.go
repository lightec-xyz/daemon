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

func (p *ProofClient) TxInEth2Prove(req *TxInEth2ProveReq) (*TxInEth2ProveResp, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) TxBlockIsParentOfCheckPointProve(req *TxBlockIsParentOfCheckPointProveReq) (*TxBlockIsParentOfCheckPointResp, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) CheckPointFinalityProve(req *CheckPointFinalityProveReq) (*CheckPointFinalityProveResp, error) {
	//TODO implement me
	panic("implement me")
}

func (p *ProofClient) GenDepositProof(req DepositRequest) (DepositResponse, error) {
	response := DepositResponse{}
	err := p.call(&response, "zkbtc_genDepositProof", req)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (p *ProofClient) GenRedeemProof(req RedeemRequest) (RedeemResponse, error) {
	response := RedeemResponse{}
	err := p.call(&response, "zkbtc_genRedeemProof", req)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (p *ProofClient) GenVerifyProof(req VerifyRequest) (VerifyResponse, error) {
	response := VerifyResponse{}
	err := p.call(&response, "zkbtc_genVerifyProof", req)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (p *ProofClient) GenSyncCommGenesisProof(req SyncCommGenesisRequest) (SyncCommGenesisResponse, error) {
	response := SyncCommGenesisResponse{}
	err := p.call(&response, "zkbtc_genSyncCommGenesisProof", req)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (p *ProofClient) GenSyncCommitUnitProof(req SyncCommUnitsRequest) (SyncCommUnitsResponse, error) {
	response := SyncCommUnitsResponse{}
	err := p.call(&response, "zkbtc_genSyncCommitUnitProof", req)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (p *ProofClient) GenSyncCommRecursiveProof(req SyncCommRecursiveRequest) (SyncCommRecursiveResponse, error) {
	response := SyncCommRecursiveResponse{}
	err := p.call(&response, "zkbtc_genSyncCommRecursiveProof", req)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (p *ProofClient) ProofInfo(proofId string) (ProofInfo, error) {
	status := ProofInfo{}
	err := p.call(&status, "zkbtc_proofInfo", proofId)
	if err != nil {
		return status, err
	}
	return status, nil
}
func (p *ProofClient) MaxNums() (int, error) {
	var result int
	err := p.call(&result, "zkbtc_maxNums")
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (p *ProofClient) CurrentNums() (int, error) {
	var result int
	err := p.call(&result, "zkbtc_currentNums")
	if err != nil {
		return 0, err
	}
	return result, nil
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
