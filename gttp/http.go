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

// Response is the response of the request.
type Response struct {
	// Status is the HTTP status code.
	StatusCode int

	// Body is the response body.
	Body []byte
}

// getClient ...
func getClient() *http.Client {
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

func getTimeout(timeout ...time.Duration) time.Duration {
	if len(timeout) == 0 {
		return 10 * time.Second
	}

	return timeout[0]
}

// Post sends a POST request to the specified URL.
func Post(url string, headers M, body any, timeout ...time.Duration) (resp *Response, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(helpers.ToBytes(body)))
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), getTimeout(timeout...))
	defer cancel()

	res, err := getClient().Do(req.WithContext(ctx))
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
func Get(url string, headers M, timeout ...time.Duration) (resp *Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), getTimeout(timeout...))
	defer cancel()

	res, err := getClient().Do(req.WithContext(ctx))
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
func Custom(url, method string, headers M, body any, timeout ...time.Duration) (resp *Response, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(helpers.ToBytes(body)))
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), getTimeout(timeout...))
	defer cancel()

	res, err := getClient().Do(req.WithContext(ctx))
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
