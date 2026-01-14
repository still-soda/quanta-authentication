package routes

import (
	_handlers "qauth-server/internal/handlers"

	"github.com/gin-gonic/gin"
)

type RegisterRouterHandlers struct {
	HealthHandler *_handlers.HealthHandler
	FileHandler   *_handlers.FileHandler
	AuthHandler   *_handlers.AuthHandler
}

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine, handlers *RegisterRouterHandlers) {
	healthHandler := handlers.HealthHandler
	fileHandler := handlers.FileHandler
	authHandler := handlers.AuthHandler

	r.GET("/health", healthHandler.Check)
	r.POST("/upload", fileHandler.Upload)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}
}
