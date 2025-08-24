// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package svc

import (
	"github.com/clin211/miniblog-v3/apps/user/api/internal/config"
	"github.com/clin211/miniblog-v3/apps/user/rpc/pb/rpc"
	"github.com/clin211/miniblog-v3/pkg/middleware"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config          config.Config
	UserRpc         rpc.UserClient
	AuthnMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:          c,
		UserRpc:         rpc.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn()),
		AuthnMiddleware: middleware.NewAuthnMiddleware().Handle,
	}
}
