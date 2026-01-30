package app_error

// 审计日志相关错误
var (
	ErrFailedToGetAuditLogs   = NewAppError(500, "failed to get audit logs", "获取审计日志失败")
	ErrAuditLogNotFound       = NewAppError(404, "audit log not found", "审计日志不存在")
	ErrFailedToCreateAuditLog = NewAppError(500, "failed to create audit log", "创建审计日志失败")
	ErrInvalidTimeRange       = NewAppError(400, "invalid time range", "时间范围无效")
)
