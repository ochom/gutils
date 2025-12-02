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
func (c *fiberClient) sendRequest(url, method string, headers M, body []byte, timeouts ...time.Duration) (*Response, error) {
	timeout := time.Hour
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	resp := c.makeRequest(url, method, headers, body, timeout)
	if len(resp.Errors) == 0 {
		return resp, nil
	}

	errStrings := []string{}
	for _, err := range resp.Errors {
		errStrings = append(errStrings, err.Error())
	}

	return resp, errors.New(strings.Join(errStrings, ", "))
}

// makeRequest sends a request to the specified URL.
func (c *fiberClient) makeRequest(url, method string, headers M, body []byte, timeout time.Duration) *Response {
	var agent *fiber.Agent
	switch method {
	case "POST":
		agent = fiber.Post(url)
	case "GET":
		agent = fiber.Get(url)
	case "DELETE":
		agent = fiber.Delete(url)
	case "PUT":
		agent = fiber.Put(url)
	case "PATCH":
		agent = fiber.Patch(url)
	default:
		err := fmt.Errorf("unknown method: %s", method)
		return NewResponse(500, []error{err}, nil)
	}

	// skip ssl verification
	agent.InsecureSkipVerify()
	agent.Timeout(timeout)

	// add request headers
	for k, v := range headers {
		agent.Add(k, v)
	}

	// add request body
	if method == "POST" || method == "PUT" || method == "PATCH" {
		agent.Body(body)
	}

	// make request
	code, content, errs := agent.Bytes()
	if code == 0 {
		code = 500
	}

	return NewResponse(code, errs, content)
}
