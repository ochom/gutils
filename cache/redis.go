package cache

import (
	"context"
	"time"

	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/logs"
	"github.com/redis/go-redis/v9"
)

// redisCache implements Cache
type redisCache struct {
	client *redis.Client
}

func newRedisCache() Cache {
	cl := redis.NewClient(&redis.Options{
		Addr:     env.Get("REDIS_URL", "localhost:6379"),
		Password: env.Get("REDIS_PASSWORD", ""),
		DB:       env.Int("REDIS_DB_INDEX", 0),
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
func (r *redisCache) set(key string, value []byte) error {
	return r.setWithExpiry(key, value, 0)
}

// setWithExpiry ...
func (r *redisCache) setWithExpiry(key string, value []byte, expiry time.Duration) error {
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

func (r *redisCache) cleanUp() {
	// TODO
}
