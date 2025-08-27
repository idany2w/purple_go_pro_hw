package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db DbConfig
}

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SSLMode  string
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file, using default config")
	}

	return &Config{
		Db: DbConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DbName:   os.Getenv("POSTGRES_DB"),
			SSLMode:  os.Getenv("POSTGRES_SSL_MODE"),
		},
	}
}
