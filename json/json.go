package json

import (
	baseJSON "encoding/json"

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

// Encode encodes the given payload into JSON format in byte slice
func Encode(payload any) byteData {
	bytesPayload, err := baseJSON.Marshal(&payload)
	if err != nil {
		logs.Error("Failed to marshal JSON: %s", err.Error())
		return nil
	}

	return bytesPayload
}

// Decode decodes JSON payload into the specified type
func Decode[T any](payload []byte) T {
	var data T
	if err := baseJSON.Unmarshal(payload, &data); err != nil {
		logs.Error("Failed to unmarshal JSON: %s", err.Error())
		return data
	}

	return data
}
