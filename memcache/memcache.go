package memcache

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"time"

	"github.com/ghongli/cache"
)

// Default prefix to prevent collision with other key stored in memcache
const PREFIX string = "mem_cache:"

type Cache struct {
	client *memcache.Client
	prefix string
}

func NewMemcacheStore(c *memcache.Client, prefix string) *Cache {
	if prefix == "" {
		prefix = PREFIX
	}
	return &Cache{client: c, prefix: prefix}
}

// Put put value to memcache.
func (mem *Cache) Put(key string, data []byte, expire time.Duration) error {
	item := &memcache.Item{
		Key:        mem.key(key),
		Value:      data,
		Expiration: int32(expire / time.Second),
	}
	return mem.client.Set(item)
}

// Get get value from memcache
func (mem *Cache) Get(key string) ([]byte, error) {
	val, err := mem.client.Get(mem.key(key))
	if err != nil {
		return nil, mem.adaptError(err)
	}

	return val.Value, nil
}

// Delete delete value from memcache.
func (mem *Cache) Delete(key string) error {
	return mem.adaptError(mem.client.Delete(mem.key(key)))
}

// IsExist check value exists in memcache.
func (mem *Cache) IsExist(key string) bool {
	_, err := mem.client.Get(mem.key(key))
	return !(err != nil)
}

// ClearAll clear all cached in memcache.
func (mem *Cache) ClearAll() error {
	return mem.client.FlushAll()
}

func (mem *Cache) adaptError(err error) error {
	switch err {
	case nil:
		return nil
	case memcache.ErrCacheMiss:
		return cache.ErrCacheMiss
	default:
		return err
	}
}

func (mem *Cache) key(key string) string {
	return fmt.Sprint(mem.prefix, key)
}

func init() {
	cache.Register("memcache", func() cache.Cache {
		return NewMemcacheStore(memcache.New("127.0.0.1:11211"), PREFIX)
	})
}
