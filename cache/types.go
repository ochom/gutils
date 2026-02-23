package cache

import (
	"time"
)

// CacheDriver represents the available cache backend types.
const (
	// Memory indicates in-memory cache storage
	Memory = iota
	// Redis indicates Redis cache storage
	Redis
)

// cacheItem represents a single item stored in the in-memory cache.
// It tracks the value, creation time, and expiration duration for automatic cleanup.
type cacheItem struct {
	value     []byte
	createdAt time.Time
	expiry    time.Duration
}

// expired checks if the cache item has exceeded its expiration duration.
// Returns false if expiry is 0 (no expiration).
func (c cacheItem) expired() bool {
	if c.expiry == 0 {
		return false
	}
	return time.Since(c.createdAt) > c.expiry
}

// Config holds configuration options for cache connections.
// Used primarily for Redis configuration.
//
// Example:
//
//	config := cache.Config{
//		Url:      "localhost:6379",
//		DbIndex:  0,
//		Password: "secret",
//	}
type Config struct {
	// Url is the Redis server address (host:port)
	Url string
	// DbIndex is the Redis database number (0-15)
	DbIndex int
	// Password is the Redis authentication password
	Password string
}
