package cache

import (
	"runtime"
	"time"

	"github.com/gofiber/storage/redis/v3"
	"github.com/ochom/gutils/helpers"
)

var store *redis.Storage

func init() {
	HOST := helpers.GetEnv("REDIS_HOST", "127.0.0.1")
	PORT := helpers.GetEnvInt("REDIS_PORT", 6379)
	USERNAME := helpers.GetEnv("REDIS_USERNAME", "")
	PASSWORD := helpers.GetEnv("REDIS_PASSWORD", "")
	DATABASE := helpers.GetEnvInt("REDIS_DATABASE", 0)

	store = initCache(HOST, PORT, USERNAME, PASSWORD, DATABASE)
}

// initCache creates a new cache instance
func initCache(host string, port int, username, password string, database int) *redis.Storage {
	return redis.New(redis.Config{
		Host:      host,
		Port:      port,
		Username:  username,
		Password:  password,
		Database:  database,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})
}

// Set a key-value pair in the cache
func Set(key string, value []byte, exp time.Duration) error {
	return store.Set(key, value, exp)
}

// Get a key from the cache
func Get(key string) ([]byte, error) {
	return store.Get(key)
}

// Has a key in the cache
func Has(key string) bool {
	val, err := store.Get(key)
	if err != nil || len(val) == 0 {
		return false
	}

	return true
}

// Delete a key from the cache
func Delete(key string) error {
	return store.Delete(key)
}

// Clear the cache
func Clear() error {
	return store.Reset()
}

// Close the cache
func Close() error {
	return store.Close()
}
