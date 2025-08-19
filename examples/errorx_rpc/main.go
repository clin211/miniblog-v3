package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clin211/miniblog-v3/pkg/errorx"
	"github.com/clin211/miniblog-v3/pkg/response"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 说明：
// 本示例模拟"RPC 网关"风格的路由，展示在服务到服务调用中如何将内部错误（模块错误/普通错误）
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
		return nil, &errorx.Errno{
			HTTP:    400,
			Code:    10002,
			Message: "用户ID非法",
			Data:    nil,
			Reason:  fmt.Sprintf("ID %s 格式不正确", id),
		}
	}
	if id == "404" {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	return &QueryUserResponse{ID: id, Name: "rpc-user"}, nil
}

func main() {
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

			// 使用统一入口处理成功或失败
			if err != nil {
				response.WriteResponse(r.Context(), w, err)
			} else {
				response.WriteResponse(r.Context(), w, resp)
			}
		},
	}})

	// 演示批量查询接口
	server.AddRoutes([]rest.Route{{
		Method: http.MethodPost,
		Path:   "/rpc/users/batch",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				IDs []string `json:"ids"`
			}
			if err := httpx.Parse(r, &req); err != nil {
				response.WriteResponse(r.Context(), w, err)
				return
			}

			// 模拟批量查询，部分成功部分失败
			var results []map[string]interface{}
			var hasError bool

			for _, id := range req.IDs {
				resp, err := callUserService(r.Context(), id)
				if err != nil {
					hasError = true
					results = append(results, map[string]interface{}{
						"id":    id,
						"error": err.Error(),
					})
				} else {
					results = append(results, map[string]interface{}{
						"id":   resp.ID,
						"name": resp.Name,
					})
				}
			}

			if hasError {
				// 部分失败的情况
				err := &errorx.Errno{
					HTTP:    207, // Multi-Status
					Code:    10004,
					Message: "批量查询部分失败",
					Data:    results,
					Reason:  "部分用户ID不存在或格式错误",
				}
				response.WriteResponse(r.Context(), w, err)
			} else {
				// 全部成功
				response.WriteResponse(r.Context(), w, results)
			}
		},
	}})

	server.Start()
}
