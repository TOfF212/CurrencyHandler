package config

import (
	"os"
)

// Config структура для хранения параметров конфигурации приложения
type Config struct {
	DatabaseURL string
	ServerPort  string
}

func LoadConfig() *Config {
	dbURL := getEnv("DATABASE_URL", "postgres://user:password@localhost:5435/mydb?sslmode=disable")
	serverPort := getEnv("SERVER_PORT", ":8080")

	return &Config{
		DatabaseURL: dbURL,
		ServerPort:  serverPort,
	}
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
