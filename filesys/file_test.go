package filesys

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"github.com/ghongli/cache"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

var (
	_          cache.Cache = MustNewFileSysStore("./../tmp", time.Second)
	fileCache  *Cache
	sampleData = []byte("file store")
	key        = "file_sys"
)

func TestMain(t *testing.M) {
	fileCache = MustNewFileSysStore("./../tmp/cache", time.Second*1)

	flag.Parse()
	os.Exit(t.Run())
}

func TestMustNewFileSysStore(t *testing.T) {
	defer func() {
		recover()
	}()

	_ = MustNewFileSysStore("/hh", time.Second*1)
}

func TestCache_Put(t *testing.T) {
	err := fileCache.Put(key, sampleData, time.Minute*1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCache_Get(t *testing.T) {
	val, err := fileCache.Get(key)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(val, sampleData) {
		t.Fatalf(`Values are not equal.. Expected %v \n
			Got %v`, sampleData, val)
	}
}

func TestCache_GetUnknownKey(t *testing.T) {
	val, err := fileCache.Get("UnknownKey")
	if err == nil {
		t.Fatal(`Expected an error for a file that doesn't exist on the filesystem`)
	}
	if val != nil {
		t.Fatalf(`Expected a nil item to be return ... Got %v instead`, val)
	}
}

func TestCache_ClearAll(t *testing.T) {
	if err := fileCache.ClearAll(); err != nil {
		t.Fatalf("The cache directory, %s could not be flushed... %v", fileCache.baseDir, err)
	}
}

func TestCache_TrashGc(t *testing.T) {
	err := fileCache.Put("xyz", []byte("mock ..."), cache.EXPIRES_DEFAULT)
	if err != nil {
		t.Fatalf("An error occurred... %v", err)
	}

	data, err := fileCache.Get("xyz")
	if err != cache.ErrCacheMiss {
		t.Fatal("Cached data is supposed to be expired")
	}
	if data != nil {
		t.Fatal("Garbage collected item is supposed to be empty")
	}
}

func TestCache_Delete(t *testing.T) {
	if err := fileCache.Delete(key); err != nil {
		t.Fatalf("Could not delete the cached data... %v", err)
	}
}

func TestCache_IsExist(t *testing.T) {
	val := fileCache.IsExist(key)
	if val {
		t.Fatalf("an item with the key is exist ? expected %v.. \n %v instead.", false, val)
	}
}

func TestFilePathForKey(t *testing.T) {
	key := "page"

	b := md5.Sum([]byte(key))
	s := hex.EncodeToString(b[:])

	path := filepath.Join(fileCache.baseDir, s[0:2], s[2:4], s[4:6], s)
	if x := fileCache.filePathForKey(key); x != path {
		t.Fatalf("Path differs.. Expected %s. Got %s instead", path, x)
	}
}

type mockSerializeHelper struct {
}

func (m *mockSerializeHelper) Serialized(i interface{}) ([]byte, error) {
	return nil, errors.New("an error occurred")
}

func (m *mockSerializeHelper) DeSerialized(data []byte, i interface{}) error {
	return errors.New("another error")
}

func TestCache_GetFails(t *testing.T) {
	fileCache.Put("test", []byte("test"), time.Second*1)

	fs := &Cache{"./../tmp/cache", &mockSerializeHelper{}}
	_, err := fs.Get("test")
	if err == nil {
		t.Fatalf("Expected a cache miss.. Got %v", err)
	}
}

func TestCache_PutFails(t *testing.T) {
	fs := &Cache{"./../tmp/cache", &mockSerializeHelper{}}
	err := fs.Put("test", []byte("test"), time.Nanosecond*4)
	if err == nil {
		t.Fatalf("Expected an error from bytes marshalling.. Got %v", err)
	}
}

func TestCache_PutFailsAnUnwriteableDirectory(t *testing.T) {
	fileCache.baseDir = "/"

	if err := fileCache.Put("test", []byte("test"), time.Nanosecond*4); err == nil {
		t.Fatal("An error was supposed to occur because the root directory isn't writeable")
	}
}

func TestCache_TrashGc2(t *testing.T) {
	store := MustNewFileSysStore("./../tmp/cache", time.Second*1)

	tests := []struct {
		key, value string
		expires    time.Duration
	}{
		{"file", "file cache", time.Microsecond},
		{"numb", "two ...", time.Microsecond},
		{"x", "cache test", time.Microsecond},
	}

	for _, v := range tests {
		store.Put(v.key, []byte(v.value), v.expires)
	}

	time.Sleep(time.Second * 2)

	var filePath string
	for _, v := range tests {
		filePath = store.filePathForKey(v.key)
		if _, err := os.Stat(filePath); err == nil {
			t.Fatal("File exists when it isn't supposed to since there was a garbage collection")
		}
	}
}
