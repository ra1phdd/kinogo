package metrics_v1

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"kinogo/pkg/cache"
	"kinogo/pkg/logger"
	"strconv"
	"time"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s Service) WriteFromDB() {
	panic("implement me")
}

func (s Service) Reset() {
	keys := []string{
		"metrics:*",
		"spent_time:*",
		"unique_users:*",
		"returning_users:*",
	}

	for _, key := range keys {
		var cursor uint64
		for {
			// Используем SCAN вместо KEYS для более эффективного поиска ключей
			keysFromKey, cursor, err := cache.Rdb.Scan(cache.Ctx, cursor, key, 1000).Result()
			if err != nil {
				logger.Warn("Ошибка сканирования ключей из кеша", zap.Error(err))
				return
			}

			if len(keysFromKey) > 0 {
				// Используем Pipeline для пакетного удаления ключей
				pipe := cache.Rdb.Pipeline()
				for _, keyFromKey := range keysFromKey {
					pipe.Del(cache.Ctx, keyFromKey)
				}
				_, err = pipe.Exec(cache.Ctx)
				if err != nil {
					logger.Warn("Ошибка удаления ключей из кеша", zap.Error(err))
					continue
				}
			}

			// SCAN возвращает курсор 0, когда итерация завершена
			if cursor == 0 {
				break
			}
		}
	}
}

func calcAvgTime() {
	cacheKeyPattern := "spent_time:*"

	// Получение всех ключей по паттерну
	keys, err := cache.Rdb.Keys(cache.Ctx, cacheKeyPattern).Result()
	if err != nil {
		logger.Warn("Ошибка получения ключей из кеша", zap.Error(err))
		return
	}

	// Если ключей нет, то возвращаемся
	if len(keys) == 0 {
		logger.Info("Нет ключей для расчета среднего времени")
		return
	}

	var totalSpentTime int64
	var count int64

	// Проходимся по всем ключам и суммируем значения
	for _, key := range keys {
		spentTimeStr, err := cache.Rdb.Get(cache.Ctx, key).Result()
		if err != nil {
			logger.Warn("Ошибка получения значения из кеша", zap.Error(err))
			continue
		}

		spentTime, err := strconv.ParseInt(spentTimeStr, 10, 64)
		if err != nil {
			logger.Warn("Ошибка преобразования строки в int", zap.Error(err))
			continue
		}

		totalSpentTime += spentTime
		count++
	}

	// Вычисляем среднее значение
	averageSpentTime := totalSpentTime / count

	_, err = cache.Rdb.Set(cache.Ctx, "metrics:avg_time", averageSpentTime, 0).Result()
	if err != nil {
		logger.Warn("Ошибка установки ключа в кеш с TTL", zap.Error(err))
	}
}

func calcBounceRate() {
	cacheKeyPattern := "spent_time:*"

	// Получение всех ключей по паттерну
	keys, err := cache.Rdb.Keys(cache.Ctx, cacheKeyPattern).Result()
	if err != nil {
		logger.Warn("Ошибка получения ключей из кеша", zap.Error(err))
		return
	}

	// Если ключей нет, то возвращаемся
	if len(keys) == 0 {
		logger.Info("Нет ключей для расчета среднего времени")
		return
	}

	var bounceRate int64
	var count int64

	// Проходимся по всем ключам и суммируем значения
	for _, key := range keys {
		spentTimeStr, err := cache.Rdb.Get(cache.Ctx, key).Result()
		if err != nil {
			logger.Warn("Ошибка получения значения из кеша", zap.Error(err))
			continue
		}

		spentTime, err := strconv.ParseInt(spentTimeStr, 10, 64)
		if err != nil {
			logger.Warn("Ошибка преобразования строки в int", zap.Error(err))
			continue
		}

		if spentTime < 60 {
			bounceRate += spentTime
		}
		count++
	}

	// Вычисляем среднее значение
	averageBounceRate := bounceRate / count

	_, err = cache.Rdb.Set(cache.Ctx, "metrics:bounce_rate", averageBounceRate, 0).Result()
	if err != nil {
		logger.Warn("Ошибка установки ключа в кеш с TTL", zap.Error(err))
	}
}

