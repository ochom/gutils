package cache

import (
	"sync"
	"time"

	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/logs"
)

var cacheWorkers int
var memoryCache map[string]cacheItem
var mut sync.Mutex

// V  is the type of the value to be stored in the cache
type V []byte

// cacheItem ...
type cacheItem struct {
	value     []byte
	createdAt time.Time
	expiry    time.Duration
	callBack  func()
}

// expired ...
func (c cacheItem) expired() bool {
	return time.Since(c.createdAt) > c.expiry
}

func init() {
	cacheWorkers = env.Int("CACHE_TOTAL_WORKERS", 10)
	memoryCache = make(map[string]cacheItem)
	go CleanUp()
}

// Set ...
func Set(key string, value V) {
	expiry := time.Hour * time.Duration(env.Int("MAX_CACHE_HOUR", 24))
	SetWithExpiry(key, value, expiry)
}

// SetWithExpiry ...
func SetWithExpiry(key string, value V, expiry time.Duration) {
	callback := func() {
		logs.Info("Session item expired: %s", key)
	}
	SetWithCallback(key, value, expiry, callback)
}

// SetWithCallback ...
func SetWithCallback(key string, value V, expiry time.Duration, callback func()) {
	mut.Lock()
	defer mut.Unlock()
	item := cacheItem{
		value:     value,
		createdAt: time.Now(),
		expiry:    expiry,
		callBack:  callback,
	}
	memoryCache[key] = item
}

// Get ...
func Get(key string) V {
	mut.Lock()
	defer mut.Unlock()

	item, ok := memoryCache[key]
	if !ok {
		return nil
	}

	return item.value
}

// Delete ...
func Delete(key string) {
	mut.Lock()
	defer mut.Unlock()
	delete(memoryCache, key)
}

// CleanUp deletes expired cache items and calls their callbacks
func CleanUp() {
	jobs := make(chan cacheItem, 100)

	//  start workers
	for i := 0; i < cacheWorkers; i++ {
		go runCallbacks(jobs)
	}

	tick(jobs)
}

// tick to run the ticker
func tick(jobs chan<- cacheItem) {
	tk := time.NewTicker(time.Second) // every second
	for range tk.C {
		// acquire
		mut.Lock()
		for key, item := range memoryCache {
			if item.expired() {
				jobs <- item
				delete(memoryCache, key)
			}
		}

		// release
		mut.Unlock()
	}
}

func runCallbacks(jobs <-chan cacheItem) {
	for item := range jobs {
		item.callBack()
	}
}
