package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clin211/miniblog-v3/pkg/errorx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 说明：
// 本示例模拟“RPC 网关”风格的路由，展示在服务到服务调用中如何将内部错误（模块错误/普通错误）
// 统一翻译为标准响应结构，供上游服务或客户端消费。这里仍以 go-zero REST 方式承载，以便在项目中直接运行。

type QueryUserRequest struct {
	ID string `path:"id"`
}

type QueryUserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// 模拟下游 RPC 调用
func callUserService(ctx context.Context, id string) (*QueryUserResponse, error) {
	if id == "0" {
		return nil, errorx.NewModuleError("user", 10002, "用户ID非法", errorx.WithReason(map[string]any{"id": id}))
	}
	if id == "404" {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	return &QueryUserResponse{ID: id, Name: "rpc-user"}, nil
}

func main() {
	// 注册模块码段
	_ = errorx.RegisterModule("user", errorx.Range{Start: 10000, End: 19999})

	server := rest.MustNewServer(rest.RestConf{Port: 8292})
	defer server.Stop()

	server.AddRoutes([]rest.Route{{
		Method: http.MethodGet,
		Path:   "/rpc/users/:id",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var path struct {
				ID string `path:"id"`
			}
			_ = httpx.Parse(r, &path)
			id := path.ID
			resp, err := callUserService(r.Context(), id)
			httpx.OkJsonCtx(r.Context(), w, errorx.Render(func() any {
				if err != nil {
					return err
				}
				return resp
			}()))
		},
	}})

	server.Start()
}
