package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/ochom/gutils/env"
	"github.com/redis/go-redis/v9"
)

// memoryCache implements the Cache interface using an in-memory map.
// It provides thread-safe operations and automatic cleanup of expired items.
//
// The memory cache is useful for:
//   - Development and testing without Redis
//   - Single-instance applications
//   - Fallback when Redis is unavailable
//
// Configuration:
//   - CACHE_TOTAL_WORKERS: Number of cleanup workers (default: 10)
type memoryCache struct {
	cacheWorkers int
	items        map[string]cacheItem
	mut          sync.Mutex
}

// newMemoryCache creates a new in-memory cache instance.
// It starts a background goroutine that periodically cleans up expired items.
func newMemoryCache() Cache {
	c := &memoryCache{
		cacheWorkers: env.Get("CACHE_TOTAL_WORKERS", 10),
		items:        make(map[string]cacheItem),
		mut:          sync.Mutex{},
	}

	go c.cleanUp()
	return c
}

// getClient ...
func (m *memoryCache) getClient() *redis.Client {
	return nil
}

// setWithExpiry ...
func (m *memoryCache) set(key string, value []byte, expiry time.Duration) error {
	m.mut.Lock()
	defer m.mut.Unlock()

	item := cacheItem{
		value:     value,
		createdAt: time.Now(),
		expiry:    expiry,
	}
	m.items[key] = item

	return nil
}

// get ...
func (m *memoryCache) get(key string) []byte {
	m.mut.Lock()
	defer m.mut.Unlock()

	item, ok := m.items[key]
	if !ok {
		return nil
	}

	return item.value
}

// delete ...
func (m *memoryCache) delete(key string) error {
	m.mut.Lock()
	defer m.mut.Unlock()

	if _, ok := m.items[key]; !ok {
		return fmt.Errorf("key %s not found", key)
	}

	delete(m.items, key)
	return nil
}

// cleanUp deletes expired cache items and calls their callbacks
func (m *memoryCache) cleanUp() {
	for {
		m.mut.Lock()
		for key, item := range m.items {
			if item.expired() {
				delete(m.items, key)
			}
		}

		m.mut.Unlock()
		<-time.After(time.Second)
	}
}
