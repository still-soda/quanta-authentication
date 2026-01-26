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

// 用户管理相关权限
const (
	// 创建用户权限
	UserCreate = "user_create"
	// 删除用户权限
	UserDelete = "user_delete"
	// 查看用户权限
	UserView = "user_view"
	// 查看用户列表权限
	UserList = "user_list"
	// 更新用户权限
	UserUpdate = "user_update"
	// 为用户分配角色权限
	UserAssignRoles = "user_assign_roles"
	// 从用户撤销角色权限
	UserRevokeRoles = "user_revoke_roles"
	// 重置用户密码权限
	UserResetPassword = "user_reset_password"
)

// 应用组权限管理相关权限
const (
	// 查看应用组权限
	AppGroupPermissionView = "app_group_permission_view"
	// 创建应用组权限
	AppGroupPermissionCreate = "app_group_permission_create"
	// 更新应用组权限
	AppGroupPermissionUpdate = "APP_GROUP_PERMISSION_UPDATE"
	// 删除应用组权限
	AppGroupPermissionDelete = "app_group_permission_delete"
)

// 应用组角色管理相关权限
const (
	// 查看应用组角色
	AppGroupRoleView = "app_group_role_view"
	// 创建应用组角色
	AppGroupRoleCreate = "APP_GROUP_ROLE_CREATE"
	// 更新应用组角色
	AppGroupRoleUpdate = "APP_GROUP_ROLE_UPDATE"
	// 删除应用组角色
	AppGroupRoleDelete = "APP_GROUP_ROLE_DELETE"
	// 为应用组角色分配权限
	AppGroupRoleAssignPermissions = "APP_GROUP_ROLE_ASSIGN_PERMISSIONS"
	// 从应用组角色撤销权限
	AppGroupRoleRevokePermissions = "app_group_role_revoke_permissions"
	// 为用户分配应用组角色
	AppGroupRoleAssignToUser = "app_group_role_assign_to_user"
	// 从用户撤销应用组角色
	AppGroupRoleRevokeFromUser = "app_group_role_revoke_from_user"
)

// 应用组管理员相关权限
const (
	// 查看应用组管理员
	AppGroupAdminView = "app_group_admin_view"
	// 添加应用组管理员
	AppGroupAdminCreate = "app_group_admin_create"
	// 移除应用组管理员
	AppGroupAdminDelete = "app_group_admin_delete"
)
