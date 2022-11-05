package storage

import (
	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func InitializeRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
