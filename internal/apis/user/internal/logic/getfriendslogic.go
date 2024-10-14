package logic

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/internal/apis/user/internal/svc"
	"github.com/n8sPxD/cowIM/internal/apis/user/internal/types"
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

func (l *GetFriendsLogic) GetFriends(req *types.GetFriendsRequest) (*types.GetFriendsResponse, error) {
	// 直接调mysql方法
	friends, err := l.svcCtx.MySQL.GetFriends(l.ctx, req.UserID)
	if err != nil {
		logx.Error("[GetFriends] Get friends failed, error: ", err)
		return nil, errors.New("获取好友列表失败")
	}

	rets := make([]types.FriendInfo, len(friends))
	for i, friend := range friends {
		rets[i].FriendID = uint32(friend.ID)
		rets[i].Username = friend.Username
		rets[i].Avatar = friend.Avatar
	}

	resp := &types.GetFriendsResponse{Friends: rets}

	return resp, nil
}
