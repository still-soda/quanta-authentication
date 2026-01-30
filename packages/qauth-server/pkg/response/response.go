package response

import (
	"errors"
	app_error "qauth-server/internal/errors"
	"qauth-server/internal/utilities"

	"github.com/gin-gonic/gin"
)

// HandlerError 处理错误响应
func HandlerError(ctx *gin.Context, err error) {
	var appErr *app_error.AppError
	if errors.As(err, &appErr) {
		ctx.JSON(appErr.Code, gin.H{
			"code":  appErr.Code,
			"msg":   appErr.SendError(),
			"error": appErr.Error(),
		})
		return
	}

	logger := ctx.MustGet("logger").(utilities.Logger)
	logger.Error("internal server error: ", err)
	ctx.JSON(500, gin.H{
		"code":  500,
		"msg":   "internal server error",
		"error": err.Error(),
	})
}

// HandlerSuccess 处理成功响应
func HandlerSuccess(ctx *gin.Context, data any) {
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": data,
	})
}
