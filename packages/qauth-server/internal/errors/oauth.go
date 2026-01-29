package app_error

import "qauth-server/pkg/app_error"

// OAuth 相关错误
var (
	ErrOAuthClientNotFound = app_error.NewAppError(404, "oauth client not found")
)
