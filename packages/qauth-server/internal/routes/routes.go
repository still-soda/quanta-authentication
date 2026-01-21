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
	RoleHandler   *_handlers.RoleHandler
}

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine, handlers *RegisterRouterHandlers) {
	healthHandler := handlers.HealthHandler
	fileHandler := handlers.FileHandler
	authHandler := handlers.AuthHandler
	oauthHandler := handlers.OAuthHandler
	oidcHandler := handlers.OIDCHandler
	roleHandler := handlers.RoleHandler

	// 健康检查端点
	r.GET("/health", healthHandler.Check)

	// OIDC 发现端点
	r.GET("/.well-known/openid-configuration", oidcHandler.GetOpenIDConfiguration)
	r.GET("/.well-known/jwks.json", oidcHandler.GetJWKS)

	// OAuth2 路由（对外）
	oauthGroup := r.Group("/v1/oauth")
	{
		// OAuth2 授权端点
		oauthGroup.GET("/authorize", oauthHandler.Authorize)
		oauthGroup.POST("/authorize", oauthHandler.Authorize)
		oauthGroup.POST("/token", oauthHandler.Token)
		oauthGroup.POST("/validate", oauthHandler.ValidateToken)
		oauthGroup.POST("/revoke", oauthHandler.RevokeToken)

		// OIDC 端点
		oauthGroup.GET("/userinfo", oauthHandler.UserInfo)
		oauthGroup.GET("/logout", oauthHandler.Logout)
		oauthGroup.POST("/logout", oauthHandler.Logout)
	}

	// system 路由（系统相关）
	systemGroup := r.Group("/system/v1")
	{
		// 认证路由
		authGroup := systemGroup.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh-token", authHandler.RefreshToken)
		}

		// 需要认证的路由
		authRequiredGroup := systemGroup.Group("/")
		{
			// 使用认证中间件
			authRequiredGroup.Use(middleware.Auth())

			// OAuth2 客户端管理
			clientGroup := authRequiredGroup.Group("/clients")
			{
				clientGroup.Use(middleware.Auth())

				clientGroup.GET("", oauthHandler.ListClients)
				clientGroup.POST("", oauthHandler.CreateClient)
				clientGroup.GET("/:id", oauthHandler.GetClient)
				clientGroup.PUT("/:id", oauthHandler.UpdateClient)
				clientGroup.DELETE("/:id", oauthHandler.DeleteClient)
			}

			// JWKS 密钥管理
			jwksGroup := authRequiredGroup.Group("/jwks")
			{
				jwksGroup.GET("/keys", oidcHandler.GetKeyRotationInfo)
				jwksGroup.POST("/rotate", oidcHandler.ForceKeyRotation)
			}

			// 角色管理
			roleGroup := authRequiredGroup.Group("/roles")
			{
				roleGroup.GET("", roleHandler.GetRoles)
				roleGroup.GET("/:id", roleHandler.GetRole)
				roleGroup.PUT("/:id", roleHandler.UpdateRole)
				roleGroup.DELETE("/:id", roleHandler.DeleteRole)
			}

			// 权限管理
			permissionGroup := authRequiredGroup.Group("/permissions")
			{
				permissionGroup.POST("/grant-to-role", roleHandler.GrantPermissionsToRole)
				permissionGroup.POST("/revoke-from-role", roleHandler.RevokePermissionsFromRole)
			}

			// 资源管理路由
			resourceGroup := authRequiredGroup.Group("/resources")
			{
				resourceGroup.POST("/upload", fileHandler.Upload)
			}
		}
	}

}
