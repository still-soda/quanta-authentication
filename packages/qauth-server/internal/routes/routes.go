package routes

import (
	"qauth-server/internal/handlers"
	"qauth-server/internal/utilities"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, s utilities.Storage) {
	health := handlers.NewHealthHandler()
	file := handlers.NewFileHandler(s)

	r.GET("/health", health.Check)
	r.POST("/upload", file.Upload)
}
