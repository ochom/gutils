package cache

import (
	"context"
	"time"

	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/logs"
	"github.com/redis/go-redis/v9"
)

// redisCache implements the Cache interface using Redis as the backend.
// It provides distributed caching suitable for multi-instance deployments.
//
// Configuration is read from environment variables:
//   - REDIS_HOST: Redis server hostname
//   - REDIS_PORT: Redis server port
//   - REDIS_URL: Full Redis URL (alternative to HOST:PORT)
//   - REDIS_PASSWORD: Authentication password
//   - REDIS_DB_INDEX: Database index (default: 0)
//
// If Redis connection fails, it automatically falls back to memory cache.
type redisCache struct {
	client *redis.Client
}

// newRedisCache creates a new Redis cache instance.
// Falls back to memory cache if Redis is unavailable.
func newRedisCache() Cache {
	host := env.Get[string]("REDIS_HOST")
	port := env.Get[string]("REDIS_PORT")

	url := ""
	if host != "" && port != "" {
		url = host + ":" + port
		goto initRedis
	}

	if url == "" {
		url = env.Get[string]("REDIS_URL")
	}

	if url == "" {
		logs.Error("newRedisCache: REDIS_HOST, REDIS_PORT or REDIS_URL must be set")
		return newMemoryCache()
	}

initRedis:
	cl := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: env.Get[string]("REDIS_PASSWORD"),
		DB:       env.Get("REDIS_DB_INDEX", 0),
	})

	if err := cl.Ping(context.Background()).Err(); err != nil {
		logs.Error("newRedisCache: %s", err.Error())
		return newMemoryCache()
	}

	logs.Info("Connected to redis")
	return &redisCache{client: cl}
}

// getClient ...
func (r *redisCache) getClient() *redis.Client {
	return r.client
}

// set ...
func (r *redisCache) set(key string, value []byte, expiry time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	if err := r.client.Set(ctx, key, value, expiry).Err(); err != nil {
		logs.Error("setWithCallback: %s", err.Error())
		return err
	}

	return nil
}

// get ...
func (r *redisCache) get(key string) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	v, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	return []byte(v)
}

// delete ...
func (r *redisCache) delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	if err := r.client.Del(ctx, key).Err(); err != nil {
		logs.Error("delete: %s", err.Error())
		return err
	}

	return nil
}
