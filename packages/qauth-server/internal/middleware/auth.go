package middleware

import (
	"fmt"
	"qauth-server/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get token from header
		authorization := ctx.GetHeader("Authorization")
		if authorization == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "authorization header is required"})
			return
		}
		// extract token
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
		// set user info to context
		ctx.Set("userInfo", userInfo)
		ctx.Next()
	}
}
