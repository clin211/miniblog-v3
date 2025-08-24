// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package main

import (
	"net/http"

	"github.com/clin211/miniblog-v3/pkg/errorx"
	"github.com/clin211/miniblog-v3/pkg/response"
	"github.com/clin211/miniblog-v3/pkg/validate"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// DemoRequest 示例请求体
type DemoRequest struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required,length(6|32)"`
}

// DemoResponse 示例响应体
type DemoResponse struct {
	Message string `json:"message"`
}

func main() {
	server := rest.MustNewServer(rest.RestConf{Port: 8188})
	defer server.Stop()

	// 注入参数校验（可选，仅用于与项目 validate 适配）
	httpx.SetValidator(validate.HttpxValidatorAdapter{})

	// 成功示例：/ok
	server.AddRoutes([]rest.Route{{
		Method: http.MethodGet,
		Path:   "/ok",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			resp := &DemoResponse{Message: "hello"}
			response.SuccessCtx(r.Context(), w, "response message", resp, "demo-trace")
		},
	}})

	// 自定义失败示例：/fail
	server.AddRoutes([]rest.Route{{
		Method: http.MethodGet,
		Path:   "/fail",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			response.FailCtx(r.Context(), w, 10001, "业务失败", "详细失败原因")
		},
	}})

	// 参数校验示例：/validate 以 POST body 触发
	server.AddRoutes([]rest.Route{{
		Method: http.MethodPost,
		Path:   "/validate",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req DemoRequest
			if err := httpx.Parse(r, &req); err != nil {
				response.WriteResponse(r.Context(), w, err)
				return
			}
			// 正常业务
			response.WriteResponse(r.Context(), w, nil)
		},
	}})

	// 直接透传对象/错误示例：/auto
	server.AddRoutes([]rest.Route{{
		Method: http.MethodGet,
		Path:   "/auto",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			// 示例：构造一个校验错误（演示自动转换为 reason 数组）
			verr := validate.ValidationErrors{
				&validate.ValidationError{Field: "Email", Message: "邮箱格式不正确", Value: "final-direct-testexample.com"},
			}
			response.WriteResponse(r.Context(), w, verr)
		},
	}})

	// 业务错误示例：/business-error
	server.AddRoutes([]rest.Route{{
		Method: http.MethodGet,
		Path:   "/business-error",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			// 构造业务错误
			err := &errorx.Errno{
				HTTP:    400,
				Code:    10002,
				Message: "用户余额不足",
				Data:    map[string]interface{}{"balance": 50, "required": 100},
				Reason:  "当前余额50元，需要100元",
			}
			response.WriteResponse(r.Context(), w, err)
		},
	}})

	// 普通错误示例：/normal-error
	server.AddRoutes([]rest.Route{{
		Method: http.MethodGet,
		Path:   "/normal-error",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			// 普通 error 会被转换为 404 Not Found
			response.WriteResponse(r.Context(), w, errorx.ErrResourceNotFound)
		},
	}})

	server.Start()
}
