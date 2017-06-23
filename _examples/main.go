package main

import (
	"bytes"
	"fmt"
	"github.com/ghongli/cache"
	"github.com/ghongli/cache/filesys"
	"github.com/ghongli/cache/memory"
	redis "github.com/ghongli/cache/redis"
	r "github.com/go-redis/redis"
	"reflect"
	"time"
)

func main() {
	serializer := cache.NewCacheSerializer()
	i := &cache.CacheItem{Data: []byte("cache"), CreatedTime: time.Now(), Expired: time.Minute * 1}

	redisCache(serializer, i)
	memoryCache(serializer, i)
	fileCache(serializer, i)
}

func redisCache(s cache.Serializer, i *cache.CacheItem) {
	opt := &r.Options{
		Addr:     "101.201.197.163:6379",
		Password: "eGd3cEn38tYCQiDBzx7PTWwO",
		DB:       6,
	}

	redisStore := redis.NewRedisCache(opt, redis.PREFIX)
	byt, _ := s.Serialized(i)
	fmt.Println(redisStore.Put("example", byt, time.Second*10))
	val, _ := redisStore.Get("example")
	fmt.Println(s.Serialized(val))
	fmt.Println(redisStore.ClearAll())
	fmt.Println(redisStore.Get("example"))
}

func memoryCache(s cache.Serializer, i *cache.CacheItem) {
	// Setting a larger interval for gc would be better. like time.Minute * 10
	// prevent doing the same thing over and over again if nothing really changed
	memoryStore := memory.NewMemoryStore(time.Minute*1, s)
	byt, _ := s.Serialized(i)
	fmt.Println(memoryStore.Put("example", byt, time.Second*10))
	fmt.Println(memoryStore.Get("example"))
	fmt.Println(memoryStore.ClearAll())
	fmt.Println(memoryStore.Put("n", []byte("n"), time.Second*5))
	fmt.Println(memoryStore.Get("n"))
	time.Sleep(time.Second * 5)
	fmt.Println(memoryStore.Get("n"))
}

func fileCache(s cache.Serializer, i *cache.CacheItem) {
	fileStore := filesys.MustNewFileSysStore("./../tmp/cache", time.Minute*5, s)
	byt, _ := s.Serialized(i)
	fmt.Println(fileStore.Put("example", byt, time.Second*20))
	fmt.Println(fileStore.Get("example"))

	u := &user{"example"}
	b, _ := s.Serialized(u)
	fmt.Println(fileStore.Put("user", b, time.Second*10))
	fmt.Println(fileStore.Get("user"))
	gu, _ := fileStore.Get("user")

	if !bytes.Equal(gu, b) {
		fmt.Errorf("user byte info is not equal ...")
	}

	nu := new(user)
	s.DeSerialized(gu, nu)
	if !reflect.DeepEqual(u, nu) {
		fmt.Errorf("user info is not equal ...")
	}

	fmt.Println(fileStore.ClearAll())
	fmt.Println(fileStore.Get("unkownKey"))

}

type user struct {
	Name string
}
