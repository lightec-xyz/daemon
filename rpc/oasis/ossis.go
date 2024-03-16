package oasis

type Client struct {
}

func NewClient(url string) (*Client, error) {
	return &Client{}, nil
}
