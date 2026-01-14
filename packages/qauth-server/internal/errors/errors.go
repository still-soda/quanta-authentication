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
	ErrBadRequest           = app_error.NewAppError(400, "bad request")
	ErrUserNotFound         = app_error.NewAppError(404, "user not found")
	ErrCreateUserConflict   = app_error.NewAppError(409, "user with the given Student ID or Email already exists")
	ErrFailedToCreateUser   = app_error.NewAppError(500, "failed to create user")
	ErrAuthenticationFailed = app_error.NewAppError(401, "authentication failed")
)
