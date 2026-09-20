package http

import (
	"net/http"

	"github.com/Nikil937/auth-service/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func NewRouter(authHandler *handler.AuthHandler) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.POST("/auth/register", authHandler.Register)

	return router
}