func (s Service) UniqueUsers(uuid string) {
	cacheKey := fmt.Sprintf("unique_users:%s:%s", time.Now().UTC().Format("2006-01-02"), uuid)

	exists, err := cache.Rdb.Exists(cache.Ctx, cacheKey).Result()
	if err != nil {
		logger.Warn("Ошибка проверки наличия ключа в кеше", zap.Error(err))
		return
	}

	if exists == 0 {
		err := cache.Rdb.Incr(cache.Ctx, "metrics:unique_users").Err()
		if err != nil {
			logger.Warn("Ошибка инкрементирования метрики metrics_unique_users", zap.Error(err))
			return
		}

		_, err = cache.Rdb.Set(cache.Ctx, cacheKey, "exists", time.Until(time.Now().UTC().AddDate(0, 0, 1))).Result()
		if err != nil {
			logger.Warn("Ошибка установки ключа в кеш с TTL", zap.Error(err))
		}
	}
}

func (s Service) NewUsers() {
	err := cache.Rdb.Incr(cache.Ctx, "metrics:new_users").Err()
	if err != nil {
		logger.Warn("Ошибка инкрементирования метрики metrics_new_users", zap.Error(err))
	}
}

func (s Service) ReturningUsers(uuid string) {
	cacheKey := fmt.Sprintf("returning_users:%s:%s", time.Now().Month(), uuid)

	exists, err := cache.Rdb.Exists(cache.Ctx, cacheKey).Result()
	if err != nil {
		logger.Warn("Ошибка проверки наличия ключа в кеше", zap.Error(err))
		return
	}

	if exists == 0 {
		err := cache.Rdb.Incr(cache.Ctx, "metrics:returning_users").Err()
		if err != nil {
			logger.Warn("Ошибка инкрементирования метрики metrics_unique_users", zap.Error(err))
			return
		}

		_, err = cache.Rdb.Set(cache.Ctx, cacheKey, "exists", time.Until(time.Now().UTC().AddDate(0, 0, 1))).Result()
		if err != nil {
			logger.Warn("Ошибка установки ключа в кеш с TTL", zap.Error(err))
		}
	}
}

func (s Service) SpentTime(userTime time.Time, uuid string) {
	cacheKey := fmt.Sprintf("spent_time:%s", uuid)

	// Попробуем получить значение из кеша
	spentTimeStr, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		logger.Warn("Ошибка получения значения из кеша", zap.Error(err))
		return
	}

	// Инициализация переменной для хранения времени
	var spentTime int64

	if err == nil {
		// Преобразуем строку в int64
		spentTime, err = strconv.ParseInt(spentTimeStr, 10, 64)
		if err != nil {
			logger.Warn("Ошибка преобразования строки в int", zap.Error(err))
			return
		}
	}

	// Обновляем значение в кеше
	newSpentTime := spentTime + userTime.Unix()
	err = cache.Rdb.Set(cache.Ctx, cacheKey, newSpentTime, 0).Err()
	if err != nil {
		logger.Warn("Ошибка обновления значения в кеше", zap.Error(err))
	}
}

func (s Service) NewComments() {
	err := cache.Rdb.Incr(cache.Ctx, "metrics:new_comments").Err()
	if err != nil {
		logger.Warn("Ошибка инкрементирования метрики metrics_comments", zap.Error(err))
	}
}

func (s Service) NewRegistrations() {
	err := cache.Rdb.Incr(cache.Ctx, "metrics:new_registrations").Err()
	if err != nil {
		logger.Warn("Ошибка инкрементирования метрики metrics_registrations", zap.Error(err))
	}
}
