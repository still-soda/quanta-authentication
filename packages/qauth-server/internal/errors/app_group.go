package app_error

// 应用组管理员相关错误
var (
	ErrFailedToCreateOwnerAdmin      = NewAppError(500, "failed to create owner admin", "创建所有者管理员失败")
	ErrFailedToAddAppGroupAdmin      = NewAppError(500, "failed to add app group admin", "添加应用组管理员失败")
	ErrFailedToRemoveAppGroupAdmin   = NewAppError(500, "failed to remove app group admin", "移除应用组管理员失败")
	ErrFailedToGetAppGroupAdmins     = NewAppError(500, "failed to get app group admins", "获取应用组管理员失败")
	ErrFailedToGetUserAppGroupAdmins = NewAppError(500, "failed to get user's app group admins", "获取用户应用组管理员失败")
	ErrFailedToCheckAppGroupAdmin    = NewAppError(500, "failed to check app group admin", "检查应用组管理员失败")
)

// 应用组权限相关错误
var (
	ErrFailedToCreateAppGroupPermission = NewAppError(500, "failed to create app group permission", "创建应用组权限失败")
	ErrFailedToFindPermission           = NewAppError(404, "failed to find permission", "权限不存在")
	ErrFailedToUpdateAppGroupPermission = NewAppError(500, "failed to update app group permission", "更新应用组权限失败")
	ErrFailedToDeleteAppGroupPermission = NewAppError(500, "failed to delete app group permission", "删除应用组权限失败")
	ErrFailedToGetAppGroupPermission    = NewAppError(500, "failed to get app group permission", "获取应用组权限失败")
	ErrFailedToGetAppGroupPermissions   = NewAppError(500, "failed to get app group permissions", "获取应用组权限列表失败")
	ErrFailedToGetPermissionByCodes     = NewAppError(500, "failed to get app group permissions by codes", "根据权限代码获取权限失败")
	ErrFailedToFindPermissions          = NewAppError(404, "failed to find permissions", "权限不存在")
	ErrPermissionNotBelongToClient      = NewAppError(400, "permission does not belong to client", "权限不属于此应用")
)

// 应用组角色相关错误
var (
	ErrFailedToCreateAppGroupRole               = NewAppError(500, "failed to create app group role", "创建应用组角色失败")
	ErrFailedToFindRole                         = NewAppError(404, "failed to find role", "角色不存在")
	ErrFailedToUpdateAppGroupRole               = NewAppError(500, "failed to update app group role", "更新应用组角色失败")
	ErrFailedToDeleteAppGroupRole               = NewAppError(500, "failed to delete app group role", "删除应用组角色失败")
	ErrFailedToGetAppGroupRole                  = NewAppError(500, "failed to get app group role", "获取应用组角色失败")
	ErrFailedToGetAppGroupRoles                 = NewAppError(500, "failed to get app group roles", "获取应用组角色列表失败")
	ErrFailedToDeleteRolePermissionAssociations = NewAppError(500, "failed to delete role permission associations", "删除角色权限关联失败")
	ErrFailedToDeleteUserRoleAssociations       = NewAppError(500, "failed to delete user role associations", "删除用户角色关联失败")
	ErrFailedToGrantPermissionToRole            = NewAppError(500, "failed to grant permission to role", "为角色授予权限失败")
	ErrFailedToRevokePermissionsFromRole        = NewAppError(500, "failed to revoke permissions from role", "撤销角色权限失败")
	ErrFailedToClearExistingPermissions         = NewAppError(500, "failed to clear existing permissions", "清除现有权限失败")
	ErrFailedToGetRolePermissions               = NewAppError(500, "failed to get role permissions", "获取角色权限失败")
	ErrFailedToCountRoleUsers                   = NewAppError(500, "failed to count role users", "统计角色用户数失败")
	ErrFailedToCountPermissions                 = NewAppError(500, "failed to count permissions", "统计权限数失败")
)

// 应用组用户角色相关错误
var (
	ErrFailedToAssignRoleToUser           = NewAppError(500, "failed to assign role to user", "为用户分配角色失败")
	ErrFailedToRevokeRoleFromUser         = NewAppError(500, "failed to revoke role from user", "撤销用户角色失败")
	ErrFailedToGetUserAppGroupRoles       = NewAppError(500, "failed to get user app group roles", "获取用户应用组角色失败")
	ErrFailedToGetUserAppGroupPermissions = NewAppError(500, "failed to get user app group permissions", "获取用户应用组权限失败")
	ErrFailedToGetRoleUsers               = NewAppError(500, "failed to get role users", "获取角色用户失败")
)

// 默认角色相关错误
var (
	ErrFailedToGetDefaultRole   = NewAppError(500, "failed to get default role", "获取默认角色失败")
	ErrFailedToClearDefaultRole = NewAppError(500, "failed to clear default role", "清除默认角色失败")
	ErrFailedToSetDefaultRole   = NewAppError(500, "failed to set default role", "设置默认角色失败")
)
