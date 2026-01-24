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

// 角色相关权限
const (
	// 创建角色权限
	RoleCreate = "role_create"
	// 删除角色权限
	RoleDelete = "role_delete"
	// 查看角色权限
	RoleView = "role_view"
	// 更新角色权限
	RoleUpdate = "role_update"
	// 角色授权权限
	RoleAssignPermissions = "role_assign_permissions"
	// 角色取消授权权限
	RoleRevokePermissions = "role_revoke_permissions"
)

// 审计日志相关权限
const (
	// 查看审计日志权限
	AuditView = "audit_view"
	// 导出审计日志权限
	AuditExport = "audit_export"
)

// 权限管理相关权限
const (
	// 创建权限
	PermissionCreate = "permission_create"
	// 删除权限
	PermissionDelete = "permission_delete"
	// 查看权限
	PermissionView = "permission_view"
	// 更新权限
	PermissionUpdate = "permission_update"
)
