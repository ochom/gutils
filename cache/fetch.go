package cache

import (
	"time"

	"github.com/ochom/gutils/jsonx"
)

// Fetch retrieves a typed value from cache or fetches and caches it when missing.
//
// When a cached value exists, it is returned immediately and the cache is refreshed
// asynchronously in the background using fetchFunc (stale-while-revalidate behavior).
// When no cached value exists, fetchFunc is executed synchronously and the result is
// stored for the specified duration before being returned.
//
// Example:
//
//	type Profile struct { Name string }
//	profile := cache.Fetch("user:42:profile", 5*time.Minute, func() Profile {
//		return fetchProfileFromDB(42)
//	})
func Fetch[T any](key string, duration time.Duration, fetchFunc func() T) T {
	fetchAndSet := func() []byte {
		result := fetchFunc()

		jsonData := jsonx.Encode(result)
		_ = Set(key, jsonData, duration)
		return jsonData
	}

	cached := Get(key)
	if cached != nil { /*  */
		go func() {
			_ = fetchAndSet()
		}()

		return jsonx.Decode[T](cached)
	}

	return jsonx.Decode[T](fetchAndSet())
}
