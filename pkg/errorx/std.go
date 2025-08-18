package errorx

// 预置的构造函数，便于快速创建常见错误。

func NewBadRequest(message string, opts ...Option) *Error {
	return New(CodeBadRequest, message, opts...)
}

func NewUnauthorized(message string, opts ...Option) *Error {
	return New(CodeUnauthorized, message, opts...)
}

func NewForbidden(message string, opts ...Option) *Error {
	return New(CodeForbidden, message, opts...)
}

func NewNotFound(message string, opts ...Option) *Error {
	return New(CodeNotFound, message, opts...)
}

func NewConflict(message string, opts ...Option) *Error {
	return New(CodeConflict, message, opts...)
}

func NewTooManyRequests(message string, opts ...Option) *Error {
	return New(CodeTooManyRequests, message, opts...)
}

func NewInternal(message string, opts ...Option) *Error {
	return New(CodeInternal, message, opts...)
}
