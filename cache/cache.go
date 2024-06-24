package cache

import (
	"fmt"
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
	cleanUp()
}

var conn Cache

// default to memory cache
func init() {
	driver := env.Get("CACHE_DRIVER", "memory")
	if driver == "memory" {
		conn = newMemoryCache()
	} else {
		url := env.Get("REDIS_URL", "localhost:6379")
		password := env.Get("REDIS_PASSWORD", "")
		dbIndex := env.Int("REDIS_DB_INDEX", 0)
		con, err := newRedisCache(&Config{
			Url:      url,
			DbIndex:  dbIndex,
			Password: password,
		})

		if err != nil {
			panic(err)
		}

		conn = con
	}

	go conn.cleanUp()
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
func Get[T any](key string) (T, error) {
	v := conn.get(key)
	if v == nil {
		return *new(T), fmt.Errorf("key %s not found", key)
	}

	return helpers.FromBytes[T](v), nil
}

// Delete ...
func Delete(key string) error {
	return conn.delete(key)
}

// CleanUp ...
func CleanUp() {
	tick := time.NewTicker(time.Second)
	for range tick.C {
		conn.cleanUp()
	}
}
