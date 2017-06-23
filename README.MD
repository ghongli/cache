## A Go caching Library
cache is a Go cache manager. It can use many cache. Default use `gob` encode/decode Serialize data.
Supported cache stores
- [x] Memory
- [x] Filesystem
- [x] Memcache
- [x] Redis

### How to install?
```go
go get -u -v  github.com/ghongli/cache
```

### How to use it?
- First must import it
```go
import (
	"time"
	"fmt"
	
	"github.com/ghongli/cache"
	"github.com/ghongli/cache/filesys"
	"github.com/ghongli/cache/memory"
	"github.com/ghongli/cache/redis"
)
```

- Memory
```go
// Setting a larger interval for gc would be better. like time.Minute * 10
// prevent doing the same thing over and over again if nothing really changed
memoryStore := memory.NewMemoryStore(time.Minute*1, nil)
fmt.Println(memoryStore.Put("example", []byte("example"), time.Second*10))
fmt.Println(memoryStore.Get("example"))
fmt.Println(memoryStore.ClearAll())
fmt.Println(memoryStore.Put("n", []byte("n"), time.Second*5))
fmt.Println(memoryStore.Get("n"))
time.Sleep(time.Second * 5)
fmt.Println(memoryStore.Get("n"))
```

- Filesystem
```go
fileStore := filesys.MustNewFileSysStore("./../tmp/cache", time.Minute*5, nil)
fmt.Println(fileStore.Put("example", []byte("example"), time.Second*20))
fmt.Println(fileStore.Get("example"))
```

- Redis
```go
opt := &r.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}

redisStore := redis.NewRedisCache(opt, redis.PREFIX)
fmt.Println(redisStore.Put("example", []byte("example"), time.Second*10))
val, _ := redisStore.Get("example")
fmt.Println(s.Serialized(val))
fmt.Println(redisStore.ClearAll())
fmt.Println(redisStore.Get("example"))
```

### LICENSE
MIT