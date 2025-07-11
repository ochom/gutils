package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/ochom/gutils/logs"
)

type ByteData []byte

// String converts byte data to string
func (b ByteData) String() string {
	if b == nil {
		return ""
	}
	return string(b)
}

// ToBytes converts provided interface to slice of bytes
func ToBytes(payload any) ByteData {
	stringValue, ok := payload.(string)
	if ok {
		return ByteData(stringValue)
	}

	numericValue, ok := payload.(float64)
	if ok {
		return ByteData(fmt.Appendf(nil, "%f", numericValue))
	}

	boolValue, ok := payload.(bool)
	if ok {
		if boolValue {
			return ByteData("true")
		}
		return ByteData("false")
	}

	if payload == nil {
		return nil
	}

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
