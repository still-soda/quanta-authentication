package app_error

var (
	ErrFailedToCreateFile = NewAppError(500, "failed to create file", "创建文件失败")
)
