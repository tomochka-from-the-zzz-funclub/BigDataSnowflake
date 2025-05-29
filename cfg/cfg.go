package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	SslMode    string
}

func LoadConfig() Config {
	err := godotenv.Load("config/docker.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		SslMode:    os.Getenv("DB_SSLMODE"),
	}
}
