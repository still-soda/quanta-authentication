package middleware

import "github.com/gin-gonic/gin"

func UseMiddleware(r *gin.Engine) {
	r.Use(Logger())
	r.Use(Recovery())
	r.Use(CORS())
}
