package gson

import (
	"encoding/json"

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

// Marshal converts provided object or array to slice of bytes
func Marshal(payload any) byteData {
	bytesPayload, err := json.Marshal(&payload)
	if err != nil {
		logs.Error("Failed to marshal JSON: %s", err.Error())
		return nil
	}

	return bytesPayload
}

// Unmarshal converts slice of bytes to provided object or array
func Unmarshal[T any](payload []byte) T {
	var data T
	if err := json.Unmarshal(payload, &data); err != nil {
		logs.Error("Failed to unmarshal JSON: %s", err.Error())
		return data
	}

	return data
}
