package config

import "os"

type Config struct {
	AppPort string
}

func Load() Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	return Config{AppPort: port}
}
