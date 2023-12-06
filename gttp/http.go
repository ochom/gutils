package gttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

// Response is the response of the request.
type Response struct {
	// Status is the HTTP status code.
	Status int

	// Body is the response body.
	Body []byte
}

var client *http.Client

// getClient ...
func getClient(timeout ...time.Duration) *http.Client {
	if client != nil {
		return client
	}

	timeOut := getTimeout(timeout...)
	client = &http.Client{
		Timeout: timeOut,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     time.Second * 90,
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
func Post(url string, headers M, body any, timeout ...time.Duration) (res *Response, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(toBytes(body)))
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), getTimeout(timeout...))
	defer cancel()

	reqDo, err := getClient(timeout...).Do(req.WithContext(ctx))
	if err != nil {
		return
	}

	defer reqDo.Body.Close()

	bodyBytes, err := io.ReadAll(reqDo.Body)
	if err != nil {
		return
	}

	return &Response{reqDo.StatusCode, bodyBytes}, nil
}

// Get sends a GET request to the specified URL.
func Get(url string, headers M, timeout ...time.Duration) (res *Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), getTimeout(timeout...))
	defer cancel()

	reqDo, err := getClient(timeout...).Do(req.WithContext(ctx))
	if err != nil {
		return
	}

	defer reqDo.Body.Close()

	bodyBytes, err := io.ReadAll(reqDo.Body)
	if err != nil {
		return
	}

	return &Response{reqDo.StatusCode, bodyBytes}, nil
}

// Custom sends a custom request to the specified URL.
func Custom(url, method string, headers M, body any, timeout ...time.Duration) (res *Response, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(toBytes(body)))
	if err != nil {
		return
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), getTimeout(timeout...))
	defer cancel()

	reqDo, err := getClient(timeout...).Do(req.WithContext(ctx))
	if err != nil {
		return
	}

	defer reqDo.Body.Close()

	bodyBytes, err := io.ReadAll(reqDo.Body)
	if err != nil {
		return
	}

	return &Response{reqDo.StatusCode, bodyBytes}, nil
}
