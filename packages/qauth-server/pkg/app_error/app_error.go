package app_error

import "fmt"

type AppError struct {
	Code    int
	Message string
}

// NewAppError 创建一个新的 AppError 实例
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	return e.Message
}

// Wrap 包装一个错误，保留原始错误信息
func (e *AppError) Wrap(err error) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message + ": " + err.Error(),
	}
}

// Wrapf 包装一个格式化的错误消息
func (e *AppError) Wrapf(format string, args ...any) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message + ": " + fmt.Sprintf(format, args...),
	}
}

// WithDetails 添加详细信息到错误消息中
func (e *AppError) WithDetails(args ...any) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message + ": " + fmt.Sprint(args...),
	}
}

// WithMessage 添加自定义消息到错误消息中
func (e *AppError) WithMessage(message string) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message + ": " + message,
	}
}
