package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseName     string `env:"DB_NAME"`
	DatabaseHost     string `env:"DB_HOST"`
	DatabasePort     string `env:"DB_PORT"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePassword string `env:"DB_PASSWORD"`
	MaxAttemptsConn  int    `env:"MAX_ATTEMPTS"`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("can't load env file")
	}

	var cnf Config
	if err := env.Parse(&cnf); err != nil {
		log.Fatal("can't parse config for postgresql")
	}

	return &cnf, nil
}
