package app_error

import "qauth-server/pkg/app_error"

// 权限和角色相关错误
var (
	ErrNoPermission        = app_error.NewAppError(403, "no permission")
	ErrPermissionNoExist   = app_error.NewAppError(400, "permission does not exist")
	ErrRoleNoExist         = app_error.NewAppError(400, "role does not exist")
	ErrFailedToGetRoles    = app_error.NewAppError(500, "failed to get roles")
	ErrFailedToGetUserRole = app_error.NewAppError(500, "failed to get user role")
)
