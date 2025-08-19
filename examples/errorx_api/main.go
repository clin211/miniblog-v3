package main

import (
	"fmt"
	"net/http"

	"github.com/clin211/miniblog-v3/pkg/errorx"
	"github.com/clin211/miniblog-v3/pkg/response"
	"github.com/clin211/miniblog-v3/pkg/validate"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// LoginRequest 演示请求
type LoginRequest struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required,length(6|32)"`
}

// LoginResponse 演示响应
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func main() {
	server := rest.MustNewServer(rest.RestConf{Port: 8291})
	defer server.Stop()

	// 与项目 validate 适配，使 httpx.Parse 自动进行标签校验
	httpx.SetValidator(validate.HttpxValidatorAdapter{})

	// 1) 登录接口：成功与参数校验失败演示
	server.AddRoutes([]rest.Route{{
		Method: http.MethodPost,
		Path:   "/api/login",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req LoginRequest
			if err := httpx.Parse(r, &req); err != nil {
				// 参数校验失败，使用统一入口自动处理
				response.WriteResponse(r.Context(), w, err)
				return
			}

			// 模拟业务：当邮箱为 notfound@example.com 时，返回用户模块的业务错误
			if req.Email == "notfound@example.com" {
				err := &errorx.Errno{
					HTTP:    404,
					Code:    10001,
					Message: "用户不存在",
					Data:    nil,
					Reason:  fmt.Sprintf("邮箱 %s 未注册", req.Email),
				}
				response.WriteResponse(r.Context(), w, err)
				return
			}

			// 成功返回
			resp := &LoginResponse{AccessToken: "at-demo", RefreshToken: "rt-demo"}
			response.WriteResponse(r.Context(), w, resp)
		},
	}})

	// 2) 用户查询接口：演示模块错误与普通 error 的差异
	server.AddRoutes([]rest.Route{{
		Method: http.MethodGet,
		Path:   "/api/users/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var path struct {
				ID string `path:"id"`
			}
			_ = httpx.Parse(r, &path)
			id := path.ID
			if id == "0" {
				// 返回模块错误：user:10002
				err := &errorx.Errno{
					HTTP:    400,
					Code:    10002,
					Message: "用户ID非法",
					Data:    nil,
					Reason:  fmt.Sprintf("ID %s 格式不正确", id),
				}
				response.WriteResponse(r.Context(), w, err)
				return
			}
			if id == "404" {
				// 返回普通 error，将按统一入口渲染为 Not Found
				response.WriteResponse(r.Context(), w, fmt.Errorf("record not found: %s", id))
				return
			}
			// 正常
			response.WriteResponse(r.Context(), w, map[string]any{"id": id, "name": "demo"})
		},
	}})

	// 3) 直接成功对象演示：/api/ok
	server.AddRoutes([]rest.Route{{
		Method: http.MethodGet,
		Path:   "/api/ok",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			response.WriteResponse(r.Context(), w, map[string]any{"hello": "world"})
		},
	}})

	// 4) 自定义成功响应演示：/api/custom
	server.AddRoutes([]rest.Route{{
		Method: http.MethodGet,
		Path:   "/api/custom",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			response.SuccessCtx(r.Context(), w, "自定义成功消息",
				map[string]any{"custom": "data"},
				"自定义备注信息")
		},
	}})

	// 5) 自定义失败响应演示：/api/fail
	server.AddRoutes([]rest.Route{{
		Method: http.MethodGet,
		Path:   "/api/fail",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			response.FailCtx(r.Context(), w, 10003, "业务失败", "详细失败原因")
		},
	}})

	server.Start()
}
