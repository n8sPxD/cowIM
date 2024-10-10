package logic

import (
	"context"
	"errors"

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

func (l *AddFriendsLogic) AddFriends(req *types.AddFriendRequest) (*types.AddFriendResponse, error) {
	// 先判定是不是已经为好友
	ok, err := l.svcCtx.MySQL.GetFriend(l.ctx, req.UserID, req.FriendID)
	if err != nil {
		logx.Error("[AddFriends] Get friend failed, error: err")
		return nil, errors.New("添加好友失败，可能是服务器出了问题")
	}
	if ok {
		// 已经是好友
		return nil, errors.New("你们已经是好友了")
	}
	// 直接调用sql方法
	err = l.svcCtx.MySQL.InsertFriend(l.ctx, req.UserID, req.FriendID)
	if err != nil {
		logx.Error("[AddFriends] Add friend failed, error: ", err)
		return nil, errors.New("添加好友失败，可能是服务器出了问题")
	}
	return nil, nil
}
