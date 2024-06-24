package cache

import (
	"fmt"
	"time"

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
	conn = newMemoryCache()

	go conn.cleanUp()
}

// NewCache ...
func Init(driver CacheDriver, url ...string) error {
	if driver == Memory {
		// cache is already running return nil
		return nil
	}

	cn, err := newRedisCache(url...)
	if err != nil {
		return err
	}

	conn = cn

	go conn.cleanUp()
	return nil
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
