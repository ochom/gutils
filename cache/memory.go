package cache

import (
	"sync"
	"time"

	"github.com/ochom/gutils/env"
	"github.com/redis/go-redis/v9"
)

// memoryCache implements Cache
type memoryCache struct {
	cacheWorkers int
	items        map[string]cacheItem
	mut          sync.Mutex
}

func newMemoryCache() Cache {
	return &memoryCache{
		cacheWorkers: env.Int("CACHE_TOTAL_WORKERS", 10),
		items:        make(map[string]cacheItem),
		mut:          sync.Mutex{},
	}
}

// getClient ...
func (m *memoryCache) getClient() *redis.Client {
	return nil
}

// set ...
func (m *memoryCache) set(key string, value []byte) {
	expiry := time.Hour * time.Duration(env.Int("MAX_CACHE_HOUR", 24))
	m.setWithExpiry(key, value, expiry)
}

// setWithExpiry ...
func (m *memoryCache) setWithExpiry(key string, value []byte, expiry time.Duration) {
	m.mut.Lock()
	defer m.mut.Unlock()
	item := cacheItem{
		value:     value,
		createdAt: time.Now(),
		expiry:    expiry,
	}
	m.items[key] = item
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
func (m *memoryCache) delete(key string) {
	m.mut.Lock()
	defer m.mut.Unlock()
	delete(m.items, key)
}

// cleanUp deletes expired cache items and calls their callbacks
func (m *memoryCache) cleanUp() {
	m.mut.Lock()
	defer m.mut.Unlock()
	for key, item := range m.items {
		if item.expired() {
			delete(m.items, key)
		}
	}
}
