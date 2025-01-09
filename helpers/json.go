package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/ochom/gutils/logs"
)

// ToBytes converts provided interface to slice of bytes
func ToBytes[T any](payload T) []byte {
	if isComparable(payload) {
		panic(fmt.Sprintf("cannot marshal to given type: %T", payload))
	}

	bytesPayload, err := json.Marshal(&payload)
	if err != nil {
		logs.Error("Failed to marshal JSON: %s", err.Error())
		return nil
	}

	return bytesPayload
}

// FromBytes converts slice of bytes to provided interface
func FromBytes[T any](payload []byte) (T, error) {
	var data T
	if isComparable(data) {
		return data, errors.New("cannot unmarshal to given type")
	}

	if err := json.Unmarshal(payload, &data); err != nil {
		return data, err
	}

	return data, nil
}

// isComparable checks for comparable and []byte
func isComparable(a any) bool {
	if _, ok := a.([]byte); ok {
		return true
	}

	switch reflect.TypeOf(a).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}
