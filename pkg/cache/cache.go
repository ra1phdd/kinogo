package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	Rdb *redis.Client
)

func Init(Addr string, Username string, Password string, DB int) error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     Addr,
		Username: Username,
		Password: Password,
		DB:       DB,
	})

	err := Rdb.Ping(Ctx).Err()
	if err != nil {
		return err
	}

	return nil
}

/*
функция для стирания кэша

нужна для дэбага
*/
func ClearCache(Rdb *redis.Client) error {
	// Удаление всего кэша из Redis
	err := Rdb.FlushAll(Ctx).Err()
	if err != nil {
		return err
	}
	return nil
}
