package logic

import (
	"context"

	"github.com/clin211/miniblog-v3/apps/user/api/internal/svc"
	"github.com/clin211/miniblog-v3/apps/user/api/internal/types"
	"github.com/clin211/miniblog-v3/apps/user/rpc/pb/rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// 调用 RPC 服务进行用户注册
	rpcResp, err := l.svcCtx.UserRpc.Register(l.ctx, &rpc.RegisterRequest{
		Username:       req.Username,
		Password:       req.Password,
		Email:          req.Email,
		Phone:          req.Phone,
		Age:            int32(req.Age),
		Gender:         int32(req.Gender),
		Avatar:         req.Avatar,
		RegisterSource: int32(req.RegisterSource),
		WechatOpenid:   req.WechatOpenid,
	})
	if err != nil {
		return nil, err
	}

	return &types.RegisterResponse{
		UserId: rpcResp.UserId,
	}, nil
}
