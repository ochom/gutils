package cache

import (
	"time"
)

const (
	Memory = iota
	Redis
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

// Config ...
type Config struct {
	Url      string
	DbIndex  int
	Password string
}
