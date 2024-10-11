package logic

import (
	"context"

	"github.com/n8sPxD/cowIM/microservices/group/api/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/group/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInviteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupInviteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInviteLogic {
	return &GroupInviteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupInviteLogic) GroupInvite(req *types.GroupInviteRequest) (resp *types.GroupInviteResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
