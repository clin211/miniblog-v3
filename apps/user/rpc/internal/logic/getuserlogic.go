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

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUser 获取用户信息
func (l *GetUserLogic) GetUser(in *rpc.GetUserRequest) (*rpc.GetUserResponse, error) {
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
		return nil, errorx.ToGRPCError(errorx.ErrUnauthorized.SetMessage("只能查看自己的用户信息"))
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

	// 构建响应
	resp := &rpc.GetUserResponse{
		UserId:    user.UserId,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Age:       int32(user.Age),
		Gender:    int32(user.Gender),
		Avatar:    user.Avatar,
		Status:    int32(user.Status),
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	l.Infow("获取用户信息成功",
		logx.Field("userId", in.UserId),
		logx.Field("username", user.Username))

	return resp, nil
}
