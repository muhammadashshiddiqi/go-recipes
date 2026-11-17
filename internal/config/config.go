package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port string
	DatabaseURL string
	JWTSecret string
	JWTTL time.Duration
}

func Load() Config  {
	hours, err := strconv.Atoi(getEnv("JWT_TTL_HOURS", "24"))
	if err != nil || hours <= 0 {
		hours = 24
	}

	return Config{
		Port: getEnv("APP_PORT", "3030"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5444/go-recipe?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", ""),
		JWTTL: time.Duration(hours) * time.Hour,
	}
}


func getEnv(key, fallback string) string  {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}