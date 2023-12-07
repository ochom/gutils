package cache

import (
	"time"

	"github.com/gofiber/storage/sqlite3/v2"
	"github.com/ochom/gutils/helpers"
)

var store *sqlite3.Storage

func init() {
	database := helpers.GetEnv("CACHE_DB_PATH", "./fiber.sqlite3")
	table := helpers.GetEnv("CACHE_TABLE_NAME", "fiber_storage")
	store = initCache(database, table)
}

// initCache creates a new cache instance
func initCache(dbPath, tableName string) *sqlite3.Storage {
	return sqlite3.New(sqlite3.Config{
		Database:        dbPath,
		Table:           tableName,
		Reset:           false,
		GCInterval:      10 * time.Second,
		MaxOpenConns:    100,
		MaxIdleConns:    100,
		ConnMaxLifetime: 1 * time.Second,
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
