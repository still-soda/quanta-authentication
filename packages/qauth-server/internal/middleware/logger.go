package middleware

import (
	"fmt"
	"qauth-server/internal/utilities"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger(logger *utilities.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Set("logger", logger)
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		utilities.GetLogger().Info(
			fmt.Sprintf("| %3d | %13v | %15s | %s | %s |",
				statusCode,
				latencyTime,
				clientIP,
				reqMethod,
				reqUri,
			),
		)
	}
}
