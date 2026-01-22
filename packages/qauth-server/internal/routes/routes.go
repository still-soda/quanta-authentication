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

	r.Use(
		middleware.Logger(),
		middleware.Recovery(),
	)

	// 公开路由（无需认证）
	openGroup := r.Group("/")
	{
		openGroup.Use(middleware.CORS())

		// 健康检查端点
		openGroup.GET("/health", healthHandler.Check)

		// OIDC 发现端点
		openGroup.GET("/.well-known/openid-configuration", oidcHandler.GetOpenIDConfiguration)
		openGroup.GET("/.well-known/jwks.json", oidcHandler.GetJWKS)
	}

	// OAuth2 路由（对外）
	oauthGroup := r.Group("/v1/oauth")
	{
		// OAuth2 授权端点
		oauthGroup.GET("/authorize", oauthHandler.AuthorizePage)
		oauthGroup.GET("/authorize/info", oauthHandler.AuthorizeInfo)
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
	systemGroup := r.Group("/_/v1")
	{
		// 认证路由
		// /system/v1/auth
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
			// /system/v1/clients
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
			// /system/v1/jwks
			jwksGroup := authRequiredGroup.Group("/jwks")
			{
				jwksGroup.GET("/keys", oidcHandler.GetKeyRotationInfo)
				jwksGroup.POST("/rotate", oidcHandler.ForceKeyRotation)
			}

			// 角色管理
			// /system/v1/roles
			roleGroup := authRequiredGroup.Group("/roles")
			{
				roleGroup.GET("", roleHandler.GetRoles)
				roleGroup.GET("/:id", roleHandler.GetRole)
				roleGroup.PUT("/:id", roleHandler.UpdateRole)
				roleGroup.DELETE("/:id", roleHandler.DeleteRole)
			}

			// 权限管理
			// /system/v1/permissions
			permissionGroup := authRequiredGroup.Group("/permissions")
			{
				permissionGroup.POST("/grant-to-role", roleHandler.GrantPermissionsToRole)
				permissionGroup.POST("/revoke-from-role", roleHandler.RevokePermissionsFromRole)
			}

			// 资源管理路由
			// /system/v1/resources
			resourceGroup := authRequiredGroup.Group("/resources")
			{
				resourceGroup.POST("/upload", fileHandler.Upload)
			}
		}
	}

}
