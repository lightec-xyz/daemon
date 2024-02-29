package rpc

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"reflect"
	"time"
)

var _ ProofAPI = (*ProofClient)(nil)
var _ ISyncCommitteeProof = (*SyncCommitteeProofClient)(nil)

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

// TODO(keep), add sync committee proof client
type SyncCommitteeProofClient struct {
	*rpc.Client
	timeout time.Duration
}

func (p *SyncCommitteeProofClient) SyncCommitteeProofInfo(period uint64, proofType SyncCommitteeProofType) (SyncCommitteeProofInfo, error) {
	status := SyncCommitteeProofInfo{}
	err := p.call(&status, "sync_committee_proofInfo", period, proofType)
	if err != nil {
		return status, err
	}
	return status, nil
}

func (p *SyncCommitteeProofClient) GenGenesisSyncCommitteeProof(request GenesisSyncCommitteeProofRequest) (SyncCommitteeProofResponse, error) {
	response := SyncCommitteeProofResponse{}
	err := p.call(&response, "sync_committee_genGenesisProof", request)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (p *SyncCommitteeProofClient) GenUnitSyncCommitteeProof(request UnitSyncCommitteeProofRequest) (SyncCommitteeProofResponse, error) {
	response := SyncCommitteeProofResponse{}
	err := p.call(&response, "sync_committee_genUintProof", request)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (p *SyncCommitteeProofClient) GenRecursiveSyncCommitteeProof(request RecursiveSyncCommitteeProofRequest) (SyncCommitteeProofResponse, error) {
	response := SyncCommitteeProofResponse{}
	err := p.call(&response, "sync_committee_genRecursiveProof", request)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (p *SyncCommitteeProofClient) call(result interface{}, method string, args ...interface{}) error {
	vi := reflect.ValueOf(result)
	if vi.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be pointer")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), p.timeout)
	defer cancelFunc()
	return p.CallContext(ctx, result, method, args...)
}

func NewSyncCommitteeProofClient(url string) (ISyncCommitteeProof, error) {
	client, err := rpc.DialHTTP(url)
	if err != nil {
		return nil, err
	}
	return &SyncCommitteeProofClient{
		Client:  client,
		timeout: 30 * time.Minute,
	}, nil
}

func NewWsSyncCommitteeProofClient(url string) (ISyncCommitteeProof, error) {
	client, err := rpc.DialWebsocket(context.Background(), url, "")
	if err != nil {
		return nil, err
	}
	return &SyncCommitteeProofClient{
		Client:  client,
		timeout: 3 * time.Hour,
	}, nil
}
