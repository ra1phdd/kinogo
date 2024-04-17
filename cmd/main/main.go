package main

import (
	"log"

	"kinogo/cmd/server"
	"kinogo/cmd/websocket"
	"kinogo/internal/config"
	"kinogo/internal/services"
	"kinogo/pkg/cache"
	"kinogo/pkg/db"
	"kinogo/pkg/logger"

	"github.com/robfig/cron"
)

func main() {
	// Подгрузка конфигурации
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	// Инициализация логгепа
	logger.Init(config.LoggerLevel)

	// Подключение к БД
	err = db.Init(config.DBUser, config.DBPassword, config.DBHost, config.DBName)
	if err != nil {
		logger.Error("Ошибка при подключении к БД", err)
	}

	// Инициализация кэша Redis
	cache.Init(config.RedisAddr, config.RedisPort, config.RedisPassword)

	go websocket.Start()

	c := cron.New()
	c.AddFunc("0 */1 * * * *", services.SaveStatsToDB)
	c.Start()

	server.Start()
}
