package permissions

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

// CURD 操作枚举
type Action int8

const (
	Create Action = iota
	Read
	Update
	Delete
)
