package app_error

import "qauth-server/pkg/app_error"

// 用户管理相关错误
var (
	ErrUserNotFound          = app_error.NewAppError(404, "user not found")
	ErrCreateUserConflict    = app_error.NewAppError(409, "user with the given Student ID or Email already exists")
	ErrFailedToCreateUser    = app_error.NewAppError(500, "failed to create user")
	ErrFailedToGetUsers      = app_error.NewAppError(500, "failed to get users")
	ErrFailedToUpdateUser    = app_error.NewAppError(500, "failed to update user")
	ErrFailedToDeleteUser    = app_error.NewAppError(500, "failed to delete user")
	ErrFailedToSetUserRoles  = app_error.NewAppError(500, "failed to set user roles")
	ErrFailedToResetPassword = app_error.NewAppError(500, "failed to reset password")
)
