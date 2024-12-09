package gttp

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
)

type fiberClient struct {
}

// Post sends a POST request to the specified URL.
func (c *fiberClient) Post(url string, headers M, body any, timeout ...time.Duration) (resp *Response, err error) {
	client := fiber.AcquireClient()
	req := client.Post(url)
	req.InsecureSkipVerify()

	for k, v := range headers {
		req.Add(k, v)
	}

	req.Body(helpers.ToBytes(body))
	code, content, errs := req.Bytes()
	if len(errs) > 0 {
		for _, er := range errs {
			logs.Error("client error: %s", er.Error())
		}

		return nil, errs[0]
	}

	return &Response{code, content}, nil
}

// Get sends a GET request to the specified URL.
func (c *fiberClient) Get(url string, headers M, timeout ...time.Duration) (resp *Response, err error) {
	client := fiber.AcquireClient()
	req := client.Get(url)
	req.InsecureSkipVerify()

	for k, v := range headers {
		req.Add(k, v)
	}

	code, content, errs := req.Bytes()
	if len(errs) > 0 {
		for _, err := range errs {
			logs.Error("client error: %s", err.Error())
		}
		return nil, errs[0]
	}

	return &Response{code, content}, nil
}

// Custom sends a custom request to the specified URL.
func (c *fiberClient) Custom(url, method string, headers M, body any, timeout ...time.Duration) (resp *Response, err error) {
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
		logs.Error("Unknown method: %s", method)
		return nil, fmt.Errorf("unknown method: %s", method)
	}

	req.InsecureSkipVerify()

	for k, v := range headers {
		req.Add(k, v)
	}

	req.Body(helpers.ToBytes(body))
	code, content, errs := req.Bytes()
	if len(errs) > 0 {
		for _, err := range errs {
			logs.Error("client error: %s", err.Error())
		}

		return nil, errs[0]
	}

	return &Response{code, content}, nil
}
