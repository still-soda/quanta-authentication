package app_error

// 角色管理相关错误
var (
	ErrRoleNotFound              = NewAppError(404, "role not found", "角色不存在")
	ErrRoleCodeAlreadyExists     = NewAppError(409, "role code already exists", "角色代码已存在")
	ErrFailedToCreateRole        = NewAppError(500, "failed to create role", "创建角色失败")
	ErrFailedToUpdateRole        = NewAppError(500, "failed to update role", "更新角色失败")
	ErrFailedToDeleteRole        = NewAppError(500, "failed to delete role", "删除角色失败")
	ErrFailedToGetRole           = NewAppError(500, "failed to get role", "获取角色失败")
	ErrFailedToListRoles         = NewAppError(500, "failed to list roles", "获取角色列表失败")
	ErrFailedToCountUsers        = NewAppError(500, "failed to count users by role", "统计角色用户数失败")
	ErrFailedToCountPermissions  = NewAppError(500, "failed to count permissions by role", "统计角色权限数失败")
	ErrFailedToGrantPermissions  = NewAppError(500, "failed to grant permissions to role", "分配权限到角色失败")
	ErrFailedToRevokePermissions = NewAppError(500, "failed to revoke permissions from role", "撤销角色权限失败")
	ErrFailedToAssignRole        = NewAppError(500, "failed to assign role to user", "分配角色到用户失败")
	ErrFailedToRevokeRole        = NewAppError(500, "failed to revoke role from user", "撤销用户角色失败")
	ErrFailedToCheckPermissions  = NewAppError(500, "failed to check role permissions", "检查角色权限失败")
	ErrFailedToDeleteRolePerms   = NewAppError(500, "failed to delete role permissions", "删除角色权限关联失败")
	ErrFailedToDeleteUserRoles   = NewAppError(500, "failed to delete user roles", "删除用户角色关联失败")
)
