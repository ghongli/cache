package redis

import (
	"fmt"
	"time"

	"github.com/ghongli/cache"
	"github.com/go-redis/redis"
)

// Default prefix to prevent collision with other key stored in redis
const PREFIX string = "r_cache:"

type Cache struct {
	client *redis.Client
	prefix string
}

func NewRedisCache(opts *redis.Options, prefix string) *Cache {
	if prefix == "" {
		prefix = PREFIX
	}
	return &Cache{redis.NewClient(opts), prefix}
}

// Put put cache to redis.
func (r *Cache) Put(key string, data []byte, expire time.Duration) error {
	return r.client.Set(r.key(key), data, expire).Err()
}

// Get cache from redis
func (r *Cache) Get(key string) ([]byte, error) {
	val, err := r.client.Get(r.key(key)).Bytes()
	if err != nil {
		return nil, err
	}

	return val, err
}

// IsExist check cache's existence in redis.
func (r *Cache) IsExist(key string) bool {
	val, err := r.client.Exists(r.key(key)).Result()
	if err != nil {
		return false
	}

	if val == -1 {
		return false
	}

	return true
}

// Delete delete cache in redis.
func (r *Cache) Delete(key string) error {
	return r.client.Del(r.key(key)).Err()
}

// ClearAll clean all cache in redis.
func (r *Cache) ClearAll() error {
	return r.client.FlushDb().Err()
}

func (r *Cache) key(key string) string {
	return fmt.Sprint(r.prefix, key)
}

func init() {
	cache.Register("redis", func() cache.Cache {
		return NewRedisCache(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		}, PREFIX)
	})
}
