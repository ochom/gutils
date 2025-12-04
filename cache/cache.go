package cache

import (
	"time"

	"github.com/ochom/gutils/env"
	"github.com/redis/go-redis/v9"
)

// Cache ...
type Cache interface {
	getClient() *redis.Client
	set(key string, value []byte, expiry time.Duration) error
	get(key string) []byte
	delete(key string) error
}

var conn Cache

// default to memory cache
func init() {
	switch env.Get("CACHE_DRIVER", "redis") {
	case "memory":
		conn = newMemoryCache()
	case "redis":
		conn = newRedisCache()
	default:
		conn = newRedisCache()
	}

}

// Client ...
func Client() *redis.Client {
	return conn.getClient()
}

// Set ...
func Set(key string, value []byte, expiry time.Duration) error {
	return conn.set(key, value, expiry)
}

// Get returns the value for the given key.
func Get(key string) []byte {
	return conn.get(key)
}

// Delete removes the value for the given key.
func Delete(key string) error {
	return conn.delete(key)
}
