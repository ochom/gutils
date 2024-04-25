package helpers

import (
	"encoding/json"

	"github.com/ochom/gutils/logs"
)

// ToJSON converts a struct to JSON
func ToJSON(payload any) []byte {
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
	bytesPayload, err := json.Marshal(&payload)
	if err != nil {
		return nil
	}

	return bytesPayload
}

// FromJSON converts json byte  to struct
func FromJSON[T any](payload []byte) T {
	var data T
	if payload == nil {
		return data
	}

	err := json.Unmarshal(payload, &data)
	if err != nil {
		logs.Error("Failed to unmarshal JSON: %s", err.Error())
		return data
	}

	return data
}
