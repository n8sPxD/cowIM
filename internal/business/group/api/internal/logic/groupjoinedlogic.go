package logic

import (
	"context"
	"errors"

	"github.com/n8sPxD/cowIM/internal/business/group/api/internal/svc"
	"github.com/n8sPxD/cowIM/internal/business/group/api/internal/types"
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

func (l *GroupJoinedLogic) GroupJoined(req *types.GroupJoinedRequest) (*types.GroupJoinedResponse, error) {
	// 直接查数据库
	groups, err := l.svcCtx.MySQL.GetGroupsJoinedBaseInfo(l.ctx, req.UserID)
	if err != nil {
		logx.Error("[GroupJoined] Get groups from MySQL failed, error: ", err)
		return nil, errors.New("获取群组信息失败")
	}
	infos := make([]types.GroupJoinedInfo, len(groups))
	for i, group := range groups {
		infos[i].GroupID = uint32(group.ID)
		infos[i].GroupName = group.GroupName
		infos[i].GroupAvatar = group.GroupAvatar
	}
	return &types.GroupJoinedResponse{Infos: infos}, nil
}
