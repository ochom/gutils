package helpers

import (
	"encoding/json"

	"github.com/ochom/gutils/logs"
)

// ToJSON converts a struct to JSON
func ToJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		logs.Error("Failed to marshal JSON: %s", err.Error())
		return nil
	}

	return b
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
