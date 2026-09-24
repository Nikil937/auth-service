package main

import (
	"context"
	"log"

	"github.com/Nikil937/auth-service/internal/config"
	"github.com/Nikil937/auth-service/internal/database"
	"github.com/Nikil937/auth-service/internal/delivery/http"
	"github.com/Nikil937/auth-service/internal/delivery/http/handler"
	"github.com/Nikil937/auth-service/internal/repository"
	"github.com/Nikil937/auth-service/internal/service"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer db.Close()

	redisClient, err := database.NewRedisClient(ctx, cfg.RedisAddr)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	defer redisClient.Close()

	refreshRepo := repository.NewRedisRefreshRepository(redisClient)

	userRepo := repository.NewPostgresUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTAccessTTL, refreshRepo, cfg.JWTRefreshTTL)
	authHandler := handler.NewAuthHandler(authService)

	router := http.NewRouter(authHandler, cfg.JWTSecret)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
