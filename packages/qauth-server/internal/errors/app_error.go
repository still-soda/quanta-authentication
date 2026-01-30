package app_error

import (
	"fmt"
)

type ILogger interface {
	Info(message string, args ...any)
	Error(message string, args ...any)
	Warn(message string, args ...any)
	With(args ...any) ILogger
}

type AppError struct {
	Code        int
	ViewMessage string
	Message     string
	Scope       string
	Cause       error
}

// NewAppError 创建一个新的 AppError 实例
func NewAppError(code int, message string, viewMsg string) *AppError {
	return &AppError{
		Code:        code,
		ViewMessage: viewMsg,
		Message:     message,
		Scope:       "",
		Cause:       nil,
	}
}

// clone 创建当前错误的副本
func (e *AppError) clone() *AppError {
	return &AppError{
		Code:        e.Code,
		ViewMessage: e.ViewMessage,
		Message:     e.Message,
		Scope:       e.Scope,
		Cause:       e.Cause,
	}
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	msg := e.Message
	if e.Scope != "" {
		msg = fmt.Sprintf("[%s] %s", e.Scope, msg)
	}
	if e.Cause != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.Cause)
	}
	return msg
}

// LogDetails 使用 logger 记录错误的详细信息
func (e *AppError) LogDetails(logger ILogger) {
	logger.Error(e.ErrorDetails())
}

// ErrorDetails 返回错误的详细信息，包括作用域和原始错误
func (e *AppError) ErrorDetails() string {
	if e.Cause != nil {
		return fmt.Sprintf("Code: %d, Message: %s, ViewMessage: %s, Scope: %s, Cause: %v", e.Code, e.Message, e.ViewMessage, e.Scope, e.Cause)
	}
	return fmt.Sprintf("Code: %d, Message: %s, ViewMessage: %s, Scope: %s", e.Code, e.Message, e.ViewMessage, e.Scope)
}

// SendError 返回适合发送给客户端的错误信息
func (e *AppError) SendError() string {
	return e.ViewMessage
}

// WithScope 添加作用域信息到错误消息中
func (e *AppError) WithScope(scope string) *AppError {
	newErr := e.clone()
	newErr.Scope = scope
	newErr.Message = fmt.Sprintf("[%s] %s", scope, e.Message)
	return newErr
}

// Wrap 包装一个错误，保留原始错误信息
func (e *AppError) Wrap(err error) *AppError {
	newErr := e.clone()
	newErr.Cause = err
	return newErr
}

// Unwrap 返回原始错误
func (e *AppError) Unwrap() error {
	return e.Cause
}

// Wrapf 包装一个格式化的错误消息
func (e *AppError) Wrapf(format string, args ...any) *AppError {
	newErr := e.clone()
	newErr.Cause = fmt.Errorf(format, args...)
	return newErr
}

// WithDetails 添加详细信息到错误消息中
func (e *AppError) WithDetails(args ...any) *AppError {
	newErr := e.clone()
	newErr.Message = fmt.Sprintf("%s: %s", e.Message, fmt.Sprint(args...))
	return newErr
}

// WithMessage 添加自定义消息到错误消息中
func (e *AppError) WithMessage(message string) *AppError {
	newErr := e.clone()
	newErr.Message = fmt.Sprintf("%s: %s", e.Message, message)
	return newErr
}

// WithViewMessage 添加自定义展示消息
func (e *AppError) WithViewMessage(viewMsg string) *AppError {
	newErr := e.clone()
	newErr.ViewMessage = viewMsg
	return newErr
}
