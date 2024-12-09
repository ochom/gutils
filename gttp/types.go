package gttp

// Response is the response of the request.
type Response struct {
	// Status is the HTTP status code.
	StatusCode int

	// Body is the response body.
	Body []byte
}

// M  is a map[string]string
type M map[string]string

// Map is a map[string]any{}
type Map map[string]any
