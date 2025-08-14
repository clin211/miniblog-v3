package logic

import (
	"context"

	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/svc"
	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RpcLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRpcLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RpcLogic {
	return &RpcLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RpcLogic) Rpc(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
