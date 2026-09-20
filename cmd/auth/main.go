package main

import (
	"context"
	"log"

	"github.com/Nikil937/auth-service/internal/config"
	"github.com/Nikil937/auth-service/internal/database"
	"github.com/Nikil937/auth-service/internal/delivery/http"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	router := http.NewRouter()

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
