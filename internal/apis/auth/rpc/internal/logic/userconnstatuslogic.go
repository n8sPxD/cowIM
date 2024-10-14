package logic

import (
	"context"
	"time"

	"github.com/n8sPxD/cowIM/internal/apis/auth/rpc/internal/svc"
	"github.com/n8sPxD/cowIM/internal/apis/auth/rpc/types/authRpc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserConnStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserConnStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserConnStatusLogic {
	return &UserConnStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserConnStatusLogic) UserConnStatus(in *authRpc.ConnRequest) (*authRpc.Empty, error) {
	err := l.svcCtx.Redis.UpdateUserRouterStatus(l.ctx, in.UserId, in.WorkId, time.Now())
	if err != nil {
		logx.Error("[UserConnStatus] Update user router status from redis failed, error: ", err)
	}
	return &authRpc.Empty{}, err
}
