package app_error

// 权限和角色相关错误
var (
	ErrNoPermission                      = NewAppError(403, "no permission", "权限不足")
	ErrPermissionNoExist                 = NewAppError(400, "permission does not exist", "权限不存在")
	ErrRoleNoExist                       = NewAppError(400, "role does not exist", "角色不存在")
	ErrFailedToGetRoles                  = NewAppError(500, "failed to get roles", "获取角色列表失败")
	ErrFailedToGetUserRole               = NewAppError(500, "failed to get user role", "获取用户角色失败")
	ErrFailedToFindPermissionsByCodes    = NewAppError(500, "failed to find permissions by codes", "根据权限代码查找权限失败")
	ErrFailedToCountPermissionsByCodes   = NewAppError(500, "failed to count permissions by codes", "统计权限失败")
	ErrFailedToFindAllPermissions        = NewAppError(500, "failed to find all permissions", "获取所有权限失败")
	ErrFailedToListPermissions           = NewAppError(500, "failed to list permissions", "列出权限失败")
	ErrFailedToFindPermissionByID        = NewAppError(500, "failed to find permission by id", "查找权限失败")
	ErrFailedToFindPermissionsByResource = NewAppError(500, "failed to find permissions by resource", "根据资源查找权限失败")
	ErrFailedToCreatePermission          = NewAppError(500, "failed to create permission", "创建权限失败")
	ErrFailedToUpdatePermission          = NewAppError(500, "failed to update permission", "更新权限失败")
	ErrFailedToDeleteRolePermissions     = NewAppError(500, "failed to delete role permissions", "删除角色权限失败")
	ErrFailedToDeletePermission          = NewAppError(500, "failed to delete permission", "删除权限失败")
)
