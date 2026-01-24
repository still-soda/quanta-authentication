package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"qauth-server/internal/config"
	"qauth-server/internal/database"
	"qauth-server/internal/handlers"
	"qauth-server/internal/handlers/business"
	"qauth-server/internal/routes"
	"qauth-server/internal/services"
	"qauth-server/internal/tasks"
	"qauth-server/internal/utilities"
	"qauth-server/pkg/jwks"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
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

	// 静态文件服务
	r.Static("/uploads", cfg.Storage.LocalDir)

	// 创建 jwks 管理器
	jwksManager, err := jwks.NewJWKSManager(nil)
	if err != nil {
		panic("failed to create jwks manager: " + err.Error())
	}
	// 启动密钥轮换
	jwksManager.StartRotation()
	defer jwksManager.StopRotation()

	// 创建服务
	storageService, err := services.NewStorageService(cfg)
	if err != nil {
		panic("failed to create storage service")
	}
	defer storageService.Close()

	cacheService := services.NewCacheService(cfg)
	defer cacheService.Close()

	permissionService := services.NewPermissionService(db)
	fileService := services.NewFileService(storageService, db)
	userService := services.NewUserService(db)
	roleService := services.NewRoleService(db, permissionService, userService)
	oauthService := services.NewOAuthService(db, cfg, jwksManager, userService)
	counterService := services.NewCounterService(db)
	auditService := services.NewAuditService(db, userService)

	// OIDC 服务
	issuer := "http://localhost:" + cfg.Server.Port
	oidcService, err := services.NewOIDCService(cfg, issuer, jwksManager)
	if err != nil {
		panic("failed to create oidc service: " + err.Error())
	}
	defer oidcService.Close()

	// 创建定时器
	cronScheduler := cron.New(cron.WithSeconds())

	// 创建定时任务
	counterTask := tasks.NewCounterTask(counterService, cacheService, userService)
	cleanupCounterTask, err := counterTask.Register(cronScheduler)
	if err != nil {
		panic("failed to register counter task: " + err.Error())
	}
	defer cleanupCounterTask()

	userTask := tasks.NewUserTask(cacheService)
	cleanupUserTask, err := userTask.Register(cronScheduler)
	if err != nil {
		panic("failed to register user task: " + err.Error())
	}
	defer cleanupUserTask()

	// 创建处理器
	healthHandler := handlers.NewHealthHandler()
	fileHandler := handlers.NewFileHandler(fileService)
	authHandler := handlers.NewAuthHandler(userService, roleService, auditService)
	oauthHandler := handlers.NewOAuthHandler(oauthService, roleService, userService, oidcService, cacheService, auditService)
	oidcHandler := handlers.NewOIDCHandler(oidcService, userService)
	roleHandler := handlers.NewRoleHandler(roleService, permissionService, auditService)
	permissionHandler := handlers.NewPermissionHandler(roleService, permissionService, auditService)
	dashboardHandler := business.NewDashboardHandler(userService, counterService, cacheService)
	auditHandler := business.NewAuditHandler(auditService, roleService)

	// 注册路由
	routes.RegisterRoutes(r, &routes.RegisterRouterHandlers{
		HealthHandler:     healthHandler,
		FileHandler:       fileHandler,
		AuthHandler:       authHandler,
		OAuthHandler:      oauthHandler,
		OIDCHandler:       oidcHandler,
		RoleHandler:       roleHandler,
		PermissionHandler: permissionHandler,
		DashboardHandler:  dashboardHandler,
		AuditHandler:      auditHandler,
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
