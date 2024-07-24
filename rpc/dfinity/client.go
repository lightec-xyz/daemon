package dfinity

import (
	"fmt"
	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/ic/ic"
	"github.com/aviate-labs/agent-go/principal"
	"time"
)

type Client struct {
	agent   *agent.Agent
	icAgent *ic.Agent
}

func (c *Client) PublicKey(canisterId string) (interface{}, error) {
	publicKey := make(map[string]interface{})
	canId, err := principal.Decode(canisterId)
	if err != nil {
		return nil, err
	}
	err = c.call(canId, "public_key", []any{}, []any{&publicKey})
	if err != nil {
		return nil, err
	}
	return &publicKey, nil
}

func (c *Client) Sign(canisterId, msg string) (interface{}, error) {
	canId, err := principal.Decode(canisterId)
	if err != nil {
		return nil, err
	}
	signature := make(map[string]interface{})
	err = c.call(canId, "sign", []any{msg}, []any{&signature})
	if err != nil {
		return nil, err
	}
	return &signature, nil
}

func (c *Client) Verify(canisterId, signature, msg, publicKey string) (bool, error) {
	canId, err := principal.Decode(canisterId)
	if err != nil {
		return false, err
	}
	var resp bool
	err = c.call(canId, "verify", []any{signature, msg, publicKey}, []any{&resp})
	if err != nil {
		return false, err
	}
	return resp, nil
}

func (c *Client) BtcUtxo(canisterId, addr string, network ...string) (interface{}, error) {
	canId, err := principal.Decode(canisterId)
	if err != nil {
		return nil, err
	}
	var resp string
	err = c.agent.Call(canId, "get_utxos", []any{addr}, []any{&resp})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) BtcBalance(canisterId, addr string, network ...string) (uint64, error) {
	canId, err := principal.Decode(canisterId)
	if err != nil {
		return 0, err
	}
	var resp uint64
	err = c.agent.Call(canId, "get_balance", []any{addr}, []any{&resp})
	if err != nil {
		return 0, err
	}
	fmt.Printf("%v \n", resp)
	return resp, nil
}

func (c *Client) call(canisterID principal.Principal, method string, args []any, rets []any) error {
	return c.agent.Call(canisterID, method, args, rets)
}
func (c *Client) query(canisterID principal.Principal, method string, args []any, rets []any) error {
	return c.agent.Query(canisterID, method, args, rets)
}

func NewClient() (*Client, error) {
	config := agent.Config{
		PollTimeout: 45 * time.Second,
	}
	agent, err := agent.New(config)
	if err != nil {
		return nil, err
	}
	canId, err := principal.Decode("0")
	icAgent, err := ic.NewAgent(canId, config)
	if err != nil {
		return nil, err
	}
	return &Client{
		agent:   agent,
		icAgent: icAgent,
	}, nil
}
