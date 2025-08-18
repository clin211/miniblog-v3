package errorx

import "fmt"

// Error 表示带业务码的错误。
//
// Error 满足 error 接口，包含：
//   - Code: 业务码
//   - Message: 面向用户或调用方的提示
//   - Reason: 结构化原因，可用于承载上下文（例如 map 或切片），用于响应 reason 字段
//   - Cause: 源错误，便于链路追踪
//   - Meta: 额外附加信息，不参与 Error() 文本输出
type Error struct {
	Code    Code
	Message string
	Reason  interface{}
	Cause   error
	Meta    map[string]interface{}
}

// Error 实现 error 接口。
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Option 为可选项配置函数。
type Option func(*Error)

// WithReason 设置结构化 reason。
func WithReason(reason interface{}) Option {
	return func(e *Error) { e.Reason = reason }
}

// WithCause 设置底层 cause。
func WithCause(cause error) Option {
	return func(e *Error) { e.Cause = cause }
}

// WithMeta 设置元数据（覆盖式）。
func WithMeta(meta map[string]interface{}) Option {
	return func(e *Error) { e.Meta = meta }
}

// New 创建新的业务错误。
func New(code Code, message string, opts ...Option) *Error {
	e := &Error{Code: code, Message: message}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Wrap 使用指定 code 与 message 包装底层错误。
func Wrap(err error, code Code, message string, opts ...Option) *Error {
	e := &Error{Code: code, Message: message, Cause: err}
	for _, opt := range opts {
		opt(e)
	}
	return e
}
