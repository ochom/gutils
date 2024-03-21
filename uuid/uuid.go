package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"io"
	"sync/atomic"
	"time"
)

var objectIDCounter = readRandomUint32()

// ObjectID ...
type ObjectID [12]byte

// New ...
func New() string {
	id := new()
	return hex.EncodeToString(id[:])
}

// new ...
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
