package middleware

import (
	"qauth-server/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// Auth 认证中间件
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := jwt.ExtractTokenFromHeader(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{
				"msg":   "authorization failed",
				"code":  401,
				"error": err.Error(),
			})
			return
		}

		userInfo, err := jwt.ParseAccessToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{
				"msg":   "invalid token",
				"code":  401,
				"error": err.Error(),
			})
			return
		}
		// 将用户信息存储在上下文中
		ctx.Set("userInfo", userInfo)
		ctx.Next()
	}
}
