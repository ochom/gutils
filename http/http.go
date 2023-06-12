package http

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
func (r *Request) Post() (resBody []byte, status int, err error) {
	req, err := http.NewRequest("POST", r.url, bytes.NewBuffer(r.body))
	if err != nil {
		return
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.client().Do(req.WithContext(ctx))
	if err != nil {
		return
	}

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	return bodyBytes, res.StatusCode, nil
}

// Get sends a GET request to the specified URL.
func (r *Request) Get() (resBody []byte, status int, err error) {
	req, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		return
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.client().Do(req.WithContext(ctx))
	if err != nil {
		return
	}

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	return bodyBytes, res.StatusCode, nil
}

// Put sends a PUT request to the specified URL.
func (r *Request) Put() (resBody []byte, status int, err error) {
	req, err := http.NewRequest("PUT", r.url, bytes.NewBuffer(r.body))
	if err != nil {
		return
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.client().Do(req.WithContext(ctx))
	if err != nil {
		return
	}

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	return bodyBytes, res.StatusCode, nil
}

// Delete sends a DELETE request to the specified URL.
func (r *Request) Delete() (resBody []byte, status int, err error) {
	req, err := http.NewRequest("DELETE", r.url, nil)
	if err != nil {
		return
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.client().Do(req.WithContext(ctx))
	if err != nil {
		return
	}

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	return bodyBytes, res.StatusCode, nil
}

// Custom sends a custom request to the specified URL.
func (r *Request) Custom(method string) (resBody []byte, status int, err error) {
	req, err := http.NewRequest(method, r.url, bytes.NewBuffer(r.body))
	if err != nil {
		return
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := r.client().Do(req.WithContext(ctx))
	if err != nil {
		return
	}

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	return bodyBytes, res.StatusCode, nil
}
