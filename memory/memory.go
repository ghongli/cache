package memory

import (
	"sync"
	"time"

	"github.com/ghongli/cache"
)

var (
	// DefaultEvery means the clock time of recycling the expired cache items in memory.
	DefaultEvery = 60 // 1 minute
)

// in memory store
// it contains a RW locker for safe map storage.
type Cache struct {
	lock  sync.RWMutex
	items map[string][]byte
	s     cache.Serializer
}

func NewMemoryStore(gcInterval time.Duration) *Cache {
	m := &Cache{
		items: make(map[string][]byte),
		s:     cache.NewCacheSerializer(),
	}

	m.TrashGc(gcInterval)

	return m
}

func (m *Cache) Put(key string, data []byte, expire time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	item := &cache.CacheItem{
		CreatedTime: time.Now(),
		Expired:     expire,
		Data:        data,
	}

	bs, err := m.s.Serialized(item)
	if err != nil {
		return err
	}

	m.items[key] = bs

	return nil
}

func (m *Cache) Get(key string) ([]byte, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	data, ok := m.items[key]
	if !ok {
		return nil, cache.ErrCacheMiss
	}

	item := new(cache.CacheItem)
	err := m.s.DeSerialized(data, item)
	if err != nil {
		return nil, err
	}

	if item.IsExpired() {
		go m.Delete(key) // Prevent a deadlock since the mutex is still locked here
		return nil, cache.ErrCacheMiss
	}

	return item.Data, nil
}

func (m *Cache) Delete(key string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.items[key]; !ok {
		return cache.ErrCacheMiss
	}

	delete(m.items, key)

	return nil
}

func (m *Cache) ClearAll() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.items = make(map[string][]byte)

	return nil
}

func (m *Cache) IsExist(key string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if v, ok := m.items[key]; ok {
		item := new(cache.CacheItem)
		err := m.s.DeSerialized(v, item)
		if err != nil {
			return false
		}

		return !item.IsExpired()
	}

	return false
}

func (m *Cache) TrashGc(gcInterval time.Duration) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if gcInterval < 1*time.Second {
		gcInterval = time.Duration(DefaultEvery) * time.Second
	}

	if len(m.items) >= 1 {
		currentItem := new(cache.CacheItem)
		for k, v := range m.items {
			err := m.s.DeSerialized(v, currentItem)
			if err == nil && currentItem.IsExpired() {
				go m.Delete(k)
			}
		}
	}

	time.AfterFunc(gcInterval, func() {
		m.TrashGc(gcInterval)
	})
}
