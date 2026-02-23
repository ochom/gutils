// Package cache provides a unified caching interface with support for Redis and in-memory storage.
//
// The package automatically selects the cache backend based on the CACHE_DRIVER environment variable:
//   - "redis": Uses Redis as the cache backend (default)
//   - "memory": Uses an in-memory cache with automatic cleanup
//
// For Redis, the following environment variables are supported:
//   - REDIS_HOST: Redis server host
//   - REDIS_PORT: Redis server port
//   - REDIS_URL: Full Redis URL (alternative to HOST:PORT)
//   - REDIS_PASSWORD: Redis password (optional)
//   - REDIS_DB_INDEX: Redis database index (default: 0)
//
// Example usage:
//
//	import (
//		"time"
//		"github.com/ochom/gutils/cache"
//		"github.com/ochom/gutils/jsonx"
//	)
//
//	// Store a string value
//	cache.Set("user:123:name", []byte("John Doe"), 5*time.Minute)
//
//	// Retrieve a value
//	name := cache.Get("user:123:name")
//	fmt.Println(string(name)) // "John Doe"
//
//	// Store a struct as JSON
//	type User struct { Name string; Email string }
//	user := User{Name: "Alice", Email: "alice@example.com"}
//	cache.Set("user:456", jsonx.Encode(user), time.Hour)
//
//	// Delete a key
//	cache.Delete("user:123:name")
package cache

import (
	"time"

	"github.com/ochom/gutils/env"
	"github.com/redis/go-redis/v9"
)

// Cache defines the interface for cache operations.
// Both Redis and in-memory implementations satisfy this interface.
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

// Client returns the underlying Redis client if available.
// Returns nil when using in-memory cache.
//
// This is useful for advanced Redis operations not covered by the simple cache interface.
//
// Example:
//
//	// Direct Redis operations
//	rdb := cache.Client()
//	if rdb != nil {
//		// Use Redis-specific features
//		rdb.HSet(ctx, "hash:key", "field", "value")
//	}
func Client() *redis.Client {
	return conn.getClient()
}

// Set stores a value in the cache with the specified key and expiration duration.
//
// The value is stored as bytes, so complex types should be serialized (e.g., using jsonx.Encode).
// Set expiry to 0 for no expiration (memory cache only; Redis requires explicit expiration).
//
// Example:
//
//	// Cache a simple string
//	cache.Set("greeting", []byte("Hello, World!"), 10*time.Minute)
//
//	// Cache JSON data
//	data := map[string]int{"count": 42}
//	cache.Set("stats", jsonx.Encode(data), time.Hour)
//
//	// Cache with no expiration (memory cache)
//	cache.Set("permanent", []byte("data"), 0)
func Set(key string, value []byte, expiry time.Duration) error {
	return conn.set(key, value, expiry)
}

// Get retrieves a value from the cache by key.
// Returns nil if the key does not exist or has expired.
//
// Example:
//
//	data := cache.Get("my-key")
//	if data == nil {
//		// Key not found or expired
//		data = fetchFromDatabase()
//		cache.Set("my-key", data, 5*time.Minute)
//	}
//
//	// Parse JSON data
//	var user User
//	if data := cache.Get("user:123"); data != nil {
//		user = jsonx.Decode[User](data)
//	}
func Get(key string) []byte {
	return conn.get(key)
}

// Delete removes a value from the cache by key.
// Returns an error if the key does not exist (memory cache only).
//
// Example:
//
//	// Invalidate user cache after update
//	updateUser(user)
//	cache.Delete("user:" + user.ID)
//
//	// Clear session data
//	cache.Delete("session:" + sessionID)
func Delete(key string) error {
	return conn.delete(key)
}
