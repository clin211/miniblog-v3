package main

import (
	"net/http"

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
			response.SuccessCtx(r.Context(), w, "response message", resp, map[string]any{"traceId": "demo-trace"})
		},
	}})

	// 未找到示例：/notfound
	server.AddRoutes([]rest.Route{{
		Method: http.MethodGet,
		Path:   "/notfound",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			response.NotFoundCtx(r.Context(), w, "resource not found")
		},
	}})

	// 参数校验示例：/validate 以 POST body 触发
	server.AddRoutes([]rest.Route{{
		Method: http.MethodPost,
		Path:   "/validate",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req DemoRequest
			if err := httpx.Parse(r, &req); err != nil {
				response.JsonBaseResponseCtx(r.Context(), w, err)
				return
			}
			// 正常业务
			response.JsonBaseResponseCtx(r.Context(), w, &DemoResponse{Message: "ok"})
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
			response.JsonBaseResponseCtx(r.Context(), w, verr)
		},
	}})

	server.Start()
}
