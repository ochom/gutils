package gttp

import "encoding/json"

// M ...
type M map[string]string

// toBytes ...
func toBytes(payload any) []byte {
	// Check if the payload is already of type []byte.
	if bytesPayload, ok := payload.([]byte); ok {
		return bytesPayload
	}

	// Check if the payload is not nil.
	if payload != nil {
		// Attempt to convert the payload to []byte.
		if bytesPayload, ok := payload.(string); ok {
			return []byte(bytesPayload)
		} else if jsonPayload, err := json.Marshal(payload); err == nil {
			return jsonPayload
		}
	}

	// Return nil if the payload cannot be converted.
	return nil
}
