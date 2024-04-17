package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func Init(RedisAddr string, RedisPort string, RedisPassword string) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     RedisAddr + ":" + RedisPort,
		Password: RedisPassword,
		DB:       0,
	})
}
