package app_error

import (
	"errors"

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
