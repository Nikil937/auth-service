package http

import (
	"net/http"

	"github.com/Nikil937/auth-service/internal/delivery/http/handler"
	"github.com/Nikil937/auth-service/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(authHandler *handler.AuthHandler, jwtSecret string) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/refresh", authHandler.Refresh)
	router.POST("/auth/logout", authHandler.Logout)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret))

	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole("admin"))

	admin.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "admin access granted",
		})
	})

	protected.GET("/me", authHandler.Me)

	return router
}
