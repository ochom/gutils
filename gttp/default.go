package gttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"github.com/ochom/gutils/helpers"
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

// Post sends a POST request to the specified URL.
func (c *defaultClient) Post(url string, headers M, body any, timeout ...time.Duration) (resp *Response, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(helpers.ToBytes(body)))
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

// Get sends a GET request to the specified URL.
func (c *defaultClient) Get(url string, headers M, timeout ...time.Duration) (resp *Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
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

// Custom sends a custom request to the specified URL.
func (c *defaultClient) Custom(url, method string, headers M, body any, timeout ...time.Duration) (resp *Response, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(helpers.ToBytes(body)))
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
