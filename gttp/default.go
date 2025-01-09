package gttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

type defaultClient struct{}

// getClient ...
func (*defaultClient) getClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return client
}

// post sends a POST request to the specified URL.
func (c *defaultClient) post(url string, headers M, body []byte, timeout ...time.Duration) (resp *Response, err error) {
	return c.sendRequest(url, "POST", headers, body, timeout...)
}

// get sends a GET request to the specified URL.
func (c *defaultClient) get(url string, headers M, timeout ...time.Duration) (resp *Response, err error) {
	return c.sendRequest(url, "GET", headers, nil, timeout...)
}

// sendRequest sends a  request to the specified URL.
func (c *defaultClient) sendRequest(url, method string, headers M, body []byte, timeout ...time.Duration) (resp *Response, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), getTimeout(timeout...))
	defer cancel()

	res, err := c.getClient().Do(req.WithContext(ctx))
	if err != nil {
		return
	}

	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	return &Response{res.StatusCode, content}, nil
}
