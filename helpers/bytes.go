package helpers

import "encoding/json"

// toBytes ...
func ToBytes(payload any) []byte {
	if payload == nil {
		return nil
	}

	// Check if the payload is already of type []byte.
	if bytesPayload, ok := payload.([]byte); ok {
		return bytesPayload
	}

	// Check if the payload is already of type string.
	if stringPayload, ok := payload.(*string); ok {
		return []byte(*stringPayload)
	}

	// Marshal the payload to JSON.
	bytesPayload, err := json.Marshal(payload)
	if err != nil {
		return nil
	}

	return bytesPayload
}

// FromBytes ...
func FromBytes[T any](payload []byte) (T, error) {
	var data T
	if payload == nil {
		return data, nil
	}

	err := json.Unmarshal(payload, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
