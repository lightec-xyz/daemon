package dfinity

import (
	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/ic/ic"
	"github.com/aviate-labs/agent-go/principal"
	"time"
)

type Client struct {
	agent      *agent.Agent
	icAgent    *ic.Agent
	canisterId principal.Principal
}

func (c *Client) PublicKey() (string, error) {
	var result string
	err := c.call(c.canisterId, "public_key", []any{}, []any{&result})
	if err != nil {
		return "", err
	}
	return result, nil
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

func (c *Client) VerifyAndSign(txRaw, receiptRaw, proof string) (*Signature, error) {
	signature := Signature{}
	err := c.call(c.canisterId, "verify_and_sign", []any{txRaw, receiptRaw, proof}, []any{&signature.Signed, &signature.Signature})
	if err != nil {
		return nil, err
	}
	return &signature, nil
}
func (c *Client) BlockHeight() (*BlockHeight, error) {
	height := BlockHeight{}
	err := c.call(c.canisterId, "block_height", []any{}, []any{&height.Hash, &height.Height, &height.Signature})
	if err != nil {
		return nil, err
	}
	return &height, nil
}

func (c *Client) call(canisterID principal.Principal, method string, args []any, rets []any) error {
	return c.agent.Call(canisterID, method, args, rets)
}
func (c *Client) query(canisterID principal.Principal, method string, args []any, rets []any) error {
	return c.agent.Query(canisterID, method, args, rets)
}

func NewClient(canId string) (*Client, error) {
	config := agent.Config{
		PollTimeout: 45 * time.Second,
	}
	canisterId, err := principal.Decode(canId)
	if err != nil {
		return nil, err
	}
	agent, err := agent.New(config)
	if err != nil {
		return nil, err
	}
	icAgent, err := ic.NewAgent(canisterId, config)
	if err != nil {
		return nil, err
	}
	return &Client{
		agent:      agent,
		icAgent:    icAgent,
		canisterId: canisterId,
	}, nil
}
