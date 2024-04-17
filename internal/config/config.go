package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Configuration struct {
	Port          string `env:"ADDRESS" envDefault:"4000"`
	LoggerLevel   string `env:"LOGGER_LEVEL" envDefault:"debug"`
	DBUser        string `env:"DB_USER,required"`
	DBPassword    string `env:"DB_PASSWORD,required"`
	DBName        string `env:"DB_NAME,required"`
	DBHost        string `env:"DB_HOST,required"`
	RedisAddr     string `env:"REDIS_ADDR,required"`
	RedisPort     string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD,required"`
	IsDebug       bool   `env:"DEBUG" envDefault:"false"`
}

func NewConfig(files ...string) (*Configuration, error) {
	err := godotenv.Load(files...)

	if err != nil {
		log.Printf("No .env file could be found %q\n", files)
	}

	cfg := Configuration{}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
