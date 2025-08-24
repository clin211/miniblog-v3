// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package logic

import (
	"context"

	"github.com/clin211/miniblog-v3/apps/user/api/internal/svc"
	"github.com/clin211/miniblog-v3/apps/user/api/internal/types"
	"github.com/clin211/miniblog-v3/apps/user/rpc/pb/rpc"
	"github.com/clin211/miniblog-v3/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 1. 调用 RPC 服务进行登录
	rpcResp, err := l.svcCtx.UserRpc.Login(l.ctx, &rpc.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		// 将 gRPC 错误转换为 errorx 错误
		return nil, errorx.FromGRPCError(err)
	}

	// 2. 构造响应
	return &types.LoginResponse{
		Token:    rpcResp.Token,
		ExpireAt: rpcResp.ExpireAt,
	}, nil
}
