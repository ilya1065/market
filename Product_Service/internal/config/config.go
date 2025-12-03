package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	HTTPPort string
	DBDSN    string
	DBDriver string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("не удалось загрузить переменные окружения")
	}
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	return &Config{
		HTTPPort: os.Getenv("HTTP_PORT"),
		DBDSN:    dsn,
		DBDriver: os.Getenv("DB_DRIVER"),
	}
}
