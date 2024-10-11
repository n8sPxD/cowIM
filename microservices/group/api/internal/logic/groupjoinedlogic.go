package logic

import (
	"context"

	"github.com/n8sPxD/cowIM/microservices/group/api/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/group/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupJoinedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupJoinedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupJoinedLogic {
	return &GroupJoinedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupJoinedLogic) GroupJoined(req *types.GroupJoinedRequest) (resp *types.GroupJoinedResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
