package gttp

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ochom/gutils/logs"
)

type fiberClient struct{}

// post sends a POST request to the specified URL.
func (c *fiberClient) post(url string, headers M, body []byte, timeouts ...time.Duration) (resp *Response, err error) {
	return c.sendRequest(url, "POST", headers, body, timeouts...)
}

// get sends a GET request to the specified URL.
func (c *fiberClient) get(url string, headers M, timeouts ...time.Duration) (resp *Response, err error) {
	return c.sendRequest(url, "GET", headers, nil, timeouts...)
}

// sendRequest sends a request to the specified URL.
func (c *fiberClient) sendRequest(url, method string, headers M, body []byte, timeouts ...time.Duration) (resp *Response, err error) {
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
			return c.makeRequest(url, method, headers, body)
		}
	}
}

// makeRequest sends a request to the specified URL.
func (c *fiberClient) makeRequest(url, method string, headers M, body []byte) (resp *Response, err error) {
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

	if method == "POST" || method == "PUT" || method == "PATCH" {
		req.Body(body)
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
