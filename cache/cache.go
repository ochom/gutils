package cache

import (
	"sync"
	"time"
)

var memoryCache map[string]cacheItem
var mut sync.Mutex

// V  is the type of the value to be stored in the cache
type V []byte

// cacheItem ...
type cacheItem struct {
	key       string
	value     V
	createdAt time.Time
	expiry    time.Duration
	callBack  func()
}

func init() {
	memoryCache = make(map[string]cacheItem)
	go CleanUp()
}

// Set ...
func Set(key string, value V) {
	SetWithExpiry(key, value, 0)
}

// SetWithExpiry ...
func SetWithExpiry(key string, value V, expiry time.Duration) {
	SetWithCallback(key, value, expiry, nil)
}

// SetWithCallback ...
func SetWithCallback(key string, value V, expiry time.Duration, callback func()) {
	mut.Lock()
	defer mut.Unlock()
	item := cacheItem{
		key:       key,
		value:     value,
		createdAt: time.Now(),
		expiry:    expiry,
		callBack:  callback,
	}
	memoryCache[key] = item
}

// Get ...
func Get(key string) V {
	item, ok := memoryCache[key]
	if !ok {
		return nil
	}

	return item.value
}

// CleanUp deletes expired cache items and calls their callbacks
func CleanUp() {
	tk := time.NewTicker(time.Millisecond * 20)
	for range tk.C {
		mut.Lock()
		for key, item := range memoryCache {
			if item.expiry > 0 && time.Since(item.createdAt) > item.expiry {
				delete(memoryCache, key)
				if item.callBack != nil {
					item.callBack()
				}
			}
		}
		mut.Unlock()
	}
}
