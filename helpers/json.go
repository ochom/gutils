package helpers

import (
	"encoding/json"

	"github.com/ochom/gutils/logs"
)

// Marshallable ...
type Marshallable interface {
	map[string]any | map[string]string |
		[]map[string]any | []map[string]string |
		struct{} | *struct{} |
		[]struct{} | []*struct{}
}

// ToBytes converts provided interface to slice of bytes
func ToBytes[T Marshallable](payload T) []byte {
	bytesPayload, err := json.Marshal(&payload)
	if err != nil {
		logs.Error("Failed to marshal JSON: %s", err.Error())
		return nil
	}

	return bytesPayload
}

// FromBytes converts slice of bytes to provided interface
func FromBytes[T Marshallable](payload []byte) (T, error) {
	var data T
	if err := json.Unmarshal(payload, &data); err != nil {
		return data, err
	}

	return data, nil
}
