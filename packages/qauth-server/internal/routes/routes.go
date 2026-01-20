package routes

import (
	_handlers "qauth-server/internal/handlers"
	"qauth-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

type RegisterRouterHandlers struct {
	HealthHandler *_handlers.HealthHandler
	FileHandler   *_handlers.FileHandler
	AuthHandler   *_handlers.AuthHandler
	OAuthHandler  *_handlers.OAuthHandler
	OIDCHandler   *_handlers.OIDCHandler
}

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine, handlers *RegisterRouterHandlers) {
	healthHandler := handlers.HealthHandler
	fileHandler := handlers.FileHandler
	authHandler := handlers.AuthHandler
	oauthHandler := handlers.OAuthHandler
	oidcHandler := handlers.OIDCHandler

	r.GET("/health", healthHandler.Check)
	r.POST("/upload", fileHandler.Upload)

	// OIDC 发现端点
	r.GET("/.well-known/openid-configuration", oidcHandler.GetOpenIDConfiguration)
	r.GET("/.well-known/jwks.json", oidcHandler.GetJWKS)

	// 认证路由
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh-token", authHandler.RefreshToken)
	}

	// OAuth2 路由
	oauthGroup := r.Group("/oauth")
	{
		// OAuth2 授权端点
		oauthGroup.GET("/authorize", oauthHandler.Authorize)
		oauthGroup.POST("/authorize", oauthHandler.Authorize)
		oauthGroup.POST("/token", oauthHandler.Token)
		oauthGroup.POST("/validate", oauthHandler.ValidateToken)
		oauthGroup.POST("/revoke", oauthHandler.RevokeToken)

		// OIDC 端点
		oauthGroup.GET("/userinfo", oidcHandler.GetUserInfo)
		oauthGroup.GET("/logout", oidcHandler.Logout)
		oauthGroup.POST("/logout", oidcHandler.Logout)

		// OAuth2 客户端管理
		clientGroup := oauthGroup.Group("/clients")
		{
			clientGroup.Use(middleware.Auth())

			clientGroup.GET("", oauthHandler.ListClients)
			clientGroup.POST("", oauthHandler.CreateClient)
			clientGroup.GET("/:id", oauthHandler.GetClient)
			clientGroup.PUT("/:id", oauthHandler.UpdateClient)
			clientGroup.DELETE("/:id", oauthHandler.DeleteClient)
		}
	}

	// 管理员路由
	adminGroup := r.Group("/admin")
	{
		// JWKS 密钥管理
		jwksGroup := adminGroup.Group("/jwks")
		{
			jwksGroup.GET("/keys", oidcHandler.GetKeyRotationInfo)
			jwksGroup.POST("/rotate", oidcHandler.ForceKeyRotation)
		}
	}
}
