package memcache

import (
	"flag"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ghongli/cache"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	_          = &Cache{}
	store      *Cache
	key        = "test"
	sampleData = []byte("memcache store")
)

func TestMain(t *testing.M) {
	//store = NewMemcacheStore(memcache.New("127.0.0.1:11211"), "memTest:")
	store = NewMemcacheStore(memcache.New("101.201.197.163:11211"), "")
	flag.Parse()
	os.Exit(t.Run())
}

func TestCache_Put(t *testing.T) {
	err := store.Put(key, sampleData, time.Minute*1)
	if err != nil {
		t.Fatalf(`An error occurred while trying to add
		some data to memcached.. \n %v`, err)
	}
}

func TestCache_Get(t *testing.T) {
	val, err := store.Get(key)
	if err != nil {
		t.Fatalf(`Could not get item with key %s when it in fact
			exists... \n %v`, key, err)
	}

	if !reflect.DeepEqual(sampleData, val) {
		t.Fatalf(`Expected %v \n ..Got %v instead`, sampleData, val)
	}
}

func TestCache_Delete(t *testing.T) {
	err := store.Delete(key)
	if err != nil {
		t.Fatalf(`The key %s could not be deleted ... %v`, key, err)
	}
}

func TestCache_ClearAll(t *testing.T) {
	err := store.ClearAll()
	// If nil, we accept the cache is flushed..
	if err != nil {
		t.Fatalf(`An error occurred while clearing the memcached db.. %v`, err)
	}

	_, err = store.Get("name")
	if err != cache.ErrCacheMiss {
		t.Fatal("All data should have been cleared from the cache")
	}
}

func TestCache_IsExist(t *testing.T) {
	val := store.IsExist(key)
	if val {
		t.Fatalf("an item with the key is exist ? expected %v.. \n %v instead.", false, true)
	}
}

func TestCache_adaptError(t *testing.T) {
	if err := store.adaptError(nil); err != nil {
		t.Fatalf("Expected %v.. Got %v", nil, err)
	}
}

func TestCache_Prefix(t *testing.T) {
	m := NewMemcacheStore(memcache.New("127.0.0.1:11211"), "")
	if !reflect.DeepEqual(m.prefix, PREFIX) {
		t.Fatalf(`Prefix doen't match.
			\n Expected %s \n.. Got %s`, PREFIX, m.prefix)
	}
}
