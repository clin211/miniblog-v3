package logic

import (
	"context"

	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/svc"
	"github.com/clin211/miniblog-v3/apps/user/rpc/pb/rpc"

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
	// todo: add your logic here and delete this line

	return &rpc.GetUserResponse{}, nil
}
