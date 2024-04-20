package internal

import (
	"context"
	"github.com/eatmoreapple/env"
	"github.com/go-redis/redis/v8"
	"sync"
)

type Cache interface {
	Add(ctx context.Context, key string, value string) error
	GetAll(ctx context.Context, key string) ([]string, error)
	DelAll(ctx context.Context, key string) error
}

type memoryCache struct {
	*sync.Map
}

func (m *memoryCache) Add(ctx context.Context, key string, value string) error {
	// fixme: not atomic
	result, err := m.GetAll(ctx, key)
	if err != nil {
		return err
	}
	result = append(result, value)
	m.Map.Store(key, result)
	return nil
}

func (m *memoryCache) GetAll(_ context.Context, key string) ([]string, error) {
	values, exists := m.Map.Load(key)
	if !exists {
		return nil, nil
	}
	r, ok := values.([]string)
	if !ok {
		return nil, nil
	}
	return r, nil
}

func (m *memoryCache) DelAll(_ context.Context, key string) error {
	m.Map.Delete(key)
	return nil
}

func NewMemoryCache(item *sync.Map) Cache {
	return &memoryCache{Map: item}
}

type redisCache struct {
	client *redis.Client
}

func (r *redisCache) Add(ctx context.Context, key string, value string) error {
	return r.client.RPush(ctx, key, value).Err()
}

func (r *redisCache) GetAll(ctx context.Context, key string) ([]string, error) {
	return r.client.LRange(ctx, key, 0, -1).Result()
}

func (r *redisCache) DelAll(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func NewRedisCache(client *redis.Client) Cache {
	return &redisCache{client: client}
}

func NewRedisCacheFromAddr(addr string) Cache {
	client := &redis.Options{
		Network:  "tcp",
		Addr:     addr,
		PoolSize: 10,
	}
	return NewRedisCache(redis.NewClient(client))
}

// CacheFromEnv creates a new Cache based on the REDIS_ADDR environment variable.
// todo set key expiration time
func CacheFromEnv() Cache {
	var cache Cache
	if addr := env.Name("REDIS_ADDR").String(); len(addr) > 0 {
		cache = NewRedisCacheFromAddr(addr)
	} else {
		cache = NewMemoryCache(new(sync.Map))
	}
	return cache
}
