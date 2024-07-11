package metrics_v1

import (
	"fmt"
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
		"metrics_unique_users",
		"metrics_new_users",
		"metrics_returning_users",
		"metrics_page_views",
		"metrics_comments",
		"metrics_registrations",
	}

	err := cache.Rdb.Del(cache.Ctx, keys...).Err()
	if err != nil {
		logger.Warn("Ошибка сброса метрики", zap.Error(err))
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

func (s Service) AvgTime(userTime time.Time) {
	cacheKey := "avg_time:values"

	fmt.Println(userTime.Unix())

	if userTime.Unix() < 10 {
		return
	}

	err := cache.Rdb.LPush(cache.Ctx, cacheKey, userTime.Unix()).Err()
	if err != nil {
		logger.Warn("Ошибка добавления значения в список", zap.Error(err))
		return
	}

	err = cache.Rdb.LTrim(cache.Ctx, cacheKey, 0, 99).Err()
	if err != nil {
		logger.Warn("Ошибка обрезки списка", zap.Error(err))
		return
	}

	values, err := cache.Rdb.LRange(cache.Ctx, cacheKey, 0, -1).Result()
	if err != nil {
		logger.Warn("Ошибка получения значений из списка", zap.Error(err))
		return
	}

	var sum int64
	for _, value := range values {
		t, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			logger.Warn("Ошибка преобразования значения из списка", zap.Error(err))
			return
		}
		sum += t
	}
	avgTime := sum / int64(len(values))

	err = cache.Rdb.Set(cache.Ctx, "metrics:avg_time", avgTime, 0).Err()
	if err != nil {
		logger.Warn("Ошибка установки среднего времени в кеш", zap.Error(err))
	}
}

func (s Service) PageViews() {
	err := cache.Rdb.Incr(cache.Ctx, "metrics:page_views").Err()
	if err != nil {
		logger.Warn("Ошибка инкрементирования метрики metrics_page_views", zap.Error(err))
	}
}

func (s Service) Comments() {
	err := cache.Rdb.Incr(cache.Ctx, "metrics:comments").Err()
	if err != nil {
		logger.Warn("Ошибка инкрементирования метрики metrics_comments", zap.Error(err))
	}
}

func (s Service) Registrations() {
	err := cache.Rdb.Incr(cache.Ctx, "metrics:registrations").Err()
	if err != nil {
		logger.Warn("Ошибка инкрементирования метрики metrics_registrations", zap.Error(err))
	}
}
