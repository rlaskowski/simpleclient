package client

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"os"
)

type Client struct {
	client  *http.Client
	headers textproto.MIMEHeader
}

func NewClient() *Client {
	return &Client{
		client:  &http.Client{},
		headers: make(textproto.MIMEHeader),
	}
}

func (c *Client) DownloadFile(filename, url string) error {
	req, err := c.GetRequest(url)
	if err != nil {
		return err
	}

	_, err = c.newResponse(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetRequest(url string) (*http.Request, error) {
	return c.newRequest(http.MethodGet, url)
}

func (c *Client) newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) newResponse(req *http.Request) (*Response, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return NewResponse(res), nil
}

/* func (c *Client) run(id int) error {
	filename := fmt.Sprintf("%v.mp4", id)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.bambulax.com/private/bvids/%s", filename), nil)

	if err != nil {
		log.Fatalf("Error when try to get resource, error type: %s", err)
	}

	auth := c.auth()

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))

	res, err := c.client.Do(req)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Could not download resource %s http status %v", filename, res.StatusCode)
	}

	path := filepath.Join("/media/rafal/My Passport/Filmy/bambulax", filename)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		log.Printf("Could not create file %s due to %s", filename, err)
	}

	log.Printf("Writing file in path: %s", path)

	b := make([]byte, 1024*1024)
	if _, err := io.CopyBuffer(file, res.Body, b); err != nil {
		return err
	}

	return file.Close()
} */

func (c *Client) AddHeader(name, val string) {
	c.headers.Set(name, val)
}

func (c *Client) copy(filename string, src io.Reader) (written int64, err error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		return 0, fmt.Errorf("Could not create file %s due to %s", filename, err)
	}

	buff := make([]byte, 32*1024)

	return io.CopyBuffer(file, src, buff)
}

func (c *Client) download(filepath, url string) error {
	return nil
}

