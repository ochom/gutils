package gttp

import (
	"time"

	"github.com/ochom/gutils/logs"
)

type ClientType int

const (
	DefaultHttp = iota
	GoFiber
)

type Client interface {
	Post(url string, headers M, body any, timeout ...time.Duration) (resp *Response, err error)
	Get(url string, headers M, timeout ...time.Duration) (resp *Response, err error)
	SendRequest(url, method string, headers M, body any, timeout ...time.Duration) (resp *Response, err error)
}

func NewClient(clientType ClientType) Client {
	switch clientType {
	case DefaultHttp:
		return &defaultClient{}
	case GoFiber:
		return &fiberClient{}
	default:
		logs.Error("Unknown provider: %d", clientType)
		return nil
	}
}
