package migrator

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"kinogo/config"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var migrationsPath string
	// Путь до папки с миграциями.
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.Parse()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка при инициализации конфига: %s", err)
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?sslmode=disable", cfg.DB.DBUser, cfg.DB.DBPassword, cfg.DB.DBHost, cfg.DB.DBName)

	m, err := migrate.New(
		"file://"+migrationsPath,
		connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No migrations to apply")
		} else {
			log.Fatal(err)
		}
	} else {
		log.Println("Migrations applied successfully")
	}
}
