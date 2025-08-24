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

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserRequest) (resp *types.UpdateUserResponse, err error) {
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

	// 调用RPC服务更新用户信息
	rpcResp, err := l.svcCtx.UserRpc.UpdateUser(rpcCtx, &rpc.UpdateUserRequest{
		UserId:   userID,
		Username: req.Username,
		Age:      int32(req.Age),
		Gender:   int32(req.Gender),
		Avatar:   req.Avatar,
	})
	if err != nil {
		logx.Errorw("调用RPC服务失败",
			logx.Field("userId", userID),
			logx.Field("error", err))
		// 将 gRPC 错误转换为 errorx 错误
		return nil, errorx.FromGRPCError(err)
	}

	// 构建响应
	resp = &types.UpdateUserResponse{
		UserId:    rpcResp.UserId,
		Username:  req.Username,
		UpdatedAt: "", // RPC响应中没有返回更新时间，这里暂时留空
	}

	logx.Infow("更新用户信息成功",
		logx.Field("userId", userID),
		logx.Field("username", req.Username))

	return resp, nil
}
