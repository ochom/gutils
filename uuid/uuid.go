// Package uuid provides unique identifier generation utilities.
//
// This package offers functions for generating:
//   - MongoDB-style 24-character hexadecimal ObjectIDs
//   - Short, human-readable IDs from integers
//
// Example usage:
//
//	// Generate a unique ID
//	id := uuid.New()
//	fmt.Println(id) // "507f1f77bcf86cd799439011" (24 hex characters)
//
//	// Generate a short ID from a sequence number
//	shortID := uuid.ShortID(12345)
//	fmt.Println(shortID) // "ABC123" (alphanumeric)
package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"io"
	"sync/atomic"
	"time"
)

// objectIDCounter provides incrementing values for ID uniqueness
var objectIDCounter = readRandomUint32()

// ObjectID is a 12-byte unique identifier, similar to MongoDB's ObjectID.
type ObjectID [12]byte

// New generates a new unique identifier as a 24-character hexadecimal string.
// The ID is time-based with random and counter components for uniqueness.
//
// The ID structure (12 bytes):
//   - Bytes 0-3: Unix timestamp (seconds)
//   - Bytes 4-8: Random/unique value
//   - Bytes 9-11: Incrementing counter
//
// Example:
//
//	// Generate IDs for database records
//	user := User{
//		ID:    uuid.New(),
//		Name:  "Alice",
//		Email: "alice@example.com",
//	}
//
//	// Generate multiple unique IDs
//	for i := 0; i < 10; i++ {
//		fmt.Println(uuid.New())
//	}
func New() string {
	id := new()
	return hex.EncodeToString(id[:])
}

// new creates a new ObjectID (internal function).
func new() ObjectID {
	var b [12]byte

	// Convert the time to a byte array
	binary.BigEndian.PutUint32(b[0:4], uint32(time.Now().Unix()))

	// Generate 5 random bytes
	unique := generateRandomByte()
	copy(b[4:9], unique[:])

	now := uint64(time.Now().UnixNano())
	binary.BigEndian.PutUint64(b[4:], now)

	// Generate 3 random bytes with counter
	putUint32(b[9:12], atomic.AddUint32(&objectIDCounter, 1))

	return b
}

func generateRandomByte() [5]byte {
	var b [5]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(err)
	}

	return b
}

func readRandomUint32() uint32 {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(err)
	}

	return (uint32(b[0]) << 0) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
}

func putUint32(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}
