package redis

import "github.com/go-redis/redis"

type Cache struct {
	client *redis.Client
	prefix string
}
