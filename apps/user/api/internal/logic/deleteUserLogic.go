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

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserRequest) (resp *types.DeleteUserResponse, err error) {
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

	// 调用RPC服务删除用户
	_, err = l.svcCtx.UserRpc.DeleteUser(rpcCtx, &rpc.DeleteUserRequest{
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
	resp = &types.DeleteUserResponse{}

	logx.Infow("删除用户成功",
		logx.Field("userId", userID))

	return resp, nil
}
