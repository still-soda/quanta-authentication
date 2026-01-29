package app_error

import "qauth-server/pkg/app_error"

var (
	ErrFailedToFindCounters = app_error.NewAppError(500, "failed to find counter")
)
