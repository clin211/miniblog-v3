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
		// 处理业务错误
		if e, ok := err.(*errorx.Errno); ok {
			return nil, e
		}
		// 处理其他错误
		logx.Errorf("调用登录RPC失败: %v", err)
		return nil, errorx.InternalServerError.SetMessage("登录失败")
	}

	// 2. 构造响应
	return &types.LoginResponse{
		Token:     rpcResp.Token,
		ExpiresAt: rpcResp.ExpireAt,
	}, nil
}
