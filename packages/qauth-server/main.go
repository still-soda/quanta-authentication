package main

import (
	"log"
	"os"
	"os/signal"
	"qauth-server/internal/config"
	"qauth-server/internal/database"
	"qauth-server/internal/middleware"
	"qauth-server/internal/models"
	"qauth-server/internal/routes"
	"qauth-server/internal/storage"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.New()

	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.Close()

	// 自动迁移数据库表
	if err := database.AutoMigrate(&models.User{}, &models.File{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库表迁移完成")

	// 初始化存储服务
	storageService, err := storage.NewCDKStorage(cfg)
	if err != nil {
		log.Fatalf("存储服务初始化失败: %v", err)
	}
	defer storageService.Close()

	r := gin.New()

	// 使用中间件
	middleware.UseMiddleware(r)

	// 静态文件服务（用于访问上传的文件）
	r.Static("/uploads", cfg.Storage.LocalDir)

	// 注册路由
	routes.SetupRoutes(r, storageService)

	// 启动服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(":" + cfg.Server.Port); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	<-quit
}
