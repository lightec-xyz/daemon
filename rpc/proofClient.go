package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/rpc/ws"
	"reflect"
	"time"
)

var _ IProof = (*ProofClient)(nil)

type ProofClient struct {
	*rpc.Client
	timeout time.Duration
	conn    *ws.Conn
	custom  bool
}

func (p *ProofClient) BtcDuperRecursiveProve(req *BtcDuperRecursiveRequest) (*ProofResponse, error) {
	var result ProofResponse
	err := p.call(&result, "zkbtc_btcDuperRecursiveProve", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProofClient) BtcDepthRecursiveProve(req *BtcDepthRequest) (*ProofResponse, error) {
	var result ProofResponse
	err := p.call(&result, "zkbtc_btcDepthRecursiveProve", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProofClient) BtcChainProve(req *BtcChainRequest) (*ProofResponse, error) {
	var result ProofResponse
	err := p.call(&result, "zkbtc_btcChainProve", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProofClient) BtcDepositProve(req *BtcDepositRequest) (*ProofResponse, error) {
	var result ProofResponse
	err := p.call(&result, "zkbtc_btcDepositProve", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProofClient) BtcChangeProve(req *BtcChangeRequest) (*ProofResponse, error) {
	var result ProofResponse
	err := p.call(&result, "zkbtc_btcChangeProve", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProofClient) BtcBaseProve(req *BtcBaseRequest) (*ProofResponse, error) {
	var result ProofResponse
	err := p.call(&result, "zkbtc_btcBaseProve", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProofClient) BtcMiddleProve(req *BtcMiddleRequest) (*ProofResponse, error) {
	var result ProofResponse
	err := p.call(&result, "zkbtc_btcMiddleProve", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProofClient) BtcUpperProve(req *BtcUpperRequest) (*ProofResponse, error) {
	var result ProofResponse
	err := p.call(&result, "zkbtc_btcUpProve", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProofClient) SupportProofType() []common.ZkProofType {
	var result []common.ZkProofType
	err := p.call(&result, "zkbtc_supportProofType")
	if err != nil {
		return nil
	}
	return result
}

func (p *ProofClient) BtcBulkProve(data *BtcBulkRequest) (*BtcBulkResponse, error) {
	var result BtcBulkResponse
	err := p.call(&result, "zkbtc_btcBulkProve", data)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProofClient) BtcPackedRequest(data *BtcPackedRequest) (*BtcPackResponse, error) {
	var result BtcPackResponse
	err := p.call(&result, "zkbtc_btcPackedRequest", data)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProofClient) TxInEth2Prove(req *TxInEth2ProveRequest) (*TxInEth2ProveResponse, error) {
	var result TxInEth2ProveResponse
	err := p.call(&result, "zkbtc_txInEth2Prove", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProofClient) BlockHeaderProve(req *BlockHeaderRequest) (*BlockHeaderResponse, error) {
	var result BlockHeaderResponse
	err := p.call(&result, "zkbtc_blockHeaderProve", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProofClient) BlockHeaderFinalityProve(req *BlockHeaderFinalityRequest) (*BlockHeaderFinalityResponse, error) {
	var result BlockHeaderFinalityResponse
	err := p.call(&result, "zkbtc_blockHeaderFinalityProve", req)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProofClient) GenRedeemProof(req *RedeemRequest) (*RedeemResponse, error) {
	response := RedeemResponse{}
	err := p.call(&response, "zkbtc_genRedeemProof", req)
	if err != nil {
		return nil, err
	}
	return &response, nil
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
	if p.custom {
		return p.customCall(result, method, args...)
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), p.timeout)
	defer cancelFunc()
	return p.CallContext(ctx, result, method, args...)
}

func (p *ProofClient) customCall(result interface{}, method string, args ...interface{}) error {
	vi := reflect.ValueOf(result)
	if vi.Kind() != reflect.Ptr {
		return fmt.Errorf("result must be pointer")
	}
	timeout, err := getTimeout(method)
	if err != nil {
		return err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	response, err := p.conn.Call(ctx, method, args...)
	if err != nil {
		return err
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProofClient) Close() error {
	if p.conn != nil {
		err := p.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
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

func NewCustomWsProofClient(conn *ws.Conn) (*ProofClient, error) {
	return &ProofClient{
		conn:    conn,
		timeout: 60 * time.Second,
		custom:  true,
	}, nil
}

func getTimeout(method string) (time.Duration, error) {
	// todo
	return 90 * time.Minute, nil
}
