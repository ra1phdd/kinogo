package cache

import (
	"context"
	"fmt"
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

func ClearCacheByPattern(pattern string) error {
	keys, err := Rdb.Keys(Ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	// Удаление всех ключей
	if len(keys) > 0 {
		if err := Rdb.Del(context.Background(), keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete keys: %w", err)
		}
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
