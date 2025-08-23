// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TokenClientInterceptor gRPC客户端拦截器，用于传递token
func TokenClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 从context中获取原始token
		if token, ok := ctx.Value("auth_token").(string); ok && token != "" {
			// 将token添加到gRPC元数据中
			md := metadata.New(map[string]string{
				"authorization": "Bearer " + token,
			})
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
