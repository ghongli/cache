package main

import (
	"fmt"
	"github.com/ghongli/cache"
	redis "github.com/ghongli/cache/redis"
	r "github.com/go-redis/redis"
	"time"
)

func main() {
	serializer := cache.NewCacheSerializer()
	i := &cache.CacheItem{Data: []byte("cache"), CreatedTime: time.Now(), Expired: time.Minute * 1}

	redisCache(serializer, i)
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
}
