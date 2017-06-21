package filesys

import (
	"flag"
	"github.com/ghongli/cache"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	_          cache.Cache = MustNewFileSysStore("./", time.Second)
	fileCache  *Cache
	sampleData = []byte("file store")
	key        = "file_sys"
)

func TestMain(t *testing.M) {
	fileCache = MustNewFileSysStore("./../cahce", time.Second*1)

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

}
