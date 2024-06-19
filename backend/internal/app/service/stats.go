package service

import (
	"fmt"
	"go.uber.org/zap"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"
	"net/http"
	"strconv"
	"strings"
)

func (s *Service) HandleView(r *http.Request, contentID int64) {
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	// Проверка, был ли пользователь зарегистрирован как просмотревший данный контент
	ipViewKey := fmt.Sprintf("ip_view:%d:%s", contentID, ipAddress)
	uaViewKey := fmt.Sprintf("ua_view:%d:%s", contentID, userAgent)

	// Если пользователь не был зарегистрирован, увеличиваем счетчик просмотров и добавляем запись о пользователе
	if cache.Rdb.Exists(cache.Ctx, ipViewKey).Val() == 0 && cache.Rdb.Exists(cache.Ctx, uaViewKey).Val() == 0 {
		// Получение текущего значения views из базы данных
		var views int64
		err := db.Conn.QueryRow(`SELECT views FROM movies WHERE id = $1`, contentID).Scan(&views)
		if err != nil {
			logger.Error("Ошибка получения значения views из базы данных", zap.Error(err))
			return
		}

		// Увеличение значения views на 1 и сохранение в Redis
		viewsKey := fmt.Sprintf("views:%d", contentID)
		cache.Rdb.SetNX(cache.Ctx, viewsKey, views+1, 0)
		cache.Rdb.SetNX(cache.Ctx, ipViewKey, "1", 0)
		cache.Rdb.SetNX(cache.Ctx, uaViewKey, "1", 0)
		logger.Debug("Просмотр зарегистрирован", zap.Int64("id", contentID), zap.String("ip", ipAddress), zap.String("user-agent", userAgent), zap.Int64("views", views+1))
	}
}

func (s *Service) HandleLike(r *http.Request, contentID int64) {
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	// Проверка, был ли пользователь зарегистрирован как оценивший данный контент
	ipLikeKey := fmt.Sprintf("ip_like:%d:%s", contentID, ipAddress)
	uaLikeKey := fmt.Sprintf("ua_like:%d:%s", contentID, userAgent)

	// Если пользователь не был зарегистрирован, увеличиваем счетчик лайков и добавляем запись о пользователе
	if cache.Rdb.Exists(cache.Ctx, ipLikeKey).Val() == 0 && cache.Rdb.Exists(cache.Ctx, uaLikeKey).Val() == 0 {
		// Получение текущего значения likes из базы данных
		var likes int64
		err := db.Conn.QueryRow(`SELECT likes FROM movies WHERE id = $1`, contentID).Scan(&likes)
		if err != nil {
			logger.Error("Ошибка получения значения likes из базы данных", zap.Error(err))
			return
		}

		// Увеличение значения likes на 1 и сохранение в Redis
		likesKey := fmt.Sprintf("likes:%d", contentID)
		cache.Rdb.SetNX(cache.Ctx, likesKey, likes+1, 0)
		cache.Rdb.SetNX(cache.Ctx, ipLikeKey, "1", 0)
		cache.Rdb.SetNX(cache.Ctx, uaLikeKey, "1", 0)
		logger.Debug("Лайк зарегистрирован", zap.Int64("id", contentID), zap.String("ip", ipAddress), zap.String("user-agent", userAgent), zap.Int64("likes", likes+1))
	}
}

func (s *Service) HandleDislike(r *http.Request, contentID int64) {
	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	// Проверка, был ли пользователь зарегистрирован как оценивший данный контент
	ipDislikeKey := fmt.Sprintf("ip_dislike:%d:%s", contentID, ipAddress)
	uaDislikeKey := fmt.Sprintf("ua_dislike:%d:%s", contentID, userAgent)

	// Если пользователь не был зарегистрирован, увеличиваем счетчик дизлайков и добавляем запись о пользователе
	if cache.Rdb.Exists(cache.Ctx, ipDislikeKey).Val() == 0 && cache.Rdb.Exists(cache.Ctx, uaDislikeKey).Val() == 0 {
		// Получение текущего значения dislikes из базы данных
		var dislikes int64
		err := db.Conn.QueryRow(`SELECT dislikes FROM movies WHERE id = $1`, contentID).Scan(&dislikes)
		if err != nil {
			logger.Error("Ошибка получения значения dislikes из базы данных", zap.Error(err))
			return
		}

		// Увеличение значения dislikes на 1 и сохранение в Redis
		dislikesKey := fmt.Sprintf("dislikes:%d", contentID)
		cache.Rdb.SetNX(cache.Ctx, dislikesKey, dislikes+1, 0)
		cache.Rdb.SetNX(cache.Ctx, ipDislikeKey, "1", 0)
		cache.Rdb.SetNX(cache.Ctx, uaDislikeKey, "1", 0)
		logger.Debug("Дизлайк зарегистрирован", zap.Int64("id", contentID), zap.String("ip", ipAddress), zap.String("user-agent", userAgent), zap.Int64("dislikes", dislikes+1))
	}
}

