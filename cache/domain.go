package cache

import (
	"time"
)

// CacheDriver ...
type CacheDriver string

const (
	RedisDriver  CacheDriver = "redis"
	MemoryDriver CacheDriver = "memory"
)

// V  is the type of the value to be stored in the cache
type V []byte

// cacheItem ...
type cacheItem struct {
	value     []byte
	createdAt time.Time
	expiry    time.Duration
}

// expired ...
func (c cacheItem) expired() bool {
	return time.Since(c.createdAt) > c.expiry
}
