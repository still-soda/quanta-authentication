package main

import (
	"context"
	"log"
	"net/http"
	"qauth-server/internal/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"

	oredis "github.com/go-oauth2/redis/v4"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
)


func main() {
	cfg := config.New()

	// 初始化数据库连接
	pgxConn, _ := pgx.Connect(context.TODO(), cfg.DatabaseURL)
	adapter := pgx4adapter.NewConn(pgxConn)
	pgStore, _ := pg.NewClientStore(adapter)
	
	// 初始化 Redis 令牌存储
	tokenStore := oredis.NewRedisStore(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	defer tokenStore.Close()

	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(pgStore)

	// 设置令牌有效期
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.SetPasswordTokenCfg(&manage.Config{AccessTokenExp: cfg.AccessTokenExpire})

	// 初始化 OAuth 服务
	srv := server.NewDefaultServer(manager)

	r := gin.Default()

	r.GET("/authorize", func(c *gin.Context) {
		err := srv.HandleAuthorizeRequest(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	r.POST("/token", func(c *gin.Context) {
		err := srv.HandleTokenRequest(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	r.GET("/userinfo", func(ctx *gin.Context) {
		ti, err := srv.ValidationBearerToken(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user_id": ti.GetUserID(),
			"client_id": ti.GetClientID(),
			"expires_in": int64(time.Until(ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn())).Seconds()),
		})
	})

	log.Println("[INFO] 🌍 Quanta Auth Server 运行在 :" + cfg.Port)
	r.Run(":" + cfg.Port)
}