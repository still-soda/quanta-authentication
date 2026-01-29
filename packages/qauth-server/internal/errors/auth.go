package app_error

import "qauth-server/pkg/app_error"

// 认证相关错误
var (
	ErrAuthenticationFailed  = app_error.NewAppError(401, "authentication failed")
	ErrInvalidRefreshToken   = app_error.NewAppError(401, "invalid refresh token")
	ErrInvalidAccessToken    = app_error.NewAppError(401, "invalid access token")
	ErrFailedToGenerateToken = app_error.NewAppError(500, "failed to generate token")
	ErrInvalidResponseType   = app_error.NewAppError(400, "invalid response type")
)
