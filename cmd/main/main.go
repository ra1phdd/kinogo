package main

import (
	"kinogo/cmd/websocket"
	"kinogo/internal/app/config"
	"kinogo/internal/pkg/app"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// Подгрузка конфигурации
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// Инициализация логгепа
	logger.Init(config.LoggerLevel)

	// Подключение к БД
	err = db.Init(config.DBUser, config.DBPassword, config.DBHost, config.DBName)
	if err != nil {
		logger.Fatal("Ошибка при подключении к БД", zap.Error(err))
	}

	// Инициализация кэша Redis
	cache.Init(config.RedisAddr, config.RedisPort, config.RedisPassword)

	go websocket.Start()

	//c := cron.New()
	//c.AddFunc("0 */1 * * * *", services.SaveStatsToDB)
	//c.Start()

	_, err = app.New()
	if err != nil {
		return
	}
}
