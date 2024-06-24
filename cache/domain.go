package cache

import (
	"time"
)

// CacheDriver ...
type CacheDriver string

const (
	Redis  CacheDriver = "redis"
	Memory CacheDriver = "memory"
)

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
