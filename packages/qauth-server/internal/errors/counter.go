package app_error

var (
	ErrFailedToFindCounters = NewAppError(500, "failed to find counter", "获取计数器失败")
)
