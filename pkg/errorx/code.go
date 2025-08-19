package errorx

import "net/http"

var (
	// OK 代表请求成功.
	OK = &Errno{HTTP: http.StatusOK, Code: 0, Message: "Success.", Data: nil, Reason: ""}

	// InternalServerError 表示所有未知的服务器端错误.
	InternalServerError = &Errno{HTTP: http.StatusInternalServerError, Code: 500001, Message: "Internal server error.", Data: nil, Reason: ""}

	// ErrPageNotFound 表示路由不匹配错误.
	ErrPageNotFound = &Errno{HTTP: http.StatusNotFound, Code: 404001, Message: "Page not found.", Data: nil, Reason: ""}

	// ErrBind 表示参数绑定错误.
	ErrBind = &Errno{HTTP: http.StatusBadRequest, Code: 400001, Message: "Error occurred while binding the request body to the struct.", Data: nil, Reason: ""}

	// ErrInvalidParameter 表示所有验证失败的错误.
	ErrInvalidParameter = &Errno{HTTP: http.StatusBadRequest, Code: 400002, Message: "Parameter verification failed.", Data: nil, Reason: ""}

	// ErrSignToken 表示签发 JWT Token 时出错.
	ErrSignToken = &Errno{HTTP: http.StatusUnauthorized, Code: 401001, Message: "Error occurred while signing the JSON web token.", Data: nil, Reason: ""}

	// ErrTokenInvalid 表示 JWT Token 格式错误.
	ErrTokenInvalid = &Errno{HTTP: http.StatusUnauthorized, Code: 401002, Message: "Token was invalid.", Data: nil, Reason: ""}

	// ErrUnauthorized 表示请求没有被授权.
	ErrUnauthorized = &Errno{HTTP: http.StatusUnauthorized, Code: 401003, Message: "Unauthorized.", Data: nil, Reason: ""}
)
