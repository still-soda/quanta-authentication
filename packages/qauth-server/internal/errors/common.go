package app_error

import "qauth-server/pkg/app_error"

// 通用错误定义
var (
	ErrBadRequest          = app_error.NewAppError(400, "bad request")
	ErrUnauthorized        = app_error.NewAppError(401, "unauthorized")
	ErrNotFound            = app_error.NewAppError(404, "not found")
	ErrInternalServerError = app_error.NewAppError(500, "internal server error")
	ErrAPINotImplemented   = app_error.NewAppError(501, "api not implemented")

	ErrInvalidParameter = app_error.NewAppError(400, "invalid parameter")
)
