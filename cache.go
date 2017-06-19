package cache

import (
	"errors"
	"fmt"
	"time"
)

const (
	EXPIRES_DEFAULT = time.Duration(0)
	EXPIRES_FOREVER = time.Duration(-1)
)

var (
	ErrCacheMiss                             = errors.New("Key not found")
	ErrCacheNotStored                        = errors.New("Data not stored")
	ErrCacheNotSupported                     = errors.New("Operation not supported")
	ErrCacheDataCannotBeIncreasedOrDecreased = errors.New(`
		Data isn't an integer/string type, it cannot be increased or decreased`)
)

// a cached piece of data
type CacheItem struct {
	Expired time.Duration
	Data    []byte
}

// Cache interface contains all behaviors for cache adapter.
type Cache interface {
	// get cached value by key.
	Get(key string) ([]byte, error)
	// set cached value with key and expire time.
	Put(key string, data []byte, expire time.Duration) error
	// delete cached value by key.
	Delete(key string) error
	// check if cached value exists or not.
	IsExist(key string) bool
	// clear all cache.
	ClearAll() error
}

// Some caches like redis automatically clear out the cache
// But for the filesystem and in memory, this cannot.
// Caches that have to manually clear out the cached data should implement this method.
// start trash gc routine based on config string settings.
type GarbageCollector interface {
	StartAndTrashGc(config string) error
	//TrashGc(interval time.Duration)
}

// Store is a function create a new Cache Instance
type Store func() Cache

var adapters = make(map[string]Store)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if adapter is nil,
// it panics.
func Register(name string, adapter Store) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}

	if _, ok := adapters[name]; !ok {
		panic("cache: Register called twice for adapter " + name)
	}

	adapters[name] = adapter
}

// NewCache Create a new cache driver by adapter name and config string.
// config need to be correct JSON as string: {"interval":360}.
// it will start gc automatically.
func NewCache(adapterName, config string) (cache Cache, err error) {
	storeFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return nil, err
	}

	cache = storeFunc()
	return cache, nil
}
