// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/clin211/miniblog-v3/pkg/known"
	"github.com/clin211/miniblog-v3/pkg/token"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestAuthnMiddleware(t *testing.T) {
	// 初始化token包
	config := token.Config{
		Secret:      "test-secret",
		IdentityKey: "user_id",
		Expiration:  1 * time.Hour,
		Issuer:      "test",
		Audience:    "test_users",
	}
	token.Init(config)

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedStatus int
		expectedUserID string
	}{
		{
			name: "valid token",
			setupRequest: func() *http.Request {
				// 签发有效token
				tokenString, _, _ := token.Sign("user_123")
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer "+tokenString)
				return req
			},
			expectedStatus: http.StatusOK,
			expectedUserID: "user_123",
		},
		{
			name: "invalid token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer invalid-token")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: "",
		},
		{
			name: "no authorization header",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试处理器
			var capturedUserID string
			handler := func(w http.ResponseWriter, r *http.Request) {
				// 直接使用token.ParseRequest获取用户ID
				if claims, err := token.ParseRequest(r); err == nil {
					capturedUserID = claims.UserID
				}
				w.WriteHeader(http.StatusOK)
			}

			// 应用中间件
			middleware := AuthnMiddlewareFunc()
			wrappedHandler := middleware(handler)

			// 创建响应记录器
			w := httptest.NewRecorder()
			req := tt.setupRequest()

			// 执行请求
			wrappedHandler.ServeHTTP(w, req)

			// 验证结果
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedUserID != "" {
				assert.Equal(t, tt.expectedUserID, capturedUserID)
			}
		})
	}
}

func TestAuthnInterceptor(t *testing.T) {
	// 初始化token包
	config := token.Config{
		Secret:      "test-secret",
		IdentityKey: "user_id",
		Expiration:  1 * time.Hour,
		Issuer:      "test",
		Audience:    "test_users",
	}
	token.Init(config)

	tests := []struct {
		name           string
		fullMethod     string
		setupContext   func() context.Context
		expectedError  bool
		expectedUserID string
	}{
		{
			name:       "login method - no auth required",
			fullMethod: "/rpc.User/Login",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedError:  false,
			expectedUserID: "",
		},
		{
			name:       "register method - no auth required",
			fullMethod: "/rpc.User/Register",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedError:  false,
			expectedUserID: "",
		},
		{
			name:       "get user method - auth required with valid token",
			fullMethod: "/rpc.User/GetUser",
			setupContext: func() context.Context {
				// 签发有效token
				tokenString, _, _ := token.Sign("user_123")
				ctx := context.Background()
				// 模拟gRPC元数据
				md := metadata.New(map[string]string{
					"authorization": "Bearer " + tokenString,
				})
				return metadata.NewIncomingContext(ctx, md)
			},
			expectedError:  false,
			expectedUserID: "user_123",
		},
		{
			name:       "get user method - auth required with invalid token",
			fullMethod: "/rpc.User/GetUser",
			setupContext: func() context.Context {
				ctx := context.Background()
				// 模拟gRPC元数据
				md := metadata.New(map[string]string{
					"authorization": "Bearer invalid-token",
				})
				return metadata.NewIncomingContext(ctx, md)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建拦截器
			interceptor := AuthnInterceptor()

			// 创建模拟的handler
			var capturedUserID string
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				if userID, ok := ctx.Value(known.XUserID).(string); ok {
					capturedUserID = userID
				}
				return "success", nil
			}

			// 创建模拟的gRPC信息
			info := &grpc.UnaryServerInfo{
				FullMethod: tt.fullMethod,
			}

			// 执行拦截器
			ctx := tt.setupContext()
			_, err := interceptor(ctx, "test-request", info, handler)

			// 验证结果
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectedUserID != "" {
					assert.Equal(t, tt.expectedUserID, capturedUserID)
				}
			}
		})
	}
}
