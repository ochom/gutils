package gttp

import (
	"context"
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
	timeout := time.Hour
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := make(chan *Response, 1)
	go func() {
		resp := c.makeRequest(url, method, headers, body)
		result <- resp
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-result:
		if len(r.Errors) == 0 {
			return r, nil
		}

		errStrings := []string{}
		for _, err := range r.Errors {
			errStrings = append(errStrings, err.Error())
		}

		return r, errors.New(strings.Join(errStrings, ", "))
	}
}

// makeRequest sends a request to the specified URL.
func (c *fiberClient) makeRequest(url, method string, headers M, body []byte) (resp *Response) {
	var req *fiber.Agent
	switch method {
	case "POST":
		req = fiber.Post(url)
	case "GET":
		req = fiber.Get(url)
	case "DELETE":
		req = fiber.Delete(url)
	case "PUT":
		req = fiber.Put(url)
	case "PATCH":
		req = fiber.Patch(url)
	default:
		err := fmt.Errorf("unknown method: %s", method)
		return NewResponse(500, []error{err}, nil)
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
	if code == 0 {
		code = 500
	}

	return NewResponse(code, errs, content)
}
