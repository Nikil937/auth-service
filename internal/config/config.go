package config

import (
	"os"
	"time"
)

type Config struct {
	AppPort      string
	DatabaseURL  string
	JWTSecret    string
	JWTAccessTTL time.Duration
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

	return Config{
		AppPort:      port,
		DatabaseURL:  databaseURL,
		JWTSecret:    jwtSecret,
		JWTAccessTTL: jwtAccessTTL,
	}
}
