package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	HTTPPort  string
	DBDSN     string
	JWTSecret string
	JWTAccessTTL int
	JWTRefreshTTL int
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	httpPort := os.Getenv("HTTP_PORT")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_USER")
	ssl := os.Getenv("DB_SSLMODE")
	JWTAccessTTL,err := strconv.Atoi( os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil{return nil, err}
	JWTRefreshTTL,err := strconv.Atoi( os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil{return nil, err}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		host, port, user, pass, name, ssl)
	secret := os.Getenv("JWT_SECRET")
	return &Config{
		HTTPPort:  httpPort,
		DBDSN:     dsn,
		JWTSecret: secret,
		JWTAccessTTL: JWTAccessTTL,
		JWTRefreshTTL: JWTRefreshTTL,
	}, nil
}
