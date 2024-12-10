package gttp

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
)

type fiberClient struct{}

func do(req *fiber.Agent) (int, []byte, error) {
	code, content, errs := req.Bytes()
	if len(errs) > 0 {
		for _, err := range errs {
			logs.Error("client error: %s", err.Error())
		}

		return 0, nil, errs[0]
	}

	return code, content, nil
}

func getClient(url, method string) *fiber.Agent {
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

// Post sends a POST request to the specified URL.
func (c *fiberClient) Post(url string, headers M, body any, timeout ...time.Duration) (resp *Response, err error) {
	req := getClient(url, "POST")
	for k, v := range headers {
		req.Add(k, v)
	}

	req.Body(helpers.ToBytes(body))
	code, content, err := do(req)
	if err != nil {
		return nil, err
	}
	return &Response{code, content}, nil
}

// Get sends a GET request to the specified URL.
func (c *fiberClient) Get(url string, headers M, timeout ...time.Duration) (resp *Response, err error) {
	req := getClient(url, "GET")
	for k, v := range headers {
		req.Add(k, v)
	}

	code, content, err := do(req)
	if err != nil {
		return nil, err
	}

	return &Response{code, content}, nil
}

// Custom sends a custom request to the specified URL.
func (c *fiberClient) Custom(url, method string, headers M, body any, timeout ...time.Duration) (resp *Response, err error) {
	req := getClient(url, method)
	for k, v := range headers {
		req.Add(k, v)
	}

	req.Body(helpers.ToBytes(body))
	code, content, err := do(req)
	if err != nil {
		return nil, err
	}

	return &Response{code, content}, nil
}
