package gson

import (
	"encoding/json"
	"fmt"

	"github.com/ochom/gutils/logs"
)

type byteData []byte

// String converts byte data to string
func (b byteData) String() string {
	if b == nil {
		return ""
	}
	return string(b)
}

// Bytes converts byte data to slice of bytes
func (b byteData) Bytes() []byte {
	return b
}

// Marshal converts provided interface to slice of bytes
func Marshal(payload any) byteData {
	stringValue, ok := payload.(string)
	if ok {
		return byteData(stringValue)
	}

	numericValue, ok := payload.(float64)
	if ok {
		return byteData(fmt.Appendf(nil, "%f", numericValue))
	}

	boolValue, ok := payload.(bool)
	if ok {
		if boolValue {
			return byteData("true")
		}
		return byteData("false")
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

// Unmarshal converts slice of bytes to provided interface
func Unmarshal[T any](payload []byte) T {
	var data T
	if err := json.Unmarshal(payload, &data); err != nil {
		logs.Error("Failed to unmarshal JSON: %s", err.Error())
		return data
	}

	return data
}
