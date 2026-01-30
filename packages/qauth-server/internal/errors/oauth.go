package app_error

// OAuth 相关错误
var (
	// 客户端错误
	ErrClientNotFound       = NewAppError(404, "OAuth client not found", "OAuth 客户端不存在")
	ErrClientAlreadyExists  = NewAppError(409, "OAuth client already exists", "OAuth 客户端已存在")
	ErrInvalidClientID      = NewAppError(400, "invalid client ID", "客户端 ID 无效")
	ErrInvalidClientSecret  = NewAppError(400, "invalid client secret", "客户端密钥无效")
	ErrClientCreationFailed = NewAppError(500, "failed to create OAuth client", "创建 OAuth 客户端失败")
	ErrClientUpdateFailed   = NewAppError(500, "failed to update OAuth client", "更新 OAuth 客户端失败")
	ErrClientDeletionFailed = NewAppError(500, "failed to delete OAuth client", "删除 OAuth 客户端失败")
	ErrClientQueryFailed    = NewAppError(500, "failed to query OAuth clients", "查询 OAuth 客户端失败")
	ErrNoFieldsToUpdate     = NewAppError(400, "no fields to update", "没有需要更新的字段")

	// 令牌错误
	ErrTokenGenerationFailed = NewAppError(500, "failed to generate token", "生成令牌失败")
	ErrTokenRevocationFailed = NewAppError(500, "failed to revoke token", "撤销令牌失败")
	ErrTokenLoadFailed       = NewAppError(500, "failed to load token info", "加载令牌信息失败")
	ErrInvalidToken          = NewAppError(401, "invalid token", "令牌无效")
	ErrTokenExpired          = NewAppError(401, "token expired", "令牌已过期")

	// 授权错误（使用 common.go 中定义的 ErrUnauthorized 和 user.go 中的 ErrUserNotFound）
	ErrAccessDenied       = NewAppError(403, "access denied", "访问被拒绝")
	ErrInvalidGrant       = NewAppError(400, "invalid grant", "授权类型无效")
	ErrInvalidScope       = NewAppError(400, "invalid scope", "权限范围无效")
	ErrInvalidRedirectURI = NewAppError(400, "invalid redirect URI", "重定向 URI 无效")
	ErrUserInactive       = NewAppError(403, "user is inactive", "用户未激活")
	ErrInvalidPassword    = NewAppError(401, "invalid password", "密码错误")

	// 登录状态记录错误
	ErrLoginStateRecordFailed = NewAppError(500, "failed to record login state", "记录登录状态失败")
	ErrErrorRecordFailed      = NewAppError(500, "failed to record error", "记录错误失败")

	// ID Token 错误
	ErrIDTokenGenerationFailed = NewAppError(500, "failed to generate ID token", "生成 ID Token 失败")
)
