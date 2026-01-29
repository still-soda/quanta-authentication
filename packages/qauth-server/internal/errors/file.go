package app_error

import "qauth-server/pkg/app_error"

var (
	ErrFailedToCreateFile = app_error.NewAppError(500, "failed to create file")
)
