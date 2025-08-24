// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package logic

import (
	"context"

	"github.com/clin211/miniblog-v3/apps/user/models"
	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/svc"
	"github.com/clin211/miniblog-v3/apps/user/rpc/pb/rpc"
	"github.com/clin211/miniblog-v3/pkg/errorx"
	"github.com/clin211/miniblog-v3/pkg/known"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteUser 删除用户
func (l *DeleteUserLogic) DeleteUser(in *rpc.DeleteUserRequest) (*rpc.DeleteUserResponse, error) {
	// 从context中获取用户ID（由拦截器设置）
	userID, ok := l.ctx.Value(known.XUserID).(string)
	if !ok {
		l.Errorw("从context中获取用户ID失败")
		return nil, errorx.ToGRPCError(errorx.ErrTokenInvalid)
	}

	// 验证用户权限
	if userID != in.UserId {
		l.Errorw("用户权限不足",
			logx.Field("currentUserID", userID),
			logx.Field("requestUserID", in.UserId))
		return nil, errorx.ToGRPCError(errorx.ErrUnauthorized.SetMessage("只能删除自己的账户"))
	}

	// 从数据库查询用户信息
	user, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		if err == models.ErrNotFound {
			l.Errorw("用户不存在", logx.Field("userId", in.UserId))
			return nil, errorx.ToGRPCError(errorx.ErrUserNotFound)
		}
		l.Errorw("查询用户信息失败",
			logx.Field("userId", in.UserId),
			logx.Field("error", err))
		return nil, errorx.ToGRPCError(errorx.InternalServerError.SetMessage("查询用户信息失败"))
	}

	// 检查用户状态
	if user.Status == 0 {
		l.Errorw("用户已被禁用", logx.Field("userId", in.UserId))
		return nil, errorx.ToGRPCError(errorx.ErrUserDisabled)
	}

	// 软删除用户（设置删除时间）
	// 注意：这里我们使用软删除，只是标记删除时间，不真正删除数据
	// 如果需要硬删除，可以使用 l.svcCtx.UserModel.Delete(l.ctx, user.Id)

	// 更新用户状态为禁用
	user.Status = 0
	err = l.svcCtx.UserModel.Update(l.ctx, user)
	if err != nil {
		l.Errorw("删除用户失败",
			logx.Field("userId", in.UserId),
			logx.Field("error", err))
		return nil, errorx.ToGRPCError(errorx.InternalServerError.SetMessage("删除用户失败"))
	}

	// 构建响应
	resp := &rpc.DeleteUserResponse{
		Success: true,
	}

	l.Infow("删除用户成功",
		logx.Field("userId", in.UserId),
		logx.Field("username", user.Username))

	return resp, nil
}
