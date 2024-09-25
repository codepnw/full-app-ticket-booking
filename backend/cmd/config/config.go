package config

import (
	"log"

	env "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	ServerPort string `env:"SERVER_PORT,required"`
	DBDriver   string `env:"DB_DRIVER,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBName     string `env:"DB_NAME,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBSSLMode  string `env:"DB_SSLMODE,required"`
	DBPort     int    `env:"DB_PORT,required"`
}

func NewEnvConfig(envFile string) *EnvConfig {
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("failed to load env file: %e", err)
	}

	config := &EnvConfig{}
	if err := env.Parse(config); err != nil {
		log.Fatalf("failed to load variable from env file: %e", err)
	}

	return config
}
