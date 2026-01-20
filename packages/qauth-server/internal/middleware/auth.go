package middleware

import (
	"fmt"
	"qauth-server/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// Auth 认证中间件
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 Authorization 头
		authorization := ctx.GetHeader("Authorization")
		if authorization == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "authorization header is required"})
			return
		}
		// 提取 token
		var accessToken string
		_, err := fmt.Sscanf(authorization, "Bearer %s", &accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization header format"})
			return
		}

		userInfo, err := jwt.ParseAccessToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "invalid token: " + err.Error()})
			return
		}
		// 将用户信息存储在上下文中
		ctx.Set("userInfo", userInfo)
		ctx.Next()
	}
}
