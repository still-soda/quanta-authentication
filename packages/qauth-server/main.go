package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"qauth-server/internal/config"
	"qauth-server/internal/database"
	"qauth-server/internal/handlers"
	"qauth-server/internal/middleware"
	"qauth-server/internal/routes"
	"qauth-server/internal/services"
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
		panic("failed to connect database: " + err.Error())
	}
	defer database.Close(db)

	// 自动迁移数据库表
	if err := database.AutoMigrate(db); err != nil {
		panic("failed to migrate database: " + err.Error())
	}
	logger.Info("successfully migrated database tables")

	// 数据库种子数据
	if err := database.SeedingDB(db); err != nil {
		panic("failed to seed database: " + err.Error())
	}
	logger.Info("successfully seeded database")

	r := gin.New()

	// 使用中间件
	r.Use(
		middleware.Logger(logger),
		middleware.Recovery(),
		middleware.CORS(),
	)

	// 静态文件服务
	r.Static("/uploads", cfg.Storage.LocalDir)

	// 创建服务
	storageService, err := services.NewStorageService(cfg)
	if err != nil {
		panic("failed to create storage service")
	}
	defer storageService.Close()

	fileService := services.NewFileService(storageService, db)
	userService := services.NewUserService(db)
	roleService := services.NewRoleService(db)
	oauthService := services.NewOAuthService(db, cfg)

	// OIDC 服务
	issuer := "http://localhost:" + cfg.Server.Port
	oidcService, err := services.NewOIDCService(cfg, issuer)
	if err != nil {
		panic("failed to create oidc service: " + err.Error())
	}
	defer oidcService.Close()

	// 创建处理器
	healthHandler := handlers.NewHealthHandler()
	fileHandler := handlers.NewFileHandler(fileService)
	authHandler := handlers.NewAuthHandler(userService, roleService)
	oauthHandler := handlers.NewOAuthHandler(oauthService)
	oidcHandler := handlers.NewOIDCHandler(oidcService)

	// 注册路由
	routes.RegisterRoutes(r, &routes.RegisterRouterHandlers{
		HealthHandler: healthHandler,
		FileHandler:   fileHandler,
		AuthHandler:   authHandler,
		OAuthHandler:  oauthHandler,
		OIDCHandler:   oidcHandler,
	})

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
