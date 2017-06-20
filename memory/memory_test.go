package memory

import (
	"flag"
	"os"
	"testing"
	"time"

	"errors"
	"github.com/ghongli/cache"
	"reflect"
)

var (
	_           cache.Cache = &Cache{}
	memoryStore *Cache
	sampleData  = []byte("memory store")
	key         = "memory"
)

func TestMain(m *testing.M) {
	memoryStore = NewMemoryStore(1 * time.Second)

	flag.Parse()
	os.Exit(m.Run())
}

func TestCache_Put(t *testing.T) {
	err := memoryStore.Put(key, sampleData, time.Minute*1)
	if err != nil {
		t.Fatalf("Data could not be stored in the memory store.. \n%v", err)
	}
}

type mockSerializeHelper struct {
}

func (mock *mockSerializeHelper) Serialized(i interface{}) ([]byte, error) {
	return nil, errors.New("Yup an error occurred")
}

func (mock *mockSerializeHelper) DeSerialized(data []byte, i interface{}) error {
	return errors.New("Yet another error")
}

func TestCache_SetError(t *testing.T) {
	m := &Cache{s: &mockSerializeHelper{}, items: make(map[string][]byte, 10)}
	err := m.Put("n", []byte("ERROR"), time.Second*2)
	if err == nil {
		t.Fatalf(`Error should be not nil as the item could not be marshalled into
			bytes.. Got %v`, err)
	}
}

func TestCache_Get(t *testing.T) {
	val, err := memoryStore.Get(key)
	if err != nil {
		t.Fatalf("Key %s should exist in the store... \n %v", key, err)
	}

	if !reflect.DeepEqual(sampleData, val) {
		t.Fatalf(`Data returned from the store does not match what was returned..
			\n.Expected %v \n.. Got %v instead`, sampleData, val)
	}
}

func TestCache_GetError(t *testing.T) {
	type d map[string][]byte
	f := make(d, 10)
	f["name"] = []byte("memory store")
	m := Cache{items: f, s: &mockSerializeHelper{}}

	val, err := m.Get("name")
	if err == nil {
		t.Fatalf(`An error is supposed to occur if bytes marshalling fails.. Got %v`, err)
	}

	if val != nil {
		t.Fatalf(`Value is supposed to be nil.. Got %v instead`, val)
	}
}

func TestCache_GetTrashGc(t *testing.T) {
	var key = "expiredItem"
	memoryStore.Put(key, []byte("just set this"), time.Nanosecond*1)
	val, err := memoryStore.Get(key)

	if val != nil {
		t.Fatalf(`Expected data to have a nil value.. Got %v instead`, val)
	}

	if err != cache.ErrCacheMiss {
		t.Fatalf(`Expected error to be a cache miss..
			\n Expected %v \n Got %v instead`, cache.ErrCacheMiss, err)
	}
}

func TestCache_GetUnknownKey(t *testing.T) {
	val, err := memoryStore.Get("unknown")
	if err != cache.ErrCacheMiss {
		t.Fatalf(`Expeted to get a cache miss error.. \n
			Got %v instead`, err)
	}

	if val != nil {
		t.Fatalf(`Expected %v to be a nil value`, val)
	}
}

func TestCache_Delete(t *testing.T) {
	err := memoryStore.Delete(key)
	if err != nil {
		t.Fatalf(`An error occurred while trying to delete the data from the store... %v`, err)
	}

	_, err = memoryStore.Get(key)
	if err != cache.ErrCacheMiss {
		t.Fatalf(`Expected an error of %v \n Got %v`, err)
	}

	err = memoryStore.Delete("unknown")
	if err != cache.ErrCacheMiss {
		t.Fatalf(`Error should be a missed cache.
			 \n. Expected %v.\n Got %v`, cache.ErrCacheMiss, err)
	}
}

func TestCache_ClearAll(t *testing.T) {
	expectedNumOfItems := 0
	expire := time.Minute * 5

	memoryStore.Put(key, []byte("memory store"), expire)
	memoryStore.Put("me", []byte("yoo"), expire)
	memoryStore.Put("something", []byte("flush"), expire)

	err := memoryStore.ClearAll()
	if err != nil {
		t.Fatalf("An error occurred while the store was being flushed... %v", err)
	}

	if x := len(memoryStore.items); x != expectedNumOfItems {
		t.Fatalf("Store was not flushed..\n Expected %d.. Got %d ", expectedNumOfItems, x)
	}
}

func TestCache_TrashGc(t *testing.T) {
	// Set garbage collection interval to every 5 second
	store := NewMemoryStore(time.Second * 5)

	tests := []struct {
		key, value string
		expires    time.Duration
	}{
		{"memory", "cache store", time.Microsecond},
		{"redis", "redis store", time.Microsecond},
		{"file", "write file", time.Microsecond},
	}

	for _, v := range tests {
		store.Put(v.key, []byte(v.value), v.expires)
	}

	// 
	time.Sleep(time.Second * 6)

	expectedNumOfItemsInCache := 0
	if x := len(store.items); x != expectedNumOfItemsInCache {
		t.Fatalf("Expected %d items in the store. %d found", expectedNumOfItemsInCache, x)
	}
}
