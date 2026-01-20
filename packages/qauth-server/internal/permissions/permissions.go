package permissions

// OAuth 客户端相关权限
const (
	// 创建 OAuth 客户端权限
	OAuthClientCreate = "oauth_client_create"
	// 删除 OAuth 客户端权限
	OAuthClientDelete = "oauth_client_delete"
	// 查看 OAuth 客户端权限
	OAuthClientView = "oauth_client_view"
	// 更新 OAuth 客户端权限
	OAuthClientUpdate = "oauth_client_update"
	// 列出 OAuth 客户端权限
	OAuthClientList = "oauth_client_list"
)

// 系统管理相关权限
const (
	// 获取密钥轮换信息权限
	SystemKeyRotationView = "system_key_rotation_view"
	// 执行密钥轮换权限
	SystemKeyRotationExecute = "system_key_rotation_execute"
)
