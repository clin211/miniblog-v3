package errorx

// Code 表示业务错误码。
//
// 约定：
// 1) 成功使用 CodeOK=0；
// 2) 参数校验失败在渲染层按现有文档保持 code=0；
// 3) 资源未找到使用 CodeNotFound=404；
// 4) 其他通用错误提供常见的 HTTP 语义对应码，便于统一映射与检索。
type Code int

const (
	// CodeOK 代表业务成功。
	CodeOK Code = 0

	// 通用错误语义（与 HTTP 含义一致，便于跨传输层映射）。
	CodeBadRequest      Code = 400
	CodeUnauthorized    Code = 401
	CodeForbidden       Code = 403
	CodeNotFound        Code = 404
	CodeConflict        Code = 409
	CodeTooManyRequests Code = 429
	CodeInternal        Code = 500
)
