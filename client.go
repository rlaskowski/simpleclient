package client

import (
	"fmt"
	"net/http"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{},
	}
}

func (c *Client) GetRequest(url string) (*http.Request, error) {
	return c.NewRequest(http.MethodGet, url)
}

func (c *Client) NewRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) Do(req *http.Request) (*Response, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad response status: %s", res.Status)
	}

	return NewResponse(res), nil
}
