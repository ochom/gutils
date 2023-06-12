package gttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

// Request sends a request to the specified URL.
type Request struct {
	url     string
	headers map[string]string
	body    []byte
	timeOut time.Duration
}

// Response is the response of the request.
type Response struct {
	Status int
	Body   []byte
}

// NewRequest ...
func NewRequest(url string, headers map[string]string, body []byte) *Request {
	return NewRequestWithTimeOut(url, headers, body, 30*time.Second)
}

// NewRequestWithTimeOut ...
func NewRequestWithTimeOut(url string, headers map[string]string, body []byte, timeOut time.Duration) *Request {
	return &Request{url, headers, body, timeOut}
}

// client ...
func (r *Request) client() *http.Client {
	return &http.Client{
		Timeout: r.timeOut,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

// Post sends a POST request to the specified URL.
func (r *Request) Post() (res *Response, err error) {
	req, err := http.NewRequest("POST", r.url, bytes.NewBuffer(r.body))
	if err != nil {
		return
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reqDo, err := r.client().Do(req.WithContext(ctx))
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
func (r *Request) Get() (res *Response, err error) {
	req, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		return
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reqDo, err := r.client().Do(req.WithContext(ctx))
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

// Put sends a PUT request to the specified URL.
func (r *Request) Put() (res *Response, err error) {
	req, err := http.NewRequest("PUT", r.url, bytes.NewBuffer(r.body))
	if err != nil {
		return
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reqDo, err := r.client().Do(req.WithContext(ctx))
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

// Delete sends a DELETE request to the specified URL.
func (r *Request) Delete() (res *Response, err error) {
	req, err := http.NewRequest("DELETE", r.url, nil)
	if err != nil {
		return
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reqDo, err := r.client().Do(req.WithContext(ctx))
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
func (r *Request) Custom(method string) (res *Response, err error) {
	req, err := http.NewRequest(method, r.url, bytes.NewBuffer(r.body))
	if err != nil {
		return
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reqDo, err := r.client().Do(req.WithContext(ctx))
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
