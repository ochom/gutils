package gttp

// Response represents an HTTP response from a request.
// It contains the status code, body content, and any errors encountered.
type Response struct {
	// StatusCode is the HTTP status code (e.g., 200, 404, 500)
	StatusCode int

	// Body contains the raw response body as bytes
	Body []byte

	// Errors contains any errors that occurred during the request
	Errors []error
}

// NewResponse creates a new Response with the given status code, errors, and body.
//
// Example:
//
//	resp := gttp.NewResponse(200, nil, []byte(`{"status": "ok"}`))
func NewResponse(code int, err []error, body []byte) *Response {
	return &Response{
		StatusCode: code,
		Errors:     err,
		Body:       body,
	}
}

// M is a convenience type alias for map[string]string.
// Used primarily for HTTP headers.
//
// Example:
//
//	headers := gttp.M{
//		"Content-Type": "application/json",
//		"Authorization": "Bearer token",
//	}
type M map[string]string

// Map is a convenience type alias for map[string]any.
// Used for generic key-value data structures.
//
// Example:
//
//	data := gttp.Map{
//		"name": "John",
//		"age": 30,
//		"active": true,
//	}
type Map map[string]any
