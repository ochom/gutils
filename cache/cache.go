package cache

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache ...
type Cache interface {
	getClient() *redis.Client
	set(key string, value V)
	setWithExpiry(key string, value V, expiry time.Duration)
	get(key string) V
	delete(key string)
	cleanUp()
}

var conn Cache

// NewCache ...
func Init(driver CacheDriver, url ...string) {
	switch driver {
	case RedisDriver:
		conn = newRedisCache(url...)
	case MemoryDriver:
		conn = newMemoryCache()
	}

	go conn.cleanUp()
}

// Client ...
func Client() *redis.Client {
	return conn.getClient()
}

// Set ...
func Set(key string, value V) {
	conn.set(key, value)
}

// SetWithExpiry ...
func SetWithExpiry(key string, value V, expiry time.Duration) {
	conn.setWithExpiry(key, value, expiry)
}

// Get ...
func Get(key string) V {
	return conn.get(key)
}

// Delete ...
func Delete(key string) {
	conn.delete(key)
}

// CleanUp ...
func CleanUp() {
	tick := time.NewTicker(time.Second)
	for range tick.C {
		conn.cleanUp()
	}
}
