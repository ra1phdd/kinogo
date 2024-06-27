package testutil

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"kinogo/pkg/logger"
	"log"
)

// SetupMocks создает моки для БД и Redis
func SetupMocks() (*sqlx.DB, *sql.DB, sqlmock.Sqlmock, *miniredis.Miniredis, *redis.Client) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// Создаем sqlx.DB из sql.DB
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	logger.Init("debug")

	return sqlxDB, mockDB, mock, mr, rdb
}
