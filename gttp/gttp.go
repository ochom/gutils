// Package gttp provides a simple HTTP client abstraction with support for multiple backends.
//
// The package automatically selects the HTTP client implementation based on the HTTP_CLIENT
// environment variable:
//   - "fiber": Uses gofiber/fiber client (default, faster for high-throughput)
//   - "default": Uses Go's standard net/http client
//
// Features:
//   - Simple API for GET and POST requests
//   - Configurable timeouts
//   - TLS/SSL support with skip verification option
//   - Header management
//
// Example usage:
//
//	import (
//		"github.com/ochom/gutils/gttp"
//		"github.com/ochom/gutils/jsonx"
//		"time"
//	)
//
//	// Simple GET request
//	resp, err := gttp.Get("https://api.example.com/users", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Status:", resp.StatusCode)
//	fmt.Println("Body:", string(resp.Body))
//
//	// GET with headers
//	headers := gttp.M{"Authorization": "Bearer token123"}
//	resp, err := gttp.Get("https://api.example.com/me", headers)
//
//	// POST with JSON body
//	payload := map[string]string{"name": "John", "email": "john@example.com"}
//	resp, err := gttp.Post(
//		"https://api.example.com/users",
//		gttp.M{"Content-Type": "application/json"},
//		jsonx.Encode(payload),
//	)
//
//	// Request with custom timeout
//	resp, err := gttp.Get("https://slow-api.example.com/data", nil, 30*time.Second)
package gttp

import (
	"time"

	"github.com/ochom/gutils/env"
)

// Client defines the interface for HTTP client implementations.
type Client interface {
	post(url string, headers M, body []byte, timeout ...time.Duration) (resp *Response, err error)
	get(url string, headers M, timeout ...time.Duration) (resp *Response, err error)
	sendRequest(url, method string, headers M, body []byte, timeout ...time.Duration) (resp *Response, err error)
}

var client Client

func init() {
	switch env.Get("HTTP_CLIENT", "fiber") {
	case "default":
		client = new(defaultClient)
	case "fiber":
		client = new(fiberClient)
	default:
		client = new(fiberClient)
	}
}

// Post sends a POST request to the specified URL with headers and body.
//
// The timeout parameter is optional (default: 10 seconds for default client, 1 hour for fiber).
//
// Example:
//
//	// Simple POST
//	resp, err := gttp.Post("https://api.example.com/data", nil, []byte(`{"key":"value"}`))
//
//	// POST with headers and timeout
//	headers := gttp.M{
//		"Content-Type": "application/json",
//		"Authorization": "Bearer token",
//	}
//	resp, err := gttp.Post("https://api.example.com/data", headers, body, 5*time.Second)
func Post(url string, headers M, body []byte, timeout ...time.Duration) (resp *Response, err error) {
	return client.post(url, headers, body, timeout...)
}

// Get sends a GET request to the specified URL with optional headers.
//
// The timeout parameter is optional (default: 10 seconds for default client, 1 hour for fiber).
//
// Example:
//
//	// Simple GET
//	resp, err := gttp.Get("https://api.example.com/users", nil)
//
//	// GET with authentication
//	resp, err := gttp.Get(
//		"https://api.example.com/protected",
//		gttp.M{"Authorization": "Bearer token"},
//	)
//
//	// GET with custom timeout
//	resp, err := gttp.Get("https://api.example.com/slow", nil, 30*time.Second)
func Get(url string, headers M, timeout ...time.Duration) (resp *Response, err error) {
	return client.get(url, headers, timeout...)
}

// SendRequest sends a custom HTTP request with the specified method.
//
// Supported methods: GET, POST, PUT, PATCH, DELETE
//
// Example:
//
//	// PUT request
//	resp, err := gttp.SendRequest(
//		"https://api.example.com/users/123",
//		"PUT",
//		gttp.M{"Content-Type": "application/json"},
//		jsonx.Encode(updateData),
//	)
//
//	// DELETE request
//	resp, err := gttp.SendRequest(
//		"https://api.example.com/users/123",
//		"DELETE",
//		nil,
//		nil,
//	)
func SendRequest(url, method string, headers M, body []byte, timeout ...time.Duration) (resp *Response, err error) {
	return client.sendRequest(url, method, headers, body, timeout...)
}
