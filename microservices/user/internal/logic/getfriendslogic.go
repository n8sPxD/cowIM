package logic

import (
	"context"

	"github.com/n8sPxD/cowIM/microservices/user/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendsLogic {
	return &GetFriendsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendsLogic) GetFriends(req *types.GetFriendsRequest) (resp *types.GetFriendsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
