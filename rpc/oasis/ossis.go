package oasis

type Client struct {
}

func NewClient(url string) (*Client, error) {
	return &Client{}, nil
}

func (c *Client) Redeem(proof string) (string, error) {
	panic(proof)
}
