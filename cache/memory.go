package cache

import (
	"fmt"
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
	c := &memoryCache{
		cacheWorkers: env.Int("CACHE_TOTAL_WORKERS", 10),
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
