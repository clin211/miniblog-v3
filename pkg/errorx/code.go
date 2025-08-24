// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package errorx

import "net/http"

// Code 定义了错误码的常量；前三位为 HTTP 状态码, 后三位为业务模块编号
var (
	// OK 代表请求成功.
	OK = &Errno{HTTP: http.StatusOK, Code: 0, Message: "Success.", Data: nil, Reason: ""}

	// InternalServerError 表示所有未知的服务器端错误.
	InternalServerError = &Errno{HTTP: http.StatusInternalServerError, Code: 500001, Message: "Internal server error.", Data: nil, Reason: ""}

	// ErrResourceNotFound 表示资源不存在.
	ErrResourceNotFound = &Errno{HTTP: http.StatusNotFound, Code: 404001, Message: "Resource not found.", Data: nil, Reason: ""}

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

// 用户模块 code 段的后三位区间为 100~199
var (
	// ErrUserAlreadyExists 表示用户已存在.
	ErrUserAlreadyExists = &Errno{HTTP: http.StatusConflict, Code: 409101, Message: "User already exists.", Data: nil, Reason: ""}

	// ErrUserNotFound 表示用户不存在.
	ErrUserNotFound = &Errno{HTTP: http.StatusNotFound, Code: 404102, Message: "User not found.", Data: nil, Reason: ""}

	// ErrPasswordIncorrect 表示密码错误.
	ErrPasswordIncorrect = &Errno{HTTP: http.StatusUnauthorized, Code: 401103, Message: "Password incorrect.", Data: nil, Reason: ""}

	// ErrUserDisabled 表示用户被禁用.
	ErrUserDisabled = &Errno{HTTP: http.StatusForbidden, Code: 403104, Message: "User is disabled.", Data: nil, Reason: ""}
)
