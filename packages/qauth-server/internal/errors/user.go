package app_error

// 用户管理相关错误
var (
	ErrUserNotFound          = NewAppError(404, "user not found", "用户不存在")
	ErrCreateUserConflict    = NewAppError(409, "user with the given Student ID or Email already exists", "学号或邮箱已被使用")
	ErrFailedToCreateUser    = NewAppError(500, "failed to create user", "创建用户失败")
	ErrFailedToGetUsers      = NewAppError(500, "failed to get users", "获取用户列表失败")
	ErrFailedToUpdateUser    = NewAppError(500, "failed to update user", "更新用户信息失败")
	ErrFailedToDeleteUser    = NewAppError(500, "failed to delete user", "删除用户失败")
	ErrFailedToSetUserRoles  = NewAppError(500, "failed to set user roles", "设置用户角色失败")
	ErrFailedToResetPassword = NewAppError(500, "failed to reset password", "重置密码失败")
)
