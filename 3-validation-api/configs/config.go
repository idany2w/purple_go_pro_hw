package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	Email  EmailConfig
}

type ServerConfig struct {
	Port string
}

type EmailConfig struct {
	Email    string
	Login    string
	Password string
	Address  string
	Port     string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		Email: EmailConfig{
			Email:    os.Getenv("MAIL_FROM"),
			Login:    os.Getenv("MAIL_SMTP_LOGIN"),
			Password: os.Getenv("MAIL_SMTP_PASSWORD"),
			Address:  os.Getenv("MAIL_SMTP_HOST"),
			Port:     os.Getenv("MAIL_SMTP_PORT"),
		},
	}
}

func LoadConfig() *Config {
	return NewConfig()
}
