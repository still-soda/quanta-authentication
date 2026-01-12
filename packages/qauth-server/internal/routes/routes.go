package routes

import (
	"qauth-server/internal/handlers"
	"qauth-server/internal/storage"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, storageService *storage.CDKStorage) {
	health := handlers.NewHealthHandler()
	user := handlers.NewUserHandler()
	file := handlers.NewFileHandler(storageService)

	r.GET("/health", health.Check)
	r.GET("/ping", health.Ping)

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/register", user.Register)
			users.POST("/login", user.Login)
			users.GET("/:id", user.GetProfile)
			users.GET("", user.List)
		}

		files := v1.Group("/files")
		{
			files.POST("/upload", file.Upload)
			files.GET("/:key", file.Download)
			files.DELETE("/:id", file.Delete)
			files.GET("", file.List)
		}
	}
}
