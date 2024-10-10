package logic

import (
	"context"

	"github.com/n8sPxD/cowIM/microservices/user/internal/svc"
	"github.com/n8sPxD/cowIM/microservices/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFriendsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFriendsLogic {
	return &AddFriendsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFriendsLogic) AddFriends(req *types.AddFriendRequest) (resp *types.AddFriendResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
