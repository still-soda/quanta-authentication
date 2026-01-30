package app_error

// 通用错误定义
var (
	ErrBadRequest          = NewAppError(400, "bad request", "请求参数错误")
	ErrUnauthorized        = NewAppError(401, "unauthorized", "未授权")
	ErrNotFound            = NewAppError(404, "not found", "资源不存在")
	ErrInternalServerError = NewAppError(500, "internal server error", "服务器内部错误")
	ErrAPINotImplemented   = NewAppError(501, "api not implemented", "接口暂未实现")

	ErrInvalidParameter = NewAppError(400, "invalid parameter", "参数格式错误")
)
