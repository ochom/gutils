package gttp

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
)

type fiberClient struct{}

// Post sends a POST request to the specified URL.
func (c *fiberClient) Post(url string, headers M, body any, timeouts ...time.Duration) (resp *Response, err error) {
	return c.SendRequest(url, "POST", headers, body, timeouts...)
}

// Get sends a GET request to the specified URL.
func (c *fiberClient) Get(url string, headers M, timeouts ...time.Duration) (resp *Response, err error) {
	return c.SendRequest(url, "GET", headers, nil, timeouts...)
}

// SendRequest sends a request to the specified URL.
func (c *fiberClient) SendRequest(url, method string, headers M, body any, timeouts ...time.Duration) (resp *Response, err error) {
	timeout := time.Hour
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return &Response{500, nil}, ctx.Err()
		default:
			return c.sendRequest(url, method, headers, body)
		}
	}
}

func (c *fiberClient) sendRequest(url, method string, headers M, body any) (resp *Response, err error) {
	client := fiber.AcquireClient()
	var req *fiber.Agent

	switch method {
	case "POST":
		req = client.Post(url)
	case "GET":
		req = client.Get(url)
	case "DELETE":
		req = client.Delete(url)
	case "PUT":
		req = client.Put(url)
	case "PATCH":
		req = client.Patch(url)
	default:
		return &Response{500, nil}, fmt.Errorf("unknown method: %s", method)
	}

	// skip ssl verification
	req.InsecureSkipVerify()

	for k, v := range headers {
		req.Add(k, v)
	}

	if method != "GET" {
		req.Body(helpers.ToBytes(body))
	}

	code, content, errs := req.Bytes()
	if len(errs) > 0 {
		for _, err := range errs {
			logs.Error("client error: %s", err.Error())
		}

		return &Response{500, content}, errs[0]
	}

	return &Response{code, content}, err
}
