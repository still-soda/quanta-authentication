package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"qauth-server/internal/config"
	"qauth-server/internal/database"
	"qauth-server/internal/middleware"
	"qauth-server/internal/routes"
	"qauth-server/internal/utilities"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.New()
	logger := utilities.GetLogger()

	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	db, err := database.InitDB(cfg)
	if err != nil {
		logger.Error("failed to initialize database: ", err)
		return
	}
	defer database.Close(db)

	// 自动迁移数据库表
	if err := database.AutoMigrate(db); err != nil {
		logger.Error("failed to migrate database: ", err)
		return
	}
	logger.Info("successfully migrated database tables")

	// 初始化存储服务
	storageService, err := utilities.NewLocalStorage(cfg)
	if err != nil {
		logger.Error("failed to initialize storage service: ", err)
		return
	}
	defer storageService.Close()

	r := gin.New()

	// 注入 orm
	r.Use(func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	})

	// 使用中间件
	r.Use(
		middleware.Logger(),
		middleware.Recovery(),
		middleware.CORS(),
		middleware.Database(db),
	)

	// 静态文件服务（用于访问上传的文件）
	r.Static("/uploads", cfg.Storage.LocalDir)

	// 注册路由
	routes.SetupRoutes(r)

	// 启动服务器
	logger.Info("server was running on port :" + cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to setup server: ", err)
		}
	}()

	// 优雅退出
	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown: ", err)
	}

	logger.Info("server exiting")
}
