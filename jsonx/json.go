// Package jsonx provides simplified JSON encoding and decoding utilities with generics support.
//
// This package wraps the standard library's encoding/json package to provide a more
// ergonomic API with less boilerplate code. It includes both error-ignoring convenience
// functions and explicit error-returning variants.
//
// Example usage:
//
//	type User struct {
//		Name  string `json:"name"`
//		Email string `json:"email"`
//	}
//
//	// Encode struct to JSON
//	user := User{Name: "Alice", Email: "alice@example.com"}
//	data := jsonx.Encode(user)
//	fmt.Println(data.String()) // {"name":"Alice","email":"alice@example.com"}
//
//	// Decode JSON to struct
//	jsonData := []byte(`{"name":"Bob","email":"bob@example.com"}`)
//	decoded := jsonx.Decode[User](jsonData)
//	fmt.Println(decoded.Name) // "Bob"
//
//	// With error handling
//	data, err := jsonx.EncodeErr(user)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	decoded, err := jsonx.DecodeErr[User](jsonData)
//	if err != nil {
//		log.Fatal(err)
//	}
package jsonx

import (
	baseJSON "encoding/json"

	"github.com/ochom/gutils/logs"
)

// byteData is a wrapper around []byte that provides convenience methods.
type byteData []byte

// String converts the byte data to a string representation.
// Returns an empty string if the data is nil.
//
// Example:
//
//	data := jsonx.Encode(map[string]int{"count": 42})
//	fmt.Println(data.String()) // {"count":42}
func (b byteData) String() string {
	if b == nil {
		return ""
	}
	return string(b)
}

// Bytes returns the underlying byte slice.
//
// Example:
//
//	data := jsonx.Encode(payload)
//	resp, err := http.Post(url, "application/json", bytes.NewReader(data.Bytes()))
func (b byteData) Bytes() []byte {
	return b
}

// Encode marshals any value to JSON bytes.
// On error, logs the error and returns nil.
// Use EncodeErr if you need to handle errors explicitly.
//
// Example:
//
//	// Encode a struct
//	type Config struct { Port int `json:"port"` }
//	data := jsonx.Encode(Config{Port: 8080})
//	// data.String() = {"port":8080}
//
//	// Encode a map
//	data := jsonx.Encode(map[string]string{"status": "ok"})
//
//	// Use with HTTP client
//	resp, _ := gttp.Post(url, headers, jsonx.Encode(payload).Bytes())
func Encode(payload any) byteData {
	bytesPayload, err := baseJSON.Marshal(&payload)
	if err != nil {
		logs.Error("Failed to marshal JSON: %s", err.Error())
		return nil
	}

	return bytesPayload
}

// EncodeErr marshals any value to JSON bytes and returns an error if marshaling fails.
// Use this when you need to handle encoding errors explicitly.
//
// Example:
//
//	data, err := jsonx.EncodeErr(payload)
//	if err != nil {
//		return fmt.Errorf("failed to encode payload: %w", err)
//	}
//	// use data...
func EncodeErr(payload any) (byteData, error) {
	bytesPayload, err := baseJSON.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	return bytesPayload, nil
}

// Decode unmarshals JSON bytes into the specified type T.
// On error, logs the error and returns the zero value of T.
// Use DecodeErr if you need to handle errors explicitly.
//
// Example:
//
//	type User struct {
//		ID   string `json:"id"`
//		Name string `json:"name"`
//	}
//
//	// Decode HTTP response
//	resp, _ := gttp.Get(url, nil)
//	user := jsonx.Decode[User](resp.Body)
//	fmt.Println(user.Name)
//
//	// Decode slice
//	users := jsonx.Decode[[]User](resp.Body)
//
//	// Decode map
//	data := jsonx.Decode[map[string]any](resp.Body)
func Decode[T any](payload []byte) T {
	var data T
	if err := baseJSON.Unmarshal(payload, &data); err != nil {
		logs.Error("Failed to unmarshal JSON: %s", err.Error())
		return data
	}

	return data
}

// DecodeErr unmarshals JSON bytes into the specified type T and returns an error if unmarshaling fails.
// Use this when you need to handle decoding errors explicitly.
//
// Example:
//
//	user, err := jsonx.DecodeErr[User](jsonData)
//	if err != nil {
//		return fmt.Errorf("invalid JSON: %w", err)
//	}
//	// use user...
//
//	// Handle API response
//	type APIResponse struct {
//		Data   []Item `json:"data"`
//		Error  string `json:"error"`
//	}
//	resp, err := jsonx.DecodeErr[APIResponse](body)
//	if err != nil {
//		log.Error("Failed to parse response")
//	}
func DecodeErr[T any](payload []byte) (T, error) {
	var data T
	if err := baseJSON.Unmarshal(payload, &data); err != nil {
		return data, err
	}

	return data, nil
}
