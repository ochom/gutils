package gttp

import "encoding/json"

// M  is a map[string]string
type M map[string]string

// Map is a map[string]interface{}
type Map map[string]interface{}

// toBytes ...
func toBytes(payload any) []byte {
	if payload == nil {
		return nil
	}

	// Check if the payload is already of type []byte.
	if bytesPayload, ok := payload.([]byte); ok {
		return bytesPayload
	}

	// Check if the payload is already of type string.
	if stringPayload, ok := payload.(string); ok {
		return []byte(stringPayload)
	}

	// Marshal the payload to JSON.
	bytesPayload, err := json.Marshal(payload)
	if err != nil {
		return nil
	}

	return bytesPayload
}