func (s *Service) SaveStatsToDB() {
	// Получение данных из Redis
	views, err := cache.Rdb.Keys(cache.Ctx, "views:*").Result()
	if err != nil {
		logger.Error("Ошибка получения данных о просмотрах из Redis", zap.Error(err))
		return
	}

	likes, err := cache.Rdb.Keys(cache.Ctx, "likes:*").Result()
	if err != nil {
		logger.Error("Ошибка получения данных о лайках из Redis", zap.Error(err))
		return
	}

	dislikes, err := cache.Rdb.Keys(cache.Ctx, "dislikes:*").Result()
	if err != nil {
		logger.Error("Ошибка получения данных о дизлайках из Redis", zap.Error(err))
		return
	}

	logger.Debug("Данные из Redis получены", zap.Strings("views", views), zap.Strings("likes", likes), zap.Strings("dislikes", dislikes))

	// Сохранение данных в базу данных
	for _, viewKey := range views {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(viewKey, "views:"), 10, 64)
		views, _ := cache.Rdb.Get(cache.Ctx, viewKey).Int64()
		// Обновление значения views в базе данных
		_, err := db.Conn.Exec(`UPDATE movies SET views = $1 WHERE id = $2`, views, contentID)
		if err != nil {
			logger.Error("Ошибка обновления значения views в базе данных", zap.Error(err))
			return
		}
	}

	for _, likeKey := range likes {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(likeKey, "likes:"), 10, 64)
		likes, _ := cache.Rdb.Get(cache.Ctx, likeKey).Int64()
		// Обновление значения likes в базе данных
		_, err := db.Conn.Exec(`UPDATE movies SET likes = $1 WHERE id = $2`, likes, contentID)
		if err != nil {
			logger.Error("Ошибка обновления значения likes в базе данных", zap.Error(err))
			return
		}
	}

	for _, dislikeKey := range dislikes {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(dislikeKey, "dislikes:"), 10, 64)
		dislikes, _ := cache.Rdb.Get(cache.Ctx, dislikeKey).Int64()
		// Обновление значения dislikes в базе данных
		_, err := db.Conn.Exec(`UPDATE movies SET dislikes = $1 WHERE id = $2`, dislikes, contentID)
		if err != nil {
			logger.Error("Ошибка обновления значения dislikes в базе данных", zap.Error(err))
			return
		}
	}

	// Очистка данных из Redis
	for _, viewKey := range views {
		cache.Rdb.Del(cache.Ctx, viewKey)
	}

	for _, likeKey := range likes {
		cache.Rdb.Del(cache.Ctx, likeKey)
	}

	for _, dislikeKey := range dislikes {
		cache.Rdb.Del(cache.Ctx, dislikeKey)
	}
}

func (s *Service) GetStatsToDB(id int) (int64, int64, int64, error) {
	var allViews int64
	var allLikes int64
	var allDislikes int64

	logger.Debug("Получение данных из базы данных", zap.Int("id", id))

	views, err := cache.Rdb.Keys(cache.Ctx, "views:*").Result()
	if err != nil {
		logger.Error("Ошибка получения данных о просмотрах из Redis", zap.Error(err))
		return 0, 0, 0, err
	}

	likes, err := cache.Rdb.Keys(cache.Ctx, "likes:*").Result()
	if err != nil {
		logger.Error("Ошибка получения данных о лайках из Redis", zap.Error(err))
		return 0, 0, 0, err
	}

	dislikes, err := cache.Rdb.Keys(cache.Ctx, "dislikes:*").Result()
	if err != nil {
		logger.Error("Ошибка получения данных о дизлайках из Redis", zap.Error(err))
		return 0, 0, 0, err
	}

	// Сохранение данных в базу данных
	for _, viewKey := range views {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(viewKey, "views:"), 10, 64)
		views, _ := cache.Rdb.Get(cache.Ctx, viewKey).Int64()
		if contentID == int64(id) {
			allViews = views
		}
	}

	for _, likeKey := range likes {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(likeKey, "likes:"), 10, 64)
		likes, _ := cache.Rdb.Get(cache.Ctx, likeKey).Int64()
		if contentID == int64(id) {
			allLikes = likes
		}
	}

	for _, dislikeKey := range dislikes {
		contentID, _ := strconv.ParseInt(strings.TrimPrefix(dislikeKey, "dislikes:"), 10, 64)
		dislikes, _ := cache.Rdb.Get(cache.Ctx, dislikeKey).Int64()
		if contentID == int64(id) {
			allDislikes = dislikes
		}
	}

	logger.Debug("Данные из Redis по фильму получены", zap.Int64("views", allViews), zap.Int64("likes", allLikes), zap.Int64("dislikes", allDislikes))

	return allViews, allLikes, allDislikes, nil
}
