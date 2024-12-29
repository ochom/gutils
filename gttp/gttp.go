package gttp

import (
	"time"

	"github.com/ochom/gutils/env"
)

type ClientType int

const (
	DefaultHttp = iota
	GoFiber
)

type Client interface {
	post(url string, headers M, body any, timeout ...time.Duration) (resp *Response, err error)
	get(url string, headers M, timeout ...time.Duration) (resp *Response, err error)
	sendRequest(url, method string, headers M, body any, timeout ...time.Duration) (resp *Response, err error)
}

var client Client

func init() {
	switch env.Int("HTTP_CLIENT", 1) {
	case DefaultHttp:
		client = new(defaultClient)
	case GoFiber:
		client = new(fiberClient)
	default:
		panic("unknown http client")
	}
}

// Post sends a POST request to the specified URL.
func Post(url string, headers M, body any, timeout ...time.Duration) (resp *Response, err error) {
	return client.post(url, headers, body, timeout...)
}

// Get sends a GET request to the specified URL.
func Get(url string, headers M, timeout ...time.Duration) (resp *Response, err error) {
	return client.get(url, headers, timeout...)
}

// SendRequest sends a request to the specified URL.
func SendRequest(url, method string, headers M, body any, timeout ...time.Duration) (resp *Response, err error) {
	return client.sendRequest(url, method, headers, body, timeout...)
}
