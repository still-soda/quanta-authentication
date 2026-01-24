package app_error

import (
	"errors"
	"qauth-server/pkg/app_error"

	"github.com/jackc/pgx/v5/pgconn"
)

// IsErrorWithPgsqlCode 判断错误是否包含特定的PostgreSQL错误代码
func IsErrorWithPgsqlCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == code {
		return true
	}
	return false
}

var (
	ErrBadRequest            = app_error.NewAppError(400, "bad request")
	ErrUnauthorized          = app_error.NewAppError(401, "unauthorized")
	ErrNotFound              = app_error.NewAppError(404, "not found")
	ErrUserNotFound          = app_error.NewAppError(404, "user not found")
	ErrCreateUserConflict    = app_error.NewAppError(409, "user with the given Student ID or Email already exists")
	ErrInternalServerError   = app_error.NewAppError(500, "internal server error")
	ErrFailedToCreateUser    = app_error.NewAppError(500, "failed to create user")
	ErrAuthenticationFailed  = app_error.NewAppError(401, "authentication failed")
	ErrFailedToGetUserRole   = app_error.NewAppError(500, "failed to get user role")
	ErrFailedToGenerateToken = app_error.NewAppError(500, "failed to generate token")
	ErrInvalidRefreshToken   = app_error.NewAppError(401, "invalid refresh token")
	ErrInvalidAccessToken    = app_error.NewAppError(401, "invalid access token")
	ErrNoPermission          = app_error.NewAppError(403, "no permission")
	ErrAPINotImplemented     = app_error.NewAppError(501, "api not implemented")
	ErrInvalidResponseType   = app_error.NewAppError(400, "invalid response type")
	ErrFailedToGetRoles      = app_error.NewAppError(500, "failed to get roles")
	ErrPermissionNoExist     = app_error.NewAppError(400, "permission does not exist")
	ErrRoleNoExist           = app_error.NewAppError(400, "role does not exist")

	// 审计日志相关错误
	ErrFailedToGetAuditLogs   = app_error.NewAppError(500, "failed to get audit logs")
	ErrAuditLogNotFound       = app_error.NewAppError(404, "audit log not found")
	ErrFailedToCreateAuditLog = app_error.NewAppError(500, "failed to create audit log")

	// 用户管理相关错误
	ErrFailedToGetUsers      = app_error.NewAppError(500, "failed to get users")
	ErrFailedToUpdateUser    = app_error.NewAppError(500, "failed to update user")
	ErrFailedToDeleteUser    = app_error.NewAppError(500, "failed to delete user")
	ErrFailedToSetUserRoles  = app_error.NewAppError(500, "failed to set user roles")
	ErrFailedToResetPassword = app_error.NewAppError(500, "failed to reset password")
)
