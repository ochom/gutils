package helpers

import (
	"github.com/goccy/go-json"

	"github.com/ochom/gutils/logs"
)

// ToBytes converts provided interface to slice of bytes
func ToBytes[T any](payload T) []byte {
	bytesPayload, err := json.Marshal(&payload)
	if err != nil {
		logs.Error("Failed to marshal JSON: %s", err.Error())
		return nil
	}

	return bytesPayload
}

// FromBytes converts slice of bytes to provided interface
func FromBytes[T any](payload []byte) T {
	var data T
	if err := json.Unmarshal(payload, &data); err != nil {
		logs.Error("Failed to unmarshal JSON: %s", err.Error())
		return data
	}

	return data
}
