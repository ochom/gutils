package cache

import (
	"time"

	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/helpers"
	"github.com/redis/go-redis/v9"
)

// Cache ...
type Cache interface {
	getClient() *redis.Client
	set(key string, value []byte) error
	setWithExpiry(key string, value []byte, expiry time.Duration) error
	get(key string) []byte
	delete(key string) error
}

var conn Cache

// default to memory cache
func init() {
	switch env.Int("CACHE_DRIVER", 0) {
	case Memory:
		conn = newMemoryCache()
	case Redis:
		conn = newRedisCache()
	default:
		panic("unknown cache driver")
	}

}

// Client ...
func Client() *redis.Client {
	return conn.getClient()
}

// Set ...
func Set[T any](key string, value T) error {
	return conn.set(key, helpers.ToBytes(value))
}

// SetWithExpiry ...
func SetWithExpiry[T any](key string, value T, expiry time.Duration) error {
	return conn.setWithExpiry(key, helpers.ToBytes(value), expiry)
}

// Get ...
func Get[T any](key string) T {
	if v := conn.get(key); v != nil {
		return helpers.FromBytes[T](v)
	}

	return *new(T)
}

// Delete ...
func Delete(key string) error {
	return conn.delete(key)
}
