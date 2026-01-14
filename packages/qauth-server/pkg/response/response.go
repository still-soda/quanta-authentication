package response

import (
	"errors"
	"qauth-server/internal/utilities"
	"qauth-server/pkg/app_error"

	"github.com/gin-gonic/gin"
)

func HandlerError(ctx *gin.Context, err error) {
	var appErr *app_error.AppError
	if errors.As(err, &appErr) {
		ctx.JSON(appErr.Code, gin.H{
			"code":  appErr.Code,
			"msg":   "error",
			"error": appErr.Message,
		})
		return
	}

	logger := ctx.MustGet("logger").(utilities.Logger)
	logger.Error("internal server error: ", err)
	ctx.JSON(500, gin.H{
		"code":  500,
		"msg":   "error",
		"error": "internal server error",
	})
}

func HandlerSuccess(ctx *gin.Context, data any) {
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": data,
	})
}
