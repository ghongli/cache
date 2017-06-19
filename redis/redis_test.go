package redis

import (
	"bytes"
	"flag"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ghongli/cache"
	"github.com/go-redis/redis"
)

var (
	_          cache.Cache = &Cache{}
	rCache     *Cache
	simpleData = []byte("redis cache")
	key        = "name"
)

const TEST_PREFIX = "r_cache_test:"

func TestMain(m *testing.M) {
	rCache = NewRedisCache(&redis.Options{
		Addr:     "101.201.197.163:6379",
		Password: "eGd3cEn38tYCQiDBzx7PTWwO",
		DB:       6,
	}, TEST_PREFIX)

	flag.Parse()
	os.Exit(m.Run())
}

func TestCache_Put(t *testing.T) {
	err := rCache.Put(key, simpleData, 1*time.Minute)
	if err != nil {
		t.Fatalf(
			`An error occurred while putting data to redis....
			 %v`, err,
		)
	}
}

func TestCache_IsExist(t *testing.T) {
	val := rCache.IsExist(key)
	if !val {
		t.Fatalf("an item with the key is exist ? expected %v.. \n %v instead.", true, false)
	}
}

func TestCache_Get(t *testing.T) {
	val, err := rCache.Get(key)
	if err != nil {
		t.Fatalf(
			`Could not get an item with the key, %s
			 due to an error %v`,
			key, err,
		)
	}

	if !reflect.DeepEqual(simpleData, val) {
		t.Fatalf("Expected %v.. \nGot %v instead", simpleData, val)
	}
}

func TestCache_Delete(t *testing.T) {
	err := rCache.Delete(key)
	if err != nil {
		t.Fatalf(`Could not delete the key %s due to an error ...%v`, "name", err)
	}
}

func TestCache_ClearAll(t *testing.T) {
	// save same data
	rCache.Put("me", []byte("gee.io"), cache.EXPIRES_DEFAULT)
	rCache.Put("animalName", []byte("gopher"), cache.EXPIRES_FOREVER)

	if err := rCache.ClearAll(); err != nil {
		t.Fatalf(`An error occured while flushing the redis database... %v`, err)
	}

	// Manually inspect all data in redis after flushing it's database
	cmd := rCache.client.Keys("*")
	if err := cmd.Err(); err != nil {
		t.Fatalf(`An error occured while trying to get all keys
			stored in redis.. %v`, err)
	}

	res, err := cmd.Result()
	if err != nil {
		t.Fatalf(`An error occured while trying to get the
			result from redis... %v`, err)
	}

	if x := len(res); x != 0 {
		t.Fatalf(`There should be no more data stored in
			REDIS since we flushed the database...
			\n Expected %d.. Got %d instead `, 0, x)
	}
}

func TestCache_GetUnknownKey(t *testing.T) {
	val, err := rCache.Get("oops")
	if err == nil {
		t.Fatalf("An Unknown key was encountered.. Yet we were able to retrieve it")
	}

	if !bytes.Equal(make([]byte, 0), val) {
		t.Fatalf(`Cache store should return a nil value. Since an unknown key was requested.. \n Got %v instead`, val)
	}
}

func TestNewRedisStore_DefaultPrefixIsUsed(t *testing.T) {
	s := NewRedisCache(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       6,
	}, "")

	if !reflect.DeepEqual(PREFIX, s.prefix) {
		t.Fatalf(`Redis store prefix is invalid..
		Expected %s \n... Got %s`, PREFIX, s.prefix)
	}
}
