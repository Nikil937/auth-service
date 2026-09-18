package main

import (
	"log"

	"github.com/Nikil937/auth-service/internal/config"
	"github.com/Nikil937/auth-service/internal/delivery/http"
)

func main() {
	cfg := config.Load()

	router := http.NewRouter()

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server %v", err)
	}

}
