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
	req := c.getClient(url, "POST")
	for k, v := range headers {
		req.Add(k, v)
	}

	timeout := time.Hour
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req.Body(helpers.ToBytes(body))
	code, content, err := c.do(ctx, req)
	return &Response{code, content}, err
}

// Get sends a GET request to the specified URL.
func (c *fiberClient) Get(url string, headers M, timeouts ...time.Duration) (resp *Response, err error) {
	req := c.getClient(url, "GET")
	for k, v := range headers {
		req.Add(k, v)
	}

	timeout := time.Hour
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	code, content, err := c.do(ctx, req)
	return &Response{code, content}, err
}

// Custom sends a custom request to the specified URL.
func (c *fiberClient) Custom(url, method string, headers M, body any, timeouts ...time.Duration) (resp *Response, err error) {
	req := c.getClient(url, method)
	for k, v := range headers {
		req.Add(k, v)
	}

	timeout := time.Hour
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req.Body(helpers.ToBytes(body))
	code, content, err := c.do(ctx, req)
	return &Response{code, content}, err
}

func (fiberClient) do(ctx context.Context, req *fiber.Agent) (int, []byte, error) {
	for {
		select {
		case <-ctx.Done():
			return 500, nil, ctx.Err()
		default:
			code, content, errs := req.Bytes()
			if len(errs) > 0 {
				for _, err := range errs {
					logs.Error("client error: %s", err.Error())
				}

				return 500, nil, errs[0]
			}

			return code, content, nil
		}
	}
}

func (fiberClient) getClient(url, method string) *fiber.Agent {
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
		panic(fmt.Errorf("unknown method: %s", method))
	}

	req.InsecureSkipVerify()
	return req
}
