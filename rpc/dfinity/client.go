package dfinity

import (
	"fmt"
	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/candid/idl"
	"github.com/aviate-labs/agent-go/ic/wallet"
	"github.com/aviate-labs/agent-go/identity"
	"github.com/aviate-labs/agent-go/principal"
	"math/big"
	"time"
)

type Client struct {
	agent           *agent.Agent
	walletAgent     *wallet.Agent
	txCanisterId    principal.Principal
	blockCanisterId principal.Principal
	timeout         time.Duration
}

func (c *Client) TxPublicKey() (string, error) {
	var result string
	err := c.call(c.txCanisterId, "public_key", []any{}, []any{&result})
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) BlockPublicKey() (string, error) {
	var result string
	err := c.call(c.blockCanisterId, "public_key", []any{}, []any{&result})
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) IcpBalance() (uint64, error) {
	if c.walletAgent == nil {
		return 0, fmt.Errorf("walletAgent is nil")
	}
	balance, err := c.walletAgent.WalletBalance()
	if err != nil {
		return 0, err
	}
	return balance.Amount, nil
}

func (c *Client) DummyAddress() (string, error) {
	var result string
	err := c.call(c.txCanisterId, "dummy_address", []any{}, []any{&result})
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) BtcTxSign(currentScRoot, ethTxHash, btcTxId, proof, minerReward string, sigHashes []string) (*TxSignature, error) {
	signature := TxSignature{}
	err := c.call(c.txCanisterId, "verify_and_sign_free", []any{currentScRoot, ethTxHash, btcTxId, minerReward, sigHashes, proof}, []any{&signature.Signed, &signature.Signature})
	if err != nil {
		return nil, err
	}
	return &signature, nil
}
func (c *Client) BtcTxSignWithCycle(currentScRoot, ethTxHash, btcTxId, proof, minerReward string, sigHashes []string) (*TxSignature, error) {
	signature := TxSignature{}
	err := c.walletCall(c.txCanisterId, 50_000_000_000, "verify_and_sign", []any{currentScRoot, ethTxHash, btcTxId, minerReward, sigHashes, proof}, []any{&signature.Signed, &signature.Signature})
	if err != nil {
		return nil, err
	}
	return &signature, nil
}

func (c *Client) BlockSignature() (*BlockSignature, error) {
	result := BlockSignature{}
	err := c.call(c.blockCanisterId, "block_height_free", []any{}, []any{&result.Height, &result.Hash, &result.Signature})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) BlockSignatureWithCycle() (*BlockSignature, error) {
	result := BlockSignature{}
	err := c.walletCall(c.blockCanisterId, 50_000_000_000, "block_height", []any{}, []any{&result.Height, &result.Hash, &result.Signature})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) WalletInfo() {

}

func (c *Client) CyclesBalance() (*big.Int, error) {
	if c.walletAgent == nil {
		return nil, fmt.Errorf("walletAgent is nil")
	}
	cycles, err := c.walletAgent.WalletBalance128()
	if err != nil {
		return nil, err
	}
	return cycles.Amount.BigInt(), nil
}

func (c *Client) call(canisterID principal.Principal, method string, args []any, rets []any) error {
	return c.agent.Call(canisterID, method, args, rets)
}

func (c *Client) walletCall(destCanId principal.Principal, cycles uint64, method string, args []any, rets []any) error {
	var input []byte
	var err error
	if len(args) != 0 {
		input, err = idl.Marshal(args)
		if err != nil {
			return err
		}
	}
	walletCallArg := WalletCallArg{
		Canister:   destCanId,
		MethodName: method,
		Args:       input,
		Cycles:     cycles,
	}
	if c.walletAgent == nil {
		return fmt.Errorf("walletAgent is nil")
	}
	res, err := c.walletAgent.WalletCall(walletCallArg)
	if err != nil {
		return err
	}
	if res.Err != nil {
		return fmt.Errorf("%v", *res.Err)
	}
	err = idl.Unmarshal(res.Ok.Return, rets)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) query(canisterID principal.Principal, method string, args []any, rets []any) error {
	return c.agent.Query(canisterID, method, args, rets)
}

func NewClient(opt *Options) (*Client, error) {
	if opt == nil {
		return nil, fmt.Errorf("opt is nil")
	}
	timeout := 2 * time.Minute
	config := agent.Config{
		PollTimeout: timeout,
	}
	txCanisterId, err := principal.Decode(opt.TxCanisterId)
	if err != nil {
		return nil, err
	}
	agent, err := agent.New(config)
	if err != nil {
		return nil, err
	}
	blockCanisterId, err := principal.Decode(opt.BlockCanisterId)
	if err != nil {
		return nil, err
	}
	var walletClient *wallet.Agent
	if opt.identity != nil && opt.WalletCanisterId != "" {
		walletClient, err = NewWalletClient(opt.WalletCanisterId, opt.identity)
		if err != nil {
			return nil, err
		}
	}
	return &Client{
		agent:           agent,
		walletAgent:     walletClient,
		txCanisterId:    txCanisterId,
		blockCanisterId: blockCanisterId,
		timeout:         timeout,
	}, nil
}

func NewWalletClient(walletCanId string, identity identity.Identity) (*wallet.Agent, error) {
	config := agent.Config{
		Identity: identity,
		//ClientConfig:                   &agent.ClientConfig{Host: host},
		FetchRootKey:                   true,
		PollTimeout:                    2 * time.Minute,
		DisableSignedQueryVerification: false,
	}
	walletCanisterId, err := principal.Decode(walletCanId)
	if err != nil {
		return nil, err
	}
	walletAgent, err := wallet.NewAgent(walletCanisterId, config)
	if err != nil {
		return nil, err
	}
	return walletAgent, nil
}
