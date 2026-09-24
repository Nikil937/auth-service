package http

import (
	"net/http"

	"github.com/Nikil937/auth-service/internal/delivery/http/handler"
	"github.com/Nikil937/auth-service/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(authHandler *handler.AuthHandler, jwtSecret string) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/refresh", authHandler.Refresh)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret))

	protected.GET("/me", authHandler.Me)

	return router
}
