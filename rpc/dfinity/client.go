package dfinity

import (
	"fmt"
	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/candid/idl"
	"github.com/aviate-labs/agent-go/ic/wallet"
	"github.com/aviate-labs/agent-go/identity"
	"github.com/aviate-labs/agent-go/principal"
	"math/big"
	"net/url"
	"time"
)

type Client struct {
	agent       *agent.Agent
	walletAgent *wallet.Agent
	canisterId  principal.Principal
	timeout     time.Duration
}

func (c *Client) PublicKey() (string, error) {
	var result string
	err := c.call(c.canisterId, "public_key", []any{}, []any{&result})
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) WalletBalance() (uint64, error) {
	balance, err := c.walletAgent.WalletBalance()
	if err != nil {
		return 0, err
	}
	return balance.Amount, nil
}

func (c *Client) DummyAddress() (string, error) {
	var result string
	err := c.call(c.canisterId, "dummy_address", []any{}, []any{&result})
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) EthDecoderCanister() (string, error) {
	var result string
	err := c.call(c.canisterId, "eth_decoder_canister", []any{}, []any{&result})
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) PlonkVerifierCanister() (string, error) {
	var result string
	err := c.call(c.canisterId, "plonk_verifier_canister", []any{}, []any{&result})
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) BtcTxSign(currentScRoot, ethTxHash, btcTxId, proof, minerReward string, sigHashes []string) (*TxSignature, error) {
	signature := TxSignature{}
	err := c.call(c.canisterId, "verify_and_sign_free", []any{currentScRoot, ethTxHash, btcTxId, minerReward, sigHashes, proof}, []any{&signature.Signed, &signature.Signature})
	if err != nil {
		return nil, err
	}
	return &signature, nil
}

func (c *Client) BlockSignature() (*BlockSignature, error) {
	result := BlockSignature{}
	err := c.call(c.canisterId, "block_height", []any{}, []any{&result.Hash, &result.Height, &result.Signature})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) BlockSignatureWithCycle() (*BlockSignature, error) {
	result := BlockSignature{}
	err := c.walletCall(28_000_000_000, "block_height", []any{}, []any{&result.Hash, &result.Height, &result.Signature})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) BtcTxSignWithCycle(currentScRoot, ethTxHash, btcTxId, proof, minerReward string, sigHashes []string) (*TxSignature, error) {
	signature := TxSignature{}
	err := c.walletCall(28_000_000_000, "verify_and_sign", []any{currentScRoot, ethTxHash, btcTxId, minerReward, sigHashes, proof}, []any{&signature.Signed, &signature.Signature})
	if err != nil {
		return nil, err
	}
	return &signature, nil
}

func (c *Client) WalletBalance128() (*big.Int, error) {
	cycles, err := c.walletAgent.WalletBalance128()
	if err != nil {
		return nil, err
	}
	return cycles.Amount.BigInt(), nil
}

func (c *Client) call(canisterID principal.Principal, method string, args []any, rets []any) error {
	return c.agent.Call(canisterID, method, args, rets)
}

func (c *Client) walletCall(cycles uint64, method string, args []any, rets []any) error {
	var input []byte
	var err error
	if len(args) != 0 {
		input, err = idl.Marshal(args)
		if err != nil {
			return err
		}
	}
	walletCallArg := WalletCallArg{
		Canister:   c.canisterId,
		MethodName: method,
		Args:       input,
		Cycles:     cycles,
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

func NewClient(canId string) (*Client, error) {
	timeout := 60 * time.Second
	config := agent.Config{
		PollTimeout: timeout,
	}
	canisterId, err := principal.Decode(canId)
	if err != nil {
		return nil, err
	}
	icpAgent, err := agent.New(config)
	if err != nil {
		return nil, err
	}
	walletAgent, err := wallet.NewAgent(canisterId, config)
	if err != nil {
		return nil, err
	}
	return &Client{
		agent:       icpAgent,
		walletAgent: walletAgent,
		canisterId:  canisterId,
		timeout:     timeout,
	}, nil
}

func NewClientWithIdentity(canId, walletId, endpoint string, identity identity.Identity) (*Client, error) {
	timeout := 60 * time.Second
	host, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	config := agent.Config{
		Identity:                       identity,
		ClientConfig:                   &agent.ClientConfig{Host: host},
		FetchRootKey:                   true,
		PollTimeout:                    timeout,
		DisableSignedQueryVerification: false,
	}
	canisterId, err := principal.Decode(canId)
	if err != nil {
		return nil, err
	}
	icpAgent, err := agent.New(config)
	if err != nil {
		return nil, err
	}
	walletCanisterId, err := principal.Decode(walletId)
	if err != nil {
		return nil, err
	}
	walletAgent, err := wallet.NewAgent(walletCanisterId, config)
	if err != nil {
		return nil, err
	}
	return &Client{
		agent:       icpAgent,
		walletAgent: walletAgent,
		canisterId:  canisterId,
		timeout:     timeout,
	}, nil
}
