package gttp

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
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
	// timeout := time.Hour
	// if len(timeouts) > 0 {
	// 	timeout = timeouts[0]
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// defer cancel()

	// select {
	// case <-ctx.Done():
	// 	return &Response{500, nil}, ctx.Err()
	// default:
	// 	return c.makeRequest(url, method, headers, body)
	// }

	return c.makeRequest(url, method, headers, body)
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
		errStrings := []string{}
		for _, err := range errs {
			errStrings = append(errStrings, err.Error())
		}

		if code == 0 {
			code = 500
		}
		return &Response{code, content}, errors.New(strings.Join(errStrings, ", "))
	}

	return &Response{code, content}, err
}
