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
	"github.com/clin211/miniblog-v3/pkg/known"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserRequest) (resp *types.GetUserResponse, err error) {
	// 从context中获取用户ID（由中间件设置）
	userID, ok := l.ctx.Value(known.XUserID).(string)
	if !ok {
		logx.Errorw("从context中获取用户ID失败")
		return nil, errorx.ErrTokenInvalid
	}

	// 从context中获取原始token
	token, ok := l.ctx.Value("auth_token").(string)
	if !ok {
		logx.Errorw("从context中获取token失败")
		return nil, errorx.ErrTokenInvalid
	}

	// 创建带token的gRPC上下文
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	rpcCtx := metadata.NewOutgoingContext(l.ctx, md)

	// 调用RPC服务获取用户信息
	rpcResp, err := l.svcCtx.UserRpc.GetUser(rpcCtx, &rpc.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		logx.Errorw("调用RPC服务失败",
			logx.Field("userId", userID),
			logx.Field("error", err))
		// 将 gRPC 错误转换为 errorx 错误
		return nil, errorx.FromGRPCError(err)
	}

	// 构建响应
	resp = &types.GetUserResponse{
		UserId:    rpcResp.UserId,
		Username:  rpcResp.Username,
		Email:     rpcResp.Email,
		Phone:     rpcResp.Phone,
		Age:       int(rpcResp.Age),
		Gender:    int(rpcResp.Gender),
		Avatar:    rpcResp.Avatar,
		Status:    int(rpcResp.Status),
		CreatedAt: rpcResp.CreatedAt,
		UpdatedAt: rpcResp.UpdatedAt,
	}

	logx.Infow("获取用户信息成功",
		logx.Field("userId", userID),
		logx.Field("username", rpcResp.Username),
	)

	return resp, nil
}
