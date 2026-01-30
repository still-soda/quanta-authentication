package app_error

// 用户管理相关错误
var (
	ErrUserNotFound             = NewAppError(404, "user not found", "用户不存在")
	ErrCreateUserConflict       = NewAppError(409, "user with the given Student ID or Email already exists", "学号或邮箱已被使用")
	ErrFailedToCreateUser       = NewAppError(500, "failed to create user", "创建用户失败")
	ErrFailedToGetUsers         = NewAppError(500, "failed to get users", "获取用户列表失败")
	ErrFailedToGetUser          = NewAppError(500, "failed to get user", "获取用户信息失败")
	ErrFailedToUpdateUser       = NewAppError(500, "failed to update user", "更新用户信息失败")
	ErrFailedToDeleteUser       = NewAppError(500, "failed to delete user", "删除用户失败")
	ErrFailedToSetUserRoles     = NewAppError(500, "failed to set user roles", "设置用户角色失败")
	ErrFailedToResetPassword    = NewAppError(500, "failed to reset password", "重置密码失败")
	ErrFailedToCountUsersTotal  = NewAppError(500, "failed to count total users", "统计用户总数失败")
	ErrFailedToGetUserRoles     = NewAppError(500, "failed to get user roles", "获取用户角色失败")
	ErrFailedToAssignRoles      = NewAppError(500, "failed to assign roles", "分配角色失败")
	ErrFailedToRevokeRoles      = NewAppError(500, "failed to revoke roles", "撤销角色失败")
	ErrFailedToGenerateSalt     = NewAppError(500, "failed to generate salt", "生成盐值失败")
	ErrFailedToFindRoles        = NewAppError(500, "failed to find roles", "查找角色失败")
	ErrFailedToCountUsersByRole = NewAppError(500, "failed to count users by role", "按角色统计用户数失败")
)
