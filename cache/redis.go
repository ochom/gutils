package cache

import (
	"context"
	"errors"
	"time"

	"github.com/ochom/gutils/logs"
	"github.com/redis/go-redis/v9"
)

// redisCache implements Cache
type redisCache struct {
	client *redis.Client
}

func newRedisCache(url ...string) (Cache, error) {
	if len(url) == 0 {
		logs.Error("newRedisCache: url is empty")
		return nil, errors.New("url is empty")
	}

	opt, err := redis.ParseURL(url[0])
	if err != nil {
		logs.Error("newRedisCache: %s", err.Error())
		return nil, err
	}

	cl := redis.NewClient(opt)
	return &redisCache{
		client: cl,
	}, nil
}

// getClient ...
func (r *redisCache) getClient() *redis.Client {
	return r.client
}

// set ...
func (r *redisCache) set(key string, value []byte) {
	r.setWithExpiry(key, value, 0)
}

// setWithExpiry ...
func (r *redisCache) setWithExpiry(key string, value []byte, expiry time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	if err := r.client.Set(ctx, key, value, expiry).Err(); err != nil {
		logs.Error("setWithCallback: %s", err.Error())
	}
}

// get ...
func (r *redisCache) get(key string) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	v, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	return []byte(v)
}

// delete ...
func (r *redisCache) delete(key string) {
	if err := r.client.Del(context.Background(), key).Err(); err != nil {
		logs.Error("delete: %s", err.Error())
	}
}

func (r *redisCache) cleanUp() {
	// TODO
}
