package main

import (
	"fmt"
	"net/http"

	"github.com/clin211/miniblog-v3/pkg/errorx"
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
	// 注册模块码段：user 10000-19999
	_ = errorx.RegisterModule("user", errorx.Range{Start: 10000, End: 19999})

	server := rest.MustNewServer(rest.RestConf{Port: 8291})
	defer server.Stop()

	// 与项目 validate 适配，使 httpx.Parse 自动进行标签校验
	httpx.SetValidator(validate.HttpxValidatorAdapter{})

	// 1) 登录接口：成功与参数校验失败演示（校验失败 -> code=0, message="参数校验失败"）
	server.AddRoutes([]rest.Route{{
		Method: http.MethodPost,
		Path:   "/api/login",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req LoginRequest
			if err := httpx.Parse(r, &req); err != nil {
				httpx.OkJsonCtx(r.Context(), w, errorx.Render(err))
				return
			}

			// 模拟业务：当邮箱为 notfound@example.com 时，返回用户模块的业务错误
			if req.Email == "notfound@example.com" {
				err := errorx.NewModuleError("user", 10001, "用户不存在", errorx.WithReason(map[string]any{"email": req.Email}))
				httpx.OkJsonCtx(r.Context(), w, errorx.Render(err))
				return
			}

			// 成功返回
			resp := &LoginResponse{AccessToken: "at-demo", RefreshToken: "rt-demo"}
			httpx.OkJsonCtx(r.Context(), w, errorx.Render(resp))
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
				httpx.OkJsonCtx(r.Context(), w, errorx.Render(
					errorx.NewModuleError("user", 10002, "用户ID非法", errorx.WithReason(map[string]any{"id": id})),
				))
				return
			}
			if id == "404" {
				// 返回普通 error，将按 docs/03 的示例渲染为 Not Found
				httpx.OkJsonCtx(r.Context(), w, errorx.Render(fmt.Errorf("record not found: %s", id)))
				return
			}
			// 正常
			httpx.OkJsonCtx(r.Context(), w, errorx.Render(map[string]any{"id": id, "name": "demo"}))
		},
	}})

	// 3) 直接成功对象演示：/api/ok
	server.AddRoutes([]rest.Route{{
		Method: http.MethodGet,
		Path:   "/api/ok",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			httpx.OkJsonCtx(r.Context(), w, errorx.Render(map[string]any{"hello": "world"}))
		},
	}})

	server.Start()
}
