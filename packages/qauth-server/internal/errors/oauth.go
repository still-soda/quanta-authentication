package app_error

import "qauth-server/pkg/app_error"

// OAuth 相关错误
var (
	// 客户端错误
	ErrClientNotFound       = app_error.NewAppError(404, "OAuth client not found")
	ErrClientAlreadyExists  = app_error.NewAppError(409, "OAuth client already exists")
	ErrInvalidClientID      = app_error.NewAppError(400, "invalid client ID")
	ErrInvalidClientSecret  = app_error.NewAppError(400, "invalid client secret")
	ErrClientCreationFailed = app_error.NewAppError(500, "failed to create OAuth client")
	ErrClientUpdateFailed   = app_error.NewAppError(500, "failed to update OAuth client")
	ErrClientDeletionFailed = app_error.NewAppError(500, "failed to delete OAuth client")
	ErrClientQueryFailed    = app_error.NewAppError(500, "failed to query OAuth clients")
	ErrNoFieldsToUpdate     = app_error.NewAppError(400, "no fields to update")

	// 令牌错误
	ErrTokenGenerationFailed = app_error.NewAppError(500, "failed to generate token")
	ErrTokenRevocationFailed = app_error.NewAppError(500, "failed to revoke token")
	ErrTokenLoadFailed       = app_error.NewAppError(500, "failed to load token info")
	ErrInvalidToken          = app_error.NewAppError(401, "invalid token")
	ErrTokenExpired          = app_error.NewAppError(401, "token expired")

	// 授权错误（使用 common.go 中定义的 ErrUnauthorized 和 user.go 中的 ErrUserNotFound）
	ErrAccessDenied       = app_error.NewAppError(403, "access denied")
	ErrInvalidGrant       = app_error.NewAppError(400, "invalid grant")
	ErrInvalidScope       = app_error.NewAppError(400, "invalid scope")
	ErrInvalidRedirectURI = app_error.NewAppError(400, "invalid redirect URI")
	ErrUserInactive       = app_error.NewAppError(403, "user is inactive")
	ErrInvalidPassword    = app_error.NewAppError(401, "invalid password")

	// 登录状态记录错误
	ErrLoginStateRecordFailed = app_error.NewAppError(500, "failed to record login state")
	ErrErrorRecordFailed      = app_error.NewAppError(500, "failed to record error")

	// ID Token 错误
	ErrIDTokenGenerationFailed = app_error.NewAppError(500, "failed to generate ID token")
)
