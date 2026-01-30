package app_error

// 认证相关错误
var (
	ErrAuthenticationFailed  = NewAppError(401, "authentication failed", "认证失败，请检查用户名或密码")
	ErrInvalidRefreshToken   = NewAppError(401, "invalid refresh token", "刷新令牌无效")
	ErrInvalidAccessToken    = NewAppError(401, "invalid access token", "访问令牌无效")
	ErrFailedToGenerateToken = NewAppError(500, "failed to generate token", "生成令牌失败")
	ErrInvalidResponseType   = NewAppError(400, "invalid response type", "响应类型无效")
)
