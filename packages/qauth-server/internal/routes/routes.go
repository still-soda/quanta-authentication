package routes

import (
	"qauth-server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	health := handlers.NewHealthHandler()

	r.GET("/health", health.Check)
}
