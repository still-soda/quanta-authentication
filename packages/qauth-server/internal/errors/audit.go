package app_error

import "qauth-server/pkg/app_error"

// 审计日志相关错误
var (
	ErrFailedToGetAuditLogs   = app_error.NewAppError(500, "failed to get audit logs")
	ErrAuditLogNotFound       = app_error.NewAppError(404, "audit log not found")
	ErrFailedToCreateAuditLog = app_error.NewAppError(500, "failed to create audit log")
	ErrInvalidTimeRange       = app_error.NewAppError(400, "invalid time range")
)
