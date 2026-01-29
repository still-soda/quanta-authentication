package app_error

import "qauth-server/pkg/app_error"

// 应用组管理员相关错误
var (
	ErrFailedToCreateOwnerAdmin      = app_error.NewAppError(500, "failed to create owner admin")
	ErrFailedToAddAppGroupAdmin      = app_error.NewAppError(500, "failed to add app group admin")
	ErrFailedToRemoveAppGroupAdmin   = app_error.NewAppError(500, "failed to remove app group admin")
	ErrFailedToGetAppGroupAdmins     = app_error.NewAppError(500, "failed to get app group admins")
	ErrFailedToGetUserAppGroupAdmins = app_error.NewAppError(500, "failed to get user's app group admins")
	ErrFailedToCheckAppGroupAdmin    = app_error.NewAppError(500, "failed to check app group admin")
)

// 应用组权限相关错误
var (
	ErrFailedToCreateAppGroupPermission = app_error.NewAppError(500, "failed to create app group permission")
	ErrFailedToFindPermission           = app_error.NewAppError(404, "failed to find permission")
	ErrFailedToUpdateAppGroupPermission = app_error.NewAppError(500, "failed to update app group permission")
	ErrFailedToDeleteAppGroupPermission = app_error.NewAppError(500, "failed to delete app group permission")
	ErrFailedToGetAppGroupPermission    = app_error.NewAppError(500, "failed to get app group permission")
	ErrFailedToGetAppGroupPermissions   = app_error.NewAppError(500, "failed to get app group permissions")
	ErrFailedToGetPermissionByCodes     = app_error.NewAppError(500, "failed to get app group permissions by codes")
	ErrFailedToFindPermissions          = app_error.NewAppError(404, "failed to find permissions")
	ErrPermissionNotBelongToClient      = app_error.NewAppError(400, "permission does not belong to client")
)

// 应用组角色相关错误
var (
	ErrFailedToCreateAppGroupRole               = app_error.NewAppError(500, "failed to create app group role")
	ErrFailedToFindRole                         = app_error.NewAppError(404, "failed to find role")
	ErrFailedToUpdateAppGroupRole               = app_error.NewAppError(500, "failed to update app group role")
	ErrFailedToDeleteAppGroupRole               = app_error.NewAppError(500, "failed to delete app group role")
	ErrFailedToGetAppGroupRole                  = app_error.NewAppError(500, "failed to get app group role")
	ErrFailedToGetAppGroupRoles                 = app_error.NewAppError(500, "failed to get app group roles")
	ErrFailedToDeleteRolePermissionAssociations = app_error.NewAppError(500, "failed to delete role permission associations")
	ErrFailedToDeleteUserRoleAssociations       = app_error.NewAppError(500, "failed to delete user role associations")
	ErrFailedToGrantPermissionToRole            = app_error.NewAppError(500, "failed to grant permission to role")
	ErrFailedToRevokePermissionsFromRole        = app_error.NewAppError(500, "failed to revoke permissions from role")
	ErrFailedToClearExistingPermissions         = app_error.NewAppError(500, "failed to clear existing permissions")
	ErrFailedToGetRolePermissions               = app_error.NewAppError(500, "failed to get role permissions")
	ErrFailedToCountRoleUsers                   = app_error.NewAppError(500, "failed to count role users")
	ErrFailedToCountPermissions                 = app_error.NewAppError(500, "failed to count permissions")
)

// 应用组用户角色相关错误
var (
	ErrFailedToAssignRoleToUser           = app_error.NewAppError(500, "failed to assign role to user")
	ErrFailedToRevokeRoleFromUser         = app_error.NewAppError(500, "failed to revoke role from user")
	ErrFailedToGetUserAppGroupRoles       = app_error.NewAppError(500, "failed to get user app group roles")
	ErrFailedToGetUserAppGroupPermissions = app_error.NewAppError(500, "failed to get user app group permissions")
	ErrFailedToGetRoleUsers               = app_error.NewAppError(500, "failed to get role users")
)

// 默认角色相关错误
var (
	ErrFailedToGetDefaultRole   = app_error.NewAppError(500, "failed to get default role")
	ErrFailedToClearDefaultRole = app_error.NewAppError(500, "failed to clear default role")
	ErrFailedToSetDefaultRole   = app_error.NewAppError(500, "failed to set default role")
)
