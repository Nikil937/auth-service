package config

import (
	"os"
	"time"
)

type Config struct {
	AppPort       string
	DatabaseURL   string
	JWTSecret     string
	JWTAccessTTL  time.Duration
	RedisAddr     string
	JWTRefreshTTL time.Duration
}

func Load() Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5433/auth?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret!@#$%^&*()password"
	}

	jwtTTLString := os.Getenv("JWT_ACCESS_TTL")
	if jwtTTLString == "" {
		jwtTTLString = "15m"
	}

	jwtAccessTTL, err := time.ParseDuration(jwtTTLString)
	if err != nil {
		jwtAccessTTL = 15 * time.Minute
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	jwtRefreshTTLString := os.Getenv("JWT_REFRESH_TTL")
	if jwtRefreshTTLString == "" {
		jwtRefreshTTLString = "168h"
	}

	jwtRefreshTTL, err := time.ParseDuration(jwtRefreshTTLString)
	if err != nil {
		jwtRefreshTTL = 168 * time.Hour
	}

	return Config{
		AppPort:       port,
		DatabaseURL:   databaseURL,
		JWTSecret:     jwtSecret,
		JWTAccessTTL:  jwtAccessTTL,
		RedisAddr:     redisAddr,
		JWTRefreshTTL: jwtRefreshTTL,
	}
}
