package request

import (
	"io"
	"net/http"
	"time"
)

// AppClient 包共用的client
type appClient struct {
	url    string
	method string
	header map[string]string
	client *http.Client
}

var Client appClient

func New(url, method string, header map[string]string) appClient {
	Client = appClient{
		url:    url,
		method: method,
		header: header,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
	return Client
}

func (c appClient) Get(payload io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.url, payload)
	if err != nil {
		return nil, err
	}
	if c.header != nil {
		for k, v := range c.header {
			req.Header.Set(k, v)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c appClient) Post(body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, c.url, body)
	if err != nil {
		return nil, err
	}
	if c.header != nil {
		for k, v := range c.header {
			req.Header.Set(k, v)
		}
	}

	resp, err := c.client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
