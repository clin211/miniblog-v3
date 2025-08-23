// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package middleware

import (
	"context"
	"net/http"

	"github.com/clin211/miniblog-v3/pkg/errorx"
	"github.com/clin211/miniblog-v3/pkg/known"
	"github.com/clin211/miniblog-v3/pkg/response"
	"github.com/clin211/miniblog-v3/pkg/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserIDKey 是存储在上下文中的用户ID键（已废弃，使用 known.XUserID）
const UserIDKey = "user_id"

// AuthnMiddleware 认证中间件结构体
type AuthnMiddleware struct {
}

// NewAuthnMiddleware 创建认证中间件实例
func NewAuthnMiddleware() *AuthnMiddleware {
	return &AuthnMiddleware{}
}

// Handle HTTP认证中间件处理方法
// 从请求头中解析JWT token，验证用户身份，并将用户ID存储到上下文中
func (m *AuthnMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析JWT token
		claims, err := token.ParseRequest(r)
		if err != nil {
			response.WriteResponse(r.Context(), w, errorx.ErrTokenInvalid)
			return
		}

		// 从请求头中提取原始token
		authHeader := r.Header.Get("Authorization")
		var originalToken string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			originalToken = authHeader[7:]
		}

		// 将用户ID和原始token存储到上下文中
		ctx := context.WithValue(r.Context(), known.XUserID, claims.UserID)
		if originalToken != "" {
			ctx = context.WithValue(ctx, "auth_token", originalToken)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// AuthnMiddlewareFunc HTTP认证中间件函数版本
// 从请求头中解析JWT token，验证用户身份，并将用户ID存储到上下文中
func AuthnMiddlewareFunc() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 解析JWT token
			claims, err := token.ParseRequest(r)
			if err != nil {
				response.WriteResponse(r.Context(), w, errorx.ErrTokenInvalid)
				return
			}

			// 将用户ID存储到上下文中
			ctx := context.WithValue(r.Context(), known.XUserID, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// AuthnInterceptor gRPC认证拦截器
// 从gRPC元数据中解析JWT token，验证用户身份，并将用户ID存储到上下文中
func AuthnInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 定义不需要认证的方法
		noAuthMethods := map[string]bool{
			"/rpc.User/Login":    true,
			"/rpc.User/Register": true,
		}

		// 检查当前方法是否需要认证
		if noAuthMethods[info.FullMethod] {
			// 不需要认证的方法直接执行
			return handler(ctx, req)
		}

		// 需要认证的方法解析JWT token
		claims, err := token.ParseRequest(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}

		// 将用户ID存储到上下文中
		ctx = context.WithValue(ctx, known.XUserID, claims.UserID)
		return handler(ctx, req)
	}
}
