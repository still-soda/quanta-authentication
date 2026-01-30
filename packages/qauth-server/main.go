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
	"qauth-server/internal/providers"
	"qauth-server/internal/repository"
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
	jwksMgr, err := jwks.NewJWKSManager(nil)
	if err != nil {
		panic("failed to create jwks manager: " + err.Error())
	}
	// 启动密钥轮换
	jwksMgr.StartRotation()
	defer jwksMgr.StopRotation()

	// 创建提供者
	loggerProvider := providers.NewRootLogger()
	cacheProvider := providers.NewRedisCache(cfg)
	storageProvider := providers.NewLocalStorage(cfg)
	oauthProvider := providers.NewGoOAuth(db, cfg)

	// 创建仓储
	appGroupAdminRepo := repository.NewAppGroupAdminRepository(db)
	appGroupPermRepo := repository.NewAppGroupPermissionRepository(db)
	appGroupRoleRepo := repository.NewAppGroupRoleRepository(db)
	appGroupUserRoleRepo := repository.NewAppGroupUserRoleRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	counterRepo := repository.NewCounterRepository(db)
	fileRepo := repository.NewFileRepository(db)
	oauthRepo := repository.NewOAuthRepository(db)
	loginStateRepo := repository.NewLoginStateRepository(db)
	errRecordRepo := repository.NewErrorRecordRepository(db)
	permRepo := repository.NewPermissionRepository(db)

	// 创建服务
	permSrv := services.NewPermissionService(
		permRepo,
		loggerProvider,
	)
	fileSrv := services.NewFileService(storageProvider, fileRepo)
	userSrv := services.NewUserService(db)
	roleSrv := services.NewRoleService(db, permSrv, userSrv)
	oauthSrv := services.NewOAuthService(
		cfg,
		oauthProvider,
		jwksMgr,
		userSrv,
		oauthRepo,
		loginStateRepo,
		errRecordRepo,
		loggerProvider,
	)
	counterSrv := services.NewCounterService(counterRepo)
	auditSrv := services.NewAuditService(
		auditRepo,
		userSrv,
		loggerProvider,
	)
	appGroupSrv := services.NewAppGroupService(
		db,
		userSrv,
		appGroupAdminRepo,
		appGroupPermRepo,
		appGroupRoleRepo,
		appGroupUserRoleRepo,
		loggerProvider,
	)

	// 设置 OAuthService 的 AppGroupService（延迟注入以避免循环依赖）
	oauthSrv.SetAppGroupService(appGroupSrv)

	// OIDC 服务
	issuer := "http://localhost:" + cfg.Server.Port
	oidcSrv, err := services.NewOIDCService(cfg, issuer, jwksMgr)
	if err != nil {
		panic("failed to create oidc service: " + err.Error())
	}
	defer oidcSrv.Close()

	// 创建定时器
	cronScheduler := cron.New(cron.WithSeconds())

	// 创建定时任务
	counterTask := tasks.NewCounterTask(cacheProvider, counterSrv, userSrv, oauthSrv)
	cleanupCounterTask, err := counterTask.Register(cronScheduler)
	if err != nil {
		panic("failed to register counter task: " + err.Error())
	}
	defer cleanupCounterTask()

	userTask := tasks.NewUserTask(cacheProvider)
	cleanupUserTask, err := userTask.Register(cronScheduler)
	if err != nil {
		panic("failed to register user task: " + err.Error())
	}
	defer cleanupUserTask()

	// 创建处理器
	healthHandler := handlers.NewHealthHandler()
	fileHandler := handlers.NewFileHandler(fileSrv)
	authHandler := handlers.NewAuthHandler(userSrv, roleSrv, auditSrv)
	oauthHandler := handlers.NewOAuthHandler(cacheProvider, oauthSrv, roleSrv, userSrv, oidcSrv, auditSrv, appGroupSrv)
	oidcHandler := handlers.NewOIDCHandler(oidcSrv, userSrv)
	roleHandler := handlers.NewRoleHandler(roleSrv, permSrv, auditSrv)
	permHandler := handlers.NewPermissionHandler(roleSrv, permSrv, auditSrv)
	userHandler := handlers.NewUserHandler(userSrv, roleSrv, auditSrv)
	appGroupHandler := handlers.NewAppGroupHandler(appGroupSrv, oauthSrv, roleSrv, auditSrv)
	dashboardHandler := business.NewDashboardHandler(cacheProvider, userSrv, counterSrv, oauthSrv)
	auditHandler := business.NewAuditHandler(auditSrv, roleSrv)

	// 注册路由
	routes.RegisterRoutes(r, &routes.RegisterRouterHandlers{
		HealthHandler:     healthHandler,
		FileHandler:       fileHandler,
		AuthHandler:       authHandler,
		OAuthHandler:      oauthHandler,
		OIDCHandler:       oidcHandler,
		RoleHandler:       roleHandler,
		PermissionHandler: permHandler,
		UserHandler:       userHandler,
		AppGroupHandler:   appGroupHandler,
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
