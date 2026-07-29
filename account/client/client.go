package client

type Client struct {
	url string
}

func NewClient(url string) (*Client, error) {
	return &Client{url: url}, nil
}

func (c *Client) Close() {}
