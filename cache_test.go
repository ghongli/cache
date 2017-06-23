package cache

import (
	"testing"
	"time"
)

func TestCacheItem_IsExpired(t *testing.T) {
	item := &CacheItem{
		CreatedTime: time.Now(),
		Expired:     time.Nanosecond * -2,
		Data:        []byte("ping"),
	}

	if !item.IsExpired() {
		t.Fatal("Item should be expired since it's expiration date is set 2 nanoseconds backwards")
	}
}

func TestFileCache(t *testing.T) {

}
